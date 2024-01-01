package fixbytespool

import "testing"

func BenchmarkByteFixPool(b *testing.B) {
	blocks := 100
	blockSize := 64
	block := make([]byte, blockSize)
	pool := NewByteFixPool(16, 0, blocks*blockSize)
	for n := 0; n < b.N; n++ {
		b := pool.Get()
		for i := 0; i < blocks; i++ {
			b = append(b, block...)
		}
		pool.Put(b)
	}
}
