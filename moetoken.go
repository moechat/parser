package moeparser

import (
	"github.com/moechat/moeparser/token"
)

// A MoeTokenClass is an extremely general token class. This will not match tags that are in the middle of a word unless AllowInWord is set.
type MoeTokenClass struct {
	Options token.TokenOptions   // Options for this matcher
	Type    token.TokenClassType // The type of this token
	Args    []string                 // The arguments

	ArgModFunc func([]string) []string // A function that modifies arguments (i.e. a function that converts a username to a user ID in @tagging)

	Tokens []token.Token // The tokens to insert
}

func (mtc *MoeTokenClass) GetOptions() token.TokenOptions {
	return mtc.Options
}

func (mtc *MoeTokenClass) GetType() token.TokenClassType {
	return mtc.Type
}

func (mtc *MoeTokenClass) GetTokens(args []string) []token.Token {
	ret := make([]token.Token, len(args))
	for i, token := range mtc.Tokens {
		ret[i] = token
		ret[i].SetArgs(args)
	}
	if mtc.ArgModFunc != nil {
		args = mtc.ArgModFunc(args)
	}
	return ret
}
