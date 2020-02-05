package main

import (
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/skip2/go-qrcode"
)

func main() {
	wac, err := whatsapp.NewConn(20 * time.Second)

	qrChan := make(chan string)
	go func() {
		qrcode.WriteFile(<-qrChan, qrcode.Medium, 256, "qr.png")
		//show qr code or save it somewhere to scan
	}()
	sess, err := wac.Login(qrChan)

	_ = err
	_ = sess
}
