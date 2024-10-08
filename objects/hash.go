// Copyright 2012 Google, Inc. All rights reserved.

// This code is based on
// https://github.com/google/gopacket/blob/master/flows.go
// which was licensed under SPDX-License-Identifier: BSD-3-Clause

package objects

const (
	fnvBasis = 14695981039346656037
	FnvPrime = 1099511628211
)

// FnvHash is used by our FastHash functions, and implements the FNV hash
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See http://isthe.com/chongo/tech/comp/fnv/.
func FnvHash(s []byte) (h uint64) {
	h = fnvBasis
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= FnvPrime
	}
	return
}
