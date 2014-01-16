package lexer_test

import (
	"."
	"fmt"
	"github.com/moechat/parser/token"
	"testing"
)

type TestMatcher struct {
	name  string
	exprs []lexer.Expression
}

func (tm *TestMatcher) Name() string {
	return tm.name
}

func (tm *TestMatcher) Exprs() []lexer.Expression {
	return tm.exprs
}

func (tm *TestMatcher) IsValid(args *token.TokenArgs, expNum int) bool {
	if tm.name == "image" {
		if args.ById(1) == "" {
			return false
		}
	}
	return true
}

func (tm *TestMatcher) BuildToken(args *token.TokenArgs, expNum int) (token.Token, token.Token) {
	if tm.name == "image" {
		url := args.ById(1)
		title := url
		if expNum == 2 {
			title = args.ById(2)
		}
		// Yes, this would be unsafe in a production environment. But it's a testing script.
		return token.TextToken{fmt.Sprintf(`<img src="%s" title="%s">`, url, title)}, nil
	} else if tm.name == "bold" {
		return token.TextToken{"<b>"}, token.TextToken{"</b>"}
	}
	return nil, nil
}

func TestLexer(*testing.T) {
	l := lexer.Must(lexer.New(
		&TestMatcher{"bold", []lexer.Expression{
			{Expr: `\[b\]`, CloseExpr: `\[/b\]`},
		}},
		&TestMatcher{"image", []lexer.Expression{
			{Expr: `\[img=(.*?)\]`},
			{Expr: `\[img]`, CloseExpr: `\[/img\]`, Flags: lexer.BodyAsArg},
			{Expr: `\[img=(.*?)\]`, CloseExpr: `\[/img\]`, Flags: lexer.BodyAsArg},
		}},
		&TestMatcher{"noparse", []lexer.Expression{
			{Expr: `\[nope\]`, CloseExpr: `\[/nope\]`, Flags: lexer.NoParseInner},
		}},
	))

	toTokenize := `[b][img=http://image.com/cool.png][/b][b]hi
[nope][img]http://fail.com/fun.png[/img][/nope]great.
`

	for _, t := range l.Tokenize(toTokenize) {
		fmt.Print(t.(token.TextToken).Body)
	}
	fmt.Println()

	// OUTPUT:
	// <b><img src="http://image.com/cool.png" title="http://image.com/cool.png"></b><b>hi
	// [img]http://fail.com/fun.png[/img]great.
	// </b>
}
