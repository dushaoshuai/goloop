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
	for i := range goloop.Range(3, 26, 5) {
		fmt.Println(i.I)
		if i.I == 18 {
			i.Break()
		}
	}

	// Output:
	// 3
	// 8
	// 13
	// 18
}

func ExampleRangeSlice() {
	for i, n := range goloop.RangeSlice[uint8](250, 255) {
		fmt.Println(i, n)
	}

	// Output:
	// 0 250
	// 1 251
	// 2 252
	// 3 253
	// 4 254
}
