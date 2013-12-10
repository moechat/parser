package moeparser_test

import (
	"fmt"
	"moeparser"
)

func ExampleAddMatcher(t *testing.T) {
	users := map[string]int{"alice": 0, "bob": 1}

	matcher := MoeParser.MoeMatcher{Options: moeparser.Single, OpenRe: regexp.MustCompile("@(\\S+)")}
	tags := MoeParser.HtmlTags{
		Tags:       []string{"span"},
		Classes:    map[string][]string{"span": {"at-tag"}},
		Attributes: map[int8][]string{0: "data-uid"},
		// No specific CssProps, but they behave in the same way as Attributes
		OutputFunc: func(args []string) {
			return "<span class=\"at-tag\" data-uid=\"" + args[0] + "\">"
		}, // This is unnecessary in this case, and here only as an example. Note that OutputFunc, if not nil, makes MoeParser ignore all options but InputModFunc.
		InputModFunc: func(args *[]string) {
			args[0] = strconv.Itoa(users[args[0]])
		},
	}

	moeparser.AddMoeMatcher(matcher, tags)

	moeparser.Lexer.BuildCharTree() // This is necessary for Tokenize to take into account changes due to AddMoeMatcher or RemoveMoeMatcher

	fmt.Println(moeparser.Parse("@alice"))
}
