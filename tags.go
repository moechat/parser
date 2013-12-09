package moeparser

import (
	"errors"
	"fmt"
)

// A type that determines what the parser will replace tags it finds with. The Attributes and CssProps are maps that assign a regexp parser group
type HtmlTags struct {
	Options      MatcherOptions        // Compatibility options for BBCode until token parsing is complete
	Tags         []string              // HTML tags
	Classes      [][]string            // Classes to give to the HTML elements
	Attributes   []map[int8]string     // HTML tag attributes
	CssProps     []map[int8]string     // CSS Properties
	OutputFunc   func([]string) string // A custom output function; this returns the string to emplace into the HTML.
	InputModFunc func(*[]string)       // A function that takes input and returns input modified (an example use case would be converting a username to a user ID in @tagging)
}

type HtmlTag struct {
	Tag string // The HTML Tag

	Classes    []string        // HTML element classes
	Attributes map[int8]string // HTML element attributes
	CssProps   map[int8]string // CSS Properties

	// A custom output function:
	// if set, the parser will ignore all fields and simply spit out the output.
	// The input is the array of strings that were captured
	// Use a $b to represent the body of the tag and $$ for a dollar sign
	OutputFunc func([]string) string
}

type MatcherOptions int

// Options for MoeMatcher/BbTag - set these bits in the Options field on a MoeMatcher or BbTag to implement this behavior
const (
	Single         MatcherOptions = 1 << iota // This makes MoeParser ignore the CloseRe and body (this is equivalent to setting CloseRe to the empty string)
	PossibleSingle                            // Interpret as single if there is no closing tag
	HtmlSingle                                // The HTML tag does not have a closing element
	NoParseInner                              // This makes MoeParser ignore any tags inside this tags body. It will be ignored if the Single bit is set.
	TagBodyAsArg                              // This makes the text inside of tag and passes it as an arg for the output. The text inside will not be parsed.
	NumberArgToPx                             // Converts a number to the number + "px" (ie 12 -> 12px)

	// Non-BBCode specific args
	AllowInWord // This makes MoeParser match the tags that don't either start with whitespace or the beginning of a line

	// BBCode specific args
	AllowTagBodyAsFirstArg // This makes the tag body become the first arg if there is no first argument (makes [name]arg0[/name] the same as [name=arg0][/name])
)

// A pair of a starting and ending tag.
// If End is the zero value, Start will be used as the closing tag.
// The lexer understands golang's regexp syntax (godoc regexp/syntax).
type TagPair struct {
	Start string // The string that starts the tag body.
	End   string // The string that closes the tag body
}

// Returns a symmetric tag, returning nil if tag == ""
func SymTag(tag string) TagPair {
	if tag == "" {
		return TagPair{}
	}
	return TagPair{Start: tag}
}

// Returns a new tag with Start = start and End = end, catching empty start (returns nil) and
// start == end (returns Start = start and End = "")
func Tag(start string, end string) TagPair {
	if start == "" {
		return TagPair{}
	}
	if start == end {
		return TagPair{Start: start}
	}
	return TagPair{start, end}
}

// Returns a TagPair that will match the bbcode [name]
func BbTag(name string) TagPair {
	if name == "" {
		return TagPair{}
	}

	start := fmt.Sprintf("[%s(=(?!]))]", name)
	end := fmt.Sprintf("[//%s]", name)

	return TagPair{start, end}
}

// A MoeMatcher is the general-purpose matcher structure. This will not match tags that are in the middle of a word unless AllowInWord is set.
type Matcher struct {
	Options MatcherOptions // Options for this matcher
	Tags    []TagPair      // Tag pairs for this MoeMatcher

	InputModFunc func(*[]string) // A function that takes input and returns input modified (an example use case would be converting a username to a user ID in @tagging)

	HtmlTags []HtmlTag // The HTML tags to insert

	id byte // An internal id set by the token builder
}

// The default matchers.
//
// This is here only for reference - changing this has no effect.
// Use AddMatcher and RemoveMatcher instead.
var DefaultMatcherMap = map[TagPair]*Matcher{
	SymTag("`"): {
		Tags:    []TagPair{SymTag("`")},
		Options: NoParseInner,
	},
	SymTag("*"): {
		Tags: []TagPair{SymTag("*")},
	},
	SymTag("**"): {
		Tags: []TagPair{SymTag("**")},
	},
}

var matcherMap = DefaultMatcherMap

var matchers map[string]*Matcher

// Adds the passed matcher to the matcher list.
// Returns matcher.Start, which can be used to remove the matcher if necessary
func AddMatcher(matcher *Matcher) error {
	for _, tagPair := range matcher.Tags {
		if _, ok := matcherMap[tagPair]; ok {
			err := errors.New("A matcher with an identical tag pair has already been inserted!")
			return err
		}
	}

	for _, tagPair := range matcher.Tags {
		matcherMap[tagPair] = matcher
	}

	return nil
}

// Adds a list of matchers to the matcher map
func AddMatchers(matchers []*Matcher) error {
	for _, matcher := range matchers {
		err := AddMatcher(matcher)
		if err != nil {
			return err
		}
	}

	return nil
}

// Removes the matcher whose start tag is startTag from the matcher list.
func RemoveMatcher(matcher *Matcher) {
	for i, m := range matcherMap {
		if m == matcher {
			delete(matcherMap, i)
		}
	}
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
