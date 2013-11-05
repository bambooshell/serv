package protoes

import proto "code.google.com/p/goprotobuf/proto"

var opfuncs = make(map[OPCODE]func() proto.Message)

//every opcode has a message prototype
func RegisterProtoFuncs() {
	opfuncs[C2S_CREATE_ACC] = func() proto.Message { return &CreateAcc{} }
	opfuncs[S2C_CREATE_ACC] = func() proto.Message { return &CreateAccRet{} }
	opfuncs[C2S_LOGIN] = func() proto.Message { return &AccLogin{} }
	opfuncs[S2C_LOGIN] = func() proto.Message { return &AccLoginRet{} }
	opfuncs[S2C_ACC_INFO] = func() proto.Message { return &AccLoinInfo{} }
}

func (op OPCODE) CreateProto() proto.Message {
	f := opfuncs[op]
	if f != nil {
		return f()
	}

	return nil
}
