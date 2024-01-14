package ch3

import "fmt"

func trace(s string) string {
	fmt.Println("entering: ", s)
	return s
}

func untrace(s string) string {
	fmt.Println("leaving: ", s)
	return s
}

func Function2() {
	defer untrace(trace("Function2"))
	fmt.Println("in Function2")
}

func Function1() {
	defer untrace(trace("Function1"))
	fmt.Println("in Function1")
	Function2()
}
