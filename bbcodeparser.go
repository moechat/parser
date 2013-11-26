package moeparser

import (
	"regexp"
	"strings"
)

var bbCodeRe = regexp.MustCompile("\\[(.*)\\]")

func BbCloseTag(name string) (*regexp.Regexp, error) {
	return regexp.Compile("\\[\\/"+name+"\\]")
}

type bbTagPair struct {
	htmlTag HtmlTags
	bbName string
}

type Stack struct {
	nodes []*bbTagPair
	count int
}

func (s *Stack) Push(n *bbTagPair) {
	s.nodes = append(s.nodes, n)
	s.count++
}

func (s *Stack) Pop() *bbTagPair {
	if s.count == 0 {
		return nil
	}
	s.count--
	return s.nodes[s.count]
}

func (s *Stack) Top() *bbTagPair {
	if s.count == 0 {
		return nil
	}
	return s.nodes[s.count]
}

func (s *Stack) TopFromN(n int) *Stack {
	return &Stack{nodes: s.nodes[n:], count: len(s.nodes)-n}
}

func closeTags(tagPairStack *Stack) string {
	endTags := ""
	for tagPair := tagPairStack.Pop(); tagPair != nil; tagPair = tagPairStack.Pop() {
		if tagPair.htmlTag.Tags != nil {
			for _, tag := range tagPair.htmlTag.Tags {
				endTags = "<" + tag + ">" + endTags
			}
		}
	}
	return endTags
}

func BbCodeParse(b []byte) ([]byte, error) {
	tagStack := &Stack{nodes: []*bbTagPair{&bbTagPair{}}}
	output := ""
	body := string(b)
	for body != "" {
		tagLoc := bbCodeRe.FindIndex([]byte(body))

		if(tagLoc == nil) {
			output += body
			body = ""
			break
		}

		tagData := strings.SplitN(body[tagLoc[0]+1:tagLoc[1]-1], "=", 2)
		htmlTags, ok := BbCodeTags[tagData[0]]
		cok := false
		if tagData[0][0] == '/' {
			_, cok = BbCodeTags[tagData[0][1:]]
		}
		if ok {
			output += body[:tagLoc[0]]
			body = body[tagLoc[1]:]

			args := make([]string, 2)
			args[0] = tagData[1]

			if htmlTags.Options ^ TagBodyAsArg != 0 || htmlTags.Options ^ AllowTagBodyAsFirstArg != 0 {
				closeTagRe, err := BbCloseTag(tagData[0])
				if err != nil {
					return nil, err
				}
				closeTagLoc := closeTagRe.FindIndex([]byte(body))
				if closeTagLoc == nil {
					if htmlTags.Options ^ PossibleSingle == 0 {
						if htmlTags.Options ^ TagBodyAsArg != 0 {
							args[1] = body[tagLoc[1]:]
						}
						if htmlTags.Options ^ AllowTagBodyAsFirstArg != 0 && args[0] == "" {
							args[0] = body[tagLoc[1]:]
						}
						body = ""
					}
				} else {
					if htmlTags.Options ^ TagBodyAsArg != 0 {
						args[1] = body[:closeTagLoc[0]]
					}
					if htmlTags.Options ^ AllowTagBodyAsFirstArg != 0 && args[0] == "" {
						args[0] = body[:closeTagLoc[0]]
					}
					body = body[closeTagLoc[1]:]
				}
			}

			if htmlTags.InputModFunc != nil {
				htmlTags.InputModFunc(&args)
			}

			if htmlTags.OutputFunc != nil {
				output += htmlTags.OutputFunc(args)
			} else {
				for i, tag := range htmlTags.Tags {
					output += "<" + tag

					if len(htmlTags.Classes) > i {
						if classes := htmlTags.Classes[i]; classes != nil {
							output += " class=\""
							for _, class := range classes {
								// TODO: use golang html/template for secure escaping
								output += " " + class
							}
							output += "\""
						}
					}

					if len(htmlTags.Attributes) > i {
						if attrs := htmlTags.Attributes[i]; attrs != nil {
							for i, attr := range attrs {
								if args[i] != "" {
									// TODO: use golang html/template for secure escaping
									output += " " + attr + "=\"" + args[i] + "\""
								}
							}
						}
					}

					if len(htmlTags.CssProps) > i {
						if cssProps := htmlTags.CssProps[i]; cssProps != nil {
							for i, cssProp := range cssProps {
								output += " style=\""
								if args[i] != "" {
									// TODO: use golang html/template for secure escaping
									output += cssProp + ": " + args[i] + ";"
								}
								output += "\""
							}
						}
					}

					output += ">"
				}
			}

			tagStack.Push(&bbTagPair{htmlTags, tagData[0]})

			if htmlTags.Options ^ HtmlSingle != 0 {
				tagStack.Pop()
			}

			if htmlTags.Options ^ PossibleSingle != 0 {
				closeTagRe, err := BbCloseTag(tagData[0])
				if err != nil {
					return nil, err
				}
				closeTagLoc := closeTagRe.FindIndex([]byte(body))
				if closeTagLoc == nil {
					output += closeTags(tagStack.TopFromN(tagStack.count-1))
				}
				tagStack.Pop()
			}
		} else if cok {
			tagStackCopy := &Stack{nodes: tagStack.nodes, count: tagStack.count}
			for tagPair := tagStackCopy.Pop(); tagPair != nil; tagPair = tagStackCopy.Pop() {
				if tagPair.bbName == tagData[0][1:] {
					output += closeTags(tagStack.TopFromN(tagStackCopy.count))
					break
				}
			}
		} else {
			output += body[:tagLoc[1]]
			body = body[tagLoc[1]:]
		}
	}

	output += closeTags(tagStack)
	return []byte(output), nil
}
