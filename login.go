package wabot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/skip2/go-qrcode"
)

// StartBot initiates and starts the bot
// - Takes in the name of the group to be run in
// func StartBot(roomName string) (whatsapp.Session, *whatsapp.Conn) { TODO Only speficied groups
func StartBot(longConnectionName, shortConnectionName string) (whatsapp.Session, *whatsapp.Conn, error) {
	// Self identification
	longConnName = longConnectionName
	shortConnName = shortConnectionName

	// Create empty functions to prevent crashing on img / sticker
	SetImageHandler(func(whatsapp.ImageMessage) {})
	SetStickerHandler(func(whatsapp.StickerMessage) {})
	SetDefaultTextHandleFunction(func(whatsapp.TextMessage) {})

	// Login
	var err error
	session, conn, err = handleLogin()

	if err != nil {
		return whatsapp.Session{}, nil, err
	}

	// Add the command handler
	conn.AddHandler(cmd{})
	// Time for the handler to get ready
	time.Sleep(time.Second)

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

	return session, conn, nil
}

// handleLogin returns a connection and session-pointer.
// If there is an error, the program will exit
func handleLogin() (whatsapp.Session, *whatsapp.Conn, error) {
	// Try to load a stored session for quick connection
	savedSession := whatsapp.Session{}
	savedData, err := ioutil.ReadFile(sessionFile)
	if err == nil {
		savedData = decryptData(savedData)
		err = json.Unmarshal(savedData, &savedSession)
	}

	// Whatsapp Session Options
	connOptions := whatsapp.Options{
		Timeout:         whatsappTimeout,
		LongClientName:  longConnName,
		ShortClientName: shortConnName,
	}

	// If there is no session stored
	if err != nil {
		// Requests token with a 20s timeout
		wac, err := whatsapp.NewConnWithOptions(&connOptions)
		wac.SetClientVersion(2, 2123, 7)
		if err != nil {
			return whatsapp.Session{}, nil, err
		}

		qrChan := make(chan string)
		scanChan := make(chan bool, 1)
		go func() {
			fmt.Println("No stored session found. Please login using the generated QR code!")
			targetDir, _ := filepath.Split(qrCodeFile)
			os.MkdirAll(targetDir, 0700)
			if err := qrcode.WriteFile(<-qrChan, qrcode.Medium, 256, qrCodeFile); err != nil {
				fmt.Println("Error saving qr code!", err.Error())
			} else {
				//Try to open the image. Makes it easier to scan
				displayQRcode(scanChan, qrCodeFile)
			}
		}()

		// Log into your session
		sess, err := wac.Login(qrChan)
		if err != nil {
			return whatsapp.Session{}, nil, err
		}
		// Save new session to quickly start the next time
		sessionJSON, _ := json.Marshal(sess)
		sessionJSON = encryptData(sessionJSON)
		ioutil.WriteFile(sessionFile, sessionJSON, 0600)
		fmt.Println("Session saved. No QR-Code needed during the next login!")
		scanChan <- true
		return sess, wac, nil
	}

	// Session loaded successfully. Use it to login
	wac, err := whatsapp.NewConnWithOptions(&connOptions)
	wac.SetClientVersion(2, 2123, 7)
	sess, err := wac.RestoreWithSession(savedSession)

	if err != nil {
		return whatsapp.Session{}, nil, err
	}

	return sess, wac, nil
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
