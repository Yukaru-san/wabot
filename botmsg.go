package wabot

import (
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

// Command struct
type Command struct {
	prefix   string
	function func(whatsapp.TextMessage)
}

// List of implemented commands
var commands []Command

// HandleBotMsg checks if a message is a command and executes the first possible command
func HandleBotMsg(message whatsapp.TextMessage) {
	// Run through command list and execute if possible
	for i := 0; i < len(commands); i++ {
		if strings.HasPrefix(strings.ToLower(message.Text), strings.ToLower(commands[i].prefix)) {
			go commands[i].function(message)
			break
		}
	}
}

// WriteTextMessage sends a given string as text
func WriteTextMessage(text string, remoteJid string) {
	// Create the msg struct
	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: remoteJid,
		},
		Text: text,
	}
	// And send it
	conn.Send(msg)
}
