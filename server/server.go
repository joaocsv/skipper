package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"os"
	"strings"

	"github.com/joaocsv/skipper/common"
)

var (
	agents          = []common.Message{}
	selectedAgentId string
)

func main() {
	log.Println("Starting server...")

	go startListener("8080")
	interfaceHandler()
}

func startListener(port string) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatal("Unable to start the server")
	}

	for {
		channel, err := listener.Accept()

		if err != nil {
			log.Println("Unable to establish a connection")
			continue
		}

		message := &common.Message{}
		gob.NewDecoder(channel).Decode(&message)

		if !agentExists(message.AgentId) {
			log.Println("Connection established: " + channel.RemoteAddr().String())
			log.Println("AgentId: " + message.AgentId)
			agents = append(agents, *message)
		} else {
			for _, command := range message.Commands {
				if len(command.Response) > 0 {
					log.Println("AgentId: " + message.AgentId)
					log.Println("Command: " + command.Command + " | Response: " + command.Response)
				}
			}

		}

		gob.NewEncoder(channel).Encode(message)

		channel.Close()
	}
}

func agentExists(agentId string) bool {
	exists := false

	for _, agent := range agents {
		if agent.AgentId == agentId {
			exists = true
		}
	}

	return exists
}

func interfaceHandler() {
	for {
		if selectedAgentId != "" {
			print(selectedAgentId + "@skipper# ")
		} else {
			print("skipper > ")
		}

		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		commandSplit := strings.Split(strings.TrimSpace(strings.TrimSuffix(command, "\n")), " ")
		commandBase := strings.TrimSpace(commandSplit[0])

		if commandBase != "" {
			switch commandBase {
			case "select":
				selectCommandHandler(commandSplit)
			case "show":
				showCommandHandler(commandSplit)
			case "exit":
				if selectedAgentId != "" {
					selectedAgentId = ""
					return
				}

				log.Println("Command doesn't exist")
			default:
				log.Println("Command doesn't exist")
			}
		}
	}
}

func selectCommandHandler(command []string) {
	if len(command) <= 1 {
		log.Println("You must inform the agent to be selected!")
		log.Println("Example: select agentId")

		return
	}

	if !agentExists(command[1]) {
		log.Println("The informed agent doesn't exist")
		return
	}

	selectedAgentId = command[1]
}

func showCommandHandler(command []string) {
	if len(command) <= 1 {
		log.Println("You must inform the parameter to be show!")
		log.Println("Example: show agents")

		return
	}

	switch command[1] {
	case "agents":
		for _, agent := range agents {
			println(agent.AgentId + "@" + agent.AgentHostName + ":" + agent.AgentCurrentWorkingDirectory)
		}
	default:
		log.Println("The informed parameter does not exist")
	}
}
