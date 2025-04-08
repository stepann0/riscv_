package main

import (
	"fmt"
	"testing"
)

type TestCase struct {
	index        int
	name         string
	instructions []uint32
	check        func(*Cpu) error
}

var TEST_CASES = []TestCase{
	{
		index: 0,
		name:  "lui/add/sub",
		instructions: []uint32{
			0x00400237, // lui x4, 1024
			0x000332b7, // lui x5, 51
			0x0001e337, // lui x6, 30
			0x00428033, // add x0, x5, x4
			0x00620533, // add x10, x4, x6
			0x40650533, // sub x10, x10, x6
			0x40530533, // sub x10, x6, x5
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				0:  0x0,
				4:  0x400000,
				5:  0x33000,
				6:  0x1e000,
				10: 0xfffffffffffeb000,
			}
			return cpu.regsMustEq(regs)
		},
	},
	{
		index: 1,
		name:  "xor/xori",
		instructions: []uint32{
			0xd8c00a13, // addi x20, x0, -628
			0x85000a93, // addi x21, x0, -1968
			0x014acb33, // xor x22, x21, x20
			0x015acbb3, // xor x23, x21, x21
			0x000b4c13, // xori x24, x22, 0
			0x7ffacc93, // xori x25, x21, 2047
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				20: 0xfffffffffffffd8c,
				21: 0xfffffffffffff850,
				22: 0x5dc,
				23: 0x0,
				24: 0x5dc,
				25: 0xffffffffffffffaf,
			}
			return cpu.regsMustEq(regs)
		},
	},
	{
		index: 2,
		name:  "or/ori",
		instructions: []uint32{
			0x7ff00a13, // addi x20, x0, 2047
			0xbd000a93, // addi x21, x0, -1072
			0x014aeb33, // or x22, x21, x20
			0x015aebb3, // or x23, x21, x21
			0x000b6c13, // ori x24, x22, 0
			0x144aec93, // ori x25, x21, 324
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				20: 0x7ff,
				21: 0xfffffffffffffbd0,
				22: 0xffffffffffffffff,
				23: 0xfffffffffffffbd0,
				24: 0xffffffffffffffff,
				25: 0xfffffffffffffbd4,
			}
			return cpu.regsMustEq(regs)
		},
	},
	{
		index: 3,
		name:  "slt/sltu",
		instructions: []uint32{
			0x01e00a13, // addi x20, x0, 30
			0x01400a93, // addi x21, x0, 20
			0xff600b13, // addi x22, x0, -10
			0x014aa0b3, // slt x1, x21, x20
			0x016a2133, // slt x2, x20, x22
			0x014ab1b3, // sltu x3, x21, x20
			0x016a3233, // sltu x4, x20, x22
			0x065aa293, // slti x5, x21, 101
			0xff1a2313, // slti x6, x20, -15
			0x00aab393, // sltiu x7, x21, 10
			0xfffa3413, // sltiu x8, x20, -1
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				20: 0x1e,
				21: 0x14,
				22: 0xfffffffffffffff6,
				1:  0x1,
				2:  0x0,
				3:  0x1,
				4:  0x1,
				5:  0x1,
				6:  0x0,
				7:  0x0,
				8:  0x1,
			}
			return cpu.regsMustEq(regs)
		},
	},
	{
		index: 4,
		name:  "lui",
		instructions: []uint32{
			0xfffff0b7, // lui x1, -1
			0xffe00137, // lui x2, -512
			0x002001b7, // lui x3, 512
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				1: 0xfffffffffffff000,
				2: 0xffffffffffe00000,
				3: 0x200000,
			}
			if cpu.pc != 0x8000000c {
				return fmt.Errorf("Program counter must be equal 0x8000000c, but equal %#x", cpu.pc)
			}
			if err := cpu.regsMustEq(regs); err != nil {
				return err
			}
			return nil
		},
	},
	{
		index: 5,
		name:  "auipc",
		instructions: []uint32{
			0x003e8097, // auipc x1, 1000
			0xfff22117, // auipc x2, -222
			0xfffff197, // auipc x3, -1
			0x00001197, // auipc x3, 1
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				1: 0x803e8000,
				2: 0x7ff22004,
				3: 0x8000100c,
			}
			if cpu.pc != 0x80000010 {
				return fmt.Errorf("Program counter must be equal 0x80000010, but equal %#x", cpu.pc)
			}
			if err := cpu.regsMustEq(regs); err != nil {
				return err
			}
			return nil
		},
	},
	{
		index: 6,
		name:  "sra",
		instructions: []uint32{
			0x14d00093, // addi x1, x0, 333
			0x00300113, // addi x2, x0, 3
			0xf9c00193, // addi x3, x0, -100
			0x4020d233, // sra x4, x1, x2
			0x4021d2b3, // sra x5, x3, x2
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				1: 0x14d,
				2: 0x3,
				3: 0xffffffffffffff9c,
				4: 0x29,
				5: 0xfffffffffffffff3,
			}
			return cpu.regsMustEq(regs)
		},
	},
	{
		index: 7,
		name:  "srai",
		instructions: []uint32{
			0x21700093, // addi x1, x0, 535
			0xf8500113, // addi x2, x0, -123
			0x20000193, // addi x3, x0, 512
			0x4050d213, // srai x4, x1, 5
			0x40b15293, // srai x5, x2, 11
			0x4021d313, // srai x6, x3, 2
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				1: 0x217,
				2: 0xffffffffffffff85,
				3: 0x200,
				4: 0x10,
				5: 0xffffffffffffffff,
				6: 0x80,
			}
			return cpu.regsMustEq(regs)
		},
	},
	{
		index: 8,
		name:  "branch",
		instructions: []uint32{
			0xff600093, // addi x1, x0, -10
			0x00009463, // bne x1, x0, 8
			0x00200113, // addi x2, x0, 2
			0x00011263, // bne x2, x0, 4
			0xfff00193, // addi x3, x0, -1
			0xfe018ee3, // beq x3, x0, -4
			0x06400213, // addi x4, x0, 100
			0x0c800293, // addi x5, x0, 200
			0x00524463, // blt x4, x5, 8
			0xfff00313, // addi x6, x0, -1
			0x0042d463, // bge x5, x4, 8
			0x00000013, // addi x0, x0, 0
			0xfff00393, // addi x7, x0, -1
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				1: 0xfffffffffffffff6,
				2: 0x0,
				3: 0xffffffffffffffff,
				4: 0x64,
				5: 0xc8,
				6: 0x0,
				7: 0xffffffffffffffff,
			}
			if cpu.pc != 0x80000034 {
				return fmt.Errorf("Program counter must be equal 0x80000034, but equal %#x", cpu.pc)
			}
			if err := cpu.regsMustEq(regs); err != nil {
				return err
			}
			return nil
		},
	},
	{
		index: 9,
		name:  "jal/jalr",
		instructions: []uint32{
			0x008000ef, // jal x1, 8
			0xfff00113, // addi x2, x0, -1
			0xfff00193, // addi x3, x0, -1
			0x0040006f, // jal x0, 4
			0xfff00213, // addi x4, x0, -1
			0xfff00293, // addi x5, x0, -1
			0x0080006f, // jal x0, 8
			0xffdff06f, // jal x0, -4
			0xfff00313, // addi x6, x0, -1
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				1: 0x80000004,
				2: 0x0,
				3: 0xffffffffffffffff,
				4: 0xffffffffffffffff,
				5: 0xffffffffffffffff,
				6: 0xffffffffffffffff,
			}
			if cpu.pc != 0x80000024 {
				return fmt.Errorf("Program counter must be equal 0x80000024, but equal %#x", cpu.pc)
			}
			if err := cpu.regsMustEq(regs); err != nil {
				return err
			}
			return nil
		},
	},
	{
		index: 10,
		name:  "mul/mulh",
		instructions: []uint32{
			0x08000093, // addi x1, x0, 128
			0x00200113, // addi x2, x0, 2
			0xfff00193, // addi x3, x0, -1
			0x0011d213, // srli x4, x3, 1
			0xb8e00a13, // addi x20, x0, -1138
			0x0141f2b3, // and x5, x3, x20
			0x0012d293, // srli x5, x5, 1
			0x023081b3, // mul x3, x1, x3
			0x022082b3, // mul x5, x1, x2
			0x02419333, // mulh x6, x3, x4
			0x0241b3b3, // mulhu x7, x3, x4
			0x0241a433, // mulhsu x8, x3, x4
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				1:  0x80,
				2:  0x2,
				3:  0xffffffffffffff80,
				4:  0x7fffffffffffffff,
				5:  0x100,
				6:  0xffffffffffffffc0,
				7:  0x7fffffffffffffbf,
				8:  0xffffffffffffffc0,
				9:  0x0,
				10: 0x0,
			}
			if cpu.pc != 0x80000030 {
				return fmt.Errorf("Program counter must be equal 0x80000030, but equal %#x", cpu.pc)
			}
			if err := cpu.regsMustEq(regs); err != nil {
				return err
			}
			return nil
		},
	},
	{
		index: 11,
		name:  "div/rem",
		instructions: []uint32{
			0x10000093, // addi x1, x0, 256
			0x00400113, // addi x2, x0, 4
			0xfff00193, // addi x3, x0, -1
			0x0220c233, // div x4, x1, x2
			0x023242b3, // div x5, x4, x3
			0x0221d333, // divu x6, x3, x2
			0xbe100393, // addi x7, x0, -1055
			0x41f00413, // addi x8, x0, 1055
			0x0283d4b3, // divu x9, x7, x8
			0x02246533, // rem x10, x8, x2
			0x0223e5b3, // rem x11, x7, x2
			0x02247633, // remu x12, x8, x2
			0x0223f6b3, // remu x13, x7, x2
		},
		check: func(cpu *Cpu) error {
			regs := map[uint]uint64{
				1:  0x100,
				2:  0x4,
				3:  0xffffffffffffffff,
				4:  0x40,
				5:  0xffffffffffffffc0,
				6:  0x3fffffffffffffff,
				7:  0xfffffffffffffbe1,
				8:  0x41f,
				9:  0x3e1e930c60174a,
				10: 0x3,
				11: 0xfffffffffffffffd,
				12: 0x3,
				13: 0x1,
			}
			if cpu.pc != 0x80000034 {
				return fmt.Errorf("Program counter must be equal 0x80000034, but equal %#x", cpu.pc)
			}
			if err := cpu.regsMustEq(regs); err != nil {
				return err
			}
			return nil
		},
	},
}

func TestAllInsts(t *testing.T) {
	cpu := NewCPU()
	cpu.reset()
	for _, test := range TEST_CASES {
		cpu.ExecuteProgram(test.instructions)
		if err := test.check(cpu); err != nil {
			t.Fatalf("Test '%s' failed with error: %v", test.name, err)
		}
		t.Logf("Test '%s' OK", test.name)
		cpu.reset()
	}
}

func BenchmarkInst(b *testing.B) {
	bench_inst := []uint32{
		0x01e00a13, 0xff1a2313, 0x00400237, 0x023181b3,
		0x01400a93, 0x00aab393, 0x00519133, 0x023212b3,
		0xff600b13, 0xfffa3413, 0xe6d13213, 0x02a4a433,
		0x014aa0b3, 0x14d00093, 0x01839413, 0x0211a133,
		0x016a2133, 0x00300113, 0x00939133, 0x005231f3,
		0x014ab1b3, 0xf9c00193, 0xfff9d4b7, 0x064294f3,
		0x016a3233, 0x4020d233, 0xffe5b637, 0x08049a73,
		0x065aa293, 0x4021d2b3, 0x0002c1b7,
	}
	cpu := NewCPU()
	for i := 0; i < b.N; i++ {
		cpu.ExecuteInst(bench_inst[i%len(bench_inst)])
	}
}
