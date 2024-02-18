package common

type Message struct {
	AgentId                      string
	AgentHostName                string
	AgentCurrentWorkingDirectory string
	Commands                     []Command
}
