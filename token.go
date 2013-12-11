package parser

import (
	"regexp"
)

type TokenClassFlag uint8

// Options for Tokens - return these bits in GetOptions() to implement this behavior
const (
	// TODO: Interpret as single if there is no closing tag
	PossibleSingle TokenClassFlag = 1 << iota
	// TODO: This makes MoeParser ignore any tags inside this tags body. It will be ignored if the Single bit is set.
	NoParseInner
	// TODO: This makes the text inside of tag and passes it as an arg for the output. The text inside will not be parsed.
	TokenBodyAsArg
	// TODO: This makes the tag body become the first arg if there is no first argument (makes [name]arg0[/name] the same as [name=arg0][/name])
	AllowTokenBodyAsFirstArg
	// TODO: remove Converts a number to the number + "px" (ie 12 -> 12px)
	NumberArgToPx
	// TODO: This makes MoeParser match the tags that don't either start with whitespace or the beginning of a line. This is only useful if the token is of the type OpenCloseClass.
	AllowMidWord
	// TODO: This makes MoeParser stop matching:
	// - OpenClass type tokens without leading whitespace/beginning of body
	// - CloseClass type tokens without trailing whitespace/end of body
	// - Single type tokens without leading and trailing whitespace
	DisallowMidWord
)

type TokenType uint8

// Token class types - return these in GetType()
const (
	// A single token with no opening or closing token
	SingleType TokenType = iota
	// A token class that starts a section
	OpenType
	// A token class that ends a section
	CloseType
	// A token class that can both begin and end a section
	SymmetricType
)

type RegexpList []*regexp.Regexp

// A utility function that automatically generates the expressions.
//
// If isRe is false, regexp.QuoteMeta will be run and the exact string will be matched.
//
// It will panic if passed an invalid regexp
func NewRegexpList(NotRe bool, exprs ...string) RegexpList {
	ret := make(RegexpList, len(exprs))
	if NotRe {
		for i, expr := range exprs {
			exprs[i] = regexp.QuoteMeta(expr)
		}
	}
	for i, expr := range exprs {
		ret[i] = regexp.MustCompile(expr)
	}
	return ret
}

// A token class is recognized by the lexer
type TokenClass interface {
	Regexps() RegexpList // Returns a list of the uncompiled regexps that are mapped to the token class

	Options() TokenClassFlag // Returns options for the token class
	Type() TokenType         // Returns the type of the token class

	Name() string // Returns a unique name for the token class. This is used as the regexp capturing group name, and will break Lexer.CompileRegexp if an invalid name is used

	// Run on args before they are passed to BuildTokens
	ModifyArgs(args []string, namesById map[string]int) ([]string, map[string]int)
	// If false, a PlainToken is returned instead of BuildTokens. This is run after ModifyArgs.
	IsValid(args *TokenArgs) bool
	BuildTokens(args *TokenArgs) []Token // Returns instances of the tokenClass with args set
}

type TokenArgs struct {
	args     []string
	idByName map[string]int

	size int
}

func NewTokenArgs(args []string, idByName map[string]int) *TokenArgs {
	return &TokenArgs{args, idByName, len(args)}
}

func (ta *TokenArgs) ById(id int) string {
	if id < ta.size {
		return ta.args[id]
	}
	return ""
}

func (ta *TokenArgs) ByName(name string) string {
	if id, ok := ta.idByName[name]; ok {
		return ta.args[id]
	}
	return ""
}

func (ta *TokenArgs) Size() int {
	return ta.size
}

// A token is returned by the lexer and recognized by the parser
type Token interface {
	Copy() Token // Get a copy of the token

	SetArgs(*TokenArgs)      // Set the args of the token
	Output() (string, error) // Get the output of the token
}

// A special case of Token used to represent text that isn't matched by any other tokens
// i.e. "hi" in <p>hi</p>
type TextToken struct {
	body string
}

func NewTextToken(body string) *TextToken {
	return &TextToken{body: body}
}

func (pt *TextToken) setBody(body string) {
	pt.body = body
}

func (pt *TextToken) Copy() Token {
	return &TextToken{body: pt.body}
}

func (pt *TextToken) SetArgs(*TokenArgs) {}

// Returns the Text's body
func (pt *TextToken) Output() (string, error) {
	return pt.body, nil
}
