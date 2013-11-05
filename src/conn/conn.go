package conn

import (
	"account"
	"fmt"
	"glog"
	"logMgr"
	"msg"
	"net"
	"time"
)

type Conn struct {
	conn     net.Conn
	acc      *account.Acc //need to know who to handle msg
	readch   chan *msg.Msg
	closerch chan bool
}

func NewConn(c net.Conn, a *account.Acc) (newConn *Conn) {
	newConn = &Conn{
		conn:     c,
		acc:      a,
		readch:   make(chan *msg.Msg, 50),
		closerch: make(chan bool, 1),
	}

	return newConn
}

func (this *Conn) GetAcc() *account.Acc {
	return this.acc
}

func (this *Conn) SetAcc(a *account.Acc) {
	this.acc = a
}

func (this *Conn) GetRCH() chan *msg.Msg {
	return this.readch
}

func (this *Conn) Close(accId uint32) {
	defer this.conn.Close()

	logMgr.PushLogicLog(glog.Linfo, fmt.Sprintf("acc:%d connection closed", accId))
}

//message read routine
func (this *Conn) AcceptMsg() {
	acc := this.GetAcc()
	defer this.Close(acc.GetAccId())

	for {
		select {
		case <-this.closerch: //heartbeat ? write error ?
			//notify SelectMsg routine
			acc.GetCloseWCH() <- true
			return
		default:
			_msg, err := msg.HandleRecv(this.conn)
			//fmt.Printf("%v\n", err)
			if err != nil { //close connect
				logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("%d: AcceptMsg() error", acc.GetAccId()))
				acc.Reset()

				//notify SelectMsg routine
				acc.GetCloseWCH() <- true

				return
			} else if _msg != nil {
				this.GetRCH() <- _msg
			}
		}
	}
}

//message handling routine
func (this *Conn) SelectMsg(ch chan *msg.ServMsg) {
	acc := this.GetAcc()
	defer this.Close(acc.GetAccId())
	saveAcc := time.Tick(1e9 * 60 * 2)

	for {
		select {
		case _msg := <-this.GetRCH(): //push msg to main message loop
			servMsg := &msg.ServMsg{this.GetAcc(), _msg}
			ch <- servMsg
		case _msg := <-acc.GetWCH(): //send out
			msg.HandleSend(this.conn, _msg.(*msg.Msg))
		case <-saveAcc:
			account.SaveAcc(acc)
		case <-acc.GetCloseWCH():
			return
		}
	}
}
