package wabot

import (
	"fmt"
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
	authorName := JidToName(message.Info.RemoteJid)

	//	fmt.Printf("---------\nNumber: %s\nNickname: %s\nAuthor: %s\n", participantNumber, userNickname, authorName)

	// Return if desired
	if useNicknames && len(userNickname) > 0 {
		return userNickname
	} else if !useContactName && message.Info.Source.Participant != nil {
		return participantNumber
	}

	// Return the custom name if there is none in the contacts
	if len(authorName) < 2 {
		return participantNumber
	}

	return authorName
}

// MessageToGroupID returns the jid of a group. If the message is from outside a group it returns ""
func MessageToGroupID(message whatsapp.TextMessage) string {

	jid := MessageToJid(message)

	groupID := strings.Split(jid, "-")[1]

	fmt.Println("MessageToName: ", MessageToName(message))

	if jid != groupID {
		return groupID
	}

	// Message is from outside a group
	return ""
}

// GetPhoneNumber returns a bots phone number
func GetPhoneNumber() string {

	// Try to find the right name
	return strings.Split(conn.Info.Wid, "@")[0]
}

// JidToGroupName returns a group's name according to the jid
func JidToGroupName(jid string) string {

	// If only a group's jid was given
	if !strings.Contains(jid, "-") {
		jid = fmt.Sprintf("%s-%s", GetPhoneNumber(), jid)
	}

	// Return the right name
	return JidToName(jid)
}

// ----------------------------------------------------------------------------------------------------- //

// JidToName returns a corresponding contact or group's name.
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
