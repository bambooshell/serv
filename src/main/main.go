package main

import (
	"account"
	"config"
	"conn"
	"fmt"
	"log"
	"logMgr"
	"msg"
	"net"
	"protoes"
)

func startServListen(port int) (net.Listener, error) {
	ip := fmt.Sprintf("localhost:%d", port)
	ls, err := net.Listen("tcp", ip)
	if err != nil {
		log.Fatalf("startServListen() failed: %s", err)
	}
	fmt.Printf("start server listenning: %s\n", ip)
	return ls, err
}

//main routine message loop,
//clients' logical messages will be dispatched to its handling function
//message handling should be in the same routine
func startMsgLoop(ch chan *msg.ServMsg) {
	for {
		s := <-ch
		msg.HandleMsg(s.Acc, s.SMsg)
	}
}

func main() {
	// start log routine
	logMgr.InitServLog()
	//logMgr.TestWriteLog()

	//database initialization stuff
	account.InitDBConnect()
	account.StartDB()

	// register functions of creating message protocol
	protoes.RegisterProtoFuncs()

	// register message handlers
	msg.InitMsgHandler()

	//start listenning port
	listen, _ := startServListen(config.Gport)

	ServMsgChan := make(chan *msg.ServMsg)
	//start server message handler
	go startMsgLoop(ServMsgChan)

	//start loop,listenning socket connection
	for {
		newConn, err := listen.Accept()
		if err != nil {
			log.Print("listen.Accept() failed")
		}
		fmt.Println("accept new connection")
		//every connection has a read/write goroutine
		cc := conn.NewConn(newConn, account.NewAcc(newConn))
		go cc.AcceptMsg()
		go cc.SelectMsg(ServMsgChan)
	}
}
