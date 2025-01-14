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
