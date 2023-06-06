package bloom

import (
	"hash/fnv"
	"unsafe"
)

type BloomFilter struct {
	size     uint64
	hashFunc uint64
	bits     []bool
}

func NewBloomFilter(size uint64, hashFunc uint64) *BloomFilter {
	return &BloomFilter{
		size:     size,
		hashFunc: hashFunc,
		bits:     make([]bool, size),
	}
}

func (b *BloomFilter) Add(data []byte) {
	for i := uint64(0); i < b.hashFunc; i++ {
		h := b.hash(data, i)
		b.bits[h%b.size] = true
	}
}

func (b *BloomFilter) AddConcurrently(data []byte) {
	// 创建一个channel用于接收处理结果
	ch := make(chan bool, b.hashFunc)

	// 启动多个goroutine处理任务
	for i := uint64(0); i < b.hashFunc; i++ {
		go func(i uint64) {
			h := b.hash(data, i)
			b.bits[h%b.size] = true
			ch <- true
		}(i)
	}

	// 等待所有任务完成
	for i := uint64(0); i < b.hashFunc; i++ {
		<-ch
	}
}

func (b *BloomFilter) Test(data []byte) bool {
	for i := uint64(0); i < b.hashFunc; i++ {
		h := b.hash(data, i)
		if !b.bits[h%b.size] {
			return false
		}
	}
	return true
}

func (b *BloomFilter) TestConcurrently(data []byte) bool {
	resultCh := make(chan bool, b.hashFunc)
	for i := uint64(0); i < b.hashFunc; i++ {
		go func(index uint64) {
			h := b.hash(data, index)
			resultCh <- b.bits[h%b.size]
		}(i)
	}

	for i := uint64(0); i < b.hashFunc; i++ {
		if !<-resultCh {
			return false
		}
	}
	return true
}

func (b *BloomFilter) hash(data []byte, seed uint64) uint64 {
	h := fnv.New64a()
	_, err := h.Write(data)
	if err != nil {
		return 0
	}
	return h.Sum64() >> 1 * (seed + 1)
}

func (b *BloomFilter) GetBites() []bool {
	return b.bits
}

func (b *BloomFilter) Size() int {
	// unsafe.Sizeof() 返回的是字节大小
	u := unsafe.Sizeof(b.size) + unsafe.Sizeof(b.hashFunc)
	return len(b.bits)*1 + int(u)
}
