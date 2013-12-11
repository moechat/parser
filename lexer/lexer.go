// This package defines the Lexer type, which is used to convert raw tags into tokens (http://en.wikipedia.org/wiki/Lexical_analysis)
//
// It is implemented using the regexp package, which is implemented using a NFA (http://en.wikipedia.org/wiki/Nondeterministic_finite_automaton), and is thus efficient in this use case.
package lexer

import (
	"errors"
	"fmt"
	"github.com/moechat/moeparser/token"
	"regexp"
	"strconv"
)

// The lexer is an object that contains the definition
type Lexer struct {
	tokenClasses   map[string]token.TokenClass // The map of TokenClasses by name
	tokenClassById map[int]token.TokenClass    // The map of TokenClasses by capture group ID

	exprs map[string]bool // The set of expressions that are mapped to some token class

	expr   string         // The main regexp expression
	regexp *regexp.Regexp // The regexp used to match tags
}

// Creates a new Lexer with the token map given.
func New(tokenClasses map[string]token.TokenClass) (*Lexer, error) {
	ret := &Lexer{
		tokenClasses: tokenClasses,
		exprs:        make(map[string]bool),
	}
	for name, tokenClass := range ret.tokenClasses {
		if tokenClass.Name() != name {
			err := errors.New("A token class whose index did not match Name() was passed!")
			return nil, err
		}
		for _, re := range tokenClass.Regexps() {
			ret.exprs[re.String()] = true
		}
	}
	err := ret.CompileRegexp()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Adds the passed token class to the token class map.
func (l *Lexer) AddTokenClass(tokenClass token.TokenClass) error {
	if _, ok := l.tokenClasses[tokenClass.Name()]; ok {
		err := errors.New("A TokenClass with the same name already exists!")
		return err
	}
	regexps := tokenClass.Regexps()
	for _, re := range regexps {
		if l.exprs[re.String()] {
			err := errors.New("Another TokenClass has an identical regexp!")
			return err
		}
	}

	for _, re := range regexps {
		l.exprs[re.String()] = true
	}

	l.tokenClasses[tokenClass.Name()] = tokenClass

	return nil
}

// Removes all instances of the token class specified from the token class map.
func (l *Lexer) RemoveTokenClass(tokenClass token.TokenClass) {
	for _, re := range tokenClass.Regexps() {
		delete(l.exprs, re.String())
	}

	delete(l.tokenClasses, tokenClass.Name())
}

// Creates the and compiles the regexp used by the Lexer.
//
// It must be run after adding or removing token classes in order for changes to take effect.
func (l *Lexer) CompileRegexp() error {
	l.tokenClassById = make(map[int]token.TokenClass)

	var err error

	l.expr = "(?:"
	for name, tokenClass := range l.tokenClasses {
		for i, re := range tokenClass.Regexps() {
			l.expr += fmt.Sprintf("(?P<_%02x%s>%s)|", i, name, re.String())
		}
	}
	l.expr = l.expr[:len(l.expr)-1]
	l.expr += ")"

	l.regexp, err = regexp.Compile(l.expr)
	if err != nil {
		return err
	}

	names := l.regexp.SubexpNames()
	usedNames := make(map[string]bool)

	for i, name := range names {
		if name != "" {
			if name[0] == '_' {
				if _, ok := l.tokenClasses[name[3:]]; usedNames[name] || !ok {
					return errors.New("lexer: capture group names starting with _ are reserved for use by the lexer! Your name is " + name)
				}

				usedNames[name] = true
				l.tokenClassById[i] = l.tokenClasses[name[3:]]
			}
		}
	}

	return nil
}

// Converts an input string into Tokens.
//
// If the tree has not been built, tokenize will run BuildCharTree().
// BuildCharTree() *must* be run if you call AddTokenClass or RemoveTokenClass between Tokenize()'s
func (l *Lexer) Tokenize(data string) []token.Token {
	ret := make([]token.Token, 0)
	subexpNames := l.regexp.SubexpNames()

	for data != "" {
		indices := l.regexp.FindStringSubmatchIndex(data)
		if indices == nil {
			ret = append(ret, token.NewText(data))
			data = ""
			break
		}

		for i, tokenClass := range l.tokenClassById {
			if indices[i*2] >= 0 {
				if indices[i*2] != 0 {
					ret = append(ret, token.NewText(data[:indices[i*2]]))
				}

				exprId, _ := strconv.ParseInt(subexpNames[i][1:3], 16, 8)
				tokenRegexp := tokenClass.Regexps()[exprId]
				args := []string(tokenRegexp.FindStringSubmatch(data[indices[i*2]:indices[i*2+1]]))

				idByName := make(map[string]int)
				for i, name := range tokenRegexp.SubexpNames() {
					if name != "" {
						idByName[name] = i
					}
				}

				args, idByName = tokenClass.ModifyArgs(args, idByName)
				tokenArgs := token.NewTokenArgs(args, idByName)

				if tokenClass.IsValid(tokenArgs) {
					ret = append(ret, tokenClass.BuildTokens(tokenArgs)...)
				} else {
					ret = append(ret, token.NewText(data[indices[i*2]:indices[i*2+1]]))
				}
				data = data[indices[i*2+1]:]
			}
		}
	}

	return ret
}
