package lexer

import (
	"errors"
	"github.com/moechat/moeparser/token"
)

// The Lexer, which tokenizes text
type Lexer struct {
	TokenClasses map[string]token.TokenClass // The possible token classes

	ExprTreeRoot *treeNode // The root node for the expression tree
}

// Creates a new Lexer with the token map given.
// BuildExprTree is run automatically.
func New(tokenClasses map[string]token.TokenClass) *Lexer {
	ret := &Lexer{tokenClasses, nil}
	ret.BuildExprTree()
	return ret
}

// A node in the possible character tree
type treeNode struct {
	expr     string               // The expression that matches this treeNode
	children map[string]*treeNode // Possible next characters

	t token.TokenClass // The Token class that the treeNode is the end of or nil if it is not the end of a token
}

// Adds the passed token class to the token class map.
func (l *Lexer) AddTokenClass(exprs []string, tokenClass token.TokenClass) error {
	for _, expr := range exprs {
		if _, ok := l.TokenClasses[expr]; ok {
			err := errors.New("A token with an identical tag pair has already been inserted!")
			return err
		}
	}

	for _, expr := range exprs {
		l.TokenClasses[expr] = tokenClass
	}

	return nil
}

// Removes the token class from the token class map.
func (l *Lexer) RemoveTokenClass(tokenClass token.TokenClass) {
	for i, t := range l.TokenClasses {
		if t == tokenClass {
			delete(l.TokenClasses, i)
		}
	}
}

// Returns an empty treeNode with Children intialized
func rootNode() *treeNode {
	return &treeNode{children: make(map[string]*treeNode)}
}

func newNode(expr string) *treeNode {
	return &treeNode{expr: expr, children: make(map[string]*treeNode)}
}

func (node *treeNode) setTokenClass(t token.TokenClass) {
	node.t = t
}

func getGroups(expr string) []string {
	ret := make([]string, len(expr))

	for i, c := range expr {
		ret[i] = string(c)
	}

	return ret
}

// This function builds an ordered tree of possible characters which is used by Tokenize().
//
// Although it is run when the Lexer is first created,
// It *must* run again in order for changes made by AddTokenClass or RemoveTokenClass to take effect.
func (l *Lexer) BuildExprTree() {
	l.ExprTreeRoot = rootNode()
	for expr, token := range l.TokenClasses {
		grps := getGroups(expr)

		currNode := l.ExprTreeRoot
		for _, grp := range grps {
			if node, ok := currNode.children[grp]; ok {
				currNode = node
			} else {
				node := newNode(grp)
				currNode.children[grp] = node
				currNode = node
			}
		}

		currNode.setTokenClass(token)
	}
}

// Converts an input []byte into Tokens.
//
// If the tree has not been built, tokenize will run BuildCharTree()
// BuildCharTree() *must* be run if you call AddTokenClass or RemoveTokenClass between Tokenize()'s
func (l *Lexer) Tokenize(data []byte) []token.Token {
	ret := make([]token.Token, 100)

	return ret
}
