package bytes_pool

import "testing"

func BenchmarkByte(b *testing.B) {
	blocks := 100
	blockSize := 64
	block := make([]byte, blockSize)
	for n := 0; n < b.N; n++ {
		var b []byte
		for i := 0; i < blocks; i++ {
			b = append(b, block...)
		}
	}
}

func BenchmarkMake(b *testing.B) {
	blocks := 100
	blockSize := 64
	block := make([]byte, blockSize)
	for n := 0; n < b.N; n++ {
		b := make([]byte, 0, blocks*blockSize)
		for i := 0; i < blocks; i++ {
			b = append(b, block...)
		}
	}
}

// 这里我们每次先从字节池拿一个字节数组Get()，使用完之后归还字节池Put()。
func BenchmarkBytePool(b *testing.B) {
	blocks := 100
	blockSize := 64
	block := make([]byte, blockSize)
	pool := NewBytePool(0, blocks*blockSize)
	for n := 0; n < b.N; n++ {
		b := pool.Get() // buffer池
		for i := 0; i < blocks; i++ {
			b.Write(block)
		}
		// 归还
		pool.Put(b)
	}
}
