package bbcode

import (
	"bytes"
	"github.com/moechat/parser"
	"html"
	"html/template"
	"regexp"
	"strconv"
	"strings"
)

var bbCodeRe = regexp.MustCompile("\\[([^\\]|^\\[]*)\\]")

func bbCloseTag(name string) (*regexp.Regexp, error) {
	return regexp.Compile("\\[\\/" + name + "\\]")
}

func bbTag(name string) (*regexp.Regexp, error) {
	return regexp.Compile("\\[" + name + "(=.*)?\\]")
}

type bbTagPair struct {
	htmlTag HtmlTags
	bbName  string
}

type stack struct {
	nodes []*bbTagPair
	count int
}

func (s *stack) push(n *bbTagPair) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

func (s *stack) pop() *bbTagPair {
	if s.count == 0 {
		return nil
	}
	s.count--
	return s.nodes[s.count]
}

func (s *stack) top() *bbTagPair {
	if s.count == 0 {
		return nil
	}
	return s.nodes[s.count-1]
}

func (s *stack) topFromN(n int) *stack {
	return &stack{nodes: s.nodes[n:], count: len(s.nodes) - n}
}

func closeNTags(tagPairStack *stack, n int) string {
	endTags := ""
	for i := 0; i < n; i++ {
		tagPair := tagPairStack.pop()
		endTags += closeTags(tagPair)
	}
	return endTags
}

func closeTags(tagPair *bbTagPair) string {
	endTags := ""
	if tagPair.htmlTag.Tags != nil {
		if tagPair.htmlTag.Options&parser.HtmlSingle == 0 {
			for _, tag := range tagPair.htmlTag.Tags {
				endTags = "</" + tag + ">" + endTags
			}
		}
	}
	return endTags
}

// Parse parses BBCode only.
// Although not used by the main Parse method, it is included in case parsing only BBCode is desired.
// Note that this function completely ignores MoeTags.
func Parse(body string) (string, error) {
	tagStack := &stack{nodes: make([]*bbTagPair, 10)}
	output := ""
	body = html.EscapeString(body)

	for body != "" {
		tagLoc := bbCodeRe.FindIndex([]byte(body))

		if tagLoc == nil {
			output += body
			body = ""
			break
		}

		tagData := strings.SplitN(body[tagLoc[0]+1:tagLoc[1]-1], "=", 2)
		htmlTags, ok := bbCodeTags[tagData[0]]
		cok := false
		if tagData[0][0] == '/' {
			_, cok = bbCodeTags[tagData[0][1:]]
		}

		if ok {
			output += body[:tagLoc[0]]
			body = body[tagLoc[1]:]

			args := make([]string, 2)
			if len(tagData) == 2 {
				args[0] = tagData[1]
			}

			if htmlTags.Options&(parser.TokenBodyAsArg|parser.AllowTokenBodyAsFirstArg) != 0 {
				closeTagRe, err := bbCloseTag(tagData[0])
				if err != nil {
					return "", err
				}
				closeTagLoc := closeTagRe.FindIndex([]byte(body))
				if closeTagLoc == nil {
					if htmlTags.Options&parser.PossibleSingle == 0 {
						if htmlTags.Options&parser.AllowTokenBodyAsFirstArg != 0 && args[0] == "" {
							args[0] = body[tagLoc[1]:]
						}
						if htmlTags.Options&parser.TokenBodyAsArg != 0 {
							args[1] = body[tagLoc[1]:]
							body = ""
						}
					}
				} else {
					tagRe, err := bbTag(tagData[0])
					if err != nil {
						return "", err
					}
					openTagLoc := tagRe.FindIndex([]byte(body))
					if htmlTags.Options&parser.PossibleSingle == 0 || openTagLoc == nil || closeTagLoc[0] < openTagLoc[0] {
						if htmlTags.Options&parser.AllowTokenBodyAsFirstArg != 0 && args[0] == "" {
							args[0] = body[:closeTagLoc[0]]
						}
						if htmlTags.Options&parser.TokenBodyAsArg != 0 {
							args[1] = body[:closeTagLoc[0]]
							body = body[closeTagLoc[1]:]
						}
					}
				}
			}

			if htmlTags.InputModFunc != nil {
				htmlTags.InputModFunc(&args)
			}

			if htmlTags.OutputFunc != nil {
				output += htmlTags.OutputFunc(args)
			} else {
				for i, tag := range htmlTags.Tags {
					templStr := "<" + tag

					if len(htmlTags.Classes) > i {
						if classes := htmlTags.Classes[i]; classes != nil {
							templStr += " class=\""
							for _, class := range classes {
								templStr += " " + class
							}
							templStr += "\""
						}
					}

					if len(htmlTags.Attributes) > i {
						if attrs := htmlTags.Attributes[i]; attrs != nil {
							for i, attr := range attrs {
								if len(args) > int(i) && args[i] != "" {
									templStr += " " + attr + "=\"{{index . " + strconv.Itoa(int(i)) + "}}\""
								}
							}
						}
					}

					if len(htmlTags.CssProps) > i {
						if cssProps := htmlTags.CssProps[i]; cssProps != nil {
							templStr += " style=\""
							for i, cssProp := range cssProps {
								if len(args) > int(i) && args[i] != "" {
									templStr += cssProp + ": {{index . " + strconv.Itoa(int(i)) + "}};"
								}
							}
							templStr += "\""
						}
					}

					templStr += ">"

					tmpl, err := template.New("elementTemplate").Parse(templStr)
					if err != nil {
						return "", err
					}

					eleBuffer := bytes.Buffer{}
					err = tmpl.Execute(&eleBuffer, args)
					if err != nil {
						return "", err
					}
					output += eleBuffer.String()
				}
			}

			tagStack.push(&bbTagPair{htmlTags, tagData[0]})

			if htmlTags.Options&parser.NoParseInner != 0 {
				closeTagRe, err := bbCloseTag(tagData[0])
				if err != nil {
					return "", err
				}
				closeTagLoc := closeTagRe.FindIndex([]byte(body))

				if closeTagLoc == nil {
					output += body
					body = ""
				} else {
					output += body[:closeTagLoc[0]]
					output += closeTags(tagStack.pop())
					body = body[closeTagLoc[1]:]
				}
			} else if htmlTags.Options&parser.PossibleSingle != 0 {
				closeTagRe, err := bbCloseTag(tagData[0])
				if err != nil {
					return "", err
				}
				closeTagLoc := closeTagRe.FindIndex([]byte(body))
				if closeTagLoc != nil {
					tagRe, err := bbTag(tagData[0])
					if err != nil {
						return "", err
					}
					openTagLoc := tagRe.FindIndex([]byte(body))
					if openTagLoc != nil && openTagLoc[0] < closeTagLoc[0] {
						if htmlTags.Options&parser.HtmlSingle == 0 {
							output += closeTags(tagStack.pop())
						}
					}
				} else {
					if htmlTags.Options&parser.HtmlSingle == 0 {
						output += closeTags(tagStack.pop())
					}
				}
			}
		} else if cok {
			tagStackCopy := &stack{nodes: tagStack.nodes, count: tagStack.count}
			foundMatch := false
			for tagPair := tagStackCopy.pop(); tagPair != nil; tagPair = tagStackCopy.pop() {
				if tagPair.bbName == tagData[0][1:] {
					output += body[:tagLoc[0]]
					output += closeNTags(tagStack, tagStack.count-tagStackCopy.count)
					body = body[tagLoc[1]:]
					foundMatch = true
					break
				}
			}

			if !foundMatch {
				output += body[tagLoc[0]:tagLoc[1]]
				body = body[tagLoc[1]:]
			}
		} else {
			output += body[:tagLoc[1]]
			body = body[tagLoc[1]:]
		}
	}

	output += closeNTags(tagStack, tagStack.count)
	return output, nil
}
