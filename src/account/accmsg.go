package account

import proto "code.google.com/p/goprotobuf/proto"
import (
	"fmt"
	"glog"
	"gnet"
	"msg"
	"protoes"
)

type LoadingAcc struct {
	accName  string
	roleName string
	sess     *gnet.Gnet
}

//create queue
var CreatingCH = make(chan *LoadingAcc, 1000)

//loading queue
var LoadingCH = make(chan *LoadingAcc, 1000)

//loaded queue
var LoadedAccCH = make(chan *Acc, 1000)

type opfunc func(acc *Acc, _msg *msg.Msg) int

var msgHandler = make(map[protoes.OPCODE]opfunc)

//register message handler
func InitMsgHandler() {
	//msgHandler[protoes.C2S_CREATE_ACC] = HandleCreateAcc
	//msgHandler[protoes.C2S_LOGIN] = HandleLogin
}

func HandleMsg(sess *gnet.Gnet, _msg *msg.Msg) int {
	accId := sess.GetId()
	defer func() {
		if err := recover(); err != nil {
			glog.PushLogicLog(glog.Lerror, fmt.Sprintf("HandleMsg:%d %s error,from accId:%d", _msg.Op, _msg.Pb, accId))
		}
	}()

	//fmt.Printf("accId[%d] HandleMsg: opcode[%d]\n", accId, _msg.Op)

	switch _msg.Op {
	case protoes.C2S_LOGIN:
		return HandleLogin(sess, _msg)
	case protoes.C2S_CREATE_ACC:
		return HandleCreateAcc(sess, _msg)
	}

	acc := GetAccById(accId)
	f := msgHandler[_msg.Op]
	if f == nil || acc == nil {
		return -1
	}

	//fmt.Printf("HandleMsg: execute opfunc[%d]\n", _msg.Op)

	return f(acc, _msg)
}

//remove account from manager
//close Gnet routine
func HandleRelogin(acc *Acc) {
	//fmt.Printf("HandleRelogin: %v\n", acc)
	RemoveAccFromMgr(acc)
	acc.sess.Close()
	acc = nil
}

//return value:1, account not exist;2, database error;otherwise, account exist accId;
func HandleCreateResult(sess *gnet.Gnet, ret uint32) {
	pb := protoes.S2C_CREATE_ACC.CreateProto()
	newPb := pb.(*protoes.CreateAccRet)
	//newPb.Ok = proto.Uint32(ret)
	retMsg := msg.NewMsg(0, protoes.S2C_CREATE_ACC)
	retMsg.Pb = newPb
	sess.WRITECH() <- retMsg

	//fmt.Printf("HandleCreateResult: %d\n", ret)
}

func HandleLoginResult(sess *gnet.Gnet, ret uint32) {
	pb := protoes.S2C_LOGIN.CreateProto()
	newPb := pb.(*protoes.AccLoginRet)
	newPb.Ok = proto.Uint32(ret)
	retMsg := msg.NewMsg(0, protoes.S2C_LOGIN)
	retMsg.Pb = newPb
	sess.WRITECH() <- retMsg

	//fmt.Printf("HandleLoginResult: %d\n", ret)
}

//////////////////////////////////////////////////////////////////
func (this *Acc) SendAccInfo() {
	pb := protoes.S2C_ACC_INFO.CreateProto()
	newPb := pb.(*protoes.AccLoinInfo)
	newPb.AccId = proto.Uint32(this.GetAccId())
	newPb.Lv = proto.Uint32(uint32(this.GetLv()))
	newPb.AccName = proto.String(this.GetAccName())
	newPb.RoleName = proto.String(this.GetRoleName())
	retMsg := msg.NewMsg(0, protoes.S2C_ACC_INFO)
	retMsg.Pb = newPb
	this.sess.WRITECH() <- retMsg

	//fmt.Println("SendAccInfo:...")
}

func HandleCreateAcc(sess *gnet.Gnet, _msg *msg.Msg) int {
	_recvPb := _msg.Pb.(*protoes.CreateAcc)
	accName := *_recvPb.AccName
	roleName := *_recvPb.RoleName

	create := &LoadingAcc{
		accName:  accName,
		roleName: roleName,
		sess:     sess,
	}

	CreatingCH <- create

	return 0
}

func HandleLogin(sess *gnet.Gnet, _msg *msg.Msg) int {
	_recvPb := _msg.Pb.(*protoes.AccLogin)
	accName := *_recvPb.AccName

	load := &LoadingAcc{
		accName:  accName,
		roleName: "",
		sess:     sess,
	}

	LoadingCH <- load

	return 0
}
