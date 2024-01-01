package bytes_pool

import (
	"bytes"
	"sync"
)

type BytePool struct {
	p sync.Pool
}

func NewBytePool(size, cap int) *BytePool {
	if size > cap {
		panic("size must be less than cap")
	}
	p := &BytePool{}
	// 设置New函数
	p.p.New = func() any {
		var b []byte
		if cap > 0 {
			b = make([]byte, size, cap)
		}
		return bytes.NewBuffer(b)
	}
	return p
}

func (p *BytePool) Get() *bytes.Buffer {
	return p.p.Get().(*bytes.Buffer)
}

// 在Put的时候重置字节数组的已用空间（这样下次才能从头开始使用）
func (p *BytePool) Put(b *bytes.Buffer) {
	// 重置已用大小
	b.Reset()
	p.p.Put(b)
}
