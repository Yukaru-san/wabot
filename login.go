package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/skip2/go-qrcode"
)

const (
	qrCodeFile      = "qr.png"
	sessionFile     = "storedSession.json"
	whatsappTimeout = 20 * time.Second
)

// HandleLogin returns a connection and session-pointer.
// If there is an error, the program will exit
func HandleLogin() (whatsapp.Session, *whatsapp.Conn) {
	savedSession := whatsapp.Session{}

	savedData, err := ioutil.ReadFile(sessionFile)
	if err == nil {
		err = json.Unmarshal(savedData, &savedSession)
	}

	// If there is no session stored
	if err != nil {
		// Requests token with a 20s timeout
		wac, err := whatsapp.NewConn(whatsappTimeout)
		if err != nil {
			fmt.Println("An error occured:", err.Error())
			os.Exit(1)
		}

		qrChan := make(chan string)
		scanChan := make(chan bool, 1)
		go func() {
			fmt.Println("No stored session found. Please login using the generated QR code!")
			if err := qrcode.WriteFile(<-qrChan, qrcode.Medium, 256, qrCodeFile); err != nil {
				fmt.Println("Error saving qr code!", err.Error())
				os.Exit(1)
			} else {
				//Try to open the image. Makes it easier to scan
				displayQRcode(scanChan, qrCodeFile)
			}
		}()

		// Log into your session
		sess, err := wac.Login(qrChan)
		if err != nil {
			println("Timeout! Exiting...")
			os.Exit(0)
		}
		// Save new session to quickly start the next time
		sessionJSON, _ := json.Marshal(sess)
		ioutil.WriteFile(sessionFile, sessionJSON, 0600)
		fmt.Println("Session saved.")
		scanChan <- true
		return sess, wac
	}

	// Session loaded successfully. Use it to login
	wac, err := whatsapp.NewConn(whatsappTimeout)
	sess, err := wac.RestoreWithSession(savedSession)

	if err != nil {
		fmt.Println("Error! Exiting...")
		os.Exit(0)
	}

	fmt.Println("Successfully logged in!")
	return sess, wac
}

func displayQRcode(ch chan bool, image string) {
	if runtime.GOOS == "linux" {
		if er := exec.Command("feh", "-v").Run(); er == nil {
			cmd := exec.Command("feh", image)
			go cmd.Run()
			if ch != nil {
				go (func() {
					<-ch
					cmd.Process.Kill()
					os.Remove(image)
				})()
			}
		} else {
			fmt.Println("no nil:", er.Error())
		}
	}
}
