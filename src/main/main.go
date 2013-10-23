package main

import (
	//"fmt"
	"logMgr"
	//"log"
)

func main() {
	// start log routine
	logMgr.InitServLog()
	logMgr.TestWriteLog()
}
