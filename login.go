package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/skip2/go-qrcode"
)

// HandleLogin returns a connection and session to the
//   web. If there is an error, the program will exit
func HandleLogin() (whatsapp.Session, *whatsapp.Conn) {

	// Try to load an old session
	savedSession := whatsapp.Session{}
	savedData, err := ioutil.ReadFile("storedSession.sess")
	err = json.Unmarshal(savedData, &savedSession)

	// If there is no session stored
	if err != nil {
		// Requests token with a 20s timeout
		wac, err := whatsapp.NewConn(20 * time.Second)

		// QR-Code Stream
		qrChan := make(chan string)
		// Write QR-Code to a file to scan it
		go func() {
			println("No stored session found. Please login using the generated QR code!")
			qrcode.WriteFile(<-qrChan, qrcode.Medium, 256, "qr.png")
			//show qr code or save it somewhere to scan
		}()
		// Log into your session
		sess, err := wac.Login(qrChan)
		// Close programm on error or timeout
		if err != nil {
			println("Timeout! Exiting...")
			os.Exit(0)
		}
		// Save new session to quickly start the next time
		sessionJSON, _ := json.Marshal(sess)
		err = ioutil.WriteFile("storedSession.sess", sessionJSON, 0600)
		println("Session saved.")

		return sess, wac
	}

	// Session loaded successfully. Use it to login
	wac, err := whatsapp.NewConn(20 * time.Second)
	sess, err := wac.RestoreWithSession(savedSession)

	// Error exit
	if err != nil {
		println("Error! Exiting...")
		os.Exit(0)
	}

	println("Successfully logged in!")
	return sess, wac

}
