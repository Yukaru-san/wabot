package wabot

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Yukaru-san/go-whatsapp"
)

// Command struct
type Command struct {
	prefix   string
	groups   []string
	function func(whatsapp.TextMessage)
}

// List of implemented commands
var commands []Command
var imageHandleFunction func(whatsapp.ImageMessage)
var stickerHandleFunction func(whatsapp.StickerMessage)
var defaultTextHandleFunction func(whatsapp.TextMessage)

// handleBotMsg checks if a message is a command and executes the first possible command
func handleBotMsg(message whatsapp.TextMessage) {

	fmt.Printf("Received: %s\n - from Group:%s\n", message.Text, fmt.Sprintf("%s@%s@g.us", GetPhoneNumber(), MessageToGroupID(message)))

	// Run through command list and execute if possible
	for i := 0; i < len(commands); i++ {
		// if no groups set || group is set && command is right
		if (len(commands[i].groups) == 0 || arrayContains(commands[i].groups, fmt.Sprintf("%s@%s@g.us", GetPhoneNumber(), MessageToGroupID(message)))) && strings.HasPrefix(strings.Split(strings.ToLower(message.Text), " ")[0], strings.ToLower(commands[i].prefix)) {
			go commands[i].function(message)
			break
		}
	}

	// Predefined function that can be turned off
	if useNicknames && strings.HasPrefix(strings.ToLower(message.Text), "/nick") && len(message.Text) > 7 {
		SetUserNickname(MessageToJid(message), message.Text[6:])
		if len(nicknameUpdateText) > 0 {
			WriteTextMessage(nicknameUpdateText, message.Info.RemoteJid)
		}
	}

	// No command found? Try to run the default code
	go defaultTextHandleFunction(message)
}

func arrayContains(array []string, group string) bool {
	for _, a := range array {
		if a == group {
			return true
		}
	}

	return false
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
