package wabot

import (
	"fmt"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

type messageHandler struct{}

// Mainly caused by another instance of Whatsapp Web being opened
func (messageHandler) HandleError(err error) {
	if strings.Contains(err.Error(), "closed") {
		// Reconnect after a given amount of time
		println("Another instance of Whatsapp Web has been opened. Waiting to try again...")
		time.Sleep(errorTimeout)
		session, conn = handleLogin()
	}
}

func (messageHandler) HandleTextMessage(message whatsapp.TextMessage) {

	if message.Info.Timestamp > startTime && JidToName(message.Info.RemoteJid) == conn.Info.Pushname {
		go handleBotMsg(message)

		if showTextMessages {
			println(fmt.Sprintf("%s: %s", JidToName(MessageToJid(message)), message.Text))
		}
	}
}

func (messageHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	//fmt.Println(message)
}

func (messageHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	//fmt.Println(message)
}

func (messageHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	//fmt.Println(message)
}

func (messageHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
	//fmt.Println(message)
}

func (messageHandler) HandleJSONMessage(message string) {
	//	fmt.Println(message)
}

func (messageHandler) HandleContactMessage(message whatsapp.ContactMessage) {
	//fmt.Println(message)
}
