package logMgr

import (
	"config"
	"fmt"
	"glog"
	"log"
	"runtime"
)

var logic = glog.NewLogger(config.Llogname, config.Loglv)
var db = glog.NewLogger(config.LDblogname, config.Loglv)

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
	if lv >= glog.Ldebug && lv <= glog.Lfatal && lv >= logic.GetLogLv() {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		} else {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		newStr := fmt.Sprintf(" %s:%d: %s %s", file, line, glog.Levels[lv], str)
		if config.PrintLog {
			log.Println(newStr)
		}

		logic.PushLog(newStr)
	}
}

func PushDbLog(lv int, str string) {
	if lv >= glog.Ldebug && lv <= glog.Lfatal && lv >= db.GetLogLv() {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		} else {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}

		newStr := fmt.Sprintf(" %s:%d: %s %s", file, line, glog.Levels[lv], str)
		if config.PrintLog {
			log.Println(newStr)
		}

		db.PushLog(newStr)
	}
}

func TestWriteLog() {
	for i := 0; i < 100000000; i++ {
		PushLogicLog(glog.Lerror, "PushLogicLogglog.  PushLogicLogglog.PushLogicLogglog.PushLogicLog glog.PushLogicLog ")
		PushDbLog(glog.Linfo, "PushDblog glog.PushDblog  glog.PushDblogglog.PushDblog   glog.PushDblogglog.PushDblog")
	}
}
