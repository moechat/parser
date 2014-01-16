package lexer

import (
	"errors"
	"fmt"
	"github.com/moechat/parser/token"
	"regexp"
)

// A Matcher pairs a set of regexps and a set of tokens.
type Matcher interface {
	Name() string // A unique name for this matcher
	Expr() string // The regular expressions that match this token

	IsValid(args *token.TokenArgs) bool
	BuildToken(args *token.TokenArgs) token.Token
}

/*
 * This is an implementation of the a Lexer, used to convert text into tokens
 * (http://en.wikipedia.org/wiki/Lexical_analysis) using the regexp package.

 * Because the regexp package is implemented using a NFA
 * (http://en.wikipedia.org/wiki/Nondeterministic_finite_automaton),
 * it's very effective for this use case.
 */
type Lexer struct {
	matchers  map[string]Matcher
	regexps   map[string]*regexp.Regexp
	subexpIds map[string]int

	expr   string         // The main regexp expression
	regexp *regexp.Regexp // The regexp used to match tags
}

func Must(l *Lexer, err error) *Lexer {
	if err != nil {
		panic(err)
	}
	return l
}

func New(matchers ...Matcher) (*Lexer, error) {
	var err error
	l := &Lexer{
		matchers:  make(map[string]Matcher),
		regexps:   make(map[string]*regexp.Regexp),
		subexpIds: make(map[string]int),
	}

	for _, matcher := range matchers {
		l.matchers[matcher.Name()] = matcher
		l.expr += fmt.Sprintf("(?P<_%s>%s)|", matcher.Name(), matcher.Expr())
		l.regexps[matcher.Name()], err = regexp.Compile(matcher.Expr())
		if err != nil {
			return nil, err
		}
	}

	l.expr = l.expr[:len(l.expr)-1] // Cut off the trailing '|'
	l.regexp, err = regexp.Compile(l.expr)
	if err != nil {
		return nil, err
	}

	for i, name := range l.regexp.SubexpNames() {
		if name != "" && name[0] == '_' {
			name = name[1:] // Remove the '_' prefix
			_, ok := l.matchers[name]
			if !ok {
				return nil, errors.New("lexer: capture group names starting with _ are reserved for use by the lexer! Your name is " + name)
			}

			l.subexpIds[name] = i
		}
	}

	return l, nil
}

/*
 * Converts an input string into Tokens.
 */
func (l *Lexer) Tokenize(data string) []token.Token {
	ret := make([]token.Token, 0)

	for data != "" {
		indices := l.regexp.FindStringSubmatchIndex(data)
		if indices == nil {
			ret = append(ret, token.TextToken{data})
			break
		}

		for name, matcher := range l.matchers {
			i := l.subexpIds[name]
			if indices[i*2] >= 0 {
				if indices[i*2] != 0 {
					ret = append(ret, token.TextToken{data[:indices[i*2]]})
				}

				args := []string(l.regexps[name].FindStringSubmatch(data[indices[0]:indices[1]]))

				tokenArgs := token.NewTokenArgs(args, l.regexps[name].SubexpNames())

				if matcher.IsValid(tokenArgs) {
					ret = append(ret, matcher.BuildToken(tokenArgs))
				} else {
					ret = append(ret, token.TextToken{data[indices[i*2]:indices[i*2+1]]})
				}
				data = data[indices[i*2+1]:]
			}
		}
	}

	return ret
}
