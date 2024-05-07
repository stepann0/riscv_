package main

import (
	"testing"
)

func TestIImm(t *testing.T) {
	words := []struct {
		inst uint32
		imm  uint32
	}{
		{0xbb748293, 0b11111111111111111111101110110111},
		{0x6de5f613, 0b000000000000000000011011011110},
		{0x25202a03, 0b000000000000000000001001010010},
	}
	for _, w := range words {
		word := InstWord(w.inst)
		if w.imm != word.iImm() {
			t.Fatalf("%032b != %032b", w.imm, word.iImm())
		}
	}
}

func TestJImm(t *testing.T) {
	words := []struct {
		inst uint32
		imm  uint32
	}{
		{0x7a100a6f, 0b000000000000000000111110100000},
		{0x00200a6f, 0b000000000000000000000000000010},
		{0xff9ffa6f, 0b11111111111111111111111111111000},
	}
	for _, w := range words {
		word := InstWord(w.inst)
		got := word.ujImm()
		if w.imm != got {
			t.Fatalf("got %032b, want %032b", got, w.imm)
		}
	}
}
func TestUImm(t *testing.T) {
	words := []struct {
		inst uint32
		imm  uint32
	}{
		{0x0327d2b7, 0b000011001001111101000000000000},
		{0xffe9c2b7, 0b11111111111010011100000000000000},
		{0xfffff617, 0b11111111111111111111000000000000},
	}
	for _, w := range words {
		word := InstWord(w.inst)
		got := word.uImm()
		if w.imm != got {
			t.Fatalf("got %032b, want %032b", got, w.imm)
		}
	}
}
