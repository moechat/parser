package bbcode_test

import (
	"."
	"fmt"
	"testing"
)

func TestBbCodeParse(t *testing.T) {
	testString := `Something [b][i][img=http://sauyon.com/blah.png][/i][/i] [b]hi[/b=what]
This seems to work very well :D [url][/url] [size=12px]something
[url]http://sauyon.com/wrappedinurl[/url] [url=http://sauyon.com] [img=//sauyon.com]What[/img] [url][/url]`

	out, err := bbcode.Parse(testString)
	if err != nil {
		fmt.Printf("Parsing failed! Error:", err)
	}
	fmt.Println("Parse succeeeded. Output is:")
	fmt.Println(out)

	testString1 := "[url=http://google.com/][img]http://www.google.com/intl/en_ALL/images/logo.gif[/img][/url]"
	out1, err := bbcode.Parse(testString1)
	if err != nil {
		fmt.Printf("Parsing failed! Error:", err)
	}
	fmt.Println("Parse succeeeded. Output is:")
	fmt.Println(out1)
}
