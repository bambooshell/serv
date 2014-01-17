package account

import (
	//"database/sql"
	"fmt"
	"glog"
)

func InitMaxAccId() {
	accId := uint32(maxAccId)
	stmt, err := dataDB.Prepare("SELECT accId FROM account ORDER BY accId DESC LIMIT 1;")
	if err != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Prepare InitMaxAccId:%s ", err.Error()))
		panic("InitMaxAccId Prepare error!")
	}
	defer stmt.Close()

	rows, err1 := stmt.Query()
	if err1 != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Query InitMaxAccId:%s", err.Error()))
		panic("InitMaxAccId Query error!")
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&accId)
		if err != nil {
			glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to rows.Next() InitMaxAccId:%s", err.Error()))
			panic("InitMaxAccId rows.Next() error!")
		}
		break
	}
	//fmt.Printf("new accId: %d\n", accId)
	SetMaxAccId(accId)
}

//return value:1, account not exist;2, database error; otherwise,account exist accId;
func LoadAccId(accName string) uint32 {
	accId := uint32(1)
	stmt, err := dataDB.Prepare("SELECT accId FROM account WHERE accName = ?;")
	if err != nil {
		glog.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to Prepare LoadAccId:%s,%s", accName, err.Error()))
		accId = 2
		return accId
	}
	defer stmt.Close()

	rows, err1 := stmt.Query(accName)
	if err1 != nil {
		glog.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to Query LoadAccId:%s,%s", accName, err.Error()))
		accId = 2
		return accId
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&accId)
		if err != nil {
			glog.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to rows.Next() LoadAccId:%s,%s", accName, err.Error()))
			accId = 2
			return accId
		}
		break
	}

	return accId
}

//return value:1, account not exist;2, database error; otherwise,account exist accId;
//same as LoadAccId() but needs two query conditions
//check accName or roleName exist
func QueryNameExist(accName, roleName string) uint32 {
	accId := uint32(1)
	stmt, err := dataDB.Prepare("SELECT accId FROM account WHERE accName = ? or roleName = ?;")
	if err != nil {
		glog.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to Prepare IsAccNameExist:%s,%s,%s", accName, roleName, err.Error()))
		accId = 2
		return accId
	}
	defer stmt.Close()

	rows, err1 := stmt.Query(accName, roleName)
	if err1 != nil {
		glog.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to Query IsAccNameExist:%s,%s,%s", accName, roleName, err.Error()))
		accId = 2
		return accId
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&accId)
		if err != nil {
			glog.PushLogicLog(glog.Lwarn, fmt.Sprintf("failed to rows.Next() IsAccNameExist:%s,%s,%s", accName, roleName, err.Error()))
			accId = 2
			return accId
		}
		break
	}

	return accId
}

//create account db
func CreateAccDB(load *LoadingAcc) (acc *Acc) {
	ret := QueryNameExist(load.accName, load.roleName)
	if ret == 1 {
		accId := GetNextMaxAccId()

		acc = NewAcc(*(load.sess))
		acc.SetAccId(accId)
		acc.sess.SetId(ret)
		acc.SetAccName(load.accName)
		acc.SetRoleName(load.roleName)
		acc.SetLv(1)

		ret = acc.insert_acc()
		acc.ResetDS()
	} else {
		acc = nil
	}

	HandleLoginResult(load.sess, ret)
	return acc
}

//load account db
func LoadAcc(load *LoadingAcc) {
	ret := LoadAccId(load.accName)
	if ret == 1 || ret == 2 { //account not exist or db error
		HandleLoginResult(load.sess, ret)
	}
	//account exist
	acc := NewAcc(*(load.sess))
	acc.SetAccId(ret)
	acc.sess.SetId(ret)
	ret = acc.query_account()
	if ret == 0 {
		LoadedAccCH <- acc
	} else {
		acc = nil
		ret = 2
		HandleLoginResult(load.sess, ret)
	}
}

//return 0,ok; 1,failed
func (this *Acc) query_account() uint32 {
	stmt, err := dataDB.Prepare("SELECT accName,roleName,lv FROM account WHERE accId = ?;")
	if err != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Prepare query_account:%d,%s", this.accId, err.Error()))
		return 1
	}
	defer stmt.Close()

	err = stmt.QueryRow(this.accId).Scan(&this.accName, &this.roleName, &this.lv)
	if err != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to QueryRow query_account:%d,%s", this.accId, err.Error()))
		return 1
	}

	return 0
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
func (this *Acc) insert_acc() uint32 {
	stmt, err := dataDB.Prepare("INSERT INTO account (accId,lv,accName,roleName)VALUES(?,?,?,?);")
	if err != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Prepare insert_acc:%d,%s,%s", this.accId, this.accName, err.Error()))
		return 1
	}
	defer stmt.Close()

	_, err1 := stmt.Exec(this.accId, this.lv, this.accName, this.roleName)
	if err1 != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Exec insert_acc:%d,%s,%s", this.accId, this.accName, err1.Error()))
		return 1
	}

	return 0
}

/*****************************************data insert end******************************************/

/*****************************************data update begin******************************************/
func (this *Acc) update_acc() {
	stmt, err := dataDB.Prepare("UPDATE account SET lv = ? WHERE accId = ?;")
	if err != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Prepare update_acc:%d,%s,%s", this.accId, this.accName, err.Error()))
		panic("update_acc error")
	}
	defer stmt.Close()

	_, err1 := stmt.Exec(this.lv, this.accId)
	if err1 != nil {
		glog.PushLogicLog(glog.Lerror, fmt.Sprintf("failed to Exec update_acc:%d,%s,%s", this.accId, this.accName, err1.Error()))
		panic("update_acc error")
	}
}

/*****************************************data update end******************************************/

/*****************************************data delete begin******************************************/

/*****************************************data delete end******************************************/
