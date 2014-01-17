package glog

import (
	"config"
	"fmt"
	"runtime"
)

var logic = NewLogger(config.GetString("Llogname", "logic"), config.GetInt32("Loglv", Lwarn))
var db = NewLogger(config.GetString("LDblogname", "db"), config.GetInt32("Loglv", Lwarn))

var is_print = false

func InitServLog() {
	is_print = config.GetBool("PrintLog", true)
	go HandleLogs()
}

//poll logs from Logger.ch
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

func PushLogicLog(lv int32, str string) {
	if lv >= Ldebug && lv <= Lfatal && lv >= logic.GetLogLv() {
		_, file, line, ok := runtime.Caller(1)
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
		newStr := fmt.Sprintf(" %s:%d: %s %s", file, line, Levels[lv], str)
		if is_print {
			fmt.Println(newStr)
		}

		logic.PushLog(newStr)
	}
}

func PushDbLog(lv int32, str string) {
	if lv >= Ldebug && lv <= Lfatal && lv >= db.GetLogLv() {
		_, file, line, ok := runtime.Caller(1)
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

		newStr := fmt.Sprintf(" %s:%d: %s %s", file, line, Levels[lv], str)
		if is_print {
			fmt.Println(newStr)
		}

		db.PushLog(newStr)
	}
}
