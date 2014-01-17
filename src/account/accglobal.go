package account

import (
	"config"
	"fmt"
	"sync"
)

var AccMgrId = make(map[uint32]*Acc)
var AccMgrName = make(map[string]*Acc)

//current max account id that had beed register by the server
var maxAccId = uint32(config.GetInt32("ServId", 1) * config.GetInt32("ServBase", 10000000))

//generate new account id, need a lock
var maxAccIdLock sync.Mutex

//require account mapping, need a lock
var accRequireLock sync.RWMutex

//for initialization
func SetMaxAccId(mid uint32) {
	if mid > maxAccId {
		maxAccId = mid
	}
}

//create unique account id
func GetNextMaxAccId() uint32 {
	maxAccIdLock.Lock()
	defer maxAccIdLock.Unlock()

	maxAccId += 1

	return maxAccId
}

//push new account to the global mapping
func AddAcc2Mgr(a *Acc) bool {
	accRequireLock.Lock()
	defer accRequireLock.Unlock()

	//fmt.Printf("AddAcc2Mgr: %v\n", a)
	if a != nil && a.accId > 0 && a.accName != "" {
		//if account already loaded, close the previous connection
		old_acc, ok := AccMgrId[a.accId]
		if ok {
			fmt.Println("warning: AddAcc2Mgr: account relogin")
			AccMgrId[a.accId] = nil
			AccMgrName[a.accName] = nil
			old_acc.sess.Close()
		}

		AccMgrId[a.accId] = a
		AccMgrName[a.accName] = a

		//send login successful message
		a.SendAccInfo()

		return true
	}

	return false
}

func RemoveAccFromMgr(a *Acc) {
	accRequireLock.Lock()
	defer accRequireLock.Unlock()

	if a != nil {
		AccMgrId[a.accId] = nil
		AccMgrName[a.accName] = nil
	}
}

func GetAccById(id uint32) (a *Acc) {
	accRequireLock.RLock()
	defer accRequireLock.RUnlock()

	a = AccMgrId[id]
	return a
}

func GetAccByName(name string) (a *Acc) {
	accRequireLock.RLock()
	defer accRequireLock.RUnlock()

	a = AccMgrName[name]
	return a
}

//account creation routine
//create one account db at a time
func AccCreation() {
	for {
		select {
		case c := <-CreatingCH:
			newAcc := CreateAccDB(c)
			if newAcc != nil {
				LoadedAccCH <- newAcc
			}
		}
	}
}

//account loading routine
//parallel query
func AccLoading() {
	for {
		select {
		case c := <-LoadingCH:
			/*
				acc := LoadAcc(c)
				LoadedAccCH <- acc
			*/
			go LoadAcc(c)
		}
	}
}
