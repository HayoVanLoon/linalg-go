package main

type Vector struct {
	values []Real
}

func IntVector(xs []int) Vector {
	v := Vector{make([]Real, cap(xs))}
	for i := range v.values {
		v.values[i] = NewReal(xs[i], 1)
	}
	return v
}

func IdentityVector(i, size int) Vector {
	v := Vector{make([]Real, size)}
	v.values[i] = NewReal(1, 1)
	return v
}

func (v Vector) String() string {
	acc := "["
	for i, x := range v.values {
		if i > 0 {
			acc += ","
		}
		acc += x.String()
	}
	return acc + "]"
}

func (v *Vector) Length() int {
	return cap(v.values)
}

func (v *Vector) Get(i int) Real {
	if i < 0 || i >= v.Length() {
		panic("out of bounds")
	}

	return v.values[i]
}

func Multiply(v1 *Vector, v2 *Vector) Real {
	if v1.Length() != v2.Length() {
		panic("different lengths; cannot multiply")
	}

	acc := NewReal(0, 1)
	for i := 0; i < v1.Length(); i += 1 {
		acc = acc.Plus(v1.Get(i).Multiply(v2.Get(i)))
	}

	return acc
}

func (v Vector) Multiply(r Real) Vector {
	v2 := Vector{make([]Real, v.Length())}
	for i, e := range v.values {
		v2.values[i] = e.Multiply(r)
	}
	return v2
}

func (v Vector) MultiplyScalar(scalar int) Vector {
	v2 := Vector{make([]Real, v.Length())}
	for i, e := range v.values {
		v2.values[i] = e.ScalarMultiply(scalar)
	}
	return v2
}

func (v Vector) Minus(v2 Vector) Vector {
	result := Vector{make([]Real, v.Length())}
	if v.Length() != v2.Length() {
		panic("different lengths")
	}
	for i := 0; i < v.Length(); i += 1 {
		result.values[i] = v.Get(i).Minus(v2.Get(i))
	}
	return result
}
