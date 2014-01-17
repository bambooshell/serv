package msg

import proto "code.google.com/p/goprotobuf/proto"

import (
	"encoding/binary"
	"fmt"
	"glog"
	"io"
	"net"
	"protoes"
)

const (
	HEADER_LEN = 4
)

type HangMsg struct {
	Sess  interface{}
	S_msg *Msg
}

type Msg struct {
	Len uint16
	Op  protoes.OPCODE
	Pb  proto.Message
}

func NewMsg(l uint16, op protoes.OPCODE) *Msg {
	_msg := &Msg{
		Len: l,
		Op:  op,
	}
	return _msg
}

//receive a completed message from connection
//return the message and push it to the main message handling routine
//the message will be eventually delivered to its handler
func HandleRecv(cc net.Conn) (*Msg, error) {
	header := make([]byte, HEADER_LEN) //include len and op
	rlen, err := io.ReadFull(cc, header)
	if err != nil || rlen != HEADER_LEN {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to read message header: %s", err.Error()))
		return nil, err
	}

	totalLen := binary.BigEndian.Uint16(header[0:2])
	op := binary.BigEndian.Uint16(header[2:4])
	opcode := protoes.OPCODE(op)
	pb := opcode.CreateProto()
	if pb == nil {
		glog.PushLogicLog(glog.Lwarn, fmt.Sprintf("opcode[%d]: undefined protocol", op))
		return nil, nil
	}

	_msg := NewMsg(totalLen, opcode)
	bodyLen := totalLen - HEADER_LEN
	body := make([]byte, bodyLen)
	_, err = io.ReadFull(cc, body)
	if err != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to read message body: %s", err.Error()))
		return nil, err
	}

	err = proto.Unmarshal(body, pb)
	if err != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Unmarshal received message[%d]: %s", op, err.Error()))
		return nil, nil
	}
	_msg.Pb = pb

	return _msg, nil
}

//send the message out
//return the length of the message
func HandleSend(cc net.Conn, _msg *Msg) (int, error) {
	l := 0
	body, err := proto.Marshal(_msg.Pb)
	if err != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Marshal message[%d]: %s", _msg.Op, err.Error()))
		return l, nil
	}
	bodyLen := len(body)

	totalLen := HEADER_LEN + bodyLen
	data := make([]byte, totalLen)
	pos := l

	//len
	binary.BigEndian.PutUint16(data[pos:pos+2], uint16(totalLen))
	pos += 2
	//opcode
	binary.BigEndian.PutUint16(data[pos:pos+2], uint16(_msg.Op))
	pos += 2
	//body
	copy(data[pos:pos+bodyLen], body)

	l, err = cc.Write(data)
	if err != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to send message[%d]: %s", _msg.Op, err.Error()))
		return l, err
	}

	return l, nil
}
