package tdqt

import "fmt"

type Rectangle struct {
	xRange Limits
	yRange Limits
}

func (r Rectangle) String() string {
	return fmt.Sprintf("%dx%d (X: %s; Y: %s)",
		r.xRange.max-r.xRange.min, r.yRange.max-r.yRange.min, // dimensions
		r.xRange.String(), r.yRange.String(), // Limits
	)
}

func (r Rectangle) Limits() (Limits, Limits) {
	return r.xRange, r.yRange
}

func (r Rectangle) cannotSubdivide() bool {
	return r.xRange.cannotSubdivide() && r.yRange.cannotSubdivide()
}

func (r Rectangle) Overlaps(b Rectangle) bool {
	return r.xRange.overlaps(b.xRange) && r.yRange.overlaps(b.yRange)
}

func (r Rectangle) xyMinMax() (int64, int64, int64, int64) {
	return r.xRange.min, r.xRange.max, r.yRange.min, r.yRange.max
}

func NewRectangle(x, y Limits) Rectangle {
	return Rectangle{x, y}
}
