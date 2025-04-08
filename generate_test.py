class Inst:
    def __init__(self, inst: str, desc=""):
        assert len(inst) == 32
        self.binary = inst
        self.desc = desc

    def dec(self):
        return int(self.binary, 2)

    def hex(self):
        return hex(self.dec())

    def bin(self):
        return self.binary

    def byte_seq(self):
        d = [self.binary[i * 8 : i * 8 + 8] for i in range(4)]
        return ", ".join(map(lambda x: hex(int(x, 2)), d))


def generate_R_inst(opcode, rd, funct3, rs1, rs2, funct7, desc=""):
    inst = "{0:07b}{1:05b}{2:05b}{3:03b}{4:05b}{5:07b}".format(
        funct7, rs2, rs1, funct3, rd, opcode
    )
    return Inst(inst, desc)


def generate_I_inst(opcode, rd, funct3, rs1, imm, desc=""):
    inst = "{0:012b}{1:05b}{2:03b}{3:05b}{4:07b}".format(imm, rs1, funct3, rd, opcode)
    return Inst(inst, desc)


def generate_S_inst(opcode, imm4_0, funct3, rs1, rs2, imm11_5, desc=""):
    inst = "{0:07b}{1:05b}{2:05b}{3:03b}{4:05b}{5:07b}".format(
        imm11_5, rs2, rs1, funct3, imm4_0, opcode
    )
    return Inst(inst, desc)


inst_list = [
    generate_I_inst(0x13, 20, 0x0, 20, 15, "addi x20, x20, 15"),
    generate_I_inst(0x13, 21, 0x4, 20, 11, "xori x21, x20, 11"),
    generate_I_inst(0x13, 22, 0x6, 20, 9, "ori  x22, x20, 9"),
    generate_I_inst(0x13, 23, 0x7, 22, 13, "addi x23, x22, 13"),
    generate_I_inst(0x13, 20, 0x1, 20, 15, "slli x24, x23, 17"),
]


def test():
    for i in inst_list:
        print(i.hex())


test()
