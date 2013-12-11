/*
 * This package implements a parser for any markup-style tags, including BBCode tags.
 */
package moeparser

import (
	"github.com/moechat/moeparser/lexer"
	"github.com/moechat/moeparser/token"
	"github.com/moechat/moeparser/token/htmltoken"
)

var Lexer *lexer.Lexer

type Parser struct {
	lexer *lexer.Lexer
}

func init() {
	tokenMap := map[string]token.TokenClass{
		"code": &TokenClass{
			name:      "code",
			regexps:   token.NewRegexpList(true, "`"),
			options:   token.NoParseInner,
			tokenType: token.SymmetricType,
			tokens: []token.Token{
				&htmltoken.Token{Name: "pre"},
				&htmltoken.Token{Name: "code"},
			},
		},
		"italic": &TokenClass{
			name:      "italic",
			regexps:   token.NewRegexpList(true, "*"),
			tokenType: token.SymmetricType,
			tokens: []token.Token{
				&htmltoken.Token{Name: "i"},
			},
		},
		"bold": &TokenClass{
			name:      "bold",
			regexps:   token.NewRegexpList(true, "**"),
			tokenType: token.SymmetricType,
			tokens: []token.Token{
				&htmltoken.Token{Name: "b"},
			},
		},
	}

	var err error
	Lexer, err = lexer.New(tokenMap)
	if err != nil {
		panic(err)
	}
}
