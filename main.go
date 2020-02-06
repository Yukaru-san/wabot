package main

import (
	"fmt"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

type cmd struct{}

var (
	botname    = "Test-Gruppe"
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
		contacList = append(contacts, c)
	}
}

// MessageToJid find the right Jid in a message
func MessageToJid(message whatsapp.TextMessage) string {

	authorID := ""

	if message.Info.Source.Participant != nil {
		authorID = *message.Info.Source.Participant
	} else {
		authorID = message.Info.RemoteJid // Personennamen
	}

	return authorID
}

// MessageToName Converts and entire message to the corresponding name
// -> allows finding the name in groups
func MessageToName(message whatsapp.TextMessage) string {

	authorID := ""

	if len(message.Info.Source.GetPushName()) > 0 {
		return message.Info.Source.GetPushName()
	} else if message.Info.Source.Participant != nil {
		authorID = *message.Info.Source.Participant
	} else {
		authorID = message.Info.RemoteJid // Personennamen
	}

	return jidToName(authorID)
}

// Converts an ID to the corresponding name
func jidToName(jid string) string {
	for _, c := range contacList {
		if c.Jid == jid {
			if c.Name == botname {
				return "Admin"
			}

			return c.Name
		}
	}
	return "{undefined}"
}
