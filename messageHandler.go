package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/bregydoc/gtranslate"
	"golang.org/x/text/language"
)

type messageHandler struct{}

func (messageHandler) HandleError(err error) {
	//fmt.Fprintf(os.Stderr, "%v", err)
}

func (messageHandler) HandleTextMessage(message whatsapp.TextMessage) {
	//	Example reaction
	if message.Info.Timestamp > startTime && strings.HasPrefix(message.Text, "!t") {

		tl, _ := gtranslate.Translate(message.Text[2:], language.German, language.English)
		msg := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: message.Info.RemoteJid,
			},
			Text: tl,
		}
		// Und schick sie ab
		time.Sleep((time.Duration)(rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(5)) * time.Second)
		_, err := conn.Send(msg)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	//fmt.Println(message)
	fmt.Printf("%s:\t%s\n", jidToName(message.Info.RemoteJid), message.Text)
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
