package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/CharlesMasson/basic-chat/chatutil"
)

func main() {
	args := os.Args
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter your name:")
	scanner.Scan()
	name := scanner.Text()

	conn, _ := net.Dial("tcp", args[1])

	sendChannel := make(chan chatutil.Message, 256)
	receiveChannel := make(chan chatutil.Message)

	go chatutil.SendMessages(conn, sendChannel)
	go chatutil.ReceiveMessages(conn, receiveChannel)

	go func() {
		for message := range receiveChannel {
			fmt.Printf("%v: %v\n", message.Sender, message.Text)
		}
	}()

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		sendChannel <- chatutil.Message{Sender: name, Text: input}
	}
}
