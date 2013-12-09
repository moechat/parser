package moeparser

import (
	"regexp"
)

// Options for HtmlTags
const (
	Single         = 1 << iota // This makes MoeParser ignore the CloseRe and body (this is equivalent to setting CloseRe to the empty string)
	PossibleSingle             // Interpret as single if there is no closing tag
	HtmlSingle                 // The HTML tag does not have a closing element
	NoParseInner               // This makes MoeParser ignore any tags inside this tags body. It will be ignored if the Single bit is set.
	TagBodyAsArg               // This makes the text inside of tag and passes it as an arg for the output. The text inside will not be parsed.
	NumberArgToPx              // Converts a number to the number + "px" (ie 12 -> 12px)

	// Non-BBCode specific args
	AllowInWord // This makes MoeParser match the tags that don't either start with whitespace or the beginning of a line

	// BBCode specific args
	AllowTagBodyAsFirstArg // This makes the tag body become the first arg if there is no first argument (makes [name]arg0[/name] the same as [name=arg0][/name])
)

// A type that determines what the parser will replace tags it finds with. The Attributes and CssProps are maps that assign a regexp parser group
type HtmlTags struct {
	Options      uint8                 // Options for this tag
	Tags         []string              // HTML tags
	Classes      [][]string            // Classes to give to the HTML elements
	Attributes   []map[int8]string     // HTML tag attributes
	CssProps     []map[int8]string     // CSS Properties
	OutputFunc   func([]string) string // A custom output function; this returns the string to emplace into the HTML.
	InputModFunc func(*[]string)       // A function that takes input and returns input modified (an example use case would be converting a username to a user ID in @tagging)
}

// A MoeMatcher is the general-purpose matcher. If CloseRe is not set, it will behave as if the paired HtmlTags object has the "Single" option set. This will not match tags that are in the middle of a word.
type MoeMatcher struct {
	OpenRe  *regexp.Regexp // Tag open regexp
	CloseRe *regexp.Regexp // Tag close regexp
}

// Tags to find and replace with HTML. One can insert use-case specific tags to this map
// Example custom tag:
/*
var users = map[string]int{"alice": 0, "bob": 1}

MoeMatcher{Options: moeparser.Single, OpenRe: regexp.MustCompile("@(\S+)")}: HtmlTags{
	Tags: []string{"span"},
	Classes: map[string][]string{"span": []string{"at-tag"}},
	Attributes: map[int8][]string{0: "data-uid"}
	// No specific CssProps, but they behave in the same way as Attributes
	OutputFunc: func(args []string) { return "<span class=\"at-tag\" data-uid=\"" + args[0] + "\">" } // This is unnecessary in this case, and here only as an example. Note that OutputFunc, if not nil, makes MoeParser ignore all options but InputModFunc.
	InputModFunc: func(args *[]string) { args[0] = strconv.Itoa(users[args[0]]) }
}
*/
var MoeTags = map[interface{}]HtmlTags{
	MoeMatcher{
		OpenRe:  regexp.MustCompile("`"),
		CloseRe: regexp.MustCompile("`"),
	}: {Options: NoParseInner, Tags: []string{"pre", "code"}},
	MoeMatcher{
		OpenRe:  regexp.MustCompile("\\*"),
		CloseRe: regexp.MustCompile("\\*"),
	}: {Tags: []string{"i"}},
	MoeMatcher{
		OpenRe:  regexp.MustCompile("\\*\\*"),
		CloseRe: regexp.MustCompile("\\*\\*"),
	}: {Tags: []string{"b"}},
}

// BBCode tags are parsed using a specific BBCode parser, but behaves in the same way as a TagMatcher with OpenRe=\[{{.Name}}(?:=(.*))?\] and CloseRe=\[\/{{.Name}}\], so one should regex escape any special characters in Name. One can insert use-case specific BBCode tags to this map.
var BbCodeTags = map[string]HtmlTags{
	"b": {Tags: []string{"b"}},
	"i": {Tags: []string{"i"}},
	"u": {
		Tags:    []string{"span"},
		Classes: [][]string{{"underline"}},
	},
	"pre":  {Tags: []string{"pre"}},
	"code": {Options: NoParseInner, Tags: []string{"pre", "code"}},
	"color": {
		Tags:     []string{"span"},
		CssProps: []map[int8]string{{0: "color"}},
	},
	"colour": {
		Tags:     []string{"span"},
		CssProps: []map[int8]string{{0: "color"}},
	},
	"size": {
		Options:  NumberArgToPx,
		Tags:     []string{"span"},
		CssProps: []map[int8]string{{0: "font-size"}},
	},
	"noparse": {Options: NoParseInner},
	"url": {
		Options:    (AllowTagBodyAsFirstArg | PossibleSingle),
		Tags:       []string{"a"},
		Attributes: []map[int8]string{{0: "href"}},
	},
	"img": {
		Options:    (AllowTagBodyAsFirstArg | TagBodyAsArg | PossibleSingle | HtmlSingle),
		Tags:       []string{"img"},
		Attributes: []map[int8]string{{0: "src", 1: "title"}},
	},
	"s":    {Tags: []string{"s"}},
	"samp": {Tags: []string{"samp"}},
	"q":    {Tags: []string{"q"}},
}
