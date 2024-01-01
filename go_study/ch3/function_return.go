package ch3

func Adder(a int) func(b int) int {
	function := func(b int) int {
		return a + b
	}
	return function
}
