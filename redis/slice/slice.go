package slice

import "log"

func SliceTest() {
	var slice = []int{1, 2, 3, 4, 5}
	newSlice := slice[1:2]

	log.Printf("newSlice=%v, len=%d, cap=%d", newSlice, len(newSlice), cap(newSlice))
	// 修改共享元素
	newSlice[0] = 10
	// slice的元素也被修改
	log.Println("after modify")
	log.Printf("slice=%v, len=%d, cap=%d", slice, len(slice), cap(slice))             // slice=[1, 10, 3, 4, 5]
	log.Printf("newSlice=%v, len=%d, cap=%d", newSlice, len(newSlice), cap(newSlice)) // len = 1, cap = 4, newSlice = [10, _, _, _]

	// newSlice扩容
	log.Println("after append")
	// append 1个元素, 容量足够会直接影响slice
	newSlice = append(newSlice, 11)
	log.Printf("slice=%v, len=%d, cap=%d", slice, len(slice), cap(slice))             // slice=[1, 10, 11, 4, 5]
	log.Printf("newSlice=%v, len=%d, cap=%d", newSlice, len(newSlice), cap(newSlice)) // len = 2, cap = 4, newSlice = [10, 11, _, _]

	// newSlice扩容
	log.Println("after append")
	// append3个元素, 容量不足会创建一个新的底层数组
	newSlice = append(newSlice, 12, 13, 14)
	log.Printf("slice=%v, len=%d, cap=%d", slice, len(slice), cap(slice))             // slice=[1, 10, 11, 4, 5]
	log.Printf("newSlice=%v, len=%d, cap=%d", newSlice, len(newSlice), cap(newSlice)) // len = 5, cap = 8, newSlice = [10, 11, 12, 13, 14]
}
