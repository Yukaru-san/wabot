package main

import (
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/skip2/go-qrcode"
)

func main() {
	// Requests token with a 20s timeout
	wac, err := whatsapp.NewConn(20 * time.Second)

	// QR-Code Stream
	qrChan := make(chan string)
	// Write QR-Code to a file to scan it
	go func() {
		qrcode.WriteFile(<-qrChan, qrcode.Medium, 256, "qr.png")
		//show qr code or save it somewhere to scan
	}()
	// Log into your session
	sess, err := wac.Login(qrChan)
	// Error check
	if err != nil {
		println(err.Error())
	}

	// Use wac TODO

	_ = sess
}
