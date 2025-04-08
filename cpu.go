package main

import (
	"fmt"
)

type PrivMode uint64

const (
	BYTE       uint8  = 8
	HALFWORD   uint8  = 16
	WORD       uint8  = 32
	DOUBLEWORD uint8  = 64
	XLEN       uint64 = 64
	FLEN       uint64 = 64

	USER_MODE       PrivMode = 0
	SUPERVISOR_MODE PrivMode = 1
	RESERVED_MODE   PrivMode = 2
	MACHINE_MODE    PrivMode = 3
)

type Cpu struct {
	pc         uint64
	privilege  PrivMode
	xregisters [32]uint64
	fregisters [32]float64 // F/D расширения
	csr        [4096]uint64
	xlen       uint64 // разрядность регистров общего назначения
	flen       uint64 // разрядность float-регистров
	memory     Dram   // доступ к памяти
}

func NewCPU() *Cpu {
	cpu := Cpu{}
	cpu.pc = DRAM_BASE
	cpu.privilege = USER_MODE
	cpu.xlen = XLEN
	cpu.flen = FLEN
	cpu.xregisters[0] = 0 // x0
	cpu.memory = InitDram(MEMORY_SIZE)
	return &cpu
}

func (cpu *Cpu) reset() {
	cpu.pc = DRAM_BASE
	cpu.privilege = USER_MODE
	cpu.xlen = XLEN
	cpu.flen = FLEN
	for i := range cpu.xregisters {
		cpu.xregisters[i] = 0
	}
}

func (cpu *Cpu) regsMustEq(xregs map[uint]uint64) error {
	for x, val := range xregs {
		if cpu.xregisters[x] != val {
			return fmt.Errorf("Регистр x%d должен быть равен %d, но равен %d", x, val, cpu.xregisters[x])
		}
	}
	return nil
}

func (cpu *Cpu) ExecuteProgram(prog []uint32) {
	for pc := cpu.pc - DRAM_BASE; int(pc/4) < len(prog); pc = cpu.pc - DRAM_BASE {
		cpu.ExecuteInst(prog[pc/4])
	}
}

func (cpu *Cpu) ExecuteInst(inst uint32) {
	legal_inst := false
	for _, i := range INSTRUCTIONS {
		if (inst & i.mask) == i.match {
			legal_inst = true
			i.execute(cpu, inst)
			cpu.pc += 4
			break
		}
	}
	if !legal_inst {
		IllegalInst(inst)
	}
}

func (cpu *Cpu) writeReg(reg uint64, val uint64) {
	if reg != 0 {
		cpu.xregisters[reg] = val
	}
}

func (cpu *Cpu) readReg(reg uint64) uint64 {
	if reg == 0 {
		return 0
	}
	return cpu.xregisters[reg]
}

func (cpu *Cpu) readCSR(csr uint64) uint64 {
	return cpu.csr[csr]
}

func (cpu *Cpu) writeCSR(csr uint64, data uint64) {
	cpu.csr[csr] = data
}

func (cpu *Cpu) dumpRegN(regs ...uint64) {
	for _, r := range regs {
		fmt.Printf("[x%d: %d]\n", r, cpu.readReg(r))
	}
}
