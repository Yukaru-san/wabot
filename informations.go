package wabot

import (
	"strings"

	"github.com/Yukaru-san/go-whatsapp"
)

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

// MessageToGroupID returns the jid of a group. If the message is from outside a group it returns ""
func MessageToGroupID(message whatsapp.TextMessage) string {

	if message.Info.Source.Participant != nil {
		return strings.Split(MessageToJid(message), "@")[1]
	}

	// Message is from outside a group
	return ""
}

// GetPhoneNumber returns a bots phone number
func GetPhoneNumber() string {

	// Try to find the right name
	return strings.Split(conn.Info.Wid, "@")[0]
}

// ----------------------------------------------------------------------------------------------------- //

// JidToName returns a corresponding contact's name. In general, you want to use MessageToName() instead
func JidToName(jid string) string {
	for _, c := range contacList {
		if c.Jid == jid {
			return c.Name
		}
	}
	return "{undefined}"
}

// NameToJid returns a corresponding contact's jid. n general, you want to use MessageToJid() instead
func NameToJid(name string) string {
	for _, c := range contacList {
		if JidToName(c.Jid) == name {
			return c.Jid
		}
	}

	return "{undefined}"
}
