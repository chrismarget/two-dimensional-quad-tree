package objects

import (
	"encoding/binary"
	"fmt"
	"image/color"

	"github.com/chrismarget/two-dimensional-quad-tree/tdqt"
)

var _ tdqt.Object = (*ColorPoint)(nil)

type ColorPoint struct {
	x     int64
	y     int64
	color color.RGBA
	hash  uint64
}

func (cp *ColorPoint) Hash() uint64 {
	return cp.hash
}

func (cp *ColorPoint) String() string {
	return fmt.Sprintf("(%d,%d): (%d,%d,%d,%d)", cp.x, cp.y, cp.color.R, cp.color.G, cp.color.B, cp.color.A)
}

func (cp *ColorPoint) computeHash() {
	bytes := make([]byte, 0, 20)
	bytes = binary.BigEndian.AppendUint64(bytes, uint64(cp.x))
	bytes = binary.BigEndian.AppendUint64(bytes, uint64(cp.y))
	bytes = append(bytes, cp.color.R, cp.color.G, cp.color.B, cp.color.A)

	cp.hash = FnvHash(bytes)
}

func (cp *ColorPoint) Overlaps(r tdqt.Rectangle) (bool, bool) {
	x, y := r.Limits()
	overlap := x.Contains(cp.x) && y.Contains(cp.y)
	if overlap {
		return true, true
	}

	return false, false
}

func NewColorPoint(x, y int64, color color.RGBA) ColorPoint {
	result := ColorPoint{
		x:     x,
		y:     y,
		color: color,
	}

	result.computeHash()

	return result
}
