package objects

import (
	"encoding/binary"
	"fmt"
	"image/color"

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

func (cl *ColorLine) Hash() uint64 {
	return cl.hash
}

func (cl *ColorLine) String() string {
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

func (cl *ColorLine) Overlaps(r tdqt.Rectangle) (bool, bool) {
	xLimits, yLimits := r.Limits()

	Point1Contained := xLimits.Contains(cl.x1) && yLimits.Contains(cl.y1)
	Point2Contained := xLimits.Contains(cl.x2) && yLimits.Contains(cl.y2)

	if Point1Contained || Point2Contained {
		// shortcut return for when one or both endpoints lie inside the rectangle
		return true, Point1Contained && Point2Contained
	}

	lineEPs := [2]geom.Coord{
		{float64(cl.x1), float64(cl.y1)},
		{float64(cl.x2), float64(cl.y2)},
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

func NewColorLine(x1, y1, x2, y2 int64, color color.RGBA) tdqt.Object {
	result := ColorLine{
		x1:    x1,
		y1:    y1,
		x2:    x2,
		y2:    y2,
		color: color,
	}

	result.computeHash()

	return &result
}
