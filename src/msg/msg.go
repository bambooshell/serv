package msg

import proto "code.google.com/p/goprotobuf/proto"

import (
	"account"
	"encoding/binary"
	"errors"
	"fmt"
	"glog"
	"io"
	"logMgr"
	"net"
	"protoes"
)

const (
	HEADER_LEN = 4
)

type Msg struct {
	Len uint16
	Op  protoes.OPCODE
	Pb  proto.Message
}

type ServMsg struct {
	Acc  *account.Acc
	SMsg *Msg
}

func NewMsg(l uint16, op protoes.OPCODE) *Msg {
	_msg := &Msg{
		Len: l,
		Op:  op,
	}
	return _msg
}

//receive a completed message from connection
//return the message,push it to the main routien
func HandleRecv(cc net.Conn) (_msg *Msg, err error) {
	rlen := 0
	header := make([]byte, HEADER_LEN) //include len and op
	rlen, err = io.ReadFull(cc, header)
	if err != nil || rlen != HEADER_LEN {
		fmt.Println(err.Error())
		return nil, err
	}

	totalLen := binary.BigEndian.Uint16(header[0:2])
	op := binary.BigEndian.Uint16(header[2:4])
	opcode := protoes.OPCODE(op)
	pb := opcode.CreateProto()
	if pb == nil {
		err = errors.New(fmt.Sprintf("opcode:%d don't has a protocol structure", op))
		fmt.Println(err.Error)
		return nil, err
	}

	_msg = NewMsg(totalLen, opcode)
	bodyLen := totalLen - HEADER_LEN
	body := make([]byte, bodyLen)
	_, err = io.ReadFull(cc, body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	err = proto.Unmarshal(body, pb)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	_msg.Pb = pb

	//fmt.Printf("recv msg: %d %s\n", op, _msg.Pb)
	return _msg, err
}

//send the message out
func HandleSend(cc net.Conn, _msg *Msg) (int, error) {
	body, err1 := proto.Marshal(_msg.Pb)
	if err1 != nil {
		return 0, err1
	}
	bodyLen := len(body)

	totalLen := HEADER_LEN + bodyLen
	data := make([]byte, totalLen)
	pos := 0

	//len
	binary.BigEndian.PutUint16(data[pos:pos+2], uint16(totalLen))
	pos += 2
	//opcode
	binary.BigEndian.PutUint16(data[pos:pos+2], uint16(_msg.Op))
	pos += 2
	//body
	copy(data[pos:pos+bodyLen], body)

	l, err := cc.Write(data)

	//fmt.Printf("send msg: %d %s\n", _msg.Op, _msg.Pb)
	return l, err
}

//belongs to main loop message routine
func HandleMsg(acc *account.Acc, _msg *Msg) int {
	defer func() {
		if err := recover(); err != nil {
			logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("HandleMsg:%d %s error,from acc:%d", _msg.Op, _msg.Pb, acc.GetAccId()))
		}
	}()

	//login account must has accId
	if protoes.C2S_CREATE_ACC == _msg.Op || protoes.C2S_LOGIN == _msg.Op {
		if acc.GetAccId() != 0 {
			acc.GetCloseWCH() <- true
		}
	} else {
		if acc.GetAccId() <= 0 {
			acc.GetCloseWCH() <- true
		}
	}

	f := msgHandler[_msg.Op]
	if f != nil {
		return f(acc, _msg)
	}
	//fmt.Printf("HandleMsg failed: %d\n", _msg.Op)
	return 0
}
