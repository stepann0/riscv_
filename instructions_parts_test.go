package main

import "testing"

func TestIImm(t *testing.T) {
	words := []struct {
		inst uint32
		imm  uint64
	}{
		{0xeaa00093, 0xfffffffffffffeaa},
		{0xbb748293, 0xfffffffffffffbb7},
		{0x6de5f613, 0x6de},
		{0x25202a03, 0x252},
	}
	for _, w := range words {
		word := InstWord(w.inst)
		got := word.iImm()
		if w.imm != got {
			t.Fatalf("Instruction %#x: got imm=%#x, want %#x", w.inst, got, w.imm)
		}
	}
}

func TestJImm(t *testing.T) {
	words := []struct {
		inst uint32
		imm  uint64
	}{
		{0x7a100fef, 0xfa0},
		{0x7a100a6f, 0xfa0},
		{0x00200a6f, 0x2},
		{0xff9ffa6f, 0xfffffffffffffff8},
	}
	for _, w := range words {
		word := InstWord(w.inst)
		got := word.ujImm()
		if w.imm != got {
			t.Fatalf("Instruction %#x: got imm=%#x, want %#x", w.inst, got, w.imm)
		}
	}
}

func TestUImm(t *testing.T) {
	words := []struct {
		inst uint32
		imm  uint64
	}{
		{0xff830297, 0xffffffffff830000},
		{0x0327d2b7, 0x327d000},
		{0xffe9c2b7, 0xffffffffffe9c000},
		{0xfffff617, 0xfffffffffffff000},
	}
	for _, w := range words {
		word := InstWord(w.inst)
		got := word.uImm()
		if w.imm != got {
			t.Fatalf("Instruction %#x: got imm=%#x, want %#x", w.inst, got, w.imm)
		}
	}
}

func TestSBImm(t *testing.T) {
	words := []struct {
		inst uint32
		imm  uint64
	}{
		{0x8a29163, 0x82},
		{0xf8e5e4e3, 0xffffffffffffff88},
		{0x4ae5d863, 0x4b0},
		{0xfee5fae3, 0xfffffffffffffff4},
	}
	for _, w := range words {
		word := InstWord(w.inst)
		got := word.sbImm()
		if w.imm != got {
			t.Fatalf("Instruction %#x got %#x, want %#x", w.inst, got, w.imm)
		}
	}
}

func TestRdRs(t *testing.T) {
	words := []struct {
		inst uint32
		rd   uint64
		rs1  uint64
		rs2  uint64
	}{
		{0x003110b3, 1, 2, 3},
		{0x0084b533, 10, 9, 8},
		{0x00620133, 2, 4, 6},
		{0x417b0ab3, 21, 22, 23},
		{0x40215133, 2, 2, 2},
	}

	for _, w := range words {
		word := InstWord(w.inst)
		if w.rd != word.rd() || w.rs1 != word.rs1() || w.rs2 != word.rs2() {
			t.Fatalf("Instruction %#x: got rd=%d; rs1=%d; rs2=%d; want rd=%d; rs1=%d; rs2=%d",
				w.inst, word.rd(), word.rs1(), word.rs2(),
				w.rd, w.rs1, w.rs2)
		}
	}
}
