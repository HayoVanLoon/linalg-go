package main

import (
	"testing"
)

func TestMinAbs(t *testing.T) {
	cases := []struct {
		input    []int
		expected int
	}{
		{[]int{0}, 0},
		{[]int{2, 6, 9}, 2},
		{[]int{2, 6, -9}, 2},
		{[]int{-2, -6, -9}, 2},
	}
	for _, c := range cases {
		if output := minAbs(&c.input); output != c.expected {
			t.Errorf("expected %v, got %v", c.expected, output)
		}
	}
}

func TestBCD(t *testing.T) {
	cases := []struct {
		input    []int
		expected int
	}{
		{[]int{2, 4}, 2},
		{[]int{6, 9, 11, 12}, 1},
		{[]int{6, 9, 36, 12}, 3},
		{[]int{6, -9, 36, 12}, 3},
		{[]int{0}, 1},
		{[]int{1}, 1},
	}
	for _, c := range cases {
		if output := bcd(&c.input); output != c.expected {
			t.Errorf("expected %v, got %v", c.expected, output)
		}
	}
}
