# WhatsApp-GroupBot

By emulating WhatsApp Web this library allows you to easiely set up and use this highly customizable WhatsApp bot.  
Since customizability is the main focus of this project, you can create as many interactions as you want.
There is also a feature to save specific informations (like settings) across users within the group.

*There is just one small problem*                                                                                      
Since there can be only one instance of WhatsApp Web at a time, the bot will pause when you open another instance               (e.g. in your browser tab)        
The bot will reconnect after a given amount of time (defined in *SetErrorTimeout()*) and therefore kick you out of your
other instance. Obviously, commands during this time-out won't have any effect and sadly, the only way around it, is completly denying another instance of WhatsApp Web being opened.                                                                       
If you wish to do so, just set the ErrorTimeout to something like 100
<br><br>
I also published [one of my private bots](https://github.com/Yukaru-san/wabot-example) so you can get some inspiration.

# Installation
Simply use:
```go
go get github.com/Yukaru-san/WhatsApp-GroupBot
```

# Usage
Using this bot should be pretty easy as soon as you know what you want to accomplish. Here is an example of a working program that will count the amount of times the command "/cmd" is used:
```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	wabot "github.com/Yukaru-san/WhatsApp-GroupBot"
	"github.com/Rhymen/go-whatsapp"
)

// Settings are saved personally for each user
type Settings struct {
	CommandCounter int
}

func main() {

	// Initialize Bot
	wabot.StartBot("My Whatsapp Bot", "bot")

	// Set a custom auto-save interval
	wabot.SetAutosaveInterval(time.Minute)

	// There is a default key, but you can change it
	wabot.SetEncryptionKey([]byte("VmDTJUQrYxTztMDB"))

	// Using default Settings can prevent errors
	wabot.CreateNewSettingsOption(Settings{
		0, // CommandCounter
	})

	// Try to load savedata and convert it into your own structure
	savedata, found := wabot.GetSaveData()
	if found {
		for i := 0; i < len(savedata.BotUsers); i++ {
			tmp := &Settings{}

			data, _ := json.Marshal(savedata.BotUsers[i].Settings)
			json.Unmarshal(data, tmp)

			savedata.BotUsers[i].Settings = *tmp
		}
		wabot.UseSaveData(savedata)
	}

	// Add and handle commands however you like!
	wabot.AddTextCommand("/cmd", func(message whatsapp.TextMessage) {
		// Get the current settings from the user who sent the message
		settings := wabot.GetUserSettings(wabot.MessageToJid(message)).(Settings)

		fmt.Println(settings)

		// Count up!
		settings.CommandCounter++

		// Update the database
		wabot.ChangeUserSettings(wabot.MessageToJid(message), settings)

		// Output
		fmt.Printf("\nUser @%s has used this command %d times!\n", wabot.MessageToName(message), settings.CommandCounter)
	})

	// Don't let the program die here!
	for {
		time.Sleep(time.Minute)
	}
}

```
Since the savable data is fully modular, at least the savedata-loading couldn't be included within the library. Therefore it has to be done inside of your own code. As long as you keep the name "Settings" you can use this exact loading part for more-or-less any other project. Also, for any aditional information regarding the struct "whatsapp.TextMessage" refer to [Rhymen's Github page](https://github.com/Rhymen/go-whatsapp) - it is quite self-explanatory though. A message's text is stored in ***message.Text***

# Functions

```go
SetSessionFilePath(string)		   // Directory to save session files in
SetUsersFilePath(string)		   // Directory to save users in
SetQRFilePath(string)			   // Directory to save qr codes in
SetErrorTimeout(time)		   	   // Timeout to reconnect on error
SetNicknameUseage(bool, msgText)	   // Allows users to choose nicknames

StartBot(string, string)                   // Initializes and starts the bot
DisplayTextMessagesInConsole(bool)         // Toggle printing new messages on / off

AddTextCommand(string, func())             // If a message starts with the given string it executes the func
AddGroupCommand(string, string, func())    // Like above, but only works in given group
SetDefaultTextHandleFunction(func())	   // Calls the given function if no command could be parsed
SetImageHandler(func())                    // Executes the given function when receiving an image  
SetStickerHandler(func())                  // Executes the given function when receiving a sticker  

WriteTextMessage(string, string)           // Writes a WhatsApp messsage to the defined Jid's owner
SendImageMessage(*os.File, string, string) // Sends the given image. Important to check imgType (e.g. "img/png")
SendStickerMessage(*os.File, string)       // Sends the given sticker. Needs to be a .webp file

SetEncryptionKey([]byte)                   // Set a 16 Byte key to encrypt your saved data

MessageToJid(whatsapp.TextMessage)         // Returns the unique ID of the message's sender
MessageToName(whatsapp.TextMessage)        // Returns the sender's name
JidToName(string)                          // Returns the name of the ID's owner
NameToJid(string)                          // Returns the corresponding Jid of the user's name

DeactivateAutoSaving()                     // Disables automatic saving
SetAutoSaveInterval(time.Duration)         // Save after every passed interval
SetErrorTimeout(time.Duration)             // Time until the program will restart after losing connection

CreateNewSettingsOption(interface)         // Standard Settings for newly added users
IsUserRegistered(string) bool              // Checks if the Jid is already registered 
AddUserByJid(string)                       // Registers the Jid's owner - if not already registered
ChangeUserSettings(string, interface)      // Updates the user's settings with the given struct
GetUserSettings(string) interface          // Returns the Jid-corresponding Settings
GetUserIndex(whatsapp.TextMessage) int     // Returns the message-owner's position in the user-slice

SaveUsersToDisk() bool                     // Returns true on successfully saving the userlist
GetSaveData() (BotUserList, bool)          // Returns the loaded savedata and true if possible (false otherwise)
UseSaveData(BotUserList)                   // Overwrites the userlist with the given one
```
*ChangeUserSettings* and *GetUserSettings* both call *AddUserByJid* if needed. You don't need to check yourself!
