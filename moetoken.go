package moeparser

import (
	"fmt"
	"github.com/moechat/moeparser/token"
)

var nextId byte = 0

// A TokenClass is an extremely general token class. This will not match tags that are in the middle of a word unless AllowInWord is set.
type TokenClass struct {
	name string // The name of this token class (must be unique and should not start with 0x)

	regexps token.RegexpList // The expressions that map to this token class

	options   token.TokenClassOptions // Options for this token class
	tokenType token.TokenType         // The type of this token class

	// A function that modifies arguments (i.e. a function that converts a username to a user ID in @tagging)
	argModFunc func(args []string, namesById map[string]int) ([]string, map[string]int)
	isValid    func(args *token.TokenArgs) bool

	tokens []token.Token // The tokens to insert
}

type TokenClassArgs struct {
	Name       string
	Options    token.TokenClassOptions
	Type       token.TokenType
	ArgModFunc func(args []string, namesById map[string]int) ([]string, map[string]int)
	IsValid    func(args *token.TokenArgs) bool
	Tokens     []token.Token
	NotRe      bool
}

func NewTokenClass(args TokenClassArgs, exprs ...string) *TokenClass {
	regexps := token.NewRegexpList(args.NotRe, exprs...)

	if args.Name == "" {
		args.Name = fmt.Sprintf("_%#x", nextId)
		nextId++
	}

	if args.IsValid == nil {
		args.IsValid = func(*token.TokenArgs) bool { return true }
	}

	return &TokenClass{args.Name, regexps, args.Options, args.Type, args.ArgModFunc, args.IsValid, args.Tokens}
}

func (mtc *TokenClass) Regexps() token.RegexpList {
	return mtc.regexps
}
func (mtc *TokenClass) Options() token.TokenClassOptions {
	return mtc.options
}
func (mtc *TokenClass) Type() token.TokenType {
	return mtc.tokenType
}
func (mtc *TokenClass) Name() string {
	return mtc.name
}

func (mtc *TokenClass) ModifyArgs(args []string, idByName map[string]int) ([]string, map[string]int) {
	return mtc.argModFunc(args, idByName)
}

func (mtc *TokenClass) IsValid(args *token.TokenArgs) bool {
	return mtc.isValid(args)
}

func (mtc *TokenClass) BuildTokens(args *token.TokenArgs) []token.Token {
	ret := make([]token.Token, len(mtc.tokens))
	for i, t := range mtc.tokens {
		ret[i] = t.Copy()
		ret[i].SetArgs(args)
		if mtc.tokenType == token.SymmetricType {
		}
	}
	return ret
}
