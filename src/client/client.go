package main

import (
	"config"
	"fmt"
	"log"
	"msg"
	"net"
	"protoes"
)

import proto "code.google.com/p/goprotobuf/proto"

func main() {

	config.InitConf()
	// register functions of creating message protocol
	protoes.RegProtoFunc()

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", config.GetInt32("Gport", 0)))
	if err != nil {
		log.Fatalf("net.Dial() failed: %s", err.Error())
	}
	defer conn.Close()

	TestLogin(conn)
}

func TestLogin(conn net.Conn) {
	fmt.Println("TestLogin..")
	pb := protoes.C2S_LOGIN.CreateProto()

	newPb := pb.(*protoes.AccLogin)
	newPb.AccName = proto.String("bambooshell")

	_msg := msg.NewMsg(0, protoes.C2S_LOGIN)
	_msg.Pb = newPb

	msg.HandleSend(conn, _msg)

	retMsg, _ := msg.HandleRecv(conn)
	testHandleRetMsg(conn, retMsg)
}

func testHandleRetMsg(conn net.Conn, retMsg *msg.Msg) {
	if retMsg == nil {
		fmt.Println("error: recv nil msg")
		return
	}
	switch retMsg.Op {
	case protoes.S2C_LOGIN:
		testRetLogin(conn, retMsg)
	case protoes.S2C_CREATE_ACC:
		testRetCreate(retMsg)
	case protoes.S2C_ACC_INFO:
		testRetAccInfo(retMsg)
	}
}

func testRetLogin(conn net.Conn, _msg *msg.Msg) {
	retPb := _msg.Pb.(*protoes.AccLoginRet)
	if *retPb.Ok == 1 { //create
		pb := protoes.C2S_CREATE_ACC.CreateProto()

		newPb := pb.(*protoes.CreateAcc)
		newPb.AccName = proto.String("bambooshell")
		newPb.RoleName = proto.String("evanchen")

		_msg2 := msg.NewMsg(0, protoes.C2S_CREATE_ACC)
		_msg2.Pb = newPb

		fmt.Printf("send create..%v\n", _msg2)
		msg.HandleSend(conn, _msg2)

		retMsg, _ := msg.HandleRecv(conn)
		testHandleRetMsg(conn, retMsg)
	} else {
		fmt.Println(_msg)
	}
}

func testRetCreate(_msg *msg.Msg) {
	fmt.Println("create failed")
}

func testRetAccInfo(_msg *msg.Msg) {
	fmt.Printf("account info: %v\n", _msg)
}
