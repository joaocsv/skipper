package main

import (
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"log"
	"net"
	"os"
	"time"

	"github.com/joaocsv/skipper/common"
)

var (
	message common.Message
)

func main() {
	log.Println("Starting agent...")
	presentSettings()
	log.Println("Identifier: " + message.AgentId)

	for {
		channel := connectServer()

		gob.NewEncoder(channel).Encode(message)
		gob.NewDecoder(channel).Decode(&message)

		channel.Close()
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func presentSettings() {
	message.AgentId = generateId()
	message.AgentHostName, _ = os.Hostname()
	message.AgentCurrentWorkingDirectory, _ = os.Getwd()
}

func generateId() string {
	hostname, _ := os.Hostname()
	timeNow := time.Now().String()
	hasher := md5.New()

	hasher.Write([]byte(hostname + timeNow))
	return hex.EncodeToString(hasher.Sum(nil))
}

func connectServer() net.Conn {
	channel, error := net.Dial("tcp", "127.0.0.1:8080")

	if error != nil {
		log.Fatal("Unable to connect to the server")
	}

	return channel
}
