package main

import (
	"math/bits"
)

func signExt(val int32, n int) int32 {
	return (val << (32 - n)) >> (32 - n)
}

func mulh(s1, s2 int64) uint64 {
	p, _ := bits.Mul64(uint64(s1), uint64(s2))
	t1 := (s1 >> 63) & s2
	t2 := (s2 >> 63) & s1
	return uint64(int64(p) - t1 - t2)
}

func mulhsu(s int64, u uint64) uint64 {
	p, _ := bits.Mul64(uint64(s), u)
	t1 := (s >> 63) & int64(u)
	return uint64(int64(p) - t1)
}

func myMulh(u1, u2 uint64) uint64 {
	n1, n2 := int64(u1), int64(u2)
	var neg1, neg2 bool
	if n1 < 0 {
		neg1, n1 = true, -n1
	}
	if n2 < 0 {
		neg2, n2 = true, -n2
	}
	v, _ := bits.Mul64(uint64(n1), uint64(n2))

	if neg1 != neg2 {
		v = -v
	}
	return v
}
