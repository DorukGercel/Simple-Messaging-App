package main

import (
	"GoChatApp/Client/services/cli_parser"
	"GoChatApp/Client/services/error_handler"
	"GoChatApp/Client/services/network"
	"GoChatApp/Server/definitions"
	"net"
)

func main() {
	// Get and check nickname, ip addr, port no given as terminal argument
	nickname, ipAddr, portNo, err := cli_parser.Parse()
	error_handler.CheckFatalError(err)
	conn, err := net.Dial(definitions.ConnType, ipAddr+":"+portNo)
	error_handler.CheckFatalError(err)
	network.InitConnection(nickname, conn)
	go network.ListenServer(conn)
	network.SendMessage(conn, nickname)
}
