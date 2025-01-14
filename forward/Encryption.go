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
	for i := 4; i < len(dst); i++ {
		dst[i] ^= dst[i&3]
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
	if len(key) > 256 {
		key = key[:256]
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return false
	}
	c.XORKeyStream(dst, src[4:])
	for i := 4; i < len(src); i++ {
		dst[i-4] ^= src[i&3]
	}
	return true
}
