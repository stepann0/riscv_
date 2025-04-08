package main

func main() {
	test_dir := "/home/stepaFedora/Dev/riscv_emulator/tests/test/riscv-tests/"

	cpu := NewCPU()
	LoadSegments(ParseElf(test_dir+"rv64ui-p-sra.elf"), cpu.memory)
}
