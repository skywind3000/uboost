package packet

import (
	"fmt"
	"testing"
)

func TestMemoryPool1(t *testing.T) {
	sizes := []int{
		8, 16, 32, 64, 128, 256, 512, 1024,
		2048, 4096, 8192, 16384, 32768, 65536,
	}
	// ...
	mp := NewMemoryPool(sizes)
	p1 := mp.Alloc(1025)
	if cap(p1) != 2048 {
		fmt.Printf("cap(p) = %d\n", cap(p1))
		t.Errorf("cap(t) != 2048")
	}
	p2 := mp.Alloc(64)
	if cap(p2) != 64 {
		t.Errorf("cap(t) != 64")
	}
	mp.Alloc(100 * 1024)
	if n := mp.DumpUsedArray(); n != nil {
		if n[3] != 1 {
			t.Errorf("n[3] != 1")
		}
		if n[8] != 1 {
			t.Errorf("n[8] != 1")
		}
	}
	mp.Free(p1)
	if mp.DumpUsedArray()[8] != 0 {
		t.Errorf("mp.DumpUsedArray()[8] != 0")
	}
	mp.Free(p2)
	if n := mp.DumpUsedArray(); n != nil {
		for i, x := range n {
			if x != 0 {
				t.Errorf("n[%d] != 0", i)
			}
		}
	}
	mp.Release()
}
