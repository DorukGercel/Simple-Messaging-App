package main

import (
	"GoChatApp/Server/services/cli_parser"
	"GoChatApp/Server/services/error_handler"
	"GoChatApp/Server/services/network"
)

func main() {
	// Get and check port no given as terminal argument
	connPortNo, err := cli_parser.Parse()
	error_handler.CheckFatalError(err)
	// Start dealing with requests
	network.AcceptConnection(connPortNo)
}
