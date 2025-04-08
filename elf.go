package main

import (
	"debug/elf"
	"encoding/binary"
)

func ParseElf(path string) *elf.File {
	exe, err := elf.Open(path)
	if err != nil {
		panic(err)
	}
	return exe
}

func loadData2Memory(mem Dram, data []byte, addr uint64) {
	// _ = m[addr+uint64(len(data))] // check length
	for i := range data {
		mem.Write8(addr+uint64(i), uint64(data[i]))
	}
}

func LoadSegments(f *elf.File, memory Dram) {
	fheader := f.FileHeader
	if fheader.ByteOrder != binary.LittleEndian ||
		fheader.Type != elf.ET_EXEC ||
		fheader.Machine != elf.EM_RISCV {
		return
	}
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_LOAD {
			data := make([]uint8, prog.Memsz)
			n, _ := prog.Open().Read(data)
			if n != int(prog.Filesz) {
				panic("Failed to read 'memsz' bytes")
			}
			loadData2Memory(memory, data, prog.Vaddr)
		}
	}
}
