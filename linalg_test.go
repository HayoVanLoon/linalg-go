package main

import (
	"testing"
	"reflect"
	"math/rand"
)

func TestMinAbs(t *testing.T) {
	cases := []struct{
		input []int
		expected int
	}{
		{[]int{0}, 0},
		{[]int{2, 6, 9}, 2},
		{[]int{2, 6, -9}, 2},
		{[]int{-2, -6, -9}, 2},
	}
	for _, c := range cases {
		if output := MinAbs(&c.input); output != c.expected {
			t.Errorf("expected %v, got %v", c.expected, output)
		}
	}
}

func TestBCD(t *testing.T) {
	cases := []struct{
		input []int
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
		if output := BCD(&c.input); output != c.expected {
			t.Errorf("expected %v, got %v", c.expected, output)
		}
	}
}

func TestVector_Simplify(t *testing.T) {
	cases := []struct{
		target   Vector
		expected Vector
	}{
		{Vector{[]int{2, 4}}, Vector{[]int{1, 2}}},
		{Vector{[]int{2, 5}}, Vector{[]int{2, 5}}},
		{Vector{[]int{2}}, Vector{[]int{1}}},
	}
	for _, c := range cases {
		c.target.Simplify()
		if !reflect.DeepEqual(c.target, c.expected) {
			t.Errorf("expected %v, got %v", c.expected, c.target)
		}
	}
}

func TestVector_Simplified(t *testing.T) {
	cases := []struct{
		input   Vector
		expected Vector
	}{
		{Vector{[]int{2, 4}}, Vector{[]int{1, 2}}},
		{Vector{[]int{2, 5}}, Vector{[]int{2, 5}}},
		{Vector{[]int{2}}, Vector{[]int{1}}},
	}
	for _, c := range cases {
		original := make([]int, cap(c.input.values))
		copy(original, c.input.values)
		output := c.input.Simplified()
		if !reflect.DeepEqual(output, c.expected) {
			t.Errorf("expected %v, got %v", c.expected, output)
		}
		if !reflect.DeepEqual(c.input.values, original) {
			t.Errorf("input was altered: now %v, was %v", c.input.values,
				original)
		}
	}
}

func TestFindScalars(t *testing.T) {
	for i := 0; i <  100; i += 1 {
		x := int(rand.Float64() * 39) - 20
		y := int(rand.Float64() * 39) - 20
		for ; x == 0 || y == 0 || x%y == 0 || y%x == 0 ||
			multiplesOfRelativePrimes(x, y); {
			x = int(rand.Float64() * 39) - 20
			y = int(rand.Float64() * 39) - 20
		}

		a, b, flip := findScalars(x, y)

		if !flip && a*x - b*y != 1 {
			t.Errorf("%d * %d - %d * %d != 1", a, x, b, y)
		}
		if flip && b*y - a*x != 1 {
			t.Errorf("%d * %d - %d * %d != 1", b, y, a, x)
		}
	}
}