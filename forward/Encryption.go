// =====================================================================
//
// Encryption.go -
//
// Created by skywind on 2025/01/14
// Last Modified: 2025/01/14 15:31:00
//
// =====================================================================
package forward

import (
	"crypto/rand"
	"crypto/rc4"
	"encoding/binary"
)

func ReverseBytes(data []byte) {
	size := len(data)
	for i := 0; i < size/2; i++ {
		data[i], data[size-i-1] = data[size-i-1], data[i]
	}
}

func CipherChaining(data []byte, reverse bool) {
	size := len(data)
	if size > 0 {
		if !reverse {
			for i := 1; i < size; i++ {
				data[i] += data[i-1]
			}
		} else {
			previous := data[0]
			for i := 1; i < size; i++ {
				p := data[i]
				data[i] -= previous
				previous = p
			}
		}
	}
}

func PacketEncrypt(dst []byte, src []byte, key []byte) bool {
	if len(dst) != len(src)+4 {
		return false
	}
	if len(key) == 0 {
		copy(dst[4:], src)
		dst[0] = 0
		dst[1] = 0
		dst[2] = 0
		dst[3] = 0
		return true
	}
	_, err := rand.Read(dst[0:4])
	if err != nil {
		dst[0] = 0
		dst[1] = 1
		dst[2] = 2
		dst[3] = 3
	}
	if len(key) > 256 {
		key = key[:256]
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return false
	}
	c.XORKeyStream(dst[4:], src)
	size := len(src)
	for i := 0; i < size; i++ {
		dst[i+4] ^= dst[i&3]
	}
	if true {
		CipherChaining(dst, false)
		ReverseBytes(dst)
	}
	return true
}

func PacketDecrypt(dst []byte, src []byte, key []byte) bool {
	if len(dst) != len(src)-4 {
		return false
	}
	if len(key) == 0 {
		copy(dst, src[4:])
		return true
	}
	if true {
		ReverseBytes(src)
		CipherChaining(src, true)
	}
	if len(key) > 256 {
		key = key[:256]
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return false
	}
	c.XORKeyStream(dst, src[4:])
	size := len(src) - 4
	for i := 0; i < size; i++ {
		dst[i] ^= src[i&3]
	}
	return true
}

func PacketEncode(pkt []byte, key []byte, seq int64) int {
	if cap(pkt) < len(pkt)+8 {
		panic("PacketEncode: buffer too small")
		return -1
	}
	size := len(pkt)
	pkt = pkt[:size+8]
	for i := size - 1; i >= 0; i-- {
		pkt[i+8] = pkt[i]
	}
	binary.LittleEndian.PutUint64(pkt[:8], uint64(seq))
	if len(key) == 0 {
		return size + 8
	}
	for i := 0; i < size; i++ {
		pkt[i+8] ^= pkt[i&7]
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return -2
	}
	c.XORKeyStream(pkt, pkt)
	return size + 8
}

func PacketDecode(pkt []byte, key []byte, seq *int64) int {
	if len(pkt) < 8 {
		return -1
	}
	size := len(pkt) - 8
	if len(key) == 0 {
		*seq = int64(binary.LittleEndian.Uint64(pkt[:8]))
		for i := 0; i < size; i++ {
			pkt[i] = pkt[i+8]
		}
		return size
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return -2
	}
	c.XORKeyStream(pkt, pkt)
	*seq = int64(binary.LittleEndian.Uint64(pkt[:8]))
	for i := 0; i < size; i++ {
		pkt[i+8] ^= pkt[i&7]
	}
	for i := 0; i < size; i++ {
		pkt[i] = pkt[i+8]
	}
	return size
}
