package moeparser

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
)

type HtmlToken struct {
	Name string // The HTML element name

	IsClose bool // Whether this element is the close token

	Args []string // Arguments passed to this HtmlToken

	Classes    []string        // HTML element classes
	Attributes map[int8]string // HTML element attributes
	CssProps   map[int8]string // CSS Properties

	// A custom output function:
	// if set, the parser will ignore all fields and simply spit out the output.
	// The input is the array of strings that were captured
	// Use a $b to represent the body of the tag and $$ for a dollar sign
	OutputFunc func([]string) string
}

func (e *HtmlToken) Copy() *HtmlToken {
	ret := *e
	return &ret
}

func (e *HtmlToken) SetArgs(args []string) {
	e.Args = args
}

func (e *HtmlToken) GetOutput() (string, error) {
	if e.OutputFunc != nil {
		return e.OutputFunc(e.Args), nil
	} else if e.IsClose {
		return fmt.Sprintf("</%s>", e.Name), nil
	} else {
		templStr := "<" + e.Name

		if e.Classes != nil {
			templStr += " class=\""
			for _, class := range e.Classes {
				templStr += " " + class
			}
			templStr += "\""
		}

		if e.Attributes != nil {
			for i, attr := range e.Attributes {
				if len(e.Args) > int(i) && e.Args[i] != "" {
					templStr += " " + attr + "=\"{{index . " + strconv.Itoa(int(i)) + "}}\""
				}
			}
		}

		if e.CssProps != nil {
			templStr += " style=\""
			for i, cssProp := range e.CssProps {
				if len(e.Args) > int(i) && e.Args[i] != "" {
					templStr += cssProp + ": {{index . " + strconv.Itoa(int(i)) + "}};"
				}
			}
			templStr += "\""
		}

		templStr += ">"

		tmpl, err := template.New("elementTemplate").Parse(templStr)
		if err != nil {
			return "", err
		}

		eleBuffer := bytes.Buffer{}
		err = tmpl.Execute(&eleBuffer, e.Args)
		if err != nil {
			return "", err
		}
		return eleBuffer.String(), nil
	}
}

/*
// Returns a symmetric tag, returning nil if tag == ""
func SymTag(tag string) TagPair {
	if tag == "" {
		return TagPair{}
	}
	return TagPair{Start: tag}
}

// Returns a new tag with Start = start and End = end, catching empty start (returns nil) and
// start == end (returns Start = start and End = "")
func Tag(start string, end string) TagPair {
	if start == "" {
		return TagPair{}
	}
	if start == end {
		return TagPair{Start: start}
	}
	return TagPair{start, end}
}

// Returns a TagPair that will match the bbcode [name]
func BbToken(name string) TagPair {
	if name == "" {
		return TagPair{}
	}

	start := fmt.Sprintf("[%s(=(?!]))]", name)
	end := fmt.Sprintf("[//%s]", name)

	return TagPair{start, end}
}
*/
