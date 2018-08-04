package main

import (
	"strconv"
	"nl/hayovanloon/bithacks"
	"fmt"
)

const maxScalarFinderIterations = 100

type Vector struct {
	values []int
}

func IdentityVector(i, size int) Vector {
	v := Vector{make([]int, size)}
	v.values[i] = 1
	return v
}

func (v Vector) String() string {
	acc := "["
	for i, x := range v.values {
		if i > 0 {
			acc += ","
		}
		acc += strconv.Itoa(x)
	}
	return acc + "]"
}

func (v *Vector) Length() int {
	return cap(v.values)
}

func (v *Vector) Get(i int) int {
	if i < 0 || i >= v.Length() {
		panic("out of bounds")
	}

	return v.values[i]
}

func (v *Vector) Simplify() {
	bcd := BCD(&v.values)
	for i, p := range v.values {
		v.values[i] = p / bcd
	}
}

func (v *Vector) Simplified() Vector {
	v2 := Vector{make([]int, cap(v.values))}
	copy(v2.values, v.values)
	v2.Simplify()
	return v2
}

func (v *Vector) CompareTo(v2 Vector) int {
	if v.Length() != v2.Length() {
		panic("different lengths; cannot compare")
	}

	// TODO

	return 0
}

func Multiply(v1 *Vector, v2 *Vector) int {
	if v1.Length() != v2.Length() {
		panic("different lengths; cannot multiply")
	}

	acc := 0
	for i := 0; i < v1.Length(); i += 1 {
		acc += v1.values[i] * v2.values[i]
	}

	return acc
}

func (v Vector) Multiply(scalar int) Vector {
	v2 := Vector{make([]int, v.Length())}
	for i, e := range v.values {
		v2.values[i] = e * scalar
	}
	return v2
}

func (v Vector) Minus(v2 Vector) Vector {
	result := Vector{make([]int, v.Length())}
	if v.Length() != v2.Length() {
		panic("different lengths")
	}
	for i := 0; i < v.Length(); i += 1 {
		result.values[i] = v.values[i] - v2.values[i]
	}
	return result
}

type Matrix struct {
	rows []Vector
}

func (m Matrix) String() string {
	acc := "["
	for i, x := range m.rows {
		if i > 0 {
			acc += "\n "
		}
		acc += x.String()
	}
	return acc + "]"
}

func (m Matrix) Dims() (int, int) {
	if rows := cap(m.rows); rows > 0 {
		return rows, m.rows[0].Length()
	} else {
		return 0, 0
	}
}

func (m Matrix) Get(i, j int) int {
	rows, cols := m.Dims()
	if i < 0 || rows <= i || j < 0 || cols <= j {
		panic(fmt.Sprintf("out of bounds: %d, %d", i, j))
	}
	return m.rows[i].values[j]
}

func (m Matrix) SwapRow(i, k int) {
	rows := cap(m.rows)
	if i < 0 || rows <= i || k < 0 || rows <= k {
		panic(fmt.Sprintf("out of bounds: %d, %d", i, k))
	}
	temp := m.rows[i]
	m.rows[i] = m.rows[k]
	m.rows[k] = temp
}

func (m Matrix) Transpose() Matrix {
	r, c := m.Dims()

	vs := make([]Vector, c)
	for ri, _ := range vs {
		vs[ri] = Vector{make([]int, r)}
		for ci, row := range m.rows {
			vs[ri].values[ci] = row.Get(ri)
		}
	}

	return Matrix{vs}
}

func MatMul(m1 *Matrix, m2 *Matrix) Matrix {
	m1Rows, m1Cols := m1.Dims()
	m2Rows, m2Cols := m2.Dims()
	if m1Rows != m2Cols || m1Cols != m2Rows {
		panic("incorrect shapes")
	}

	t2 := m2.Transpose()

	vs := make([]Vector, m1Rows)
	for i, rv := range m1.rows {
		vs[i] = Vector{make([]int, m2Cols)}
		for j, cv := range t2.rows {
			vs[i].values[j] = Multiply(&rv, &cv)
		}
	}

	return Matrix{vs}
}

func (m Matrix) Simplify() Matrix {
	return m
}

func MinAbs(xs *[]int) int {
	min := int(^uint(0) >> 1)
	for _, x := range *xs {
		absX := int(bithacks.Abs64(int64(x)))
		if absX < min {
			min = absX
		}
	}
	return min
}

func BCD(xs *[]int) int {
	candidate := MinAbs(xs)
	for ; candidate > 1; candidate -= 1 {
		ok := true
		for i := 0; ok && i < cap(*xs); i += 1 {
			ok = ok && (*xs)[i] % candidate == 0
		}
		if ok {
			return candidate
		}
	}
	return 1
}

func multiplesOfRelativePrimes(x, y int) bool {
	bcd := BCD(&[]int{x, y})
	return bcd > 1 && BCD(&[]int{x / bcd, y / bcd}) == 1
}

func findScalars(x, y int) (int, int, bool) {
	if x == 0 || y == 0 {
		panic("arguments cannot be zero")
	}
	if y%x == 0 || x%y == 0 {
		panic("arguments cannot be multiples mof each other")
	}
	if multiplesOfRelativePrimes(x, y) {
		panic("multiples of relative primes")
	}

	sx := bithacks.Sign(int64(x))
	sy := bithacks.Sign(int64(y))

	a := sx
	b := sy

	// TODO: discover relation between input and required number of iterations
	for i := 0; i < maxScalarFinderIterations; i += 1 {
		if a * x - b *y == 1 {
			return a, b, false
		}
		if  b *y - a * x == 1 {
			return a, b, true
		}
		if a *x > b *y {
			b += sy
		} else {
			b = sy
			a += sx
		}
	}

	panic(fmt.Sprintf("maximum iterations exceeded for %d and %d", x, y))
}

func (m Matrix) GaussReduction() Matrix {
	dimRow, dimCol := m.Dims()

	for done := 0; done < dimCol; done += 1 {
		found := false
		for i := done; !found && i < dimRow - 1; i += 1 {
			x := m.Get(i, done)
			if x == 0 {
				continue
			} else if bithacks.Abs64(int64(x)) == 1 {
				found = true
				m.SwapRow(done, i)
			} else {
				for k := i + 1; !found && k < dimRow; k += 1 {
					y := m.Get(k, done)
					if y == 0 {
						break
					}

					if !multiplesOfRelativePrimes(x, y) {
						found = true
						a, b, flip := findScalars(x, y)
						top, other := i, k
						if flip {
							top, other = k, i
							a, b = b, a
						}

						v1 := m.rows[top].Multiply(a)
						v2 := m.rows[other].Multiply(b)

						m.rows[top] = v1.Minus(v2)
						m.SwapRow(top, done)
					}
				}
			}
		}

		for i := done + 1; found && i < dimRow; i += 1 {
			m.rows[i] = m.rows[i].Minus(m.rows[done].Multiply(m.Get(i, done)))
		}
	}

	return m
}

func (m Matrix) GaussJordan() Matrix {
	m2 := m.GaussReduction()
	rows, cols := m2.Dims()

	for j := 1; j < rows && j < cols; j += 1 {
		if j == rows - 1 && j == cols - 1 {
			m2.rows[j] = IdentityVector(j, cols)
		}
		v := m2.rows[j]
		for i := 0; i < rows && i < cols; i += 1 {
			if j == i {
				continue
			}

			v2 := v.Multiply(m2.Get(i, j))
			m2.rows[i] = m2.rows[i].Minus(v2)
		}
	}
	return m2
}

func main() {
	fmt.Println(findScalars(3, 16))

	v1 := Vector{[]int{1, 2, 4}}
	v2 := Vector{[]int{9, 3, 8}}
	inner := Multiply(&v1, &v2)
	fmt.Println(inner)

	m2 := Matrix{[]Vector{{[]int{2, 3, 1}}, {[]int{5, 4, 2}}, {[]int{4, 2, 7}}}}
	fmt.Println(m2.GaussReduction())

	fmt.Println()

	m3 := Matrix{[]Vector{{[]int{2, 3}}, {[]int{5, 4}}, {[]int{4, 2}}}}
	fmt.Println(m3.GaussReduction())

	fmt.Println()

	m4 := Matrix{[]Vector{{[]int{2, 3, 1, 5}}, {[]int{5, 4, 2, 9}},
		{[]int{4, 2, 7, 8}}}}
	fmt.Println(m4.GaussReduction())

	fmt.Println()

	m5 := Matrix{[]Vector{{[]int{0, 1, 0}}, {[]int{1, 0, 1}}, {[]int{0, 0, 1}}}}
	fmt.Println(m5.GaussReduction())

	fmt.Println()
	fmt.Println(m2.GaussJordan())
	fmt.Println()
	fmt.Println(m4.GaussJordan())
}