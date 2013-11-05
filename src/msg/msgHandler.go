package msg

import (
	"account"
	//"fmt"
	"protoes"
)
import proto "code.google.com/p/goprotobuf/proto"

var msgHandler = make(map[protoes.OPCODE]func(acc *account.Acc, _msg *Msg) int)

//message handling functions
//logic requests handled in the same routine
func InitMsgHandler() {
	msgHandler[protoes.C2S_CREATE_ACC] = HandleCreateAcc
	msgHandler[protoes.C2S_LOGIN] = HandleLogin
}

func SendAccInfo(acc *account.Acc) {
	pb := protoes.S2C_ACC_INFO.CreateProto()
	newPb := pb.(*protoes.AccLoinInfo)
	newPb.AccId = proto.Uint32(acc.GetAccId())
	newPb.Lv = proto.Uint32(uint32(acc.GetLv()))
	newPb.AccName = proto.String(acc.GetAccName())
	newPb.RoleName = proto.String(acc.GetRoleName())
	retMsg := NewMsg(0, protoes.S2C_ACC_INFO)
	retMsg.Pb = newPb
	acc.GetWCH() <- retMsg
}

func HandleCreateAcc(acc *account.Acc, _msg *Msg) int {
	_recvPb := _msg.Pb.(*protoes.CreateAcc)
	accName := *_recvPb.AccName
	roleName := *_recvPb.RoleName

	ret := acc.CreateAccDB(accName, roleName)
	if ret == 0 { //failed
		pb := protoes.S2C_CREATE_ACC.CreateProto()
		newPb := pb.(*protoes.CreateAccRet)
		newPb.Ok = proto.Int32(int32(ret))
		retMsg := NewMsg(0, protoes.S2C_CREATE_ACC)
		retMsg.Pb = newPb
		acc.GetWCH() <- retMsg
	} else {
		//send account login info
		SendAccInfo(acc)
	}
	return 0
}

func HandleLogin(acc *account.Acc, _msg *Msg) int {
	_recvPb := _msg.Pb.(*protoes.AccLogin)
	ret := account.LoadAccId(*_recvPb.AccName)
	if ret > 2 { //account exist
		acc.SetAccId(ret)
		acc.LoadAcc()

		//send account login info
		SendAccInfo(acc)

	} else if ret == 2 { //db error
		acc.GetCloseWCH() <- true
	} else if ret == 1 { // need to create new account
		pb := protoes.S2C_LOGIN.CreateProto()
		newPb := pb.(*protoes.AccLoginRet)
		newPb.Ok = proto.Uint32(ret)
		retMsg := NewMsg(0, protoes.S2C_LOGIN)
		retMsg.Pb = newPb
		acc.GetWCH() <- retMsg
	}

	return 0
}
