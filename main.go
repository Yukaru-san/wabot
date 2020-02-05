package main

import (
	"fmt"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

type cmd struct{}

var (
	contacList []whatsapp.Contact
	sess       whatsapp.Session
	conn       *whatsapp.Conn
	startTime  = uint64(time.Now().Unix())
)

func main() {

	sess, conn = HandleLogin()
	_ = sess

	// Add a complete MSG Handler
	conn.AddHandler(cmd{})

	go (func() {
		conn.AddHandler(messageHandler{})
	})()

	for {
		time.Sleep(time.Second)
	}
}

func (cmd) HandleError(err error) {
	fmt.Println(err.Error())
}

func (cmd) HandleContactList(contacts []whatsapp.Contact) {
	for _, c := range contacts {
		fmt.Println(c.Name)
		contacList = append(contacts, c)
	}
}

func jidToName(jid string) string {
	for _, c := range contacList {
		if c.Jid == jid {
			return c.Name
		}
	}
	return ""
}
