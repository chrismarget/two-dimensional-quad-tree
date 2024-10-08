package tdqt

import "fmt"

// Limits define upper and lower bounds in one dimension. Limits work like a
// slice index, so min: 0 and max: 5 covers 5 values: 0, 1, 2, 3, 4. Value 5 is
// out of bounds.
type Limits struct {
	min          int64
	max          int64
	midpointFunc func(int64, int64) int64
}

func (l *Limits) Contains(i int64) bool {
	return l.min <= i && i < l.max
}

func (l *Limits) Max() int64 {
	return l.max
}

func (l *Limits) Min() int64 {
	return l.min
}

func (l *Limits) String() string {
	return fmt.Sprintf("%d-%d", l.min, l.max)
}

func (l *Limits) cannotSubdivide() bool {
	return l.min+1 == l.max
}

func (l *Limits) overlaps(b Limits) bool {
	if l.min > b.max {
		return false // l is too far to the right
	}

	if l.max <= b.min {
		return false // l is too far to the left
	}

	return true
}

func (l *Limits) midpoint() int64 {
	return l.midpointFunc(l.min, l.max)
}

func NewLimits(min, max int64) Limits {
	return Limits{
		min: min,
		max: max,
		midpointFunc: func(a, b int64) int64 {
			return (a | b) - ((a ^ b) >> 1)
		},
	}
}
