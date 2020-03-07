package wabot

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Yukaru-san/go-whatsapp"
)

type messageHandler struct{}

// Mainly caused by another instance of Whatsapp Web being opened
func (messageHandler) HandleError(err error) {
	if strings.Contains(err.Error(), "closed") {

		if errorTimeout == -1 {
			println("Exit due to connection break")
			os.Exit(0)
		}

		// Reconnect after a given amount of time
		println("Another instance of Whatsapp Web has been opened. Waiting to try again...")
		time.Sleep(errorTimeout)
		conn = nil
		session, conn = handleLogin()
		println("Connected again!")
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
	if message.Info.Timestamp > startTime && JidToName(message.Info.RemoteJid) == conn.Info.Pushname {
		go imageHandleFunction(message)
	}
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

func (messageHandler) HandleStickerMessage(message whatsapp.StickerMessage) {
	if message.Info.Timestamp > startTime && JidToName(message.Info.RemoteJid) == conn.Info.Pushname && strings.Contains(message.Info.Source.Message.String(), "url:") {
		go stickerHandleFunction(message)
	}
}
