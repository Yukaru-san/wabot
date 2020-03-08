package wabot

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yukaru-san/go-whatsapp"
)

type cmd struct{}

/*
	TODO

	Access to User Handle
	Allow Adding more stuff

*/

var (
	botname        = ""
	consoleWriteTo = ""

	contacList []whatsapp.Contact

	session whatsapp.Session
	conn    *whatsapp.Conn

	startTime        = uint64(time.Now().Unix())
	errorTimeout     = time.Minute * 1
	enableAutosaving = true
	autosaveInterval = time.Minute * 3

	encrypKey        = []byte("r4gyXrWSPXzvpBZJ")
	showTextMessages = true
	useContactName   = false

	useNicknames       = false
	nicknameUpdateText = "Your nickname has been updated!"

	qrCodeFile      = "qr.png"
	sessionFile     = "storedSession.dat"
	usersFile       = "storedUsers.dat"
	whatsappTimeout = 20 * time.Second
)

// SetEncryptionKey replaces the standard encryption Key
//  - Key has to be 16 Byte long
//  - Returns false if the key is bad
func SetEncryptionKey(key []byte) bool {
	if len(key) < 16 || len(key) > 16 {
		return false
	}

	encrypKey = key
	return true
}

// SetSessionFilePath changes the name a file should be saved in ([folder/]filename)
func SetSessionFilePath(path string) {
	sessionFile = path
}

// SetUsersFilePath changes the location the users will be saved in ([folder/]filename)
func SetUsersFilePath(path string) {
	usersFile = path
}

// SetQRFilePath changes the location the users will be saved in ([folder/]filename)
func SetQRFilePath(path string) {
	qrCodeFile = path
}

// UseContactNames tells the bot to use names saved in contacts (or not)
func UseContactNames(use bool) {
	useContactName = use
}

// DeactivateAutoSaving disables automatic saving
func DeactivateAutoSaving() {
	enableAutosaving = false
}

// SetAutosaveInterval - interval of userdata saving
func SetAutosaveInterval(interval time.Duration) {
	autosaveInterval = interval
}

// SetErrorTimeout sets the default time to reconnect after
// an error caused the program to disconnect
func SetErrorTimeout(timeout time.Duration) {
	errorTimeout = timeout
}

// AddTextCommand add a comand the program will listen to
// Inputs will be determined by their prefix
// - Always requires a function to take whatsapp.TextMessage as parameter
func AddTextCommand(cmd string, functionToExecute func(whatsapp.TextMessage)) { // TODO implement another checking method
	commands = append(commands, Command{prefix: cmd, function: functionToExecute})
}

// SetImageHandler calls the given function when receiving an img
func SetImageHandler(functionToExecute func(whatsapp.ImageMessage)) {
	imageHandleFunction = functionToExecute
}

// SetStickerHandler calls the given function when receiving an img
func SetStickerHandler(functionToExecute func(whatsapp.StickerMessage)) {
	stickerHandleFunction = functionToExecute
}

// SetDefaultTextHandleFunction calls the given function when no commands have fit
func SetDefaultTextHandleFunction(functionToExecute func(whatsapp.TextMessage)) {
	defaultTextHandleFunction = functionToExecute
}

// DisplayTextMessagesInConsole toggles visibility in console
func DisplayTextMessagesInConsole(display bool) {
	showTextMessages = display
}

// StartBot initiates and starts the bot
// - Takes in the name of the group to be run in
func StartBot(roomName string) (whatsapp.Session, *whatsapp.Conn) {
	// Self identification
	botname = roomName

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

// HandleError prints potential errors
func (cmd) HandleError(err error) {
	fmt.Println(err.Error())
}

// HandleContactList fills the contacList on load
func (cmd) HandleContactList(contacts []whatsapp.Contact) {
	for _, c := range contacts {
		contacList = append(contacts, c)
	}
}

// MessageToJid find the right Jid in a message
func MessageToJid(message whatsapp.TextMessage) string {

	authorID := ""

	if message.Info.Source.Participant != nil {
		authorID = *message.Info.Source.Participant
	} else {
		authorID = message.Info.RemoteJid // Personennamen
	}

	return authorID
}

// MessageToName Converts and entire message to the corresponding name
// -> allows finding the name in groups
func MessageToName(message whatsapp.TextMessage) string {

	// Try to find the right name
	participantNumber := strings.Split(MessageToJid(message), "@")[0]
	userNickname := GetUserNickname(MessageToJid(message))

	// Return if desired
	println(userNickname)
	if useNicknames && len(userNickname) > 0 {
		return userNickname
	} else if !useContactName && message.Info.Source.Participant != nil {
		return participantNumber
	}

	authorName := JidToName(message.Info.RemoteJid)

	// Return the custom name if there is none in the contacts
	if len(authorName) < 2 {
		return participantNumber
	}

	return authorName
}

// JidToName returns a corresponding contact's name. In general, you want to use MessageToName() instead
func JidToName(jid string) string {
	for _, c := range contacList {
		if c.Jid == jid {
			if c.Name == botname {
				return conn.Info.Pushname
			}

			return c.Name
		}
	}
	return "{undefined}"
}

// NameToJid finds the corresponding Jid to a name
// Returns {undefined} on error
func NameToJid(name string) string {
	for _, c := range contacList {
		if JidToName(c.Jid) == name {
			return c.Jid
		}
	}

	return "{undefined}"

}

// SetNicknameUseage changes the users abbility to use nicknames (false by default) - string will be the output. "" for no output
func SetNicknameUseage(allowNicknames bool, output string) {
	useNicknames = allowNicknames
	nicknameUpdateText = output
}
