package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type matcherWrap struct {
	Matcher

	id      int              // A unique ID for this matcher
	regexps []*regexp.Regexp // The compiled regexps
}

/*
This is an implementation of the a Lexer, used to convert text into tokens
(http://en.wikipedia.org/wiki/Lexical_analysis) using the regexp package.

Because the regexp package is implemented using a NFA
(http://en.wikipedia.org/wiki/Nondeterministic_finite_automaton),
it's very effective for this use case.
*/
type Lexer struct {
	matchersById map[int]*matcherWrap // All matchers by their ID
	matchers     map[*matcherWrap]int // All matchers and their corresponding capture group ID

	exprs map[string]bool // The set of regexps that are matched

	expr   string         // The main regexp expression
	regexp *regexp.Regexp // The regexp used to match tags

	nextId int // The ID of the next matcher to be added
}

// Creates a new Lexer.
func NewLexer() *Lexer {
	return &Lexer{matchers: make(map[*matcherWrap]int), exprs: make(map[string]bool)}
}

func MustCompile(matchers ...Matcher) *Lexer {
	l, err := Compile(matchers...)
	if err != nil {
		panic(err)
	}

	return l
}

func Compile(matchers ...Matcher) (*Lexer, error) {
	ret := NewLexer()
	err := ret.AddMatchers(matchers...)
	if err != nil {
		return nil, err
	}

	err = ret.Compile()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Adds the passed matcher/matchers to the lexer
func (l *Lexer) AddMatchers(matchers ...Matcher) error {
	for _, matcher := range matchers {
		wrap := &matcherWrap{Matcher: matcher, id: l.nextId}
		l.nextId++
		for _, expr := range matcher.Exprs() {
			if l.exprs[expr] {
				err := errors.New("Another matcher has an identical regexp!")
				return err
			}
		}

		for _, expr := range matcher.Exprs() {
			l.exprs[expr] = true
		}

		l.matchersById[wrap.id] = wrap
		l.matchers[wrap] = -1
	}

	return nil
}

// Removes all instances of the token class specified from the token class map.
func (l *Lexer) RemoveMatchers(matchers ...Matcher) {
	for _, matcher := range matchers {
		for wrap := range l.matchers {
			if Matcher(wrap) == matcher {
				for _, expr := range wrap.Exprs() {
					delete(l.exprs, expr)
				}

				delete(l.matchersById, wrap.id)
				delete(l.matchers, wrap)
			}
		}
	}
}

func (l *Lexer) MustCompile() {
	err := l.Compile()
	if err != nil {
		panic(err)
	}
}

// Creates the and compiles the regexp used by the Lexer.
//
// It must be run after adding or removing token classes in order for changes to take effect.
func (l *Lexer) Compile() error {
	var err error

	l.expr = ""
	for matcher := range l.matchers {
		matcher.regexps = make([]*regexp.Regexp, len(matcher.Exprs()))
		for i, expr := range matcher.Exprs() {
			matcher.regexps[i], err = regexp.Compile(expr)
			if err != nil {
				return err
			}
			l.expr += fmt.Sprintf("(?P<_%02x%x>%s)|", i, matcher.id, expr)
		}
	}
	l.expr = l.expr[:len(l.expr)-1]

	l.regexp, err = regexp.Compile(l.expr)
	if err != nil {
		return err
	}

	names := l.regexp.SubexpNames()
	usedNames := make(map[string]bool)

	for i, name := range names {
		if name != "" {
			if name[0] == '_' {
				matcherId64, _ := strconv.ParseInt(name[3:], 16, 0)
				matcherId := int(matcherId64)
				matcher, ok := l.matchersById[matcherId]
				if !ok {
					return errors.New("lexer: capture group names starting with _ are reserved for use by the lexer! Your name is " + name)
				}

				usedNames[name] = true
				l.matchers[matcher] = i
			}
		}
	}

	return nil
}

// Converts an input string into Tokens.
//
// If the tree has not been built, tokenize will run BuildCharTree().
// BuildCharTree() *must* be run if you call AddTokenClass or RemoveTokenClass between Tokenize()'s
func (l *Lexer) Tokenize(data string) []Token {
	ret := make([]Token, 0)
	subexpNames := l.regexp.SubexpNames()

	for data != "" {
		indices := l.regexp.FindStringSubmatchIndex(data)
		if indices == nil {
			ret = append(ret, NewTextToken(data))
			data = ""
			break
		}

		for matcher, i := range l.matchers {
			if indices[i*2] >= 0 {
				if indices[i*2] != 0 {
					ret = append(ret, NewTextToken(data[:indices[i*2]]))
				}

				exprId, _ := strconv.ParseInt(subexpNames[i][1:3], 16, 8)
				currRe := matcher.regexps[exprId]
				args := []string(currRe.FindStringSubmatch(data[indices[i*2]:indices[i*2+1]]))

				idByName := make(map[string]int)
				for i, name := range currRe.SubexpNames() {
					if name != "" {
						idByName[name] = i
					}
				}

				args, idByName = matcher.ArgModFunc(args, idByName)
				tokenArgs := NewTokenArgs(args, idByName)

				if matcher.IsValid(tokenArgs) {
					ret = append(ret, matcher.BuildTokens(tokenArgs)...)
				} else {
					ret = append(ret, NewTextToken(data[indices[i*2]:indices[i*2+1]]))
				}
				data = data[indices[i*2+1]:]
			}
		}
	}

	return ret
}
