// =====================================================================
//
// utils_test.go -
//
// Created by skywind on 2024/11/17
// Last Modified: 2024/11/17 02:16:13
//
// =====================================================================
package forward

import (
	"fmt"
	"testing"
)

func TestRC4(t *testing.T) {
	key := []byte("hello")
	src := []byte("hello world")
	// key = []byte("")
	dst := make([]byte, len(src))
	if !EncryptRC4(dst, src, key) {
		t.Fatal("rc4 encrypt failed")
	}
	if !EncryptRC4(src, dst, key) {
		t.Fatal("rc4 decrypt failed")
	}
	fmt.Printf("%v\n", HexDump(src, true, 0))
	fmt.Printf("%v\n", HexDump(dst, true, 0))
	if string(src) != "hello world" {
		t.Fatal("rc4 decrypt failed")
	}
}

func TestAddress(t *testing.T) {
	addr := AddressResolve("1.2.3.4:5678")
	if addr == nil {
		t.Fatal("resolve address failed")
	}
	if AddressString(addr) != "1.2.3.4:5678" {
		t.Fatal("address string failed")
	}
	addr = AddressResolve("8086")
}
