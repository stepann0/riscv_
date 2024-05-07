package main

import (
	"fmt"
)

const (
	BYTE       uint8  = 8
	HALFWORD   uint8  = 16
	WORD       uint8  = 32
	DOUBLEWORD uint8  = 64
	XLEN       uint64 = 64
	FLEN       uint64 = 64
)

type Cpu struct {
	pc         uint64
	xregisters [32]uint64
	fregisters [32]float64 // F/D extension
	csr        [4096]uint64
	xlen       uint64 // only 32 or 64
	flen       uint64 // only 32 or 64
	memory     Dram
}

func NewCPU() *Cpu {
	cpu := Cpu{}
	cpu.xlen = XLEN
	cpu.flen = FLEN
	cpu.xregisters[0] = 0           // x0
	cpu.xregisters[2] = MEMORY_SIZE // x2 stack pointer
	return &cpu
}

func (cpu *Cpu) ExecuteInst(inst uint32) {
	for _, i := range INSTRUCTIONS {
		if (inst & i.mask) == i.match {
			i.op(cpu, inst)
			cpu.pc += 4
		}
	}
	panic("Illegal instruction")
}

func (cpu *Cpu) writeReg(reg uint32, val uint64) {
	if reg != 0 {
		cpu.xregisters[reg] = val
	}
}

func (cpu *Cpu) readReg(reg uint32) uint64 {
	if reg == 0 {
		return 0
	}
	return cpu.xregisters[reg]
}

func (cpu *Cpu) readCSR(csr uint32) uint64 {
	return cpu.csr[csr]
}

func (cpu *Cpu) writeCSR(csr uint32, data uint64) {
	cpu.csr[csr] = data
}

func (cpu *Cpu) dumpRegN(regs ...uint32) {
	for _, r := range regs {
		fmt.Printf("[x%d: %d]\n", r, cpu.readReg(r))
	}
}
