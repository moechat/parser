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

	Classes    []string            // HTML element classes
	ModClasses func(classes []string, args *token.TokenArgs) []string // Change classes to output

	AttributesById map[int]string // HTML element attributes by argument ID
	AttributesByName map[string]string // HTML element attributes by argument name
	CssPropsById   map[int]string // CSS Properties by argument ID
	CssPropsByName   map[string]string // CSS Properties by argument name

	// A custom output function:
	// if set, the parser will ignore all fields and simply spit out the output.
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
		var classes, attrs, styles string

		if e.ModClasses != nil {
			e.Classes = e.ModClasses(e.Classes, e.Args)
		}
		if e.Classes != nil {
			classes += ` class="`
			for _, class := range e.Classes {
				classes += " " + class
			}
			classes += `"`
		}

		for argId, attr := range e.AttributesById {
			if e.Args.ById(argId) != "" {
				attrs += fmt.Sprintf(` %s="{{.ById %d}}"`, attr, argId)
			}
		}
		for argName, attr := range e.AttributesByName {
			if e.Args.ByName(argName) != "" {
				attrs += fmt.Sprintf(` %s="{{.ByName %s}}"`, attr, argName)
			}
		}

		if e.CssPropsById != nil || e.CssPropsByName != nil {
			styles += ` style="`
			for argId, cssProp := range e.CssPropsById {
				if e.Args.ById(argId) != "" {
					styles += fmt.Sprintf(`%s:{{.ById %d}};`, cssProp, argId)
				}
			}
			for argName, cssProp := range e.CssPropsByName {
				if e.Args.ByName(argName) != "" {
					styles += fmt.Sprintf(`%s:{{.ByName %s}};`, cssProp, argName)
				}
			}
			styles += `"`
		}

		templStr := fmt.Sprintf(`<%s%s%s%s>`, e.Name, classes, attrs, styles)

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
