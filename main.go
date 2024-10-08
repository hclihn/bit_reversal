package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"time"
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

// use uint64 to reverse bytes from 1-8.
func Reverse(buf []byte, start, length uint) {
	if length == 0 {
		return
	} else if length > 8 {
		ReverseAllBits(buf)
		return
	}
	// Inspired by math/bits/Reverse64()
	// We are not using the lookup table approch from Reverse8() just to save some spaces
	_ = buf[start+length-1] // bounds check hint to compiler
	x := uint64(buf[start]) // pull bytes from buf
	for i := uint(1); i < length; i++ { 
		x = (x << 8) | uint64(buf[start+i])
	}
	switch length { // shift left to mkae it symmetric to the center of 64 (for 5-7)/32 (for 3) bits
	case 3, 7:
		x <<= 4
	case 6:
		x <<= 8
	case 5:
		x <<= 12
	}
	x = (x >> 1) & m0 | (x & m0) << 1 // swap odd and even bits
	x = (x >> 2) & m1 | (x & m1) << 2 // swap pairs of bits
	if length == 1 {
		x = (x >> 4) | (x << 4) // swap nibbles of the lowest byte
		goto done
	}
	x = (x >> 4) & m2 | (x & m2) << 4 // swap nibbles
	if length == 2 {
		x = (x >> 8) | (x << 8) // swap nibbles of the lowest 2 bytes
		goto done
	}
	x = (x >> 8) & m3 | (x & m3) << 8 // swap bytes
	if length < 5 { // lengths 3 and 4
		x = (x >> 16) | (x << 16) // swap words of the lowest 2 words
		goto done
	}
	// lengths 5, 6, 7, and 8
	x = (x >> 16) & m4 | (x & m4) << 16 // swap words
	x = (x >> 32) | (x << 32) // swap double words
done:
	switch length { // shift right to restore the bytes
	case 3, 7:
		x >>= 4
	case 6:
		x >>= 8
	case 5:
		x >>= 12
	}
	// push bytes back to buf reversed
	for i := length-1; i > 0; i-- { // since i is uint, we need to separate the case for 0!
		buf[start+i] = byte(x)
		x >>= 8
	}
	buf[start] = byte(x)
}

func ReverseAllBits(buf []byte) {
	l := len(buf) 
	if l == 0 {
		return
	}
	l2, next, end := uint(l/2), uint(0), uint(l)
	if (end-next) > 15 { // reverse paris of 8-byte groups
		for i := next; i <= l2-8; i += 8 {
			j := end - 8
			Reverse(buf, i, 8)
			Reverse(buf, j, 8)
			for k := uint(0); k < 8; k++ { // swap the two groups
				buf[i+k], buf[j+k] = buf[j+k], buf[i+k]
			} 
			end -= 8
			next += 8
		}
	}
	if (end-next) > 8 { // lengths 9-15
		// Break it into two halves, each less than 8 bytes.
		// If end-next is even, we have the equal halves; otherwise, the first half is 1 byte longer.
		fmt.Printf("reverse x2 %d\n", end-next)
		r := (end-next) % 2 // remainder
		i, j, l0 := next, (end+next+r)/2, (end-next)/2
		Reverse(buf, i, l0+r)
		Reverse(buf, j, l0)
		for k := uint(0); k < l0; k++ { // swap the two halves
			buf[i+r+k], buf[j+k] = buf[j+k], buf[i+r+k]
		} 
		if r > 0 { // move the extra byte to the center
			for k := i; k < i+l0; k++ {
				buf[k], buf[k+1] = buf[k+1], buf[k]
			}
		}
		end -= l0
		next += l0+r
	} else if end-next > 0 { // lengths 1-8
		fmt.Printf("reverse x1 %d\n", end-next)
		Reverse(buf, next, end-next)
		next, end = l2, l2
	}
	// we should not get here!
	if end-next > 0 {
		panic("Reverse did not finish!")
	}
}

func testReverse(buf []byte) {
	Reverse(buf, 0, uint(len(buf)))
}

func testFor(buf []byte, szs []int, name string, fn func([]byte)) {
	failed := false
	start := time.Now()
	for _, sz := range szs {
		fmt.Printf("** Reverse %d bytes using %s...\n", sz, name)
		buf_1 := make([]byte, sz)
		copy(buf_1, buf)
		buf_2 := make([]byte, sz)
		copy(buf_2, buf)
		// create a reference
		for i := 0; i < sz; i++ {
			buf_2[i] = bits.Reverse8(buf[sz-i-1])
		}
		fmt.Printf("Original: %08b\n", buf_1)
		fn(buf_1)
		fmt.Printf("Reversed: %08b\n", buf_1)
		fmt.Printf("Reversed (%d): %v\n", len(buf_1), bytes.Equal(buf_2, buf_1))
		for i := 0; i < len(buf_1); i++ {
			if buf_2[i] != buf_1[i] {
				fmt.Printf("-> buf_2[%d] = %08b vs buf_1[%d] = %08b\n", i, buf_2[i], i, buf_1[i])
			}
		}
		if !bytes.Equal(buf_2, buf_1) {
			failed = true
		}
		fn(buf_1)
		fmt.Printf("Reversed twice (%d): %v\n", len(buf_1), bytes.Equal(buf[:sz], buf_1))
		for i := 0; i < len(buf_1); i++ {
			if buf[i] != buf_1[i] {
				fmt.Printf("-> buf[%d] = %08b vs buf_1[%d] = %08b\n", i, buf[i], i, buf_1[i])
			}
		}
		if !bytes.Equal(buf[:sz], buf_1) {
			failed = true
		}
		fmt.Println()
	}
	fmt.Printf("** %s Failed: %v, %s\n\n", name, failed, time.Since(start))
}

func main() {
	buf := []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e}
  szs := []int{31, 30, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

	testFor(buf, szs, "ReverseBits", ReverseBits)
	testFor(buf, szs, "ReverseAllBits", ReverseAllBits)
	testFor(buf, szs, "Reverse", testReverse)
}
