package fixbytespool

// 限制对象池的大小, 通过channel+select, 申请一个固定长度缓冲区的channel
// Put: channel不满则put，否则default分支丢弃
// Get: channel不空则get, 否则default分支申请新对象
type ByteFixPool struct {
	cache chan []byte
	size  int
	cap   int
}

// cacheSize: 字节池缓存长度
// size: 字节数组长度
// cap: 字节数组容量
func NewByteFixPool(cacheSize, size, cap int) *ByteFixPool {
	if size > cap {
		panic("size must be less than cap")
	}
	return &ByteFixPool{cache: make(chan []byte, cacheSize),
		size: size,
		cap:  cap}
}

func (p *ByteFixPool) Get() []byte {
	select {
	// 从channel读
	case b := <-p.cache:
		return b
	default:
		return make([]byte, p.size, p.cap)
	}
}

func (p *ByteFixPool) Put(b []byte) {
	// 重置已用大小
	b = b[:0]
	select {
	// 放入channel
	case p.cache <- b:
	// channel满了则丢弃
	default:
	}
}
