package protoes

type OPCODE uint16

const (
	C2S_CREATE_ACC = OPCODE(1)
	S2C_CREATE_ACC = OPCODE(2)
	C2S_LOGIN      = OPCODE(3)
	S2C_LOGIN      = OPCODE(4)
	S2C_ACC_INFO   = OPCODE(5)
)
