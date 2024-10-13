package objects

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"log"

	"github.com/chrismarget/two-dimensional-quad-tree/tdqt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/lineintersector"
)

var _ tdqt.Object = (*ColorLine)(nil)

type ColorLine struct {
	x1    int64
	y1    int64
	x2    int64
	y2    int64
	color color.RGBA
	hash  uint64
}

func (cl ColorLine) Hash() uint64 {
	return cl.hash
}

func (cl ColorLine) String() string {
	return fmt.Sprintf("(%d,%d)<->(%d,%d): (%d,%d,%d,%d)", cl.x1, cl.y1, cl.x2, cl.y2, cl.color.R, cl.color.G, cl.color.B, cl.color.A)
}

func (cl *ColorLine) computeHash() {
	bytes := make([]byte, 0, 36)
	bytes = binary.BigEndian.AppendUint64(bytes, uint64(cl.x1))
	bytes = binary.BigEndian.AppendUint64(bytes, uint64(cl.y1))
	bytes = binary.BigEndian.AppendUint64(bytes, uint64(cl.x2))
	bytes = binary.BigEndian.AppendUint64(bytes, uint64(cl.y2))
	bytes = append(bytes, cl.color.R, cl.color.G, cl.color.B, cl.color.A)

	cl.hash = FnvHash(bytes)
}

func (cl ColorLine) Overlaps(r tdqt.Rectangle) (bool, bool) {
	oi1 := octothorpeInfo(r, cl.x1, cl.y1)
	oi2 := octothorpeInfo(r, cl.x2, cl.y2)
	if oi1 == (row2&col2) || oi2 == (row2&col2) {
		// at least one point is in the rectangle
		return true, oi1&oi2 == row2&col2
	}

	bothOiAnded := oi1 & oi2
	if bothOiAnded == (row2 | col2) {
		return true, true // both points are in the rectangle
	}

	switch bothOiAnded & rowBits {
	case row1: // both points in top row
		return false, false // too high
	case row3: // both points in bottom row
		return false, false // too low
	case row2: // both points in center row
		if (oi1|oi2)&colBits == col1|col3 {
			return true, false // line crosses left/right through rectangle
		}
	}

	switch bothOiAnded & colBits {
	case col1:
		return false, false // both points in col1 (too left)
	case col3:
		return false, false // both points in col3 (too right)
	case col2:
		if (oi1|oi2)&rowBits == row1|row3 {
			return true, false // line crosses top/bottom through rectangle
		}
	}

	// If we got here, the line must have:
	//  - One endpoint in the center row and one in the max/min column
	//  - One endpoint in the center column and one in the max/min row
	// No shortcuts available. See if the line intersects the rectangle by
	// checking for intersections with each side of the rectangle.

	xLimits, yLimits := r.Limits()

	lineEPs := [2]geom.Coord{
		{float64(cl.x1), float64(cl.y1)}, {float64(cl.x2), float64(cl.y2)},
	}

	for _, sideEPs := range [4][2]geom.Coord{ // all 4 sides of the rectangle
		{{float64(xLimits.Min()), float64(yLimits.Min())}, {float64(xLimits.Max()), float64(yLimits.Min())}},
		{{float64(xLimits.Max()), float64(yLimits.Min())}, {float64(xLimits.Max()), float64(yLimits.Max())}},
		{{float64(xLimits.Max()), float64(yLimits.Max())}, {float64(xLimits.Min()), float64(yLimits.Max())}},
		{{float64(xLimits.Min()), float64(yLimits.Max())}, {float64(xLimits.Min()), float64(yLimits.Min())}},
	} {
		li := lineintersector.LineIntersectsLine(
			lineintersector.NonRobustLineIntersector{},
			lineEPs[0], lineEPs[1], sideEPs[0], sideEPs[1],
		)
		if li.HasIntersection() {
			// possible false positive here when an endpoint lands on a maximum boundary, but we'll live with it
			return true, false
		}
	}

	return false, false
}

func NewColorLine(x1, y1, x2, y2 int64, color color.RGBA) ColorLine {
	result := ColorLine{
		x1:    x1,
		y1:    y1,
		x2:    x2,
		y2:    y2,
		color: color,
	}

	result.computeHash()

	return result
}

const (
	col1    = 1 << 0
	col2    = 1 << 1
	col3    = 1 << 2
	colBits = col1 | col2 | col3

	row1    = 1 << 3
	row2    = 1 << 4
	row3    = 1 << 5
	rowBits = row1 | row2 | row3
)

// octothorpeInfo takes an (x,y) coordinate pair returns a byte with exactly two
// bits set. The lower bit indicates the coordinate pair's relationship with the
// rectangle's x-axis limits. The higher bit indicates the coordinate pair's
// relationship with the rectangle's y-axis limits.
//
// One of these nine values will be returned: [9, 10, 12, 17, 18, 20, 33, 34, 36]
//
//	  col1                 col2                 col3
//	00000001    x-min    00000010    x-max    00000100
//	              |                    |
//	00001001      |      00001010      |      00001100       00001000
//	    9         |         10         |         12            row1
//	              |                    |
//
// --------------------+--------------------+-------------------- y-max
//
//	              |                    |
//	00010001      |      00010010      |      00010100       00010000
//	   17         |         18         |         20            row2
//	              |                    |
//
// --------------------+--------------------+-------------------- y-min
//
//	              |                    |
//	00100001      |      00100010      |      00100100       00100000
//	   33         |         34         |         36            row3
//	              |                    |
func octothorpeInfo(r tdqt.Rectangle, x, y int64) byte {
	xRange, yRange := r.Limits()
	xMin := xRange.Min()
	xMax := xRange.Max()
	yMin := yRange.Min()
	yMax := yRange.Max()

	var result byte
	switch {
	case x >= xMax:
		result = col3
	case x >= xMin:
		result = col2
	case x < xMin:
		result = col1
	}

	switch {
	case y >= yMax:
		return result | row1
	case y >= yMin:
		return result | row2
	case y < yMin:
		return result | row3
	}

	if result&(col1|col2|col3) == 0 {
		log.Panic("impossible situation in octothorpeInfo x-axis")
	}

	if result&(row1|row2|row3) == 0 {
		log.Panic("impossible situation in octothorpeInfo y-axis")
	}

	log.Panic("impossible situation in octothorpeInfo")
	return 0
}
