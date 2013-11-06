package account

import (
	//"log"
	//"logMgr"
	"config"
	"net"
	"sync"
)

//data state
const (
	DATA_STATUS_NONE = iota
	DATA_STATUS_UPDATE
	DATA_STATUS_INSERT
)

type Acc struct {
	accId    uint32
	lv       uint8
	accName  string
	roleName string
	ds       uint8
	conn     net.Conn //need to know which conn to send to
	writech  chan interface{}
	closewch chan bool
}

var AccMgrId map[uint32]*Acc
var AccMgrName map[string]*Acc

var maxAccId = uint32(config.ServId * config.ServBase)
var maxAccIdLock sync.Mutex

//for initialization
func SetMaxAccId(mid uint32) {
	if mid > maxAccId {
		maxAccId = mid
	}
}

//for create unique accId
func GetNextMaxAccId() uint32 {
	maxAccIdLock.Lock()
	defer maxAccIdLock.Unlock()

	maxAccId += 1

	return maxAccId
}

func NewAcc(c net.Conn) (acc *Acc) {
	acc = &Acc{
		accId:    0,
		lv:       1,
		accName:  "",
		roleName: "",
		ds:       DATA_STATUS_NONE,
		conn:     c,
		writech:  make(chan interface{}, 50),
		closewch: make(chan bool, 1),
	}

	return acc
}

func AddAcc2Mgr(a *Acc) bool {
	if a != nil && a.accId > 0 && a.accName != "" {
		AccMgrId[a.accId] = a
		AccMgrName[a.accName] = a
		return true
	}

	return false
}

func (this *Acc) UpdateDS() {
	if this.ds < DATA_STATUS_UPDATE {
		this.ds = DATA_STATUS_UPDATE
	}
}

func (this *Acc) ResetDS() {
	this.ds = DATA_STATUS_NONE
}

func (this *Acc) GetWCH() chan interface{} {
	return this.writech
}

func (this *Acc) GetCloseWCH() chan bool {
	return this.closewch
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

func (this *Acc) SetAccId(accId uint32) {
	this.accId = accId
}

func (this *Acc) GetLv() uint8 {
	return this.lv
}

func (this *Acc) SetLv(lv uint8) {
	this.lv = lv
}

func (this *Acc) GetAccName() string {
	return this.accName
}

func (this *Acc) SetAccName(name string) {
	this.accName = name
}

func (this *Acc) GetRoleName() string {
	return this.roleName
}

func (this *Acc) SetRoleName(name string) {
	this.roleName = name
}

func (this *Acc) GetConn() net.Conn {
	return this.conn
}

func (this *Acc) Reset() {
	this.accId = 0
	this.lv = 0
	this.accName = ""
	this.roleName = ""
	this.conn = nil
}
