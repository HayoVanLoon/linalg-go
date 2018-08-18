package main

import (
	"fmt"
	"strconv"
)

type Real interface {
	N() int
	D() int
	ToInt() (int, int)
	ToFloat() float64
	Plus(Real) Real
	Minus(Real) Real
	Multiply(Real) Real
	Divide(Real) Real
	ScalarMultiply(int) Real
	ScalarDivide(int) Real
	Equals(Real) bool
	IsZero() bool
	IsOne() bool
	Abs() Real
	Compare(Real) int
	Simplify() Real
	fmt.Stringer
}

type realImpl struct {
	n int
	d int
}

func NewReal(n, d int) Real {
	if d == 0 {
		panic("division by zero")
	} else if d < 0 {
		return realImpl{-n, -d}
	} else {
		return realImpl{n, d}
	}
}

func (r realImpl) IsZero() bool {
	return r.n == 0
}

func (r realImpl) IsOne() bool {
	return r.n == r.d
}

func (r realImpl) N() int {
	return r.n
}

func (r realImpl) D() int {
	return r.d
}

func (r realImpl) ToInt() (int, int) {
	return r.n / r.d, r.n % r.d
}

func (r realImpl) ToFloat() float64 {
	return float64(r.n) / float64(r.d)
}

func (r realImpl) cheapSimplify() Real {
	if r.n%r.d == 0 {
		r.n = r.n / r.d
		r.d = 1
	} else if r.d%r.n == 0 {
		r.n = r.n / r.n
		r.d = r.d / r.n
	}
	return r
}

func (r realImpl) Plus(r2 Real) Real {
	return realImpl{r.n*r2.D() + r2.N()*r.d, r.d * r2.D()}.cheapSimplify()
}

func (r realImpl) Minus(r2 Real) Real {
	return realImpl{r.n*r2.D() - r2.N()*r.d, r.d * r2.D()}.cheapSimplify()
}

func (r realImpl) Multiply(r2 Real) Real {
	return realImpl{r.n * r2.N(), r.d * r2.D()}.cheapSimplify()
}

func (r realImpl) Divide(r2 Real) Real {
	if r2.N() == 0 {
		panic("division by zero")
	}
	return realImpl{r.n * r2.D(), r.d * r2.N()}.cheapSimplify()
}

func (r realImpl) ScalarMultiply(i int) Real {
	return r.Multiply(realImpl{i, 1})
}

func (r realImpl) ScalarDivide(i int) Real {
	if i == 0 {
		panic("division by zero")
	}
	return r.Multiply(realImpl{1, i})
}

func (r realImpl) Equals(r2 Real) bool {
	return r.Minus(r2).N() == 0
}

func (r realImpl) Abs() Real {
	return realImpl{abs(r.n), r.d}
}

func (r realImpl) Compare(r2 Real) int {
	return r.Minus(r2).N()
}

func (r realImpl) Simplify() Real {
	f := bcd(&[]int{r.n, r.d})
	return realImpl{r.n / f, r.d / f}
}

func (r realImpl) String() string {
	if r.d == 1 {
		return strconv.Itoa(r.n)
	} else {
		r2 := r.Simplify()
		return fmt.Sprintf("%v/%v", r2.N(), r2.D())
	}
}

func abs(x int) int {
	// from: http://graphics.stanford.edu/~seander/bithacks.html
	mask := x >> 63
	return (x + mask) ^ mask
}

func minAbs(xs *[]int) int {
	min := int(^uint(0) >> 1)
	for _, x := range *xs {
		absX := int(abs(x))
		if absX < min {
			min = absX
		}
	}
	return min
}

func bcd(xs *[]int) int {
	candidate := minAbs(xs)
	for ; candidate > 1; candidate -= 1 {
		ok := true
		for i := 0; ok && i < cap(*xs); i += 1 {
			ok = ok && (*xs)[i]%candidate == 0
		}
		if ok {
			return candidate
		}
	}
	return 1
}
