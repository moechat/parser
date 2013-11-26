package moeparser

import (
	"fmt"
	"testing"
)

func TestBbCodeParse(t *testing.T) {
	testString :=
		`Something [b][i][img=http://sauyon.com/blah.png][/i][/i] [b]hi[/b=what]
This seems to work very well :D [url][/url] [size=12px]something [url]http://sauyon.com[/url] [url=http://sauyon.com] [img] [url][/url]`
	out, err := BbCodeParse([]byte(testString))
	if err != nil {
		fmt.Printf("Parsing failed! Error: %v", err)
	}
	fmt.Printf("Parse succeeeded. Output is: %s\n", out)
}
