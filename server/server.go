package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/CharlesMasson/basic-chat/chatutil"
)

func main() {
	args := os.Args

	var sendChannels [](chan chatutil.Message)

	l, err := net.Listen("tcp", ":"+args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		receiveChannel := make(chan chatutil.Message, 256)
		sendChannel := make(chan chatutil.Message, 256)

		sendChannels = append(sendChannels, sendChannel)

		// Forwards received messages to channel
		go chatutil.ReceiveMessages(conn, receiveChannel)
		go chatutil.SendMessages(conn, sendChannel)

		// Handles received messages
		go func(channel chan chatutil.Message) {
			for message := range channel {
				fmt.Printf("%v sent \"%v\"\n", message.Sender, message.Text)
				for _, channel := range sendChannels {
					if channel != sendChannel {
						channel <- message
					}
				}
			}
		}(receiveChannel)
	}
}
