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
)

func ReverseBytes(data []byte) {
	size := len(data)
	for i := 0; i < size/2; i++ {
		data[i], data[size-i-1] = data[size-i-1], data[i]
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
		for i := 1; i < len(dst); i++ {
			dst[i] += dst[i-1]
		}
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
		var previous byte = src[0]
		for i := 1; i < len(src); i++ {
			p := src[i]
			src[i] -= previous
			previous = p
		}
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
