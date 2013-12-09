package moeparser

import (
	"errors"
)

// A type that determines what the parser will replace tags it finds with. The Attributes and CssProps are maps that assign a regexp parser group
type HtmlTags struct {
	Options      TagOptions            // Compatibility options for BBCode until token parsing is complete
	Tags         []string              // HTML tags
	Classes      [][]string            // Classes to give to the HTML elements
	Attributes   []map[int8]string     // HTML tag attributes
	CssProps     []map[int8]string     // CSS Properties
	OutputFunc   func([]string) string // A custom output function; this returns the string to emplace into the HTML.
	InputModFunc func(*[]string)       // A function that takes input and returns input modified (an example use case would be converting a username to a user ID in @tagging)
}

type TagOptions int

// Options for MoeMatcher/BbTag - set these bits in the Options field on a MoeMatcher or BbTag to implement this behavior
const (
	Single         TagOptions = 1 << iota // This makes MoeParser ignore the CloseRe and body (this is equivalent to setting CloseRe to the empty string)
	PossibleSingle                        // Interpret as single if there is no closing tag
	HtmlSingle                            // The HTML tag does not have a closing element
	NoParseInner                          // This makes MoeParser ignore any tags inside this tags body. It will be ignored if the Single bit is set.
	TagBodyAsArg                          // This makes the text inside of tag and passes it as an arg for the output. The text inside will not be parsed.
	NumberArgToPx                         // Converts a number to the number + "px" (ie 12 -> 12px)

	// Non-BBCode specific args
	AllowInWord // This makes MoeParser match the tags that don't either start with whitespace or the beginning of a line

	// BBCode specific args
	AllowTagBodyAsFirstArg // This makes the tag body become the first arg if there is no first argument (makes [name]arg0[/name] the same as [name=arg0][/name])
)

// A MoeMatcher is the general-purpose matcher. If CloseRe is not set, it will behave as if the paired HtmlTags object has the "Single" option set. This will not match tags that are in the middle of a word unless AllowInWord is set.
type MoeMatcher struct {
	Options TagOptions // Options for this matcher
	Start   string     // The string that starts the tag body. This should be a unique value, and is used as they key for internal maps.
	End     string     // The string that closes the tag body
}

var moeMatchers = map[string]*MoeMatcher{
	"`": {
		Start:   "`",
		End:     "`",
		Options: NoParseInner,
	},
	"*": {
		Start: "*",
		End:   "*",
	},
	"**": {
		Start: "**",
		End:   "**",
	},
}

var moeTags = map[string]*HtmlTags{
	"`":  {Tags: []string{"pre", "code"}},
	"*":  {Tags: []string{"i"}},
	"**": {Tags: []string{"b"}},
}

// Adds the passed matcher to the matcher list.
// Returns matcher.Start, which can be used to remove the matcher if necessary
func AddMoeMatcher(matcher *MoeMatcher, htmlTags *HtmlTags) (string, error) {
	if _, ok := moeMatchers[matcher.Start]; ok {
		err := errors.New("A matcher starting with " + matcher.Start + " already exists!")
		return "", err
	}

	moeTags[matcher.Start] = htmlTags
	moeTags[matcher.Start] = htmlTags
	return matcher.Start, nil
}

func AddMoeMatchers(matcherMap map[*MoeMatcher]*HtmlTags) (map[*MoeMatcher]string, error) {
	ret := make(map[*MoeMatcher]string, len(matcherMap))
	for matcher, htmlTags := range matcherMap {
		key, err := AddMoeMatcher(matcher, htmlTags)
		if err != nil {
			return nil, err
		}

		ret[matcher] = key
	}

	return ret, nil
}

// Removes the matcher whose start tag is startTag from the matcher list.
func RemoveMoeMatcher(startTag string) {
	delete(moeTags, startTag)
	delete(moeMatchers, startTag)
}

var bbCodeTags = map[string]HtmlTags{
	"b": {Tags: []string{"b"}},
	"i": {Tags: []string{"i"}},
	"u": {
		Tags:    []string{"span"},
		Classes: [][]string{{"underline"}},
	},
	"pre":  {Options: NoParseInner, Tags: []string{"pre"}},
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

// One can insert use-case specific BBCode tags by using this function.
// BBCode tags are parsed using a specific BBCode parser.
//
// IMPORTANT: This is ignored by moeparser.Parse - you should use AddMoeMatcher instead!
// Only use this function if you plan on using moeparser.BbCodeParse()!
//
// This will be deprecated in the future after BBCode functionality is added to AddMoeMatcher.
func AddBbTag(tag string, htmlTags HtmlTags) {
	bbCodeTags[tag] = htmlTags
}
