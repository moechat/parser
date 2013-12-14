/*
 * This package implements a parser for any markup-style tags, including BBCode tags.
 */
package parser

import (
	"github.com/moechat/parser/lexer"
	//"github.com/moechat/parser/token"
)

var Lexer *lexer.Lexer

type Parser struct {
	lexer *lexer.Lexer
}

/*
func init() {
	tokenMap := map[string]Matcher{
		"code": &Matcher{
			name:      "code",
			exprs:     []string{"`"},
			options:   token.NoParseInner,
			tokenType: token.SymmetricType,
			tokens: []token.Token{
				&htmltoken.Token{Name: "pre"},
				&htmltoken.Token{Name: "code"},
			},
		},
		"italic": &Matcher{
			name:      "italic",
			exprs:     []string{`\*`},
			tokenType: token.SymmetricType,
			tokens: []token.Token{
				&htmltoken.Token{Name: "i"},
			},
		},
		"bold": &Matcher{
			name:      "bold",
			exprs:     []string{`\*\*`},
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
*/
