package lexer_test

import (
	"."
	"fmt"
	"github.com/moechat/moeparser"
	"github.com/moechat/moeparser/token"
	"github.com/moechat/moeparser/token/htmltoken"
)

func ExampleLexer_AddTokenClass() {
	// The ID's are strings because I'm too lazy to use strconv (and not lazy enough to document the reason, apparently)
	users := map[string]string{"alice": "0", "bob": "1"}

	matcher := moeparser.NewTokenClass(
		moeparser.TokenClassArgs{
			ArgModFunc: func(args []string, idByName map[string]int) ([]string, map[string]int) {
				uid := users[args[idByName["username"]]]
				uidIndex := len(args)
				args = append(args, uid)
				idByName["uid"] = uidIndex
				return args, idByName
			},
			IsValid: func(args *token.TokenArgs) bool {
				return args.ByName("uid") != ""
			},

			Tokens: []token.Token{
				&htmltoken.Token{
					Name:       "span",
					Type:       token.SingleType, // This is the default

					Prefix: "{{.ById 1}}",
					Classes:    []string{"at-tag", "{{.ByName uid}}"},
					Attributes: map[string]string{
						"data-user": "{{.ById 2}}",
						"data-uid": "{{.ByName uid}}",
					},
					// CssProps is not necessary here, but behaves in the same way as Attributes

					// This is overkill in this case, and here only as an example. Note that OutputFunc, if not nil, makes MoeParser ignore all other options.
					OutputFunc: func(args *token.TokenArgs) string {
						output := args.ById(1)
						output += `<span class="at-tag user-` + args.ByName("uid") + `"`
						output += `data-uid="` + args.ByName("uid") + `" `
						output += `data-user="` + args.ById(2) + `">`
						output += args.ById(2)
						output += `</span>`
						return output
					},
				},
			},
		},
		`(^|\s)@(?P<username>\S+)`)

	l, _ := lexer.New(make(map[string]token.TokenClass))

	l.AddTokenClass(matcher)

	l.CompileRegexp() // This is necessary for Tokenize to take into account changes due to AddMoeMatcher or RemoveMoeMatcher

	for _, token := range l.Tokenize("@alice  @bob some text (non-matching)) \n@charlie\n\t@hi\n\r\r\t\t\tarstarts@bob doesn't match.\t@alice") {
		t, err := token.Output()
		if err != nil {
			panic(err)
		}
		fmt.Print(t)
	}

	fmt.Println()
}
