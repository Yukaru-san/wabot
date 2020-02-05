package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

type messageHandler struct{}

func (messageHandler) HandleError(err error) {
	//fmt.Fprintf(os.Stderr, "%v", err)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (messageHandler) HandleTextMessage(message whatsapp.TextMessage) {
	//	Example reaction
	if !message.Info.FromMe && message.Info.RemoteJid == "__JID__" && message.Info.Timestamp > startTime {

		s := trans(message.Text)
		msg := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: message.Info.RemoteJid,
			},
			Text: s,
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
