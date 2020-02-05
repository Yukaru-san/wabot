package main

import (
	"fmt"

	"github.com/Rhymen/go-whatsapp"
)

type messageHandler struct{}

func (messageHandler) HandleError(err error) {
	//fmt.Fprintf(os.Stderr, "%v", err)
}

func (messageHandler) HandleTextMessage(message whatsapp.TextMessage) {

	/*	Example reaction
			if message.Info.FromMe && jidToName(message.Info.RemoteJid) == "____NAME____" && message.Info.Timestamp > startTime {

			// Erstelle ne Nachricht
			msg := whatsapp.TextMessage{
				Info: whatsapp.MessageInfo{
					RemoteJid: message.Info.RemoteJid,
				},
				Text: message.Text,
			}
			// Und schick sie ab
			_, err := conn.Send(msg)

			if err != nil {
				println(err.Error())
			}

		}
	*/

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
