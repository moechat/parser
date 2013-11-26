package moeparser

import (
	"html"
)

func Parse(b string) (string, error) {
	body := html.EscapeString(html.UnescapeString(b)) // Unescape first for uniform strings (EscapeString only escapes <, >, &, ', and "

	body, err := BbCodeParse(body)
	if err != nil {
		return "", err
	}

	return body, nil
}
