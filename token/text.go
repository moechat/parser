package token

// An implementation of Token used to represent text that isn't matched by any other tokens
// i.e. "hi" in <p>hi</p>
type Text struct {
	body string
}

func NewText(body string) *Text {
	return &Text{body: body}
}

func (pt *Text) setBody(body string) {
	pt.body = body
}

func (pt *Text) Copy() Token {
	return &Text{body: pt.body}
}

func (pt *Text) SetArgs(*TokenArgs) {}

// Returns the Text's body
func (pt *Text) Output() (string, error) {
	return pt.body, nil
}
