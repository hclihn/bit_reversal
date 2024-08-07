package main

import (
	"bytes"
	"fmt"
	"math/bits"
)

const (
	m0 = 0x5555555555555555 // 01010101 ...: adjecnet bits
	m1 = 0x3333333333333333 // 00110011 ...: adjecent pair of bits
	m2 = 0x0f0f0f0f0f0f0f0f // 00001111 ...: adjecent nibbles
	m3 = 0x00ff00ff00ff00ff // adjecent bytes
	m4 = 0x0000ffff0000ffff // adjecent words
)

func Reverse8(buf []byte, start int) {
	// Inspired by math/bits/Reverse64()
	// We are not using the lookup table approch from Reverse8() just to save some spaces
	const m = 1<<8 - 1
	x := uint8(buf[start])
	x = (x >> 1) & (m0&m) | (x & (m0&m)) << 1 // swap odd and even bits
	x = (x >> 2) & (m1&m) | (x & (m1&m)) << 2 // swap pairs of bits
	x = (x >> 4) | (x << 4) // swap nibbles
	buf[start] = byte(x)
}

func Reverse16(buf []byte, start int) {
	// Inspired by math/bits/Reverse64()
	// We are not using the lookup table approch from Reverse8() just to save some spaces
	const m = 1<<16 - 1
	_ = buf[start+1] // bounds check hint to compiler
	x := uint16(buf[start]) + (uint16(buf[start+1]) << 8)
	x = (x >> 1) & (m0&m) | (x & (m0&m)) << 1 // swap odd and even bits
	x = (x >> 2) & (m1&m) | (x & (m1&m)) << 2 // swap pairs of bits
	x = (x >> 4) & (m2&m) | (x & (m2&m)) << 4 // swap nibbles
	x = (x >> 8) | (x << 8) // swap bytes
	buf[start] = byte(x)
	buf[start+1] = byte(x >> 8)
}

func Reverse32(buf []byte, start int) {
	// Inspired by math/bits/Reverse64()
	// We are not using the lookup table approch from Reverse8() just to save some spaces
	const m = 1<<32 - 1
	_ = buf[start+3] // bounds check hint to compiler
	x := uint32(buf[start])
	for i := 1; i < 4; i++ {
		x = (x << 8) | uint32(buf[start+i])
	}
	x = (x >> 1) & (m0&m) | (x & (m0&m)) << 1 // swap odd and even bits
	x = (x >> 2) & (m1&m) | (x & (m1&m)) << 2 // swap pairs of bits
	x = (x >> 4) & (m2&m) | (x & (m2&m)) << 4 // swap nibbles
	x = (x >> 8) & (m3&m) | (x & (m3&m)) << 8 // swap bytes
	x = (x >> 16) | (x << 16) // swap words
	for i := 3; i >= 0; i-- {
		buf[start+i] = byte(x)
		x >>= 8
	}
}

func Reverse64(buf []byte, start int) {
	// Inspired by math/bits/Reverse64()
	// We are not using the lookup table approch from Reverse8() just to save some spaces
	const m = 1<<64 - 1
	_ = buf[start+7] // bounds check hint to compiler
	x := uint64(buf[start])
	for i := 1; i < 8; i++ {
		x = (x << 8) | uint64(buf[start+i])
	}
	x = (x >> 1) & (m0&m) | (x & (m0&m)) << 1 // swap odd and even bits
	x = (x >> 2) & (m1&m) | (x & (m1&m)) << 2 // swap pairs of bits
	x = (x >> 4) & (m2&m) | (x & (m2&m)) << 4 // swap nibbles
	x = (x >> 8) & (m3&m) | (x & (m3&m)) << 8 // swap bytes
	x = (x >> 16) & (m4&m) | (x & (m4&m)) << 16 // swap words
	x = (x >> 32) | (x << 32) // swap double words
	for i := 7; i >= 0; i-- {
		buf[start+i] = byte(x)
		x >>= 8
	}
}

func ReverseBits(buf []byte) {
	l := len(buf) 
	if l == 0 {
		return
	}
	l2, next, end := l/2, 0, l
	if (end-next) > 15 { // reverse paris of 8-byte groups
		for i := next; i < l2-8; i += 8 {
			j := end - 8
			Reverse64(buf, i)
			Reverse64(buf, j)
			for k := 0; k < 8; k++ {
				buf[i+k], buf[j+k] = buf[j+k], buf[i+k]
			} 
			end -= 8
			next += 8
		}
	}
	if (end-next) > 7 { // reverse a pair of 4-byte groups (length 8-15)
		i, j := next, end - 4
		Reverse32(buf, i)
		Reverse32(buf, j)
		for k := 0; k < 4; k++ {
			buf[i+k], buf[j+k] = buf[j+k], buf[i+k]
		} 
		end -= 4
		next += 4
	}
	if (end-next) > 3 { // reverse a pair of wordss (length 4-7)
		i, j := next, end - 2
		Reverse16(buf, i)
		Reverse16(buf, j)
		for k := 0; k < 2; k++ {
			buf[i+k], buf[j+k] = buf[j+k], buf[i+k]
		} 
		end -= 2
		next += 2
	}
	if (end-next) > 1 { // reverse a pair of bytes (length 2 and 3)
		i, j := next, end - 1
		Reverse8(buf, i)
		Reverse8(buf, j)
		buf[i], buf[j] = buf[j], buf[i]
		end -= 1
		next += 1
	}
	if (end-next) > 0 { // reverse the remaining byte (length 1)
		Reverse8(buf, l2)
	}
}

func main() {
	buf := []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e}
  szs := []int{31, 30, 15, 14, 7, 6, 3, 2, 1, 0}

	for _, sz := range szs {
		fmt.Printf("Reverse %d bytes...\n", sz)
		buf_1 := make([]byte, sz)
		copy(buf_1, buf)
		buf_2 := make([]byte, sz)
		copy(buf_2, buf)
		// create a reference
		for i := 0; i < sz; i++ {
			buf_2[i] = bits.Reverse8(buf[sz-i-1])
		}
		fmt.Printf("Original: %08b\n", buf_1)
		ReverseBits(buf_1)
		fmt.Printf("Reversed: %08b\n", buf_1)
		fmt.Printf("Reversed (%d): %v\n", len(buf_1), bytes.Equal(buf_2, buf_1))
		ReverseBits(buf_1)
		fmt.Printf("Reversed twice (%d): %v\n", len(buf_1), bytes.Equal(buf[:sz], buf_1))
		for i := 0; i < len(buf_1); i++ {
			if buf[i] != buf_1[i] {
				fmt.Printf("-> buf[%d] = %08b vs buf_1[%d] = %08b\n", i, buf[i], i, buf_1[i])
			}
		}
		fmt.Println()
	}
}
