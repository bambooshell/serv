package account

import (
	_ "code.google.com/mysql"
	"config"
	"database/sql"
	"fmt"
	"log"
)

//there are 3 underlying database connection

//configture database
var confDB *sql.DB

//account database
var dataDB *sql.DB
var dataCH chan string

//account operation log database
var logDB *sql.DB
var logCH chan string

func InitDBConnect() {
	//format "user:password@tcp(ip:port)/database"
	dbstuff := fmt.Sprintf("%s:%s@tcp(%s:%d)", config.DBuser, config.DBpw, config.DBip, config.DBport)

	cname := dbstuff + "/" + config.DBnameConf + "?charset=utf8"
	//fmt.Println("connecting to db:" + cname)
	dbconf, err1 := sql.Open("mysql", cname)
	if err1 != nil {
		log.Panicf("failed to open db: %s ", cname)
	} else if dbconf.Ping() != nil {
		log.Panicf("failed to connect to db: %s ", cname)
	}
	confDB = dbconf

	dname := dbstuff + "/" + config.DBnameData + "?charset=utf8"
	//fmt.Println("connecting to db:" + dname)
	dbdata, err2 := sql.Open("mysql", dname)
	if err2 != nil {
		log.Panicf("failed to open db: %s ", dname)
	} else if dbdata.Ping() != nil {
		log.Panicf("failed to connect to db: %s ", dname)
	}
	dataDB = dbdata

	lname := dbstuff + "/" + config.DBlog + "?charset=utf8"
	//fmt.Println("connecting to db:" + lname)
	dblog, err3 := sql.Open("mysql", lname)
	if err3 != nil {
		log.Panicf("failed to open db: %s ", lname)
	} else if dblog.Ping() != nil {
		log.Panicf("failed to connect to db: %s ", lname)
	}
	logDB = dblog
}

func StartDB() {
	InitRoleConf()
	InitMaxAccId()
}
