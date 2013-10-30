package msg

import (
	"account"
	//"fmt"
	"protoes"
)
import proto "code.google.com/p/goprotobuf/proto"

var msgHandler = make(map[protoes.OPCODE]func(acc *account.Acc, _msg *Msg) int)

//message handling function
func InitMsgHandler() {
	msgHandler[protoes.C2S_CREATE_ACC] = HandleCreateAcc
	msgHandler[protoes.C2S_LOGIN] = HandleLogin
}

func HandleCreateAcc(acc *account.Acc, _msg *Msg) int {
	_recvPb := _msg.Pb.(*protoes.CreateAcc)

	pb := protoes.S2C_CREATE_ACC.CreateProto()

	newPb := pb.(*protoes.CreateAccRet)
	newPb.Ok = proto.Int32(1)
	newPb.AccId = proto.Uint32(1000000)
	newPb.RoleName = proto.String(*(_recvPb.RoleName))

	retMsg := NewMsg(0, protoes.S2C_CREATE_ACC)
	retMsg.Pb = newPb

	HandleSend(acc.GetConn(), retMsg)
	return 0
}

func HandleLogin(acc *account.Acc, _msg *Msg) int {

	return 0
}
