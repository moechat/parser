package bbcode

import (
	"github.com/moechat/parser"
)

// A type that determines what the parser will replace tags it finds with. The Attributes and CssProps are maps that assign a regexp parser group
type HtmlTags struct {
	Options      int                   // Compatibility options for BBCode until token parsing is complete
	Tags         []string              // HTML tags
	Classes      [][]string            // Classes to give to the HTML elements
	Attributes   []map[int8]string     // HTML tag attributes
	CssProps     []map[int8]string     // CSS Properties
	OutputFunc   func([]string) string // A custom output function; this returns the string to emplace into the HTML.
	InputModFunc func(*[]string)       // A function that takes input and returns input modified (an example use case would be converting a username to a user ID in @tagging)
}

var bbCodeTags = map[string]HtmlTags{
	"b": {Tags: []string{"b"}},
	"i": {Tags: []string{"i"}},
	"u": {
		Tags:    []string{"span"},
		Classes: [][]string{{"underline"}},
	},
	"pre":  {Options: parser.NoParseInner, Tags: []string{"pre"}},
	"code": {Options: parser.NoParseInner, Tags: []string{"pre", "code"}},
	"color": {
		Tags:     []string{"span"},
		CssProps: []map[int8]string{{0: "color"}},
	},
	"colour": {
		Tags:     []string{"span"},
		CssProps: []map[int8]string{{0: "color"}},
	},
	"size": {
		Options:  parser.NumberArgToPx,
		Tags:     []string{"span"},
		CssProps: []map[int8]string{{0: "font-size"}},
	},
	"noparse": {Options: parser.NoParseInner},
	"url": {
		Options:    (parser.AllowTokenBodyAsFirstArg | parser.PossibleSingle),
		Tags:       []string{"a"},
		Attributes: []map[int8]string{{0: "href"}},
	},
	"img": {
		Options: (parser.AllowTokenBodyAsFirstArg |
			parser.TokenBodyAsArg |
			parser.PossibleSingle |
			parser.HtmlSingle),
		Tags:       []string{"img"},
		Attributes: []map[int8]string{{0: "src", 1: "title"}},
	},
	"s":    {Tags: []string{"s"}},
	"samp": {Tags: []string{"samp"}},
	"q":    {Tags: []string{"q"}},
}

// One can insert use-case specific BBCode tags by using this function.
// BBCode tags are parsed using a specific BBCode parser.
//
// IMPORTANT: This is ignored by parser.Parse - you should use AddTokenClass instead!
// Only use this function if you plan on using parser.BbCodeParse()!
//
// This will be deprecated in the future after BBCode functionality is added to AddMatcher.
func AddBbToken(name string, htmlTags HtmlTags) {
	bbCodeTags[name] = htmlTags
}
