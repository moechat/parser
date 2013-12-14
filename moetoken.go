package parser

import (
	"fmt"
	"github.com/moechat/parser/token"
	"regexp"
)

var nextId byte = 0

type HtmlTokenBuilder struct {
	HtmlElements []string
}

func (htb *HtmlTokenBuilder) Build(args *token.TokenArgs) token.Token {

}

// A Matcher is an extremely general lexer.Matcher. This will not match tags that are in the middle of a word unless AllowInWord is set.
type Matcher struct {
	name string // The name of this token class (must be unique and should not start with 0x)

	exprs []string // The expressions that map to this token class

	options   int // Options for this token class
	tokenType int // The type of this token class

	// A function that modifies arguments (i.e. a function that converts a username to a user ID in @tagging)
	argModFunc func(args []string, namesById map[string]int) ([]string, map[string]int)
	isValid    func(args *token.TokenArgs) bool

	tokenBuilders []token.TokenBuilder // The token builders to use when matched
}

type MatcherArgs struct {
	Name          string
	Options       int
	Type          int
	ArgModFunc    func(args []string, namesById map[string]int) ([]string, map[string]int)
	IsValid       func(args *token.TokenArgs) bool
	TokenBuilders []token.TokenBuilder
	NotRe         bool
}

func NewMatcher(args MatcherArgs, exprs ...string) *Matcher {
	if args.NotRe {
		for i, expr := range exprs {
			exprs[i] = regexp.QuoteMeta(expr)
		}
	}

	if args.Name == "" {
		args.Name = fmt.Sprintf("_%#x", nextId)
		nextId++
	}

	if args.IsValid == nil {
		args.IsValid = func(*token.TokenArgs) bool { return true }
	}

	return &Matcher{args.Name, exprs, args.Options, args.Type, args.ArgModFunc, args.IsValid, args.TokenBuilders}
}

func (m *Matcher) Exprs() []string {
	return m.exprs
}

func (m *Matcher) Options() int {
	return m.options
}

func (m *Matcher) Type() int {
	return m.tokenType
}

func (m *Matcher) Name() string {
	return m.name
}

func (m *Matcher) ModifyArgs(args []string, idByName map[string]int) ([]string, map[string]int) {
	return m.argModFunc(args, idByName)
}

func (m *Matcher) IsValid(args *token.TokenArgs) bool {
	return m.isValid(args)
}

func (m *Matcher) BuildTokens(args *token.TokenArgs) []token.Token {
	ret := make([]token.Token, len(m.tokenBuilders))
	for i, t := range m.tokenBuilders {
		ret[i] = t.Build(args)
		if m.Type() == token.SymmetricToken {
		}
	}
	return ret
}
