package lexer

import (
	"errors"
	"fmt"
	"github.com/moechat/parser/token"
	"regexp"
)

type Flags int

const (
	NoParseInner Flags = 1 << iota
	BodyAsArg
	NoNewline
	RequireClose
)

type Expression struct {
	Expr      string
	CloseExpr string

	Flags Flags
}

// A Matcher pairs a set of regexps and a set of tokens.
type Matcher interface {
	Name() string        // A unique name for this matcher
	Exprs() []Expression // The regular expressions that match this token

	IsValid(args *token.TokenArgs, expNum int) bool
	BuildToken(args *token.TokenArgs, expNum int) (openToken token.Token, closeToken token.Token)
}

/*
 * This is an implementation of the a Lexer, used to convert text into tokens
 * (http://en.wikipedia.org/wiki/Lexical_analysis) using the regexp package.
 *
 * Because the regexp package is implemented using a NFA
 * (http://en.wikipedia.org/wiki/Nondeterministic_finite_automaton),
 * it's very effective for this use case.
 */
type Lexer struct {
	matchers   map[string]Matcher
	regexps    map[string][]*regexp.Regexp
	subexpIds  map[string][]int
	bodyExpIds map[string][]int

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
		matchers:   make(map[string]Matcher),
		regexps:    make(map[string][]*regexp.Regexp),
		subexpIds:  make(map[string][]int),
		bodyExpIds: make(map[string][]int),
	}

	for _, matcher := range matchers {
		l.matchers[matcher.Name()] = matcher

		numExprs := len(matcher.Exprs())
		l.regexps[matcher.Name()] = make([]*regexp.Regexp, numExprs, numExprs)
		l.subexpIds[matcher.Name()] = make([]int, numExprs, numExprs)
		l.bodyExpIds[matcher.Name()] = make([]int, numExprs, numExprs)
		for i, expr := range matcher.Exprs() {
			var argExpr, realExpr string
			if expr.CloseExpr != "" {
				closeExpr := expr.CloseExpr
				if expr.Flags&RequireClose == 0 {
					closeExpr = "(?:" + expr.CloseExpr + ")|$"
				}

				var bodyExpr, rBodyExpr string
				if expr.Flags&BodyAsArg != 0 {
					bodyExpr += "?:"
				}
				if expr.Flags&NoNewline == 0 {
					bodyExpr += "(?s)"
					rBodyExpr += "(?s)"
				}

				argExpr = fmt.Sprintf("(?:%s)(%s.*?)(?:%s)", expr.Expr, bodyExpr, closeExpr)
				realExpr = fmt.Sprintf("(?:%s)(?P<_i%02x%s>%s.*?)(?:%s)",
					expr.Expr, i, matcher.Name(), rBodyExpr, closeExpr)
			} else {
				argExpr = "(?:" + expr.Expr + ")(?-imsU)"
				realExpr = "(?:" + expr.Expr + ")(?-imsU)"
			}

			l.regexps[matcher.Name()][i], err = regexp.Compile(argExpr)
			if err != nil {
				return nil, err
			}

			for _, subexpName := range l.regexps[matcher.Name()][i].SubexpNames() {
				if subexpName != "" && subexpName[0] == '_' {
					return nil, errors.New("lexer: capture group names starting with _ are reserved for use by the lexer! Your name is " + subexpName)
				}
			}

			l.expr += fmt.Sprintf("(?P<_%02x%s>%s)|", i, matcher.Name(), realExpr)
		}
	}

	l.expr = l.expr[:len(l.expr)-1] // Cut off the trailing '|'
	l.regexp, err = regexp.Compile(l.expr)
	if err != nil {
		return nil, err
	}

	for i, name := range l.regexp.SubexpNames() {
		if name != "" && name[0] == '_' {
			var expNum int
			var matcherName string
			var err error
			if name[1] == 'i' {
				_, err = fmt.Sscanf(name, "_i%02x%s", &expNum, &matcherName)
				if err != nil {
					return nil, errors.New("Something went wrong while compiling D:")
				}

				l.bodyExpIds[matcherName][expNum] = i
			} else {
				_, err = fmt.Sscanf(name, "_%02x%s", &expNum, &matcherName)
				if err != nil {
					return nil, errors.New("Something went wrong while compiling D:")
				}

				l.subexpIds[matcherName][expNum] = i
			}
		}
	}

	return l, nil
}

/*
 * Converts an input string into Tokens.
 */
func (l *Lexer) Tokenize(data string) []token.Token {
	ret := make([]token.Token, 0)
	toAppend := ""

	for data != "" {
		indices := l.regexp.FindStringSubmatchIndex(data)
		if indices == nil {
			ret = append(ret, token.TextToken{toAppend + data})
			break
		}

		for name, matcher := range l.matchers {
			for expNum, i := range l.subexpIds[name] {
				if i != 0 && indices[i*2] >= 0 {
					if indices[i*2] != 0 {
						toAppend += data[:indices[i*2]]
					}

					args := l.regexps[name][expNum].FindStringSubmatch(data[indices[0]:indices[1]])

					tokenArgs := token.NewTokenArgs(args, l.regexps[name][expNum].SubexpNames())

					if matcher.IsValid(tokenArgs, expNum) {
						openToken, closeToken := matcher.BuildToken(tokenArgs, expNum)

						if openToken != nil {
							if toAppend != "" {
								ret = append(ret, token.TextToken{toAppend})
							}
							ret = append(ret, openToken)
						}

						bodyExpId := l.bodyExpIds[name][expNum]
						if matcher.Exprs()[expNum].Flags&(NoParseInner|BodyAsArg) == 0 && bodyExpId != 0 {
							ret = append(ret, l.Tokenize(data[indices[bodyExpId*2]:indices[bodyExpId*2+1]])...)
						} else {
							toAppend += data[indices[bodyExpId*2]:indices[bodyExpId*2+1]]
						}

						if closeToken != nil {
							if toAppend != "" {
								ret = append(ret, token.TextToken{toAppend})
								toAppend = ""
							}
							ret = append(ret, closeToken)
						}

						data = data[indices[i*2+1]:]
					} else {
						toAppend += data[:1]
						data = data[1:]
					}
				}
			}
		}
	}

	return ret
}
