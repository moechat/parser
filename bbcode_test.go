package moeparser

import (
	"fmt"
	"testing"
)

func TestBbCodeParse(t *testing.T) {
	testString := "[b][i][img=http://sauyon.com/blah.png][/i]"
	out, err := BbCodeParse([]byte(testString))
	if err != nil {
		fmt.Printf("Parsing failed! Error: %v", err)
	}
	fmt.Printf("Parse succeeeded. Output is: %s\n", out)
}
