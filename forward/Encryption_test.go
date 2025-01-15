package forward

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"
)

func TestEncryption(t *testing.T) {
	key := []byte("foobar")
	src := []byte("hello")
	dst := make([]byte, len(src)+4)
	PacketEncrypt(dst, src, key)
	if bytes.Equal(dst, src) {
		t.Error("encrypt failed")
	}
	dec := make([]byte, len(src))
	PacketDecrypt(dec, dst, key)
	if !bytes.Equal(dec, src) {
		t.Error("decrypt failed")
	}
	src = make([]byte, 65536)
	dst = make([]byte, len(src)+4)
	dec = make([]byte, len(src))
	rand.Read(src)
	ts := time.Now().Nanosecond()
	for i := 0; i < 100; i++ {
		PacketEncrypt(dst, src, key)
		PacketDecrypt(dec, dst, key)
		// if !bytes.Equal(dec, src) {
		// 	t.Error("decrypt failed")
		// }
	}
	ts = time.Now().Nanosecond() - ts

	t.Log("encrypt/decrypt 100 times cost", ts, "ns")
}

func TestPacketEncode(t *testing.T) {
	key := []byte("foobar")
	src := []byte("hello")
	dst := make([]byte, len(src)+8)
	copy(dst, src)
	h1 := PacketEncode(dst[:len(src)], key, 1234)
	if h1 != len(src)+8 {
		t.Error("encode failed")
	}
	if bytes.Equal(dst, src) {
		t.Error("encrypt failed")
	}
	var seq int64
	h2 := PacketDecode(dst, key, &seq)
	if h2 != len(src) {
		t.Error("decode failed")
	}
	if !bytes.Equal(dst[:len(src)], src) {
		t.Error("decrypt failed")
	}
	if seq != 1234 {
		t.Error("seq not match")
	}
	src = make([]byte, 65536)
	dst = make([]byte, len(src)+8)
	rand.Read(src)
	copy(dst, src)
	ts := time.Now().Nanosecond()
	for i := 0; i < 100; i++ {
		PacketEncode(dst[:len(src)], key, 1234)
		PacketDecode(dst, key, &seq)
		if seq != 1234 {
			t.Error("seq not match")
		}
		if !bytes.Equal(dst[:len(src)], src) {
			t.Error("decrypt failed")
		}
	}
	ts = time.Now().Nanosecond() - ts

	t.Log("encode/decode 100 times cost", ts, "ns")
}
