package ch3

import "fmt"

func RecursivePrintK(k int) {
	if k == 1 {
		fmt.Print(k)
		return
	}
	RecursivePrintK(k - 1)
	fmt.Println(k)
}
