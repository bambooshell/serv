package gnet

import (
	"fmt"
	"glog"
	"msg"
	"net"
	"time"
)

//conn wrapper,
//binded with id
type Gnet struct {
	conn    net.Conn
	id      uint32
	readch  chan *msg.Msg
	writech chan *msg.Msg
	werr    chan bool
}

func NewSess(c net.Conn) (s *Gnet) {
	s = &Gnet{
		conn:    c,
		id:      0,
		readch:  make(chan *msg.Msg, 50),
		writech: make(chan *msg.Msg, 50),
		werr:    make(chan bool),
	}

	return s
}

func (this *Gnet) GetId() uint32 {
	return this.id
}

func (this *Gnet) SetId(id uint32) {
	this.id = id
}

func (this *Gnet) READCH() chan *msg.Msg {
	return this.readch
}

func (this *Gnet) WRITECH() chan *msg.Msg {
	return this.writech
}

//socket connection close
func (this *Gnet) Close() {
	defer this.conn.Close()

	//do cleanup stuff
	glog.PushLogicLog(glog.Linfo, fmt.Sprintf("id[%d]: connection closed", this.GetId()))
}

//read
func (this *Gnet) read(ch chan bool) {
	defer close(this.readch)
	fmt.Println("read exist..")

	for {
		_msg, err := msg.HandleRecv(this.conn)
		if err != nil {
			ch <- true
			return
		} else if _msg != nil {
			this.readch <- _msg
		}
	}
}

//message read routine
func (this *Gnet) Poll(ch chan *msg.HangMsg) {
	defer this.Close()
	defer fmt.Println("Poll exist..")

	breakch := make(chan bool)
	go this.read(breakch)

	for {
		select {
		case <-breakch:
			this.werr <- true
			return
		case _msg := <-this.READCH():
			//fmt.Printf("Poll: %v\n", _msg)
			_h_msg := &msg.HangMsg{this, _msg}
			ch <- _h_msg
		}
	}
}

//message handling routine
//1.push msg to main message loop
//2.send out message
//3.account saving
func (this *Gnet) Push() {
	defer this.Close()
	defer close(this.writech)
	defer fmt.Println("Push exist..")

	saveTick := time.Tick(1e9 * 60 * 2)

	for {
		select {
		case <-this.werr:
			return
		case _msg := <-this.WRITECH():
			//fmt.Printf("Push: %v\n", _msg)
			_, err := msg.HandleSend(this.conn, _msg)
			if err != nil {
				return
			}
		case <-saveTick:
			//fmt.Println("tick")
		}
	}
}
