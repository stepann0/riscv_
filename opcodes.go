package main

type Instruction struct {
	mask    uint32
	match   uint32
	execute func(*Cpu, uint32)
}

var INSTRUCTIONS = [...]Instruction{
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x33,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.add(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x13,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.addi(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0x707f,
		match: 0x1b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.addiw(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfe00707f,
		match: 0x3b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.addw(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x7033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.and(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x7013,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.andi(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x7f,
		match: 0x17,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.auipc(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x63,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.beq(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x5063,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.bge(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x7063,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.bgeu(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x4063,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.blt(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x6063,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.bltu(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x1063,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.bne(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0x707f,
		match: 0x3073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.csrrc(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0x707f,
		match: 0x7073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.csrrci(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0x707f,
		match: 0x2073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.csrrs(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0x707f,
		match: 0x6073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.csrrsi(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0x707f,
		match: 0x1073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.csrrw(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0x707f,
		match: 0x5073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.csrrwi(InstWord(inst))
		},
	},
	Instruction{
		// RVM extension
		mask:  0xfe00707f,
		match: 0x2004033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.div(InstWord(inst))
		},
	},
	Instruction{
		// RVM extension
		mask:  0xfe00707f,
		match: 0x2005033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.divu(InstWord(inst))
		},
	},
	Instruction{
		// RV64M extension
		mask:  0xfe00707f,
		match: 0x200503b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.divuw(InstWord(inst))
		},
	},
	Instruction{
		// RV64M extension
		mask:  0xfe00707f,
		match: 0x200403b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.divw(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xffffffff,
		match: 0x100073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.ebreak(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xffffffff,
		match: 0x73,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.ecall(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0xf,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.fence(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfff0707f,
		match: 0x8330000f,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.fence_tso(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfffff07f,
		match: 0x302073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.frcsr(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfffff07f,
		match: 0x102073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.frflags(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfffff07f,
		match: 0x202073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.frrm(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfff0707f,
		match: 0x301073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.fscsr(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfff0707f,
		match: 0x101073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.fsflags(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfff0707f,
		match: 0x105073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.fsflagsi(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfff0707f,
		match: 0x201073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.fsrm(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfff0707f,
		match: 0x205073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.fsrmi(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x7f,
		match: 0x6f,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.jal(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x67,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.jalr(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x3,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.lb(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x4003,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.lbu(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0x707f,
		match: 0x3003,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.ld(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x1003,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.lh(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x5003,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.lhu(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x7f,
		match: 0x37,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.lui(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x2003,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.lw(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0x707f,
		match: 0x6003,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.lwu(InstWord(inst))
		},
	},
	Instruction{
		// RVM extension
		mask:  0xfe00707f,
		match: 0x2000033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.mul(InstWord(inst))
		},
	},
	Instruction{
		// RVM extension
		mask:  0xfe00707f,
		match: 0x2001033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.mulh(InstWord(inst))
		},
	},
	Instruction{
		// RVM extension
		mask:  0xfe00707f,
		match: 0x2002033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.mulhsu(InstWord(inst))
		},
	},
	Instruction{
		// RVM extension
		mask:  0xfe00707f,
		match: 0x2003033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.mulhu(InstWord(inst))
		},
	},
	Instruction{
		// RV64M extension
		mask:  0xfe00707f,
		match: 0x200003b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.mulw(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x6033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.or(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x6013,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.ori(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xffffffff,
		match: 0x100000f,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.pause(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfffff07f,
		match: 0xc0002073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.rdcycle(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfffff07f,
		match: 0xc8002073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.rdcycleh(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfffff07f,
		match: 0xc0202073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.rdinstret(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfffff07f,
		match: 0xc8202073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.rdinstreth(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfffff07f,
		match: 0xc0102073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.rdtime(InstWord(inst))
		},
	},
	Instruction{
		// RVZICSR extension
		mask:  0xfffff07f,
		match: 0xc8102073,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.rdtimeh(InstWord(inst))
		},
	},
	Instruction{
		// RVM extension
		mask:  0xfe00707f,
		match: 0x2006033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.rem(InstWord(inst))
		},
	},
	Instruction{
		// RVM extension
		mask:  0xfe00707f,
		match: 0x2007033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.remu(InstWord(inst))
		},
	},
	Instruction{
		// RV64M extension
		mask:  0xfe00707f,
		match: 0x200703b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.remuw(InstWord(inst))
		},
	},
	Instruction{
		// RV64M extension
		mask:  0xfe00707f,
		match: 0x200603b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.remw(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x23,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sb(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0x707f,
		match: 0x3023,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sd(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x1023,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sh(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x1033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sll(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfc00707f,
		match: 0x1013,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.slli(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfe00707f,
		match: 0x101b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.slliw(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfe00707f,
		match: 0x103b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sllw(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x2033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.slt(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x2013,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.slti(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x3013,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sltiu(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x3033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sltu(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x40005033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sra(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfc00707f,
		match: 0x40005013,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.srai(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfe00707f,
		match: 0x4000501b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sraiw(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfe00707f,
		match: 0x4000503b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sraw(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x5033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.srl(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfc00707f,
		match: 0x5013,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.srli(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfe00707f,
		match: 0x501b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.srliw(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfe00707f,
		match: 0x503b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.srlw(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x40000033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sub(InstWord(inst))
		},
	},
	Instruction{
		// RV64I extension
		mask:  0xfe00707f,
		match: 0x4000003b,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.subw(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x2023,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.sw(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0xfe00707f,
		match: 0x4033,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.xor(InstWord(inst))
		},
	},
	Instruction{
		// RVI extension
		mask:  0x707f,
		match: 0x4013,
		execute: func(cpu *Cpu, inst uint32) {
			cpu.xori(InstWord(inst))
		},
	},
}
