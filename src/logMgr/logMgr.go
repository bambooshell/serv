package logMgr

import (
	"config"
	"glog"
)

var logic = glog.NewLogger(config.Logic_logname, config.Loglv)
var db = glog.NewLogger(config.Db_logname, config.Loglv)

func InitServLog() {
	go HandleLogs()
}

//poll logs from Logger.ch, if the log written size >= MAX_LOGFILE_SIZE,
//close old file, create a new log.Logger
func HandleLogs() {
	for {
		select {
		case str := <-logic.CH():
			logic.PrintLog(str)
		case str := <-db.CH():
			db.PrintLog(str)
		}
	}
}

func PushLogicLog(lv int, str string) {
	logic.PushLog(lv, str)
}

func PushDbLog(lv int, str string) {
	db.PushLog(lv, str)
}

func TestWriteLog() {
	for i := 0; i < 100000000; i++ {
		PushLogicLog(glog.Lerror, "PushLogicLogglog.  PushLogicLogglog.PushLogicLogglog.PushLogicLog glog.PushLogicLog ")
		PushDbLog(glog.Linfo, "PushDblog glog.PushDblog  glog.PushDblogglog.PushDblog   glog.PushDblogglog.PushDblog")
	}
}
