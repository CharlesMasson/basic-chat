package chatutil

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"time"
)

// Message is a structure that represents a message.
type Message struct {
	Sender string
	Text   string
}

// ReadNBytes reads conn until n bytes are received.
// This is a blocking function
func ReadNBytes(conn net.Conn, n int) []byte {
	bytes := make([]byte, n)
	readBytes := 0
	for readBytes < n {
		r, _ := conn.Read(bytes[readBytes:])
		readBytes += r
		time.Sleep(time.Millisecond)
	}
	return bytes
}

// SendMessages forwards messages that are received via channel to conn.
func SendMessages(conn net.Conn, channel chan Message) {
	for message := range channel {
		bytes, err := json.Marshal(message)
		if err != nil {
			log.Fatal(err)
		}
		lengthBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(lengthBytes, uint16(len(bytes)))
		conn.Write(lengthBytes)
		conn.Write(bytes)
	}
}

// ReceiveMessages forwards messages that are received via conn to channel.
func ReceiveMessages(conn net.Conn, channel chan Message) {
	var message Message
	var length int
	var err error
	for {
		length = int(binary.BigEndian.Uint16(ReadNBytes(conn, 2)))
		err = json.Unmarshal(ReadNBytes(conn, length), &message)
		if err != nil {
			log.Fatal(err)
		}
		channel <- message
	}
}
