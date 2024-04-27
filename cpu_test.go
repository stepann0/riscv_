package main

import (
	"testing"
)

func Test0x13(t *testing.T) {
	code := []uint32{
		0x00500793, //addi x15, x0, 5
		0x00a7c813, //xori x16, x15, 10
		0x00b7e893, //ori x17, x15, 11
		0x0137f913, //andi x18, x15, 19
		0x00579993, //slli x19, x15, 5
		0x0037da13, //srli x20, x15, 3
		0x4067da93, //srai x21, x15, 6
		0x00f7ab13, //slti x22, x15, 15
		0x00b7bb93, //sltiu x23, x15, 11
	}
	cpu := NewCPU(32)
	for _, i := range code {
		rd := (i >> 7) & 0x1f
		cpu.exec(i)
		cpu.dumpRegN(rd)
	}
}
