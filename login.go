package wabot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/Yukaru-san/go-whatsapp"
	"github.com/skip2/go-qrcode"
)

// StartBot initiates and starts the bot
// - Takes in the name of the group to be run in
// func StartBot(roomName string) (whatsapp.Session, *whatsapp.Conn) { TODO Only speficied groups
func StartBot(longConnectionName, shortConnectionName string) (whatsapp.Session, *whatsapp.Conn) {
	// Self identification
	longConnName = longConnectionName
	shortConnName = shortConnectionName

	// Create empty functions to prevent crashing on img / sticker
	SetImageHandler(func(whatsapp.ImageMessage) {})
	SetStickerHandler(func(whatsapp.StickerMessage) {})
	SetDefaultTextHandleFunction(func(whatsapp.TextMessage) {})

	// Login
	session, conn = handleLogin()

	// Add the command handler
	conn.AddHandler(cmd{})

	// Saves in intervals
	go (func() {
		for enableAutosaving {
			time.Sleep(autosaveInterval)
			SaveUsersToDisk()
		}
	})()
	//
	// Add the message handler
	go (func() {
		conn.AddHandler(messageHandler{})
	})()

	return session, conn
}

// handleLogin returns a connection and session-pointer.
// If there is an error, the program will exit
func handleLogin() (whatsapp.Session, *whatsapp.Conn) {
	// Try to load a stored session for quick connection
	savedSession := whatsapp.Session{}
	savedData, err := ioutil.ReadFile(sessionFile)
	if err == nil {
		savedData = decryptData(savedData)
		err = json.Unmarshal(savedData, &savedSession)
	}

	// If there is no session stored
	if err != nil {
		// Requests token with a 20s timeout
		wac, err := whatsapp.NewConn(whatsappTimeout, longConnName, shortConnName)
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
		sessionJSON = encryptData(sessionJSON)
		ioutil.WriteFile(sessionFile, sessionJSON, 0600)
		fmt.Println("Session saved. No QR-Code needed during the next login!")
		scanChan <- true
		return sess, wac
	}

	// Session loaded successfully. Use it to login
	wac, err := whatsapp.NewConn(whatsappTimeout, longConnName, shortConnName)
	sess, err := wac.RestoreWithSession(savedSession)

	if err != nil {
		fmt.Printf("Error connecting to WhatsApp:\n%s\n", err.Error())
		os.Exit(0)
	}

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
			fmt.Println("To display the QR-code directly on linux, you need to install feh!\nPlease scan the generated QR code in data/qr.png")
		}
	}
}
