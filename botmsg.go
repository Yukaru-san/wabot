package wabot

import (
	"bytes"
	"io"
	"io/ioutil"
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

// SendImageMessage takes img type like "image/png"  / If caption len is 0 there won't be a text
func SendImageMessage(caption string, img io.Reader, imgType string, remoteJid string) {
	var msg whatsapp.ImageMessage

	// Create a Thumbnail
	var buffer bytes.Buffer
	tee := io.TeeReader(img, &buffer)

	thumbnail, _ := ioutil.ReadAll(&buffer)

	if len(caption) > 0 {
		// Create the struct
		msg = whatsapp.ImageMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: remoteJid,
			},
			Caption:   caption,
			Type:      imgType,
			Content:   tee,
			Thumbnail: thumbnail,
		}
	} else {
		// Create the struct
		msg = whatsapp.ImageMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: remoteJid,
			},
			Type:      imgType,
			Content:   tee,
			Thumbnail: thumbnail,
		}
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
