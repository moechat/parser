package moeparser

import (
	"html"
)

func Parse(b []byte) ([]byte, error) {
	body := []byte(html.EscapeString(html.UnescapeString(string(b)))) // Unescape first for uniform strings (EscapeString only escapes <, >, &, ', and "

	body, err := BbCodeParse(body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
