package ch3

import (
	"fmt"
	"strings"
)

func MapFunction(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 'a' + 'A'
	} else {
		return r
	}
}

func ChangeString() {
	var str string = "This is a hello world sentence"
	news := strings.Map(MapFunction, str)
	fmt.Println("After map function: ", news)
}
