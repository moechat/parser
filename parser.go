// Package parser (github.com/moechat/parser) provides interfaces for a formal language parser that matches
package parser

// The Parser type is used to analyse a list of Tokens into
type Parser interface {
	Parse(text string) string
}
