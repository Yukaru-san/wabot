package main

import (
	"time"
)

func main() {

	sess, conn := HandleLogin()
	_ = sess

	// Add a complete MSG Handler
	conn.AddHandler(messageHandler{})

	for {
		time.Sleep(time.Second)
	}
}
