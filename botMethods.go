package main

import (
	"strings"

	"github.com/Rhymen/go-whatsapp"
	"github.com/bregydoc/gtranslate"
	"golang.org/x/text/language"
)

var (
	helpMsg = "Possible commands:\n" +
		"/t\n   -   translate using default langs\n" +
		"/t-{from}-{to}\n   -   translate using the given langs\n" +
		"      e.g:   /t-de-es\n" +
		"/help\n   -   Hi.\n" +
		"/tf {language}\n   -   set your default {from} lang\n" +
		"/tt {language}\n   -   set your default {to} lang\n" +
		"/printSettings\n   -   prints your current settings\n" +
		"--------------------------------------------------------\n" +
		"commands can be initiated with\n \"!\" or \"/\"\n" +
		"--------------------------------------------------------"
)

// HandleTranslateRequest handles a translation request
func HandleTranslateRequest(message whatsapp.TextMessage) {

	langFrom := GetUserSettings(message.Info.RemoteJid).StdTranslationFrom
	langTo := GetUserSettings(message.Info.RemoteJid).StdTranslationTo

	commandPretext := 2

	// If a translation direction was given use it instead
	if strings.HasPrefix(message.Text, "/t-") || strings.HasPrefix(message.Text, "!t-") {
		// Splitting
		wantedLangFromString := message.Text[3:5]
		wantedLangToString := message.Text[6:8]

		// Try to parse it
		wantedLangFrom, err := language.Parse(wantedLangFromString)
		wantedLangTo, err := language.Parse(wantedLangToString)

		// Success!
		if err == nil {
			langFrom = wantedLangFrom
			langTo = wantedLangTo
			commandPretext = 9 // Start at the end of the command
		}
	}

	var txt string

	// If the message quotes another one, take that one
	if message.ContextInfo.QuotedMessage != nil {
		txt = message.ContextInfo.QuotedMessage.GetConversation()
	} else if len(message.Text) > 2 { // Else, use the given one
		txt = message.Text[commandPretext:] // Ignore the command-part
	}

	// Translate the message
	if len(txt) > 0 {
		tl, _ := gtranslate.Translate(txt, langFrom, langTo)
		WriteTextMessage(tl, message.Info.RemoteJid)
	}
}

// HandleSettingsRequest handels messages concerning a users settings.
// TODO add more if needed
func HandleSettingsRequest(message whatsapp.TextMessage) {

	// Translate From xx
	if strings.HasPrefix(message.Text, "/tf") {
		newTranslation, err := language.Parse(message.Text[4:])

		// Input error
		if err != nil {
			WriteTextMessage("Sorry!\n Couldn't parse your Input.\nExample: /tf de -> translate from german", message.Info.RemoteJid)
			return
		}

		// Search User and change the entry
		AddUserByJid(MessageToJid(message))
		users.BotUsers[GetUserIndex(message)].Settings.StdTranslationTo = newTranslation

		// Send a reply
		WriteTextMessage("Changed @"+MessageToName(message)+"'s {From} to "+newTranslation.String(), message.Info.RemoteJid)
		return
	}

	// Translate to xx
	if strings.HasPrefix(message.Text, "/tt") {
		newTranslation, err := language.Parse(message.Text[4:])

		// Input error
		if err != nil {
			WriteTextMessage("Sorry!\n Couldn't parse your Input.\nExample: /tt eng -> translate to english", message.Info.RemoteJid)
			return
		}

		// Search User and change the entry
		AddUserByJid(MessageToJid(message))
		users.BotUsers[GetUserIndex(message)].Settings.StdTranslationTo = newTranslation

		// Send a reply
		WriteTextMessage("Changed @"+MessageToName(message)+"'s {To} to "+newTranslation.String(), message.Info.RemoteJid)
		return
	}

	// Print Settings
	if strings.HasPrefix(strings.ToLower(message.Text), "/printsettings") { // TODO Update if needed

		settings := GetUserSettings(message.Info.RemoteJid)

		infos := "@" + MessageToName(message) + "'s settings:"
		infos += "\nStdTranslationFrom: " + settings.StdTranslationFrom.String()
		infos += "\nStdTranslationTo: " + settings.StdTranslationTo.String()

		WriteTextMessage(infos, message.Info.RemoteJid)
	}

	// TODO Add more here
}
