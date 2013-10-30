package account

import (
	//"log"
	//"logMgr"
	"net"
)

type Acc struct {
	accId uint32
	lv    uint8
	name  string
	conn  net.Conn //need to know which conn to send to
}

var AccMgrId map[uint32]*Acc
var AccMgrName map[string]*Acc

func NewAcc(c net.Conn) (acc *Acc) {
	acc = &Acc{
		accId: 0,
		lv:    0,
		name:  "",
		conn:  c,
	}

	return acc
}

func AddAcc2Mgr(a *Acc) bool {
	if a != nil && a.accId > 0 && a.name != "" {
		AccMgrId[a.accId] = a
		AccMgrName[a.name] = a
		return true
	}

	return false
}

func GetAccById(id uint32) (a *Acc) {
	a = AccMgrId[id]
	return a
}

func GetAccByName(name string) (a *Acc) {
	a = AccMgrName[name]
	return a
}

func (this *Acc) GetAccId() uint32 {
	return this.accId
}

func (this *Acc) GetLv() uint8 {
	return this.lv
}

func (this *Acc) GetName() string {
	return this.name
}

func (this *Acc) GetConn() net.Conn {
	return this.conn
}

func (this *Acc) Reset() {
	this.accId = 0
	this.lv = 0
	this.name = ""
	this.conn = nil
}
