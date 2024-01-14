package bytesbuffer

import (
	"fmt"
)

func BytesLenTest() {
	slice := make([]byte, 0, 10)
	slice = append(slice, []byte("abc")...)
	fmt.Printf("slice len: %d, cap: %d, slice= %s\n", len(slice), cap(slice), slice) // len = 3, cap = 10, abc

	newSlice := slice[len(slice):]
	fmt.Printf("newSlice len: %d, cap: %d, slice= %s\n", len(newSlice), cap(newSlice), newSlice) // len = 0, cap = 10, []

	// copy
	// copy(slice, "newslice")                                                          // copy函数拷贝数据到src[:len(src)]
	// fmt.Printf("slice len: %d, cap: %d, slice= %s\n", len(slice), cap(slice), slice) // len = 0, cap = 10, []

	// append
	newSlice = append(newSlice, []byte("newslice1")...)
	fmt.Println("after append over cap slice") // slice and newslice don't share the common array
	fmt.Printf("slice len: %d, cap: %d, slice= %s\n", len(slice), cap(slice), slice)
	fmt.Printf("newSlice len: %d, cap: %d, slice= %s\n", len(newSlice), cap(newSlice), newSlice)
}
