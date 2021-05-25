package wabot

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

type messageHandler struct{}

/*

	TODO Use custom Handlers from custom Bots

*/

// Mainly caused by another instance of Whatsapp Web being opened
func (messageHandler) HandleError(err error) {

	if !strings.Contains(err.Error(), "error processing data") {

		if errorTimeout == -1 {
			fmt.Printf("Exit due to connection break. Error:%s\n", err.Error())
			os.Exit(0)
		}

		// Reconnect after a given amount of time
		fmt.Printf("WhatsApp crashed with the following Message:\n%s\nTrying again in: %d\n", err.Error(), errorTimeout)
		time.Sleep(errorTimeout)
		conn.Disconnect()
		err = nil
		session, conn, err = handleLogin()
		if err != nil {
			fmt.Println("--- Connected again! ---")
		} else {
			fmt.Println("Reconnect error: ", err.Error())
			os.Exit(1)
		}
	}
}

func (messageHandler) HandleTextMessage(message whatsapp.TextMessage) {

	if message.Info.Timestamp > startTime {
		go handleBotMsg(message)

		if showTextMessages {
			fmt.Printf("%s: %s\n", JidToGroupName(message.Info.RemoteJid), message.Text)
		}
	}
}

func (messageHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	if message.Info.Timestamp > startTime {
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
	if message.Info.Timestamp > startTime && strings.Contains(message.Info.Source.Message.String(), "url:") {
		go stickerHandleFunction(message)
	}
}
