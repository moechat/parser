/*
 This package implements a parser for any markup-style tags, including BBCode tags.
*/
package moeparser

import (
	"github.com/moechat/moeparser/lexer"
	"github.com/moechat/moeparser/token"
	"github.com/moechat/moeparser/token/htmltoken"
)

// The default moeparser tokens.
//
// This is here only for reference - changing this has no effect.
// Use Lexer.AddTokenClass and Lexer.RemoveTokenClass instead.
var DefaultTokenMap = map[string]token.TokenClass{
	"`": &MoeTokenClass{
		Options: token.NoParseInner,
		Type:    token.OpenClose,
		Tokens: []*htmltoken.Token{
			{Name: "pre"},
			{Name: "code"},
		},
	},
	"*": &MoeTokenClass{
		Type: token.OpenClose,
		Tokens: []*htmltoken.Token{
			{Name: "i"},
		},
	},
	"**": &MoeTokenClass{
		Type: token.OpenClose,
		Tokens: []*htmltoken.Token{
			{Name: "b"},
		},
	},
}

var Lexer = lexer.New(DefaultTokenMap)

/*
func Parse(b string) (string, error) {
	body, err := BbCodeParse(body)
	if err != nil {
		return "", err
	}

	return body, nil
}
*/
