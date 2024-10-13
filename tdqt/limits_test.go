package tdqt

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLimits_Midpoint_int64(t *testing.T) {
	type testCase struct {
		min      int64
		max      int64
		expected int64
	}

	testCases := []testCase{
		// symmetric
		{min: 0, max: 0, expected: 0},
		{min: -1, max: 1, expected: 0},
		{min: -2, max: 2, expected: 0},

		// symmetric negative
		{min: -1, max: 0, expected: 0},
		{min: -2, max: 1, expected: 0},
		{min: -3, max: 2, expected: 0},
		{min: math.MinInt8, max: math.MaxInt8, expected: 0},

		// same value
		{min: 1, max: 1, expected: 1},
		{min: -1, max: -1, expected: -1},
		{min: 2, max: 2, expected: 2},
		{min: -2, max: -2, expected: -2},
		{min: 20, max: 20, expected: 20},
		{min: -20, max: -20, expected: -20},
		{min: math.MinInt8, max: math.MinInt8, expected: math.MinInt8},
		{min: math.MaxInt8, max: math.MaxInt8, expected: math.MaxInt8},
		{min: math.MinInt64, max: math.MinInt64, expected: math.MinInt64},
		{min: math.MaxInt64, max: math.MaxInt64, expected: math.MaxInt64},

		// one positive
		{min: 0, max: 1, expected: 1},
		{min: 0, max: 2, expected: 1},
		{min: 0, max: 3, expected: 2},
		{min: 0, max: 4, expected: 2},
		{min: 0, max: math.MaxInt8, expected: (math.MaxInt8 / 2) + 1},
		{min: 0, max: math.MaxInt64, expected: (math.MaxInt64 / 2) + 1},

		// one negative
		{min: -1, max: 0, expected: 0},
		{min: -2, max: 0, expected: -1},
		{min: -3, max: 0, expected: -1},
		{min: -4, max: 0, expected: -2},
		{min: math.MinInt8, max: 0, expected: math.MinInt8 / 2},
		{min: math.MinInt64, max: 0, expected: math.MinInt64 / 2},

		// both positive
		{min: 1, max: 2, expected: 2},
		{min: 1, max: 3, expected: 2},
		{min: 2, max: 3, expected: 3},
		{min: 2, max: 4, expected: 3},

		// both negative
		{min: -2, max: -1, expected: -1},
		{min: -3, max: -1, expected: -2},
		{min: -4, max: -1, expected: -2},
		{min: math.MinInt8, max: -1, expected: math.MinInt8 / 2},
		{min: math.MinInt8 + 1, max: -1, expected: math.MinInt8 / 2},
		{min: math.MinInt64, max: -1, expected: math.MinInt64 / 2},
		{min: math.MinInt64 + 1, max: -1, expected: math.MinInt64 / 2},

		// others
		{min: 10, max: 20, expected: 15},
		{min: 10, max: 19, expected: 15},
		{min: -20, max: -10, expected: -15},
		{min: -20, max: -11, expected: -15},
		{min: -5, max: 11, expected: 3},
		{min: -5, max: 10, expected: 3},
		{min: -10, max: 6, expected: -2},
		{min: -10, max: 5, expected: -2},
		{min: -100, max: 50, expected: -25},
		{min: -100, max: -50, expected: -75},
		{min: -50, max: 100, expected: 25},
		{min: 50, max: 100, expected: 75},
	}

	for i, tCase := range testCases {
		t.Run(fmt.Sprintf("test_case_%d", i), func(t *testing.T) {
			nl := NewLimits(tCase.min, tCase.max)
			actual := nl.midpoint()
			require.Equalf(t, tCase.expected, actual, "values: %d and %d. expected %d, got %d", tCase.min, tCase.max, tCase.expected, actual)
			t.Logf("range: %d <-> %d; mid: %d", tCase.min, tCase.max, actual)
		})
	}
}

// func TestLimits_Midpoint_uint64(t *testing.T) {
// 	type testCase struct {
// 		min      uint64
// 		max      uint64
// 		expected uint64
// 	}
//
// 	testCases := []testCase{
// 		// symmetric
// 		{min: 0, max: 0, expected: 0},
//
// 		// same value
// 		{min: 1, max: 1, expected: 1},
// 		{min: 2, max: 2, expected: 2},
// 		{min: 20, max: 20, expected: 20},
// 		{min: math.MaxInt8, max: math.MaxInt8, expected: math.MaxInt8},
// 		{min: math.MaxInt64, max: math.MaxInt64, expected: math.MaxInt64},
//
// 		// one positive
// 		{min: 0, max: 1, expected: 1},
// 		{min: 0, max: 2, expected: 1},
// 		{min: 0, max: 3, expected: 2},
// 		{min: 0, max: 4, expected: 2},
// 		{min: 0, max: math.MaxInt8, expected: (math.MaxInt8 / 2) + 1},
// 		{min: 0, max: math.MaxInt64, expected: (math.MaxInt64 / 2) + 1},
//
// 		// both positive
// 		{min: 1, max: 2, expected: 2},
// 		{min: 1, max: 3, expected: 2},
// 		{min: 2, max: 3, expected: 3},
// 		{min: 2, max: 4, expected: 3},
//
// 		// others
// 		{min: 10, max: 20, expected: 15},
// 		{min: 10, max: 19, expected: 15},
// 		{min: 50, max: 100, expected: 75},
// 		{min: 50, max: 99, expected: 75},
// 	}
//
// 	for i, tCase := range testCases {
// 		t.Run(fmt.Sprintf("test_case_%d", i), func(t *testing.T) {
// 			actual := NewLimits(tCase.min, tCase.max).midpoint()
// 			require.Equalf(t, tCase.expected, actual, "values: %d and %d. expected %d, got %d", tCase.min, tCase.max, tCase.expected, actual)
// 			t.Logf("range: %d <-> %d; mid: %d", tCase.min, tCase.max, actual)
// 		})
// 	}
//
// }
