package main

import (
	"fmt"
)

const maxScalarFinderIterations = 100

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

func (m Matrix) Get(i, j int) Real {
	rows, cols := m.Dims()
	if i < 0 || rows <= i || j < 0 || cols <= j {
		panic(fmt.Sprintf("out of bounds: %d, %d", i, j))
	}
	return m.rows[i].Get(j)
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
	for i := range vs {
		vs[i] = Vector{make([]Real, r)}
		for j, row := range m.rows {
			vs[i].values[j] = row.Get(i)
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
		vs[i] = Vector{make([]Real, m2Cols)}
		for j, cv := range t2.rows {
			vs[i].values[j] = Multiply(&rv, &cv)
		}
	}

	return Matrix{vs}
}

func (m Matrix) Simplify() Matrix {
	return m
}

func (m Matrix) Copy() Matrix {
	m2 := Matrix{make([]Vector, cap(m.rows))}
	for i, v := range m.rows {
		m2.rows[i] = v.Copy()
	}
	return m2
}

func (m Matrix) GaussReduction() Matrix {
	m2 := m.Copy()
	rows, cols := m2.Dims()

	for current := 0; current < cols; current += 1 {
		found := false
		for i := current; i < rows; i += 1 {
			if m2.Get(i, current).IsOne() {
				m2.SwapRow(i, current)
				found = true
			}
		}
		if found {
			for j := current + 1; j < rows; j += 1 {
				v := m2.rows[current].Multiply(m2.Get(j, current))
				m2.rows[j] = m2.rows[j].Minus(v)
			}
		} else {
			for i := current + 1; i < rows; i += 1 {
				x := m2.Get(current, current)
				y := m2.Get(i, current)
				for !x.IsZero() && !y.IsZero() {
					if x.Abs().Compare(y.Abs()) >= 0 {
						v := m2.rows[i].Multiply(x.Divide(y))
						m2.rows[current] = m2.rows[current].Minus(v)
						x = m2.Get(current, current)
					} else {
						v := m2.rows[current].Multiply(y.Divide(x))
						m2.rows[i] = m2.rows[i].Minus(v)
						y = m2.Get(i, current)
					}
				}
				if x.Abs().Compare(y.Abs()) < 0 {
					m2.SwapRow(current, i)
				}
			}
		}
	}

	return m2
}

func (m Matrix) GaussJordan() Matrix {
	m2 := m.GaussReduction()
	rows, cols := m2.Dims()

	for i := rows - 1; i >= 0; i -= 1 {
		j := 0
		for ; j < cols; j += 1 {
			if !m2.Get(i, j).IsZero() {
				break
			}
		}
		x := m2.Get(i, j)

		if !x.IsZero() {
			if !x.IsOne() {
				m2.rows[i] = m2.rows[i].Multiply(NewReal(x.D(), x.N()))
			}

			for k := 0; k < i; k += 1 {
				y := m2.Get(k, j)
				if !(y.IsZero()) {
					v := m2.rows[i].Multiply(y)
					m2.rows[k] = m2.rows[k].Minus(v)
				}
			}
		}
	}

	return m2
}

func main() {
	v1 := IntVector([]int{1, 2, 4})
	v2 := IntVector([]int{9, 3, 8})
	inner := Multiply(&v1, &v2)
	fmt.Println(inner)

	m2 := Matrix{[]Vector{IntVector([]int{2, 3, 1}), IntVector([]int{5, 4, 2}),
		IntVector([]int{4, 2, 7})}}
	fmt.Println(m2.GaussReduction())

	fmt.Println()

	m3 := Matrix{[]Vector{IntVector([]int{2, 3}), IntVector([]int{5, 4}),
		IntVector([]int{4, 2})}}
	fmt.Println(m3.GaussReduction())

	fmt.Println()

	//	m4 := Matrix{[]Vector{{[]int{2, 3, 1, 5}}, {[]int{4, 4, 2, 9}},
	//		{[]int{4, 2, 7, 8}}}}
	m4 := Matrix{[]Vector{IntVector([]int{10, 3, 1, 5}),
		IntVector([]int{6, 4, 2, 9}), IntVector([]int{14, 2, 7, 8})}}
	fmt.Println(m4.GaussJordan())
	fmt.Println()
	fmt.Println(m4)
	fmt.Println()

	m42 := Matrix{[]Vector{IntVector([]int{1, 1, 1, 1}),
		IntVector([]int{1, 2, 2, 2}), IntVector([]int{1, 1, 1, 6})}}
	fmt.Println(m42.GaussReduction())
	fmt.Println()

	m5 := Matrix{[]Vector{IntVector([]int{0, 1, 0}), IntVector([]int{1, 0, 1}),
		IntVector([]int{0, 0, 1})}}
	fmt.Println(m5)

	fmt.Println()
	//	fmt.Println(m2.GaussJordan())
	fmt.Println()
	fmt.Println(m42.GaussJordan())
}
