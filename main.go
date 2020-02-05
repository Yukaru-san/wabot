package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
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
	fmt.Println(trans("hey wie gehts dir"))
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

func trans(in string) string {
	file := "tmp" + strconv.FormatInt(time.Now().Unix(), 10)
	err := ioutil.WriteFile(file, []byte(in), 0777)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	out, err := exec.Command("trans", ":en", "file://"+file).Output()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	os.Remove(file)
	return (string)(out)
}

func (cmd) HandleError(err error) {
	fmt.Println(err.Error())
}

func (cmd) HandleContactList(contacts []whatsapp.Contact) {
	for _, c := range contacts {
		fmt.Println(c.Name, c.Jid)
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
