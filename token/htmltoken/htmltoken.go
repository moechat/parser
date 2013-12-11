package htmltoken

import (
	"bytes"
	"fmt"
	"github.com/moechat/moeparser/token"
	"html/template"
)

type Token struct {
	Name string // The HTML element name

	Type token.TokenType

	Args *token.TokenArgs // Arguments passed to this HTML Token

	// The following fields are run through html/template.
	//
	// To get a capture group by id, use {{.ById id}}.
	// Analogously, use {{.ByName name}} to get a capture group by name.
	//
	// See http://golang.org/pkg/html/template/ for more details
	Prefix     string            // A string to insert before the element
	Suffix     string            // A string to insert after the element
	Classes    []string          // HTML element classes
	Attributes map[string]string // HTML element attributes by argument ID
	CssProps   map[string]string // CSS Properties by argument ID

	// A custom output function
	// if set, the parser will ignore all fields and simply spit out the return value.
	// The input is the array of strings that were captured
	// Use a $b to represent the body of the tag and $$ for a dollar sign
	OutputFunc func(*token.TokenArgs) string
}

func (e *Token) Copy() token.Token {
	ret := *e
	return &ret
}

func (e *Token) SetArgs(args *token.TokenArgs) {
	e.Args = args
}

func (e *Token) Output() (string, error) {
	if e.OutputFunc != nil {
		return e.OutputFunc(e.Args), nil
	} else if e.Type == token.CloseType {
		return fmt.Sprintf("</%s>", e.Name), nil
	} else {
		templStr := "<" + e.Name

		if e.Classes != nil {
			templStr += ` class="`
			for i, class := range e.Classes {
				if i != 0 {
					templStr += " "
				}
				templStr += class
			}
			templStr += `"`
		}

		for attr, value := range e.Attributes {
			templStr += fmt.Sprintf(` %s="%s"`, attr, value)
		}

		if e.CssProps != nil {
			templStr += ` style="`
			for cssProp, value := range e.CssProps {
				templStr += cssProp + ":" + value
			}
			templStr += `"`
		}

		tmpl, err := template.New("elementTemplate").Parse(templStr)
		if err != nil {
			return "", err
		}

		eleBuffer := bytes.Buffer{}
		err = tmpl.Execute(&eleBuffer, e.Args)
		if err != nil {
			return "", err
		}
		openTag := eleBuffer.String()

		switch e.Type {
		case token.SingleType:
			return openTag + fmt.Sprintf("</%s>", e.Name), nil
		default:
			return openTag, nil
		}
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
