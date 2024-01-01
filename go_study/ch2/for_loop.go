package ch2

import "fmt"

func PrintCharator(line int, ch string) {
	for i := 1; i <= line; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%s ", ch)
		}
		fmt.Println()
	}
}
