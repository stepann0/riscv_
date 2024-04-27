package main

import (
	"fmt"
)

func _main() {
	cpu := NewCPU()
	cpu.ExecuteInst(0xfa0a13)
	cpu.ExecuteInst(0xba4a93)
	cpu.ExecuteInst(0x9a6b13)
	cpu.ExecuteInst(0xdb7b93)
	cpu.ExecuteInst(0xfa1a13)
	cpu.ExecuteInst(0x5a5c93)
	fmt.Println(cpu.readReg(20))
	fmt.Println(cpu.readReg(21))
	fmt.Println(cpu.readReg(22))
	fmt.Println(cpu.readReg(23))
	fmt.Println(cpu.readReg(24))
	fmt.Println(cpu.readReg(25))
	a := uint64(17005786465874477057)
	b := uint64(16568759573321351169)
	hi := myMulh(a, b)
	fmt.Printf("%064b\n", hi)
}

func main() {
	test_dir := "/home/stepaFedora/Dev/riscv_emulator/tests/test/riscv-tests/"
	mem := InitDram(MEMORY_SIZE)
	LoadSegments(ParseElf(test_dir+"rv64ui-p-sra.elf"), mem)
	for i, b := range mem {
		fmt.Printf("%#x ", b)
		if i >= 512 {
			break
		}
	}
}
