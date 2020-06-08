package oanda

import "strconv"

// Float64 is a helper type for float64
type Float64 float64

func (f Float64) String() string {
	return strconv.FormatFloat(float64(f), 'f', 10, 64)
}

// Int is a helper type for int
type Int int

func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// Bool is a helper type for bool
type Bool bool

func (b Bool) String() string {
	return strconv.FormatBool(bool(b))
}
