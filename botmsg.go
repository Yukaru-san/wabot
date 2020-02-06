package main

import (
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

// HandleBotMsg handles messages sent to the bot
func HandleBotMsg(message whatsapp.TextMessage, conn *whatsapp.Conn) {
	// Function 1: Translate
	if (!strings.HasPrefix(message.Text, "/tf") && !strings.HasPrefix(message.Text, "/tt")) && (strings.HasPrefix(message.Text, "/t") || strings.HasPrefix(message.Text, "!t")) {
		println("Handling a new translation request...")
		go HandleTranslateRequest(message)

		// Function 2: Help
	} else if strings.HasPrefix(message.Text, "/help") || strings.HasPrefix(message.Text, "!help") {
		println("Sending a help message...")
		go WriteTextMessage(helpMsg, message.Info.RemoteJid)

		// Function 3: Settings
	} else {
		// Always run this code on unknown messages. Spares time evaluating a msg
		go HandleSettingsRequest(message)
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
