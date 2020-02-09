# WhatsApp-GroupBot
By emulating a WhatsApp session this bot allows you to easiely set up and use this highly customizable WhatsApp bot.  
Since customizability is the main focus of this project, you can create as many interactions as you want.
There is also a feature to save specific informations (like settings) across users within the group.

*There is just one small problem*                                                                                      
Since there can be only one instance of WhatsApp Web at a time, the bot will pause when you open another instance               (e.g. in your browser tab)        
The bot will reconnect after a given amount of time (defined in *SetErrorTimeout()*) and therefore kick you out of your
other instance. Obviously, commands during this time-out won't have any effect and sadly, the only way arount it, is completly denying another instance of WhatsApp Web being opened.                                                                       
If you wish to do so, just set the ErrorTimeout to something like 100

# Installation
The bot depends on Rhymen's WhatsApp API but since I forked his original work to add more features, you will need to install these two libraries:
```go
go get github.com/Yukaru-san/go-whatsapp
```
```go
go get github.com/Yukaru-san/WhatsApp-GroupBot
```

# Usage
To use this bot you mainly need these two functions
```go
wabot.StartBot()    // Initializes and starts the bot
wabot.AddCommand()  // Add unlimited amount of commands
```
but since this is a bit on the short side, here is an example of a working program:
```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Rhymen/go-whatsapp"
	wabot "github.com/Yukaru-san/WhatsApp-GroupBot"
)

// Settings are saved personally for each user
type Settings struct {
	CommandCounter int
}

func main() {

	// Initialize Bot
	wabot.StartBot("My Whatsapp Bot")

	// Set a custom auto-save interval
	wabot.SetAutosaveInterval(time.Minute)

	// There is a default key, but you can change it
	wabot.SetEncryptionKey([]byte("VmDTJUQrYxTztMDB"))

	// Setting a default Settings can prevent many errors
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
	wabot.AddCommand("/cmd", func(message whatsapp.TextMessage) {
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
Since the savable data is fully modular, at least the savedata-loading couldn't be included within the library. Therefore it has to be done inside of your own code. As long as you keep the name "Settings" you can use this exact loading part for more-or-less any other project. Also, for any aditional information regarding the struct "whatsapp.TextMessage" refer to This site was built using [Rhymen's Github page](https://github.com/Rhymen/go-whatsapp) - it is quite self-explanatory though. A message's text is stored in ***message.Text***

# Functions

```go
StartBot(string)                           // Initializes and starts the bot
DisplayTextMessagesInConsole(bool          // Toggle printing new messages on / off

AddCommand(string, func())                 // If a message starts with the given string it executes the func
SetImageHandler(func())                    // Executes the given function when receiving a image  
SetStickerHandler(func())                  // Executes the given function when receiving a sticker  

WriteTextMessage(string, string)           // Writes a WhatsApp messsage to the defined Jid's owner
SendImageMessage(*os.File, string, string) // Sends the given image. Important to check imgType (e.g. "img/png")
SendStickerMessage(*os.File, string)       // Sends the given sticker. Needs to be a .webp file

SetEncryptionKey([]byte)                   // Set a 16 Byte key to encrypt your saved data

MessageToJid(whatsapp.TextMessage)         // Returns the unique ID of the message sender
MessageToName(whatsapp.TextMessage)        // Returns the sender's name
JidToName(string)                          // Returns the name of the ID's owner
NameToJid(string)                          // Returns the corresponding Jid of the user's name

DeactivateAutoSaving()                     // Disables automatic saving
SetAutoSaveInterval(time.Duration)         // Save after every passed interval
SetErrorTimeout(time.Duration)             // Time until the program will restart after losing connection

CreateNewSettingsOption(interface)         // Standard Settings for newly added users
IsUserRegistered(string) bool              // Checks if the Jid is registered already
AddUserByJid(string)                       // Registers the jid's owner - if not already registered
ChangeUserSettings(string, interface)      // Updates the user's settings with the given struct
GetUserSettings(string) interface          // Returns the Jid-corresponding Settings
GetUserIndex(whatsapp.TextMessage) int     // Returns the message-owner's position in the user-slice

SaveUsersToDisk() bool                     // Returns true on successful saving
GetSaveData() (BotUserList, bool)          // Returns the loaded savedata and true if possible (false otherwise)
UseSaveData(BotUserList)                   // Overwrites the userlist with the given one
```
*ChangeUserSettings* and *GetUserSettings* both call *AddUserByJid* if needed. You don't need to check yourself!
