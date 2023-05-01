package goloop_test

import (
	"fmt"

	"github.com/dushaoshuai/goloop"
)

func ExampleRepeat() {
	for i := range goloop.Repeat(3) {
		fmt.Println("Repeat", i)
	}
	// Output:
	// Repeat 0
	// Repeat 1
	// Repeat 2
}

func ExampleRepeatWithBreak() {
	for i := range goloop.RepeatWithBreak(3) {
		fmt.Println("Repeat", i.I)
		if i.I == 1 {
			i.Break()
		}
	}
	// Output:
	// Repeat 0
	// Repeat 1
}

func ExampleRange() {
	for i := range goloop.Range[int8](13, -15, 5) {
		fmt.Println(i.I)
		if i.I <= -7 {
			i.Break()
		}
	}

	// Output:
	// 13
	// 8
	// 3
	// -2
	// -7
}

func ExampleRangeSlice() {
	for i, n := range goloop.RangeSlice[uint8](250, 255) {
		fmt.Println(i, n)
		if n >= 253 {
			break
		}
	}

	// Output:
	// 0 250
	// 1 251
	// 2 252
	// 3 253
}
