package forward

import "testing"

func TestPacketReduce(t *testing.T) {
	r := NewPacketReduce(10)
	if r.Push(1) != true {
		t.Error("r.Push(1) != true")
	}
	if r.Push(1) != false {
		t.Error("r.Push(1) != false")
	}
	if r.Push(5) != true {
		t.Error("r.Push(5) != true")
	}
	if r.Push(10) != true {
		t.Error("r.Push(10) != true")
	}
	if r.Push(13) != true {
		t.Error("r.Push(10) != true")
	}
	if r.Push(13) != false {
		t.Error("r.Push(10) != false")
	}
	for i := 14; i < 100; i++ {
		// println("test: ", i)
		if r.Push(int64(i)) != true {
			t.Error("r.Push(", i, ") != true")
		}
		if r.Push(int64(i)) != false {
			t.Error("r.Push(", i, ") != false")
		}
	}
}
