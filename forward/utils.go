// =====================================================================
//
// utils.go -
//
// Created by skywind on 2024/11/17
// Last Modified: 2024/11/17 01:12:48
//
// =====================================================================
package forward

import (
	"crypto/rc4"
	"net"
	"strings"
)

func AddressSet(dst *net.UDPAddr, src *net.UDPAddr) *net.UDPAddr {
	if len(dst.IP) != len(src.IP) {
		dst.IP = make(net.IP, len(src.IP))
	}
	copy(dst.IP, src.IP)
	dst.Port = src.Port
	dst.Zone = src.Zone
	return dst
}

func AddressClone(addr *net.UDPAddr) *net.UDPAddr {
	// fmt.Printf("AddressClone(%v)\n", addr)
	naddr := &net.UDPAddr{
		IP:   make(net.IP, len(addr.IP)),
		Port: addr.Port,
		Zone: addr.Zone,
	}
	copy(naddr.IP, addr.IP)
	return naddr
}

func AddressParse(addr *net.UDPAddr, ip string, port int) {
	addr.IP = net.ParseIP(ip)
	addr.Port = port
}

func AddressResolve(address string) *net.UDPAddr {
	if !strings.Contains(address, ":") {
		address = ":" + address
	}
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil
	}
	return addr
}

func AddressString(addr *net.UDPAddr) string {
	return addr.String()
}

func EncryptRC4(dst []byte, src []byte, key []byte) bool {
	if len(dst) < len(src) {
		panic("encrypt buffer too small")
	}
	if len(key) == 0 {
		copy(dst, src)
		return true
	}
	if len(key) > 256 {
		key = key[:256]
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return false
	}
	c.XORKeyStream(dst, src)
	return true
}

func HexDump(p []byte, char_visible bool, limit int) string {
	hex := []rune("0123456789ABCDEF")
	if len(p) > limit && limit > 0 {
		p = p[:limit]
	}
	size := len(p)
	count := (size + 15) / 16
	var buffer []rune = make([]rune, 100)
	offset := 0
	output := []string{}
	for i := 0; i < count; i++ {
		length := min(16, len(p))
		line := buffer[:]
		for j := 0; j < len(line); j++ {
			line[j] = ' '
		}
		line[0] = hex[(offset>>12)&15]
		line[1] = hex[(offset>>8)&15]
		line[2] = hex[(offset>>4)&15]
		line[3] = hex[(offset>>0)&15]
		for j := 0; j < length; j++ {
			start := 6 + j*3
			line[start+0] = hex[(p[j]>>4)&15]
			line[start+1] = hex[(p[j]>>0)&15]
			if j == 8 {
				line[start-1] = '-'
			}
			if char_visible {
				c := '.'
				if p[j] >= 32 && p[j] < 127 {
					c = rune(p[j])
				}
				line[6+16*3+2+j] = c
			}
		}
		if len(p) >= 16 {
			p = p[16:]
		}
		s := string(line)
		s = strings.TrimRight(s, " ")
		output = append(output, s)
	}
	return strings.Join(output, "\n")
}
