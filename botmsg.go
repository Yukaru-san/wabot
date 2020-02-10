package wabot

import (
	"os"
	"strings"

	"github.com/Yukaru-san/go-whatsapp"
)

// Command struct
type Command struct {
	prefix   string
	function func(whatsapp.TextMessage)
}

// List of implemented commands
var commands []Command
var imageHandleFunction func(whatsapp.ImageMessage)
var stickerHandleFunction func(whatsapp.StickerMessage)
var defaultTextHandleFunction func(whatsapp.TextMessage)

// handleBotMsg checks if a message is a command and executes the first possible command
func handleBotMsg(message whatsapp.TextMessage) {
	// Run through command list and execute if possible
	for i := 0; i < len(commands); i++ {
		if strings.HasPrefix(strings.ToLower(message.Text), strings.ToLower(commands[i].prefix)) {
			go commands[i].function(message)
			break
		}
	}

	// No command found? Try to run the default code
	go defaultTextHandleFunction(message)
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

// SendImageMessage takes img type like "image/png"
func SendImageMessage(img *os.File, imgType string, remoteJid string) {
	// Create the struct
	msg := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: remoteJid,
		},
		Type:    "image/png",
		Content: img,
	}

	// And send it
	conn.Send(msg)
}

// SendStickerMessage only uses .webp files
func SendStickerMessage(img *os.File, remoteJid string) {
	// Create the struct
	msg := whatsapp.StickerMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: remoteJid,
		},

		Type:    "image/webp",
		Content: img,
	}
	// And send it
	conn.Send(msg)
}
