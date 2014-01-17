package main

import (
	"account"
	"config"
	"fmt"
	"glog"
	"gnet"
	"log"
	"msg"
	"net"
	"protoes"
)

func startServListen(port int32) (net.Listener, error) {
	ip := fmt.Sprintf("localhost:%d", port)
	ls, err := net.Listen("tcp", ip)
	if err != nil {
		log.Fatalf("startServListen() failed: %s", err.Error())
	}
	fmt.Printf("start server listenning: %s\n", ip)
	return ls, err
}

//main routine message loop,
//all messages should be handled in the same routine
//account manager holds all of the loaded accounts
func startMsgLoop(msgch chan *msg.HangMsg) {
	for {
		select {
		case acc := <-account.LoadedAccCH:
			account.AddAcc2Mgr(acc)
		case s := <-msgch:
			account.HandleMsg(s.Sess.(*gnet.Gnet), s.S_msg)
		}
	}
}

func main() {
	//server config
	config.InitConf()

	// start log routine
	glog.InitServLog()
	// glog.TestWriteLog()
	// database initialization stuff
	account.StartDB()
	// account creation and loading stuff
	go account.AccCreation()
	go account.AccLoading()
	// register account message handlers
	account.InitMsgHandler()
	// register functions of creating message protocol
	protoes.RegProtoFunc()
	// every client message is dispatched via this channel
	_s_msg_ch := make(chan *msg.HangMsg, 10000)
	go startMsgLoop(_s_msg_ch)

	// start listenning socket connection
	listen, _ := startServListen(config.GetInt32("Gport", 0))
	for {
		c, err := listen.Accept()
		if err != nil {
			log.Print("listen.Accept() failed: %s", err.Error())
			continue
		}

		fmt.Println("main(): new connection accepted..")

		cc := gnet.NewSess(c)

		go cc.Poll(_s_msg_ch)
		go cc.Push()
	}
}
