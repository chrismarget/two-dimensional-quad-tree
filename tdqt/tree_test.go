package tdqt_test

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/chrismarget/two-dimensional-quad-tree/objects"
	"github.com/chrismarget/two-dimensional-quad-tree/tdqt"
	"github.com/stretchr/testify/require"
)

func TestTree_Insert_Search(t *testing.T) {
	records := 1000 * 1000

	var insertions uint64
	var maxDepth uint8

	insertCallback := func(tree *tdqt.Tree, depth uint8) {
		insertions++
		maxDepth = max(maxDepth, depth)
	}

	randomColor := func() color.RGBA {
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, rand.Uint32())
		return color.RGBA{R: bytes[0], G: bytes[1], B: bytes[2], A: bytes[3]}
	}

	tree := tdqt.NewTree(0, math.MaxInt64, 0, math.MaxInt64, 400)
	tree.SetInsertCallback(insertCallback)

	start := time.Now()
	for i := range records {
		randColor := randomColor()
		if i%2 == 0 {
			x, y := rand.Int64(), rand.Int64()
			tree.Insert(objects.NewColorPoint(x, y, randColor))
		} else {
			x1, y1 := rand.Int64(), rand.Int64()
			var x2, y2 int64

			if x1 > 0 {
				x2 = x1 - rand.Int64N(1000)
			} else {
				x2 = x1 + rand.Int64N(1000)
			}

			if y1 > 0 {
				y2 = y1 - rand.Int64N(1000)
			} else {
				y2 = y1 + rand.Int64N(1000)
			}

			tree.Insert(objects.NewColorLine(x1, y1, x2, y2, randColor))
		}
	}
	duration := time.Now().Sub(start)
	fmt.Printf("Inserted %d records in %s (%d ips)\n", records, duration, time.Duration(records)*time.Second/duration)
	fmt.Printf("Total (re)insertions: %d\nMax depth: %d\n", insertions, maxDepth)

	half := int64(math.MaxInt64 / 2)
	x1, y1 := rand.Int64N(half), rand.Int64N(half)
	var x2, y2 int64
	if x1 > 0 {
		x2 = x1
		x1 = x1 - rand.Int64N(half)
	} else {
		x2 = x1 + rand.Int64N(half)
	}
	if y1 > 0 {
		y2 = y1
		y1 = y2 - rand.Int64N(half)
	} else {
		y2 = y1 + rand.Int64N(half)
	}

	xLimits := tdqt.NewLimits(x1, x2)
	yLimits := tdqt.NewLimits(y1, y2)

	start = time.Now()
	found := tree.Search(tdqt.NewRectangle(xLimits, yLimits))
	duration = time.Now().Sub(start)
	fmt.Printf("Found %d objects in %s", len(found), duration)
}

func TestTree_Search(t *testing.T) {
	records := int64(1000 * 1000)

	randomColor := func() color.RGBA {
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, rand.Uint32())
		return color.RGBA{R: bytes[0], G: bytes[1], B: bytes[2], A: bytes[3]}
	}

	tree := tdqt.NewTree(0, math.MaxInt64, 0, math.MaxInt64, 400)

	for i := range records {
		randColor := randomColor()
		tree.Insert(objects.NewColorPoint(i, i, randColor))
	}

	xLimits := tdqt.NewLimits(0, math.MaxInt64)
	yLimits := tdqt.NewLimits(0, math.MaxInt64)

	found := tree.Search(tdqt.NewRectangle(xLimits, yLimits))
	fmt.Printf("Found %d objects", len(found))
}

func TestSearchingForLines(t *testing.T) {
	lineCount := 5

	tree := tdqt.NewTree(math.MinInt64, math.MaxInt64, math.MinInt64, math.MaxInt64, 400)

	for i := range lineCount {
		tree.Insert(objects.NewColorLine(int64(i+1), 10, int64(i+1), -10, color.RGBA{}))
	}

	leftHalf := tree.Search(
		tdqt.NewRectangle(
			tdqt.NewLimits(math.MinInt64, -0),
			tdqt.NewLimits(math.MinInt64, math.MaxInt64),
		),
	)
	require.Empty(t, leftHalf)

	topHalf := tree.Search(
		tdqt.NewRectangle(
			tdqt.NewLimits(math.MinInt64, math.MaxInt64),
			tdqt.NewLimits(0, math.MaxInt64),
		),
	)
	require.Len(t, topHalf, lineCount)

	bottomHalf := tree.Search(
		tdqt.NewRectangle(
			tdqt.NewLimits(math.MinInt64, math.MaxInt64),
			tdqt.NewLimits(math.MinInt64, 0),
		),
	)
	require.Len(t, bottomHalf, lineCount)

	rightHalf := tree.Search(
		tdqt.NewRectangle(
			tdqt.NewLimits(0, math.MaxInt64),
			tdqt.NewLimits(math.MinInt64, math.MaxInt64),
		),
	)
	require.Len(t, rightHalf, lineCount)
}
