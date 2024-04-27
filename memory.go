package main

type RAM interface {
	Write(uint64, uint64, uint8)
	Write8(uint64, uint64)
	Write16(uint64, uint64)
	Write32(uint64, uint64)
	Write64(uint64, uint64)
	Read(uint64, uint8) uint64
	Read8(uint64, uint64) uint64
	Read16(uint64, uint64) uint64
	Read32(uint64, uint64) uint64
	Read64(uint64, uint64) uint64
}

const (
	MEMORY_SIZE uint64 = 10 * 1024 * 1024 // 10Mb
	DRAM_BASE   uint64 = 0x80000000       // starting from 2Gb
)

type Dram []byte

func InitDram(size uint64) Dram {
	return make([]byte, size)
}

func (m *Dram) Read(addr uint64, size uint8) uint64 {
	switch size {
	case 8:
		return m.Read8(addr)
	case 16:
		return m.Read16(addr)
	case 32:
		return m.Read32(addr)
	case 64:
		return m.Read64(addr)
	default:
		panic("Store access fault")
	}
}

func (m Dram) Read8(addr uint64) uint64 {
	addr -= DRAM_BASE
	return uint64(m[addr])
}

func (m Dram) Read16(addr uint64) uint64 {
	addr -= DRAM_BASE
	return uint64(m[addr]) |
		(uint64(m[addr+1]) << 8)
}

func (m Dram) Read32(addr uint64) uint64 {
	addr -= DRAM_BASE
	return uint64(m[addr]) |
		(uint64(m[addr+1]) << 8) |
		(uint64(m[addr+2]) << 16) |
		(uint64(m[addr+3]) << 24)
}

func (m Dram) Read64(addr uint64) uint64 {
	addr -= DRAM_BASE
	return uint64(m[addr]) |
		(uint64(m[addr+1]) << 8) |
		(uint64(m[addr+2]) << 16) |
		(uint64(m[addr+3]) << 24) |
		(uint64(m[addr+4]) << 32) |
		(uint64(m[addr+5]) << 40) |
		(uint64(m[addr+6]) << 48) |
		(uint64(m[addr+7]) << 56)
}

func (m Dram) Write(addr uint64, value uint64, size uint8) {
	switch size {
	case 8:
		m.Write8(addr, value)
	case 16:
		m.Write16(addr, value)
	case 32:
		m.Write32(addr, value)
	case 64:
		m.Write64(addr, value)
	default:
		panic("Store access fault")
	}
}
func (m Dram) Write8(addr uint64, val uint64) {
	addr -= DRAM_BASE
	m[addr] = uint8(val)
}

func (m Dram) Write16(addr uint64, val uint64) {
	addr -= DRAM_BASE
	m[addr] = uint8(val)
	m[addr+1] = uint8(val >> 8)
}

func (m Dram) Write32(addr uint64, val uint64) {
	addr -= DRAM_BASE
	m[addr] = uint8(val)
	m[addr+1] = uint8(val >> 8)
	m[addr+2] = uint8(val >> 16)
	m[addr+3] = uint8(val >> 24)
}

func (m Dram) Write64(addr uint64, val uint64) {
	addr -= DRAM_BASE
	m[addr] = uint8(val)
	m[addr+1] = uint8(val >> 8)
	m[addr+2] = uint8(val >> 16)
	m[addr+3] = uint8(val >> 24)
	m[addr+4] = uint8(val >> 32)
	m[addr+5] = uint8(val >> 40)
	m[addr+6] = uint8(val >> 48)
	m[addr+7] = uint8(val >> 56)
}
