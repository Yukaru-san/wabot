package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

type cmd struct{}

var (
	botname    = "Yukaru-Bot"
	contacList []whatsapp.Contact
	sess       whatsapp.Session
	conn       *whatsapp.Conn
	startTime  = uint64(time.Now().Unix())
	running    = true
	encrypKey  = []byte("r4gyXrWSPXzvpBZJ")
)

func main() {
	// Login
	sess, conn = HandleLogin()
	_ = sess

	// Add the command handler
	conn.AddHandler(cmd{})

	// Add the message handler
	go (func() {
		conn.AddHandler(messageHandler{})
	})()

	// Autosave
	go (func() {
		for running {
			time.Sleep(time.Minute * 2)
			SaveUsersToDisk(usersFile)
		}
	})()

	// Handle console input
	go HandleConsoleInput()

	for running {
		time.Sleep(time.Minute)
	}
}

// HandleConsoleInput is used to execute Admin-commands in console
func HandleConsoleInput() {
	reader := bufio.NewReader(os.Stdin)

	for running {
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", 1)

		if strings.ToLower(input) == "save" {
			println("saving...")
			SaveUsersToDisk(usersFile)
		} else if strings.ToLower(input) == "quit" {
			SaveUsersToDisk(usersFile)
			running = false
		}
	}
}

// HandleError prints potential errors
func (cmd) HandleError(err error) {
	fmt.Println(err.Error())
}

// HandleContactList fills the contacList on load
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
