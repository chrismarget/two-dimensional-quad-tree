package objects_test

import (
	"image/color"
	"strconv"
	"testing"

	"github.com/chrismarget/two-dimensional-quad-tree/objects"
	"github.com/chrismarget/two-dimensional-quad-tree/tdqt"
	"github.com/stretchr/testify/require"
)

func TestColorLine_Overlaps(t *testing.T) {
	limitsTenTwenty := tdqt.NewLimits(10, 20)
	testRectangle := tdqt.NewRectangle(limitsTenTwenty, limitsTenTwenty)

	type testCase struct {
		x1, y1, x2, y2 int64
		r              tdqt.Rectangle
		overlap        bool
		fullyContained bool
	}

	// these test cases refer to line endpoints in 9 octothorpe regions:
	//
	//   o1 | o2 | o3
	//   ---+----+----
	//   o4 | o5 | o6
	//   ---+----+----
	//   o7 | o8 | o9
	//
	testCases := map[string]testCase{
		// horizontal test cases
		"o1_o2": {
			x1: 5,
			y1: 25,
			x2: 15,
			y2: 25,
			r:  testRectangle,
		},
		"o1_o3": {
			x1: 5,
			y1: 25,
			x2: 25,
			y2: 25,
			r:  testRectangle,
		},
		"o2_o3": {
			x1: 15,
			y1: 25,
			x2: 25,
			y2: 25,
			r:  testRectangle,
		},
		"o4_o5": {
			x1:      5,
			y1:      15,
			x2:      15,
			y2:      15,
			r:       testRectangle,
			overlap: true,
		},
		"o4_o6": {
			x1:      5,
			y1:      15,
			x2:      25,
			y2:      15,
			r:       testRectangle,
			overlap: true,
		},
		"o5_o6": {
			x1:      15,
			y1:      15,
			x2:      25,
			y2:      15,
			r:       testRectangle,
			overlap: true,
		},
		"o7_o8": {
			x1: 5,
			y1: 5,
			x2: 15,
			y2: 5,
			r:  testRectangle,
		},
		"o7_09": {
			x1: 5,
			y1: 5,
			x2: 25,
			y2: 5,
			r:  testRectangle,
		},
		"o8_o9": {
			x1: 15,
			y1: 5,
			x2: 25,
			y2: 5,
			r:  testRectangle,
		},

		// vertical test cases
		"o1_o4": {
			x1: 5,
			y1: 25,
			x2: 5,
			y2: 15,
			r:  testRectangle,
		},
		"o1_o7": {
			x1: 5,
			y1: 25,
			x2: 5,
			y2: 5,
			r:  testRectangle,
		},
		"o4_o7": {
			x1: 5,
			y1: 15,
			x2: 5,
			y2: 5,
			r:  testRectangle,
		},
		"o2_o5": {
			x1:      15,
			y1:      25,
			x2:      15,
			y2:      15,
			r:       testRectangle,
			overlap: true,
		},
		"o2_o8": {
			x1:      15,
			y1:      25,
			x2:      15,
			y2:      5,
			r:       testRectangle,
			overlap: true,
		},
		"o5_o8": {
			x1:      15,
			y1:      15,
			x2:      15,
			y2:      5,
			r:       testRectangle,
			overlap: true,
		},
		"o3_o6": {
			x1: 25,
			y1: 25,
			x2: 25,
			y2: 15,
			r:  testRectangle,
		},
		"o3_o9": {
			x1: 25,
			y1: 25,
			x2: 25,
			y2: 5,
			r:  testRectangle,
		},
		"o6_09": {
			x1: 25,
			y1: 15,
			x2: 25,
			y2: 5,
			r:  testRectangle,
		},

		// both points in same region test cases
		"o1": {
			x1: 5,
			y1: 25,
			x2: 6,
			y2: 26,
			r:  testRectangle,
		},
		"o2": {
			x1: 15,
			y1: 25,
			x2: 16,
			y2: 26,
			r:  testRectangle,
		},
		"o3": {
			x1: 25,
			y1: 25,
			x2: 26,
			y2: 26,
			r:  testRectangle,
		},
		"o4": {
			x1: 5,
			y1: 15,
			x2: 6,
			y2: 16,
			r:  testRectangle,
		},
		"o5": {
			x1:             15,
			y1:             15,
			x2:             16,
			y2:             16,
			r:              testRectangle,
			overlap:        true,
			fullyContained: true,
		},
		"o6": {
			x1: 25,
			y1: 15,
			x2: 26,
			y2: 16,
			r:  testRectangle,
		},
		"o7": {
			x1: 5,
			y1: 5,
			x2: 6,
			y2: 6,
			r:  testRectangle,
		},
		"o8": {
			x1: 15,
			y1: 5,
			x2: 16,
			y2: 6,
			r:  testRectangle,
		},
		"o9": {
			x1: 25,
			y1: 5,
			x2: 26,
			y2: 6,
			r:  testRectangle,
		},

		// center with corner test cases
		"o5_o1": {
			x1:      15,
			y1:      15,
			x2:      5,
			y2:      25,
			r:       testRectangle,
			overlap: true,
		},
		"o5_o3": {
			x1:      15,
			y1:      15,
			x2:      25,
			y2:      25,
			r:       testRectangle,
			overlap: true,
		},
		"o5_o7": {
			x1:      15,
			y1:      15,
			x2:      5,
			y2:      5,
			r:       testRectangle,
			overlap: true,
		},
		"o5_o9": {
			x1:      15,
			y1:      15,
			x2:      25,
			y2:      5,
			r:       testRectangle,
			overlap: true,
		},

		// diagonal test cases
		"o7_o3": { // positive slope diagonal
			x1:      5,
			y1:      5,
			x2:      25,
			y2:      25,
			r:       testRectangle,
			overlap: true,
		},
		"o1_o9": { // negative slope diagonal
			x1:      5,
			y1:      25,
			x2:      25,
			y2:      5,
			r:       testRectangle,
			overlap: true,
		},

		// diagonal no intersection test cases
		"o4_o2_no_intersection": {
			x1: 4,
			y1: 16,
			x2: 14,
			y2: 26,
			r:  testRectangle,
		},
		"o2_o6_no_intersection": {
			x1: 16,
			y1: 26,
			x2: 26,
			y2: 16,
			r:  testRectangle,
		},
		"o6_o8_no_intersection": {
			x1: 26,
			y1: 14,
			x2: 16,
			y2: 4,
			r:  testRectangle,
		},
		"o8_o4_no_intersection": {
			x1: 14,
			y1: 4,
			x2: 4,
			y2: 14,
			r:  testRectangle,
		},

		// diagonal single point intersection test cases
		"o4_o2_corner": {
			x1:      5,
			y1:      15,
			x2:      15,
			y2:      25,
			r:       testRectangle,
			overlap: true, // top left corner intersection, prefer this not match, but it does
		},
		"o2_o6_corner": {
			x1:      15,
			y1:      25,
			x2:      25,
			y2:      15,
			r:       testRectangle,
			overlap: true, // top right corner intersection, prefer this not match, but it does
		},
		"o6_o8_corner": {
			x1:      25,
			y1:      15,
			x2:      15,
			y2:      5,
			r:       testRectangle,
			overlap: true, // bottom right corner intersection, prefer this not match, but it does
		},
		"o8_o4_corner": {
			x1:      15,
			y1:      5,
			x2:      5,
			y2:      15,
			r:       testRectangle,
			overlap: true, // bottom left corner intersection
		},

		// diagonal triangular intersection test cases
		"o4_o2_triangle": {
			x1:      6,
			y1:      14,
			x2:      16,
			y2:      24,
			r:       testRectangle,
			overlap: true, // top left corner triangular intersection
		},
		"o2_o6_triangle": {
			x1:      14,
			y1:      24,
			x2:      24,
			y2:      14,
			r:       testRectangle,
			overlap: true, // top right corner triangular intersection
		},
		"o6_o8_triangle": {
			x1:      24,
			y1:      16,
			x2:      14,
			y2:      6,
			r:       testRectangle,
			overlap: true, // bottom right corner triangular intersection
		},
		"o8_o4_triangle": {
			x1:      16,
			y1:      6,
			x2:      4,
			y2:      16,
			r:       testRectangle,
			overlap: true, // bottom left corner triangular intersection
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			lines := make([]objects.ColorLine, 2)
			lines[0] = objects.NewColorLine(tCase.x1, tCase.y1, tCase.x2, tCase.y2, color.RGBA{})
			lines[1] = objects.NewColorLine(tCase.x2, tCase.y2, tCase.x1, tCase.y1, color.RGBA{})

			for i, line := range lines {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					overlap, fullyContained := line.Overlaps(tCase.r)
					require.Equalf(t, tCase.overlap, overlap, "line %s should overlap rectangle %s", line.String(), tCase.r.String())
					require.Equalf(t, tCase.fullyContained, fullyContained, "line %s should be fully contained by rectangle %s", line.String(), tCase.r.String())
				})
			}
		})
	}
}
