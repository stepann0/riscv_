package main

import (
	"fmt"
	"math/bits"
	"syscall"
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
	// execute := decode(inst)
	// execute(cpu, inst)
	cpu.execRV64IM(inst)
	cpu.pc += 4
}

func (cpu *Cpu) execRV64IM(inst uint32) {
	var (
		opcode         uint32
		rs1, rs2, rd   uint32
		funct3, funct7 uint32
	)
	opcode = inst & 0x7f
	rd = (inst >> 7) & 0x1f
	rs1 = (inst >> 15) & 0x1f
	rs2 = (inst >> 20) & 0x1f
	funct3 = (inst >> 12) & 7
	switch opcode {
	case 0x3:
		offset := uint64(int64(int32(inst)) >> 20)
		addr := cpu.readReg(rs1) + offset
		var val int64
		switch funct3 {
		case 0x0: // lb
			val = int64(cpu.memory.Read(addr, BYTE))
		case 0x1: // lh
			val = int64(cpu.memory.Read(addr, HALFWORD))
		case 0x2: // lw
			val = int64(cpu.memory.Read(addr, WORD))
		case 0x3: // ld
			val = int64(cpu.memory.Read(addr, DOUBLEWORD))
		case 0x4: // lbu
			val := cpu.memory.Read(addr, HALFWORD)
			cpu.writeReg(rd, val)
			goto ALREADY_WRITE_REG
		case 0x5: // lhu
			val := cpu.memory.Read(addr, HALFWORD)
			cpu.writeReg(rd, val)
			goto ALREADY_WRITE_REG
		case 0x6: // lwu
			val := cpu.memory.Read(addr, WORD)
			cpu.writeReg(rd, val)
			goto ALREADY_WRITE_REG
		default:
		}
		cpu.writeReg(rd, uint64(val))
	ALREADY_WRITE_REG:
		break
	case 0x13:
		funct7 = inst >> 20
		a := cpu.xregisters[rs1]
		imm := uint64(funct7)
		shamt := imm & (cpu.xlen - 1)
		switch funct3 {
		case 0x0: // addi
			a += imm
		case 0x1: // slli
			a = a << shamt
		case 0x2: // slti
			if int64(a) < int64(imm) {
				a = 1
			} else {
				a = 0
			}
		case 0x3: // sltiu
			if a < imm {
				a = 1
			} else {
				a = 0
			}
		case 0x4: // xori
			a ^= imm
		case 0x5: // srli, srai
			if (imm & 0x400) == 1 {
				a = a >> shamt
			} else {
				a = uint64(int64(a) >> int64(shamt))
			}
		case 0x6: // ori
			a |= imm
		case 0x7: // andi
			a &= imm
		default:
		}
		cpu.writeReg(rd, a)
	case 0x1b:
		funct7 = inst >> 20
		a := cpu.xregisters[rs1]
		imm := uint64(int32(funct7))
		switch funct3 {
		case 0x0: // addiw
			a = uint64(int32(a + imm))
		case 0x1: // slliw
			a = uint64(int32(a << (imm & 31)))
		case 0x5: // srliw, sraiw
			if imm&0x400 == 1 {
				a = uint64(int32(a) >> (imm & 31))
			} else {
				a = uint64(uint32(a) >> (imm & 31))
			}
		default:
		}
	case 0x23:
		funct7 = rd | ((inst >> (25 - 5)) & 0xfe0)
		funct7 = (funct7 << 20) >> 20
		addr := cpu.xregisters[rs1] + uint64(funct7)
		val := cpu.xregisters[rs2]
		switch funct3 {
		case 0x0: // sb
			cpu.memory.Write8(addr, val)
		case 0x1: // sh
			cpu.memory.Write16(addr, val)
		case 0x2: // sW
			cpu.memory.Write32(addr, val)
		case 0x3: // sd
			cpu.memory.Write64(addr, val)
		}
	case 0x33:
		funct7 = inst >> 25
		a := cpu.readReg(rs1)
		b := cpu.readReg(rs2)
		switch funct7 {
		case 0x0:
			switch funct3 {
			case 0x0: // add
				a += b
			case 0x1: // sll
				a = a << b
			case 0x2: // slt
				if int64(rs1) < int64(rs2) {
					a = 1
				} else {
					a = 0
				}
			case 0x3: // sltu
				if rs1 < rs2 {
					a = 1
				} else {
					a = 0
				}
			case 0x4: // xor
				a ^= b
			case 0x5: // srl
				a = a >> b
			case 0x6: // or
				a |= b
			case 0x7: // and
				a &= b
			default:
			}
		case 0x1:
			// RV64M
			hi_bits, low_bits := bits.Mul64(a, b)
			switch funct3 {
			case 0x0: // mul
				a = low_bits
			case 0x1: // mulh
				a = mulh(int64(a), int64(b))
			case 0x2: // mulhsu
				a = mulhsu(int64(a), b)
			case 0x3: // mulhu
				a = hi_bits
			case 0x4: // div
				a = uint64(int64(a) / int64(b))
			case 0x5: // divu
				a /= b
			case 0x6: // rem
				a = uint64(int64(a) % int64(b))
			case 0x7: // remu
				a = a % b
			default:
			}
		case 0x20:
			switch funct3 {
			case 0x0: // sub
				a -= b
			case 0x5: // sra
				a = uint64(int64(a) >> b)
			default:
			}
		}
		cpu.writeReg(rd, a)
	case 0x37: // lui
		cpu.writeReg(rd, uint64(inst&0xfffff000))
	case 0x17: // auipc
		cpu.writeReg(rd, cpu.pc+uint64(inst&0xfffff000))
	case 0x3b:
		funct7 = inst >> 25
		a := cpu.readReg(rs1)
		b := cpu.readReg(rs2)
		switch funct7 {
		case 0x0:
			switch funct3 {
			case 0x0: // addw
				a = uint64(int32(a + b))
			case 0x1: // sllw
				a = uint64(int32(a << (b & 31)))
			case 0x5: // srlw
				a = uint64(int32(a >> (b & 31)))
			default:
			}
		case 0x1:
			switch funct3 {
			case 0x0: // mulw
				a = uint64(int32(uint32(a) * uint32(b)))
			case 0x4: // divw
				a = uint64(int32(a) / int32(b))
			case 0x5: // divuw
				a = uint64(int32(uint32(a) / uint32(b)))
			case 0x6: // remw
				a = uint64(int32(a) % int32(b))
			case 0x7: // remuw
				a = uint64(int32(uint32(a) % uint32(b)))
			default:
			}
			cpu.writeReg(rd, a)
		case 0x20:
			switch funct3 {
			case 0x0: // subw
				a = uint64(int32(a - b))
			case 0x5: //sraw
				a = uint64(int32(a) >> (b & 31))
			default:
			}
		}
	case 0x63:
		imm := ((inst >> (31 - 12)) & (1 << 12)) |
			((inst >> (25 - 5)) & 0x7e0) |
			((inst >> (8 - 1)) & 0x1e) |
			((inst << (11 - 7)) & (1 << 11))
		imm = (imm << 19) >> 19
		switch funct3 {
		case 0x0: // beq
			if rs1 == rs2 {
				cpu.pc += uint64(imm)
			}
		case 0x1: // bne
			if rs1 != rs2 {
				cpu.pc += uint64(imm)
			}
		case 0x4: // blt
			if int64(rs1) < int64(rs2) {
				cpu.pc += uint64(imm)
			}
		case 0x5: // bge
			if int64(rs1) >= int64(rs2) {
				cpu.pc += uint64(imm)
			}
		case 0x6: // bltu
			if rs1 < rs2 {
				cpu.pc += uint64(imm)
			}
		case 0x7: // bgeu
			if rs1 >= rs2 {
				cpu.pc += uint64(imm)
			}
		default:
		}
	case 0x67: // jalr
		cpu.writeReg(rd, cpu.pc+4)
		imm := int64(inst) >> 20
		x := int64(^1) // can't use literal const
		reg := int64(cpu.readReg(rs1))
		cpu.pc = uint64((reg + imm) & x)
	case 0x6f: // jal
		funct7 = ((inst >> (31 - 20)) & (1 << 20)) |
			((inst >> (21 - 1)) & 0x7fe) |
			((inst >> (20 - 11)) & (1 << 11)) |
			(inst & 0xff000)
		funct7 = (funct7 << 11) >> 11
		cpu.writeReg(rd, cpu.pc+4)
		cpu.pc += uint64(int64(funct7))
	case 0x73: // system
		// csr_addr := uint16((inst >> 20) & 0xfff)
		switch funct3 {
		case 0x0: // ecall
			a1 := uintptr(cpu.readReg(11))
			a2 := uintptr(cpu.readReg(12))
			a3 := uintptr(cpu.readReg(13))
			a4 := uintptr(cpu.readReg(14))
			a5 := uintptr(cpu.readReg(15))
			a6 := uintptr(cpu.readReg(16))
			a7 := uintptr(cpu.readReg(17))
			syscall.Syscall6(a7, a1, a2, a3, a4, a5, a6)
		case 0x1: // csrrw

		case 0x2: // csrrs
		case 0x3: // csrrc
		case 0x5: // csrrwi
		case 0x6: // csrrsi
		case 0x7: // csrrci
		default:
		}
	default:
		panic("Illegal instruction")
	}
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

func (cpu *Cpu) dumpRegN(regs ...uint32) {
	for _, r := range regs {
		fmt.Printf("[x%d: %d]\n", r, cpu.readReg(r))
	}
}
