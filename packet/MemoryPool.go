// =====================================================================
//
// # MemoryPool.go - Memory Pool
//
// Last Modified: 2024/08/01 16:41:12
//
// =====================================================================
package packet

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

// MemoryPool is an array of sync.Pool, each pool is
// for a specific block size.
type MemoryPool struct {
	blockSize []int          // a list of block size
	blockPool []*sync.Pool   // block allocator
	blockUsed []atomic.Int64 // block count
	sizeUp    []int          // size round-up -> index
	sizeDown  []int          // size round-down -> index
	sizeLimit int            // max block size
	usedBlock atomic.Int64
	usedSize  atomic.Int64
	usedLarge atomic.Int64
}

// default block size
var defaultBlockSize = []int{
	8, 16, 32, 64, 128, 256, 512, 1024,
	2048, 4096, 8192, 16384, 32768, 65536}

func roundUp(n int) int {
	return (n + 7) & (^int(7))
}

func roundDown(n int) int {
	return n & (^int(7))
}

func NewMemoryPool(sizes []int) *MemoryPool {
	if sizes == nil {
		sizes = defaultBlockSize
	}
	sizelist := []int{}
	for _, size := range sizes {
		if size > 0 {
			sizelist = append(sizelist, size)
		}
	}
	sort.Ints(sizelist)
	sizes = sizelist
	self := &MemoryPool{
		blockSize: make([]int, len(sizes)),
		blockPool: make([]*sync.Pool, len(sizes)),
		blockUsed: make([]atomic.Int64, len(sizes)),
		sizeUp:    nil,
		sizeDown:  nil,
		sizeLimit: 0,
	}
	for i, size := range sizes {
		if size > self.sizeLimit {
			self.sizeLimit = size
		}
		self.blockSize[i] = size
		self.blockPool[i] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, size)
			},
		}
		self.blockUsed[i].Store(0)
	}
	for size := 0; size <= self.sizeLimit; size += 8 {
		if size <= 0 {
			self.sizeUp = append(self.sizeUp, -1)
			self.sizeDown = append(self.sizeDown, -1)
		} else {
			var found int
			found = -1
			for i := 0; i < len(self.blockSize); i++ {
				if size <= self.blockSize[i] {
					found = i
					break
				}
			}
			if found < 0 {
				panic("MemoryPool: invalid size")
			}
			self.sizeUp = append(self.sizeUp, found)
			found = -1
			for i := len(self.blockSize) - 1; i >= 0; i-- {
				if size >= self.blockSize[i] {
					found = i
					break
				}
			}
			self.sizeDown = append(self.sizeDown, found)
		}
	}
	self.usedBlock.Store(0)
	self.usedSize.Store(0)
	self.usedLarge.Store(0)
	return self
}

func (self *MemoryPool) Release() {
	self.blockSize = nil
	self.blockPool = nil
	self.sizeUp = nil
	self.sizeDown = nil
	self.sizeLimit = 0
}

func (self *MemoryPool) Alloc(size int) []byte {
	if size <= 0 {
		return nil
	}
	if size > self.sizeLimit {
		self.usedLarge.Add(int64(size))
		return make([]byte, size)
	}
	idx := self.sizeUp[roundUp(size)/8]
	if idx < 0 {
		return make([]byte, size)
	}
	obj := self.blockPool[idx].Get().([]byte)
	if cap(obj) < size {
		panic("MemoryPool: invalid block index")
	}
	self.usedSize.Add(int64(cap(obj)))
	self.usedBlock.Add(1)
	self.blockUsed[idx].Add(1)
	return obj[:size]
}

func (self *MemoryPool) Free(obj []byte) {
	if obj == nil {
		return
	}
	size := cap(obj)
	if size <= 0 {
		return
	}
	if size > self.sizeLimit {
		self.usedLarge.Add(int64(-size))
		return
	}
	idx := self.sizeDown[roundDown(size)/8]
	if idx < 0 {
		return
	}
	self.blockPool[idx].Put(obj)
	self.blockUsed[idx].Add(-1)
	self.usedSize.Add(-int64(cap(obj)))
	self.usedBlock.Add(-1)
}

func (self *MemoryPool) String() string {
	used := self.usedSize.Load()
	block := self.usedBlock.Load()
	large := self.usedLarge.Load()
	s := fmt.Sprintf("MemoryPool(allocated=%d, block=%d, large=%d)",
		used, block, large)
	return s
}

func (self *MemoryPool) DumpSizeArray() []int {
	var m []int = nil
	for i := 0; i < len(self.blockSize); i++ {
		m = append(m, self.blockSize[i])
	}
	return m
}

func (self *MemoryPool) DumpUsedArray() []int {
	var m []int = nil
	for i := 0; i < len(self.blockUsed); i++ {
		m = append(m, int(self.blockUsed[i].Load()))
	}
	return m
}

func (self *MemoryPool) DumpInfo() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", self.String())
	for i, size := range self.blockSize {
		count := int(self.blockUsed[i].Load())
		fmt.Fprintf(&buf, "  %d: size=%d count=%d\n", i, size, count)
	}
	return buf.String()
}
