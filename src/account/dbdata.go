package account

import (
	//"database/sql"
	"fmt"
	"glog"
	"logMgr"
)

func InitMaxAccId() {
	accId := uint32(0)
	stmt, err := dataDB.Prepare("SELECT accId FROM account ORDER BY accId DESC LIMIT 1;")
	if err != nil {
		logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Prepare InitMaxAccId:%s ", err.Error()))
		panic("InitMaxAccId Prepare error!")
	}
	defer stmt.Close()

	rows, err1 := stmt.Query()
	if err1 != nil {
		logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Query InitMaxAccId:%s", err.Error()))
		panic("InitMaxAccId Query error!")
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&accId)
		if err != nil {
			logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to rows.Next() InitMaxAccId:%s", err.Error()))
			panic("InitMaxAccId rows.Next() error!")
		}
		break
	}

	SetMaxAccId(accId)
}

//return value: >2, ok;1, account not exist;2, database error
func LoadAccId(accName string) uint32 {
	accId := uint32(1)
	stmt, err := dataDB.Prepare("SELECT accId FROM account WHERE accName = ?;")
	if err != nil {
		logMgr.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to Prepare LoadAccId:%s,%s", accName, err.Error()))
		accId = 2
		return accId
	}
	defer stmt.Close()

	rows, err1 := stmt.Query(accName)
	if err1 != nil {
		logMgr.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to Query LoadAccId:%s,%s", accName, err.Error()))
		accId = 2
		return accId
	}
	defer rows.Close()

	for rows.Next() {
		accId = 0
		err = rows.Scan(&accId)
		if err != nil {
			logMgr.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to rows.Next() LoadAccId:%s,%s", accName, err.Error()))
			accId = 2
			return accId
		}
		break
	}

	return accId
}

//check accName or roleName exist
//same as LoadAccId() but needs two query conditions
func QueryNameExist(accName, roleName string) uint32 {
	accId := uint32(1)
	stmt, err := dataDB.Prepare("SELECT accId FROM account WHERE accName = ? or roleName = ?;")
	if err != nil {
		logMgr.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to Prepare IsAccNameExist:%s,%s,%s", accName, roleName, err.Error()))
		accId = 2
		return accId
	}
	defer stmt.Close()

	rows, err1 := stmt.Query(accName, roleName)
	if err1 != nil {
		logMgr.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to Query IsAccNameExist:%s,%s,%s", accName, roleName, err.Error()))
		accId = 2
		return accId
	}
	defer rows.Close()

	for rows.Next() {
		accId = 0
		err = rows.Scan(&accId)
		if err != nil {
			logMgr.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to rows.Next() IsAccNameExist:%s,%s,%s", accName, roleName, err.Error()))
			accId = 2
			return accId
		}
		break
	}

	return accId
}

//if return 0, creation failed
func (this *Acc) CreateAccDB(accName, roleName string) uint32 {
	ret := QueryNameExist(accName, roleName)
	accId := uint32(0)
	if ret == 1 {
		accId = GetNextMaxAccId()
		this.SetAccId(accId)
		this.SetAccName(accName)
		this.SetRoleName(roleName)
		this.SetLv(1)

		this.insert_acc()

		this.ResetDS()
	}

	return accId
}

/*****************************************data load begin******************************************/
func (this *Acc) LoadAcc() {
	if this.accId <= 0 {
		logMgr.PushLogicLog(glog.Lerror, "no account accId")
		panic("loading error")
	}

	this.query_account()
}

func (this *Acc) query_account() {
	stmt, err := dataDB.Prepare("SELECT accName,roleName,lv FROM account WHERE accId = ?;")
	if err != nil {
		logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Prepare query_account:%d,%s", this.accId, err.Error()))
		panic("loading error")
	}
	defer stmt.Close()

	err = stmt.QueryRow(this.accId).Scan(&this.accName, &this.roleName, &this.lv)
	if err != nil {
		logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to QueryRow query_account:%d,%s", this.accId, err.Error()))
		panic("loading error")
	}
}

/*****************************************data load end******************************************/

/*****************************************data save begin******************************************/
//data saving includes insert,update and delete operations,according to data status
//*data status must reset after saving*
func SaveAcc(this *Acc) {
	//account data
	if DATA_STATUS_UPDATE == this.ds {
		this.update_acc()
		this.ResetDS()
	}

	//....

}

/*****************************************data save end******************************************/

/*****************************************data insert begin******************************************/
func (this *Acc) insert_acc() {
	stmt, err := dataDB.Prepare("INSERT INTO account (accId,lv,accName,roleName)VALUES(?,?,?,?);")
	if err != nil {
		logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Prepare insert_acc:%d,%s,%s", this.accId, this.accName, err.Error()))
		panic("insert_acc error")
	}
	defer stmt.Close()

	_, err1 := stmt.Exec(this.accId, this.lv, this.accName, this.roleName)
	if err1 != nil {
		logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Exec insert_acc:%d,%s,%s", this.accId, this.accName, err1.Error()))
		panic("insert_acc error")
	}
}

/*****************************************data insert end******************************************/

/*****************************************data update begin******************************************/
func (this *Acc) update_acc() {
	stmt, err := dataDB.Prepare("UPDATE account SET lv = ? WHERE accId = ?;")
	if err != nil {
		logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Prepare update_acc:%d,%s,%s", this.accId, this.accName, err.Error()))
		panic("update_acc error")
	}
	defer stmt.Close()

	_, err1 := stmt.Exec(this.lv, this.accId)
	if err1 != nil {
		logMgr.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Exec update_acc:%d,%s,%s", this.accId, this.accName, err1.Error()))
		panic("update_acc error")
	}
}

/*****************************************data update end******************************************/

/*****************************************data delete begin******************************************/

/*****************************************data delete end******************************************/
