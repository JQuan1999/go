package ch3

import "fmt"

func Clouse1() {
	var local_param int = 1
	function := func(param, k int) {
		for i := 1; i <= k; i++ {
			param += i
		}
		fmt.Println("param: ", param)
	}
	function(local_param, 10)
	fmt.Println("local_param: ", local_param)
}

func Clouse2() {
	var local_param int = 1
	function := func(k int) {
		for i := 1; i <= k; i++ {
			local_param += i
		}
		fmt.Println("param: ", local_param)
	}
	function(10)
	fmt.Println("local_param: ", local_param)
}
