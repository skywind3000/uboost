// =====================================================================
//
// PacketReduce.go -
//
// Created by skywind on 2025/01/09
// Last Modified: 2025/01/09 16:53:45
//
// =====================================================================
package forward

type PacketReduce struct {
	packets map[int64]bool
	size    int64
	seqMax  int64
	seqMin  int64
}

func NewPacketReduce(size int) *PacketReduce {
	self := &PacketReduce{
		packets: make(map[int64]bool),
		size:    int64(max(size, 1)),
		seqMax:  -1,
		seqMin:  -1,
	}
	return self
}

func (self *PacketReduce) Clear() {
	for key := range self.packets {
		delete(self.packets, key)
	}
	self.seqMax = -1
	self.seqMin = -1
}

func (self *PacketReduce) Add(seq int64) {
	if len(self.packets) == 0 {
		self.seqMax = seq
		self.seqMin = seq
	}
	boundary := self.size*2 + 200
	if seq > self.seqMax+boundary || seq < self.seqMin-boundary {
		self.Clear()
		self.seqMax = seq
		self.seqMin = seq
	}
	self.packets[seq] = true
	self.seqMax = max(seq, self.seqMax)
	self.seqMin = min(seq, self.seqMin)
	for self.seqMax-self.seqMin > self.size {
		_, ok := self.packets[self.seqMin]
		if ok {
			delete(self.packets, self.seqMin)
			// println("reduce: ", self.seqMin)
		}
		self.seqMin++
	}
}

func (self *PacketReduce) Exist(seq int64) bool {
	_, ok := self.packets[seq]
	return ok
}

func (self *PacketReduce) Push(seq int64) bool {
	exists := self.Exist(seq)
	if exists == false {
		self.Add(seq)
	}
	return !exists
}
