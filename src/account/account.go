package account

import (
	"gnet"
)

//account data state
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
	sess     gnet.Gnet
}

func NewAcc(c gnet.Gnet) (acc *Acc) {
	acc = &Acc{
		accId:    0,
		lv:       1,
		accName:  "",
		roleName: "",
		ds:       DATA_STATUS_NONE,
		sess:     c,
	}

	return acc
}

func (this *Acc) UpdateDS() {
	if this.ds < DATA_STATUS_UPDATE {
		this.ds = DATA_STATUS_UPDATE
	}
}

func (this *Acc) ResetDS() {
	this.ds = DATA_STATUS_NONE
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
	this.UpdateDS()
}

func (this *Acc) GetAccName() string {
	return this.accName
}

func (this *Acc) SetAccName(name string) {
	this.accName = name
	this.UpdateDS()
}

func (this *Acc) GetRoleName() string {
	return this.roleName
}

func (this *Acc) SetRoleName(name string) {
	this.roleName = name
	this.UpdateDS()
}

func (this *Acc) GetConn() gnet.Gnet {
	return this.sess
}
