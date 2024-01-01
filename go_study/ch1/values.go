package ch1

import "fmt"

func value() {
	fmt.Println("go" + "lang")
	fmt.Println("1+1=", 1+1)
	fmt.Println("7.0/3.0 =", 7.0/3.0)

	// 布尔型，以及常见的布尔操作。
	fmt.Println(true && false)
	fmt.Println(true || false)
	fmt.Println(!true)
}
