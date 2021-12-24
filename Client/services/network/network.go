package network

import (
	"GoChatApp/Client/definitions"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	list       = "List of users: "
	incoming   = "\nIncoming message: "
	typesmt    = "Type Something: "
	badMsgForm = "Bad message format!"
)

func InitConnection(nickname string, conn net.Conn) {
	fmt.Fprintf(conn, formInitMessage(nickname))
	msg, _ := bufio.NewReader(conn).ReadString(definitions.MsgEndChar)
	msg = msg[:len(msg)-1]
	readMessage(msg, list)
}

// ListenServer Handles incoming messages.
func ListenServer(conn net.Conn) {
	for {
		msg, err := bufio.NewReader(conn).ReadString(definitions.MsgEndChar)
		if err != nil {
			handleServerShutdown()
			break
		}
		msg = msg[:len(msg)-1]
		readMessage(msg, incoming)
		fmt.Print(typesmt)
	}
	conn.Close()
	os.Exit(0)
}

func SendMessage(conn net.Conn, nickname string) {
	fmt.Print(typesmt)
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = formatInput(text)
		if isValidMsg(text) {
			if isValidQueryReq(text) {
				fmt.Fprintf(conn, formQueryMessage(text))
			} else if isValidSendMsg(text, nickname) {
				fmt.Fprintf(conn, formSendMessage(text))
				fmt.Print(typesmt)
			} else {
				printInputError(badMsgForm)
			}
		} else {
			printInputError(badMsgForm)
		}
	}
}

func isValidMsg(text string) bool {
	if text[0] != ' ' {
		return true
	}
	return false
}

func isValidSendMsg(text string, nickname string) bool {
	textSlice := strings.Split(text, " ")
	if len(textSlice) >= 2 && len(textSlice[0]) >= definitions.MinNickLen && textSlice[0] != nickname {
		return true
	}
	return false
}

func isQueryReq(text string) bool {
	if len(text) >= len(definitions.QueryStrWoutLim) && text[0] == definitions.Query && text[1] == definitions.Delim {
		return true
	}
	return false
}

func isValidQueryReq(text string) bool {
	if !isQueryReq(text) {
		return false
	}
	if !isFromMeQueryReq(text) && !isToMeQueryReq(text) {
		return false
	}
	if len(text) == len(definitions.QueryStrWLim) {
		textSlice := strings.Split(text, " ")
		if !containsLimitValue(textSlice) {
			return false
		}
	}
	return true
}

func isFromMeQueryReq(text string) bool {
	if text[2] == definitions.FromMe && text[1] == definitions.Delim {
		return true
	}
	return false
}

func isToMeQueryReq(text string) bool {
	if text[2] == definitions.ToMe && text[1] == definitions.Delim {
		return true
	}
	return false
}

func containsLimitValue(textSlice []string) bool {
	if len(textSlice) == 3 {
		val, err := strconv.Atoi(textSlice[2])
		if err != nil {
			fmt.Println(val, err)
			return false
		}
		return true
	}
	return false
}

func formInitMessage(nickname string) string {
	msg := string(definitions.InitConn) + nickname + string(definitions.MsgEndChar)
	return msg
}

func formSendMessage(text string) string {
	msg := string(definitions.SendMsg) + text + string(definitions.MsgEndChar)
	return msg
}

func formQueryMessage(text string) string {
	queryMessage := string(definitions.QueryMsg)
	if isFromMeQueryReq(text) {
		queryMessage += string(definitions.FromMeMsg)
	} else if isToMeQueryReq(text) {
		queryMessage += string(definitions.ToMeMsg)
	}
	textSlice := strings.Split(text, " ")
	if containsLimitValue(textSlice) {
		queryMessage += textSlice[2]
	} else {
		queryMessage += string(definitions.MsqEmptyLimit)
	}
	return queryMessage + string(definitions.MsgEndChar)
}

func handleServerShutdown() {
	// Close the connection when server is down.
	fmt.Println("Server closed!")
}

func readMessage(msg string, clientExtra string) {
	if strings.Contains(msg, string(definitions.MsgSenderDelim)+string(definitions.ListItemDelim)) {
		// Check if list text
		readListMessage(msg, clientExtra)
		return
	}
	readNormalMessage(msg, clientExtra)
}

func readNormalMessage(msg string, clientExtra string) {
	fmt.Println(clientExtra + msg)
	fmt.Println()
}

func readListMessage(msg string, clientExtra string) {
	listTxt := msg[strings.Index(msg, string(definitions.MsgSenderDelim)+string(definitions.ListItemDelim))+2:]
	list := strings.Split(listTxt, string(definitions.ListItemDelim))
	fmt.Println(clientExtra)
	for _, val := range list {
		fmt.Println(val)
	}
	fmt.Println()
}

func printInputError(warning string) {
	fmt.Println(warning + "\n")
	fmt.Print(typesmt)
}

func formatInput(text string) string {
	index := -1
	for i, c := range text {
		if c != ' ' && c != '\n' {
			index = i
		}
	}
	cleanText := text[:index+1]
	return cleanText
}
