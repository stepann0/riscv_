package main

func main() {
	RISCV_CPU := NewCPU()
	RISCV_CPU.ExecuteInst(0x02a00093)
	RISCV_CPU.dumpRegN(0, 1)
}
