package main

import "math/bits"

type InstWord uint32

func (i InstWord) rd() uint32    { return i.x(7, 5) }
func (i InstWord) rs1() uint32   { return i.x(15, 5) }
func (i InstWord) rs2() uint32   { return i.x(20, 5) }
func (i InstWord) rs3() uint32   { return i.x(27, 5) }
func (i InstWord) rm() uint32    { return i.x(12, 3) }
func (i InstWord) csr() uint32   { return i.x(20, 12) }
func (i InstWord) iImm() uint32  { return i.xs(20, 12) }
func (i InstWord) shamt() uint32 { return i.x(20, 6) }
func (i InstWord) sImm() uint32  { return i.x(7, 5) + (i.xs(25, 7) << 5) }
func (i InstWord) uImm() uint32  { return i.xs(12, 20) << 12 }
func (i InstWord) sbImm() uint32 {
	return (i.x(8, 4) << 1) + (i.x(25, 6) << 5) + (i.x(7, 1) << 11) + (i.immSign() << 12)
}

func (i InstWord) ujImm() uint32 {
	return (i.x(21, 10) << 1) + (i.x(20, 1) << 11) + (i.x(12, 8) << 12) + (i.immSign() << 20)
}

func (i InstWord) x(start, len int) uint32 {
	return (uint32(i) >> start) & ((uint32(1) << len) - 1)
}
func (i InstWord) xs(start, len int) uint32 {
	return uint32(int32(i) << (32 - start - len) >> (32 - len))
}

func (i InstWord) immSign() uint32 { return i.xs(31, 1) }
func signExtend(val int64, bit uint) int64 {
	return val << (64 - bit) >> (64 - bit)
}

func (cpu *Cpu) add(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), rs1+rs2)
}

func (cpu *Cpu) addi(inst InstWord) {
	rs1 := cpu.readReg(inst.rs1())
	cpu.writeReg(inst.rd(), rs1+uint64(inst.iImm()))
}

func (cpu *Cpu) addiw(inst InstWord) {
	rs1 := cpu.readReg(inst.rs1())
	sum := int32(rs1 + uint64(inst.iImm()))
	cpu.writeReg(inst.rd(), uint64(sum))
}

func (cpu *Cpu) addw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(int32(rs1+rs2)))
}

func (cpu *Cpu) and(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), rs1&rs2)
}

func (cpu *Cpu) andi(inst InstWord) {
	rs1, imm := cpu.readReg(inst.rs1()), uint64(inst.iImm())
	cpu.writeReg(inst.rd(), rs1&imm)
}

func (cpu *Cpu) auipc(inst InstWord) {
	cpu.writeReg(inst.rd(), cpu.pc+uint64(inst&0xfffff000))
}

func (cpu *Cpu) beq(inst InstWord) {
	if cpu.readReg(inst.rs1()) == cpu.readReg(inst.rs2()) {
		cpu.pc += uint64(inst.sbImm())
	}
}

func (cpu *Cpu) bge(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	if int64(rs1) >= int64(rs2) {
		cpu.pc += uint64(inst.sbImm())
	}
}

func (cpu *Cpu) bgeu(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	if rs1 >= rs2 {
		cpu.pc += uint64(inst.sbImm())
	}
}

func (cpu *Cpu) blt(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	if int64(rs1) < int64(rs2) {
		cpu.pc += uint64(inst.sbImm())
	}
}

func (cpu *Cpu) bltu(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	if rs1 < rs2 {
		cpu.pc += uint64(inst.sbImm())
	}
}

func (cpu *Cpu) bne(inst InstWord) {
	if cpu.readReg(inst.rs1()) != cpu.readReg(inst.rs2()) {
		cpu.pc += uint64(inst.sbImm())
	}
}

func (cpu *Cpu) csrrc(inst InstWord) {
	csr_data := cpu.readCSR(inst.csr())
	rs_data := cpu.readReg(inst.rs1())
	cpu.writeReg(inst.rd(), csr_data)
	if inst.rs1() != 0 {
		cpu.writeCSR(inst.csr(), csr_data&(^rs_data))
	}
}

func (cpu *Cpu) csrrci(inst InstWord) {
	csr_data := cpu.readCSR(inst.csr())
	cpu.writeReg(inst.rd(), csr_data)
	if rs := uint64(inst.rs1()); rs != 0 {
		cpu.writeCSR(inst.csr(), csr_data&(^rs))
	}
}

func (cpu *Cpu) csrrs(inst InstWord) {
	csr_data := cpu.readCSR(inst.csr())
	rs_data := cpu.readReg(inst.rs1())
	cpu.writeReg(inst.rd(), csr_data)
	if inst.rs1() != 0 {
		cpu.writeCSR(inst.csr(), csr_data|rs_data)
	}
}

func (cpu *Cpu) csrrsi(inst InstWord) {
	csr_data := cpu.readCSR(inst.csr())
	cpu.writeReg(inst.rd(), csr_data)
	if rs := uint64(inst.rs1()); rs != 0 {
		cpu.writeCSR(inst.csr(), csr_data|rs)
	}
}

func (cpu *Cpu) csrrw(inst InstWord) {
	// if inst.rd() == 0 {
	// 	return
	// } ???
	csr_data := cpu.readCSR(inst.csr())
	rs_data := cpu.readReg(inst.rs1())
	cpu.writeReg(inst.rd(), csr_data)
	cpu.writeCSR(inst.csr(), rs_data)
}

func (cpu *Cpu) csrrwi(inst InstWord) {
	// if inst.rd() == 0 {
	// 	return
	// } ???
	csr_data := cpu.readCSR(inst.csr())
	cpu.writeReg(inst.rd(), csr_data)
	cpu.writeCSR(inst.csr(), uint64(inst.rs1()))
}

func (cpu *Cpu) div(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(int64(rs1)/int64(rs2)))
}

func (cpu *Cpu) divu(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), rs1/rs2)
}

func (cpu *Cpu) divuw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(int32(rs1)/int32(rs2)))
}

func (cpu *Cpu) divw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(uint32(rs1)/uint32(rs2)))
}

func (cpu *Cpu) ebreak(inst InstWord) {

}

func (cpu *Cpu) ecall(inst InstWord) {

}

func (cpu *Cpu) fence(inst InstWord) {

}

func (cpu *Cpu) fence_tso(inst InstWord) {

}

func (cpu *Cpu) frcsr(inst InstWord) {

}

func (cpu *Cpu) frflags(inst InstWord) {

}

func (cpu *Cpu) frrm(inst InstWord) {

}

func (cpu *Cpu) fscsr(inst InstWord) {

}

func (cpu *Cpu) fsflags(inst InstWord) {

}

func (cpu *Cpu) fsflagsi(inst InstWord) {

}

func (cpu *Cpu) fsrm(inst InstWord) {

}

func (cpu *Cpu) fsrmi(inst InstWord) {

}

func (cpu *Cpu) jal(inst InstWord) {
	cpu.writeReg(inst.rd(), cpu.pc+4)
	cpu.pc += uint64(inst.ujImm())
}

func (cpu *Cpu) jalr(inst InstWord) {
	cpu.writeReg(inst.rd(), cpu.pc+4)
	imm := int64(inst) >> 20
	x := int64(^1)
	reg := int64(cpu.readReg(inst.rs1()))
	cpu.pc = uint64((reg + imm) & x)
}

func (cpu *Cpu) lb(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	cpu.writeReg(inst.rd(), cpu.memory.Read8(addr))
}

func (cpu *Cpu) lbu(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	cpu.writeReg(inst.rd(), cpu.memory.Read8(addr))
}

func (cpu *Cpu) ld(inst InstWord) { // ???
	// addr := uint64(inst.iImm() + inst.rs1())
	// cpu.writeReg(inst.rd(), cpu.memory.Read8(addr))
}

func (cpu *Cpu) lh(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	data := cpu.memory.Read16(addr)
	cpu.writeReg(inst.rd(), uint64(int64(int16(data))))
}

func (cpu *Cpu) lhu(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	cpu.writeReg(inst.rd(), cpu.memory.Read16(addr))
}

func (cpu *Cpu) lui(inst InstWord) {
	cpu.writeReg(inst.rd(), uint64(inst&0xfffff000))
}

func (cpu *Cpu) lw(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	data := cpu.memory.Read32(addr)
	cpu.writeReg(inst.rd(), uint64(int64(int32(data))))
}

func (cpu *Cpu) lwu(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	cpu.writeReg(inst.rd(), cpu.memory.Read32(addr))
}

func (cpu *Cpu) mul(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	_, low_bits := bits.Mul64(rs1, rs2)
	cpu.writeReg(inst.rd(), low_bits)
}

func (cpu *Cpu) mulh(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), mulh(int64(rs1), int64(rs2)))
}

func (cpu *Cpu) mulhsu(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), mulhsu(int64(rs1), rs2))
}

func (cpu *Cpu) mulhu(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	high_bits, _ := bits.Mul64(rs1, rs2)
	cpu.writeReg(inst.rd(), high_bits)
}

func (cpu *Cpu) mulw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(int32(uint32(rs1)*uint32(rs2))))
}

func (cpu *Cpu) or(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), rs1|rs2)
}

func (cpu *Cpu) ori(inst InstWord) {
	rs1, imm := cpu.readReg(inst.rs1()), uint64(inst.iImm())
	cpu.writeReg(inst.rd(), rs1|imm)
}

func (cpu *Cpu) pause(inst InstWord) {

}

func (cpu *Cpu) rdcycle(inst InstWord) {

}

func (cpu *Cpu) rdcycleh(inst InstWord) {

}

func (cpu *Cpu) rdinstret(inst InstWord) {

}

func (cpu *Cpu) rdinstreth(inst InstWord) {

}

func (cpu *Cpu) rdtime(inst InstWord) {

}

func (cpu *Cpu) rdtimeh(inst InstWord) {

}

func (cpu *Cpu) rem(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(int64(rs1)%int64(rs2)))
}

func (cpu *Cpu) remu(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), rs1%rs2)
}

func (cpu *Cpu) remuw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(int32(uint32(rs1)%uint32(rs2))))
}

func (cpu *Cpu) remw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(int32(rs1)%int32(rs2)))
}

func (cpu *Cpu) sb(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	cpu.memory.Write8(addr, cpu.readReg(inst.rs2()))
}

func (cpu *Cpu) sh(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	cpu.memory.Write16(addr, cpu.readReg(inst.rs2()))
}

func (cpu *Cpu) sw(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	cpu.memory.Write32(addr, cpu.readReg(inst.rs2()))
}

func (cpu *Cpu) sd(inst InstWord) {
	addr := uint64(inst.iImm() + inst.rs1())
	cpu.memory.Write64(addr, cpu.readReg(inst.rs2()))
}

func (cpu *Cpu) sll(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), rs1<<(rs2&0x1f))
}

func (cpu *Cpu) srl(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), rs1>>(rs2&0x1f))
}

func (cpu *Cpu) slliw(inst InstWord) {
	rs1 := cpu.readReg(inst.rs1())
	shift := int32(uint32(rs1) << inst.shamt())
	cpu.writeReg(inst.rd(), uint64(shift)) // ???
}

func (cpu *Cpu) sllw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	shift := int32(uint32(rs1) << (rs2 & 0x1f))
	cpu.writeReg(inst.rd(), uint64(shift))
}

func (cpu *Cpu) srlw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	shift := int32(uint32(rs1) << (rs2 & 0x1f))
	cpu.writeReg(inst.rd(), uint64(shift))
}

func (cpu *Cpu) sraw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(int32(rs1)>>(rs2&0x1f)))
}

func (cpu *Cpu) slt(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	res := uint64(0)
	if int64(rs1) < int64(rs2) {
		res = 1
	}
	cpu.writeReg(inst.rd(), res)
}

func (cpu *Cpu) sltu(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	res := uint64(0)
	if rs1 < rs2 {
		res = 1
	}
	cpu.writeReg(inst.rd(), res)
}

func (cpu *Cpu) slti(inst InstWord) {
	rs1, imm := cpu.readReg(inst.rs1()), uint64(inst.iImm())
	res := uint64(0)
	if int64(rs1) < int64(imm) {
		res = 1
	}
	cpu.writeReg(inst.rd(), res)
}

func (cpu *Cpu) sltiu(inst InstWord) {
	rs1, imm := cpu.readReg(inst.rs1()), uint64(inst.iImm())
	res := uint64(0)
	if rs1 < imm {
		res = 1
	}
	cpu.writeReg(inst.rd(), res)
}

func (cpu *Cpu) sra(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	if cpu.xlen == 32 {
		cpu.writeReg(inst.rd(), uint64(int64(rs1)>>rs2))
	} else {
		cpu.writeReg(inst.rd(), uint64(int64(rs1)>>(rs2&0x1f)))
	}
}

func (cpu *Cpu) srai(inst InstWord) {
	rs1 := cpu.readReg(inst.rs1())
	cpu.writeReg(inst.rd(), uint64(int64(rs1)>>(inst.iImm()&0x1f)))
}

func (cpu *Cpu) sraiw(inst InstWord) {
	rs1 := cpu.readReg(inst.rs1())
	cpu.writeReg(inst.rd(), uint64(int32(rs1)>>(inst.iImm()&0x1f)))
}

func (cpu *Cpu) srliw(inst InstWord) {
	rs1 := cpu.readReg(inst.rs1())
	cpu.writeReg(inst.rd(), uint64(uint32(rs1)>>(inst.iImm()&0x1f)))
}

func (cpu *Cpu) srli(inst InstWord) {
	rs1 := cpu.readReg(inst.rs1())
	cpu.writeReg(inst.rd(), uint64(rs1>>(inst.iImm()&0x1f)))
}

func (cpu *Cpu) slli(inst InstWord) {
	rs1 := cpu.readReg(inst.rs1())
	cpu.writeReg(inst.rd(), uint64(rs1<<(inst.iImm()&0x1f)))
}

func (cpu *Cpu) sub(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), rs1-rs2)
}

func (cpu *Cpu) subw(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), uint64(int32(rs1-rs2)))
}

func (cpu *Cpu) xor(inst InstWord) {
	rs1, rs2 := cpu.readReg(inst.rs1()), cpu.readReg(inst.rs2())
	cpu.writeReg(inst.rd(), rs1^rs2)
}

func (cpu *Cpu) xori(inst InstWord) {
	rs1, imm := cpu.readReg(inst.rs1()), uint64(inst.iImm())
	cpu.writeReg(inst.rd(), rs1^imm)
}
