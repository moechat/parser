package token

type TokenOptions uint8

// Options for Tokens - return these bits in GetOptions() to implement this behavior
const (
	PossibleSingle         TokenOptions = 1 << iota // Interpret as single if there is no closing tag
	HtmlSingle                                      // The HTML tag does not have a closing element
	NoParseInner                                    // This makes MoeParser ignore any tags inside this tags body. It will be ignored if the Single bit is set.
	TagBodyAsArg                                    // This makes the text inside of tag and passes it as an arg for the output. The text inside will not be parsed.
	NumberArgToPx                                   // Converts a number to the number + "px" (ie 12 -> 12px)
	AllowInWord                                     // This makes MoeParser match the tags that don't either start with whitespace or the beginning of a line
	AllowTagBodyAsFirstArg                          // This makes the tag body become the first arg if there is no first argument (makes [name]arg0[/name] the same as [name=arg0][/name])
)

type TokenClassType uint8

// Token class types - return these in GetType()
const (
	Single     TokenClassType = iota // A single token with no opening or closing token
	PossSingle                       // An open token that can be single if there is no matching closing token
	Open                             // A token that starts a section
	Close                            // A token that ends a section
	OpenClose                        // A token that both begins and ends a section
)

// A token class recognized by the lexer
type TokenClass interface {
	GetOptions() TokenOptions // Returns options for the token class
	GetType() TokenClassType  // Returns the type of the token class

	GetTokens([]string, bool) []Token // Returns instances of the tokenClass with args set
}

// A token returned by the lexer and recognized by the parser
type Token interface {
	Copy() Token // Get a copy of the token

	SetArgs([]string)           // Set the args of the token
	GetOutput() (string, error) // Get the output of the token
}

// An implementation of Token used to represent text that isn't matched by any other tokens
// i.e. "hi" in <p>hi</p>
type plainText struct {
	Body string
}

func (pt *plainText) setBody(body string) {
	pt.Body = body
}

func (pt *plainText) Copy() *plainText {
	return &plainText{Body: pt.Body}
}

// Returns the PlainText's body
func (pt *plainText) GetOutput() (string, error) {
	return pt.Body, nil
}
