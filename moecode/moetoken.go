package moecode

import (
	"github.com/moechat/moeparser"
)

// The default moecode tokens.
//
// This is here only for reference - changing this has no effect.
// Use Addmoeparser.TokenClass and Removemoeparser.TokenClass instead.
var DefaultTokenMap = moeparser.TokenClassMap{
	"`": &MoeTokenClass{
		Options: moeparser.NoParseInner,
	},
	"*":  &MoeTokenClass{},
	"**": &MoeTokenClass{},
}

// A Moemoeparser.TokenClass is an extremely general token class. This will not match tags that are in the middle of a word unless AllowInWord is set.
type MoeTokenClass struct {
	Options moeparser.TokenOptions   // Options for this matcher
	Type    moeparser.TokenClassType // The type of this token
	Args    []string                 // The arguments

	ArgModFunc func([]string) []string // A function that modifies arguments (i.e. a function that converts a username to a user ID in @tagging)

	Tokens []moeparser.Token // The tokens to insert
}

func (mtc *MoeTokenClass) GetOptions() moeparser.TokenOptions {
	return mtc.Options
}

func (mtc *MoeTokenClass) GetType() moeparser.TokenClassType {
	return mtc.Type
}

func (mtc *MoeTokenClass) GetTokens(args []string) []moeparser.Token {
	ret := make([]moeparser.Token, len(args))
	for i, token := range mtc.Tokens {
		ret[i] = token
		ret[i].SetArgs(args)
	}
	if mtc.ArgModFunc != nil {
		args = mtc.ArgModFunc(args)
	}
	return ret
}
