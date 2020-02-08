package wabot

import (
	"fmt"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

type cmd struct{}

/*
	TODO

	Access to User Handle
	Allow Adding more stuff

*/

var (
	botname        = "Yukaru-Bot"
	consoleWriteTo = ""

	contacList []whatsapp.Contact

	session whatsapp.Session
	conn    *whatsapp.Conn

	startTime        = uint64(time.Now().Unix())
	errorTimeout     = time.Minute * 1
	autosaveInterval = time.Minute * 3

	encrypKey        = []byte("r4gyXrWSPXzvpBZJ")
	showTextMessages = true

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

// SetAutosaveInterval - interval of userdata saving
func SetAutosaveInterval(interval time.Duration) {
	autosaveInterval = interval
}

// SetErrorTimeout sets the default time to reconnect after
// an error caused the program to disconnect
func SetErrorTimeout(timeout time.Duration) {
	errorTimeout = timeout
}

// AddCommand add a comand the program will listen to
// Inputs will be determined by their prefix
// - Always requires a function to take whatsapp.TextMessage as parameter
func AddCommand(cmd string, functionToExecute func(whatsapp.TextMessage)) { // TODO implement another checking method
	commands = append(commands, Command{prefix: cmd, function: functionToExecute})
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

	// Login
	session, conn = HandleLogin()

	// Add the command handler
	conn.AddHandler(cmd{})

	// Saves in intervals
	go (func() {
		for {
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

	authorID := ""

	if len(message.Info.Source.GetPushName()) > 0 {
		return message.Info.Source.GetPushName()
	} else if message.Info.Source.Participant != nil {
		authorID = *message.Info.Source.Participant
	} else {
		authorID = message.Info.RemoteJid // Personennamen
	}

	return jidToName(authorID)
}

// Converts an ID to the corresponding name
func jidToName(jid string) string {
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
	botname = "Yukaru-Bot"

	for _, c := range contacList {
		if jidToName(c.Jid) == name {
			return c.Jid
		}
	}

	return "{undefined}"

}
