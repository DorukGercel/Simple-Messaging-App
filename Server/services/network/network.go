package network

import (
	"GoChatApp/Server/definitions"
	"GoChatApp/Server/services/db_handler"
	"GoChatApp/Server/services/error_handler"
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

const server = "Server"

func AcceptConnection(connPortNo string) {
	// Listen for incoming connections.
	l, err := net.Listen(definitions.ConnType, definitions.ConnHost+":"+connPortNo)
	error_handler.CheckFatalError(err)
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + definitions.ConnHost + ":" + connPortNo)
	// Init connections map
	connMap := map[string]net.Conn{}
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		error_handler.CheckFatalError(err)
		fmt.Println("New Conn")
		go HandleRequest(conn, connMap)
	}
}

// HandleRequest Handles incoming requests.
func HandleRequest(currentConn net.Conn, connMap map[string]net.Conn) {
	var currentUserNick string
	for {
		req, err := bufio.NewReader(currentConn).ReadString(definitions.MsgEndChar)
		if err != nil {
			handleClientShutdown(currentUserNick, connMap)
			break
		}
		req = req[:len(req)-1]
		oppCode, msg := definitions.GetMsgOppCodeAndMessage(req)
		if oppCode == definitions.InitConn {
			currentUserNick, err = initConn(msg, currentConn, connMap)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Fprintf(currentConn, formMsg(server, formTextForList(getUsersList(connMap))))
		} else if oppCode == definitions.SendMsg {
			sendMessage(msg, currentUserNick, currentConn, connMap)
		} else if oppCode == definitions.QueryMsg {
			var records []string
			queryOppCode, queryLimit := definitions.GetQueryOppCode(msg)
			if queryOppCode == definitions.ToMeMsg {
				records = db_handler.ExecQuery(currentUserNick, db_handler.TO, queryLimit)
			} else {
				records = db_handler.ExecQuery(currentUserNick, db_handler.FROM, queryLimit)
			}
			fmt.Fprintf(currentConn, formMsg(server, formTextForList(records)))
		}
	}
	currentConn.Close()
}

func initConn(msg string, currentConn net.Conn, connMap map[string]net.Conn) (string, error) {
	fmt.Println("Name: " + msg)
	if _, ok := connMap[msg]; !ok {
		currentUserNick := msg
		connMap[currentUserNick] = currentConn
		return currentUserNick, nil
	}
	fmt.Println("User already exists")
	return "", errors.New(definitions.ErrUexists)
}

func sendMessage(msg string, currentUserNick string, currentConn net.Conn, connMap map[string]net.Conn) {
	msgSlice := getSendMsgSlice(msg)
	if checkSendMsgSlice(msgSlice) {
		sentUserNick, chatMsg := getSendUserAndMsg(msgSlice)
		if currentUserNick == sentUserNick {
			fmt.Println("You can't send message to yourself")
			fmt.Fprint(currentConn, formMsg(server, definitions.ErrSendSame))
		} else if sentConn, ok := connMap[sentUserNick]; ok {
			go db_handler.WriteRecord([3]string{currentUserNick, sentUserNick, chatMsg})
			fmt.Fprint(sentConn, formMsg(currentUserNick, chatMsg))
		} else {
			fmt.Fprint(currentConn, formMsg(server, definitions.ErrUNotexist))
		}
	} else {
		fmt.Fprint(currentConn, formMsg(server, definitions.ErrNarg))
	}
}

func getUsersList(connMap map[string]net.Conn) []string {
	var list []string
	for key := range connMap {
		list = append(list, key)
	}
	return list
}

func formMsg(sender string, msg string) string {
	if sender == "" {
		return msg + string(definitions.MsgEndChar)
	}
	return sender + string(definitions.MsgSenderDelim) + msg + string(definitions.MsgEndChar)
}

func handleClientShutdown(currentUserNick string, connMap map[string]net.Conn) {
	// Close the connection when you're done with it.
	delete(connMap, currentUserNick)
	fmt.Println("User " + currentUserNick + " deleted")
}

func getSendMsgSlice(msg string) []string {
	return strings.SplitN(msg, string(definitions.PersonMsgDelim), 2)
}

func checkSendMsgSlice(msgSlice []string) bool {
	return len(msgSlice) == 2
}

func getSendUserAndMsg(msgSlice []string) (sentUserNick string, chatMsg string) {
	sentUserNick = msgSlice[0][0:len(msgSlice[0])]
	chatMsg = msgSlice[1]
	return
}

func formTextForList(list []string) string {
	listText := string(definitions.ListItemDelim)
	noRecords := len(list)
	for i, val := range list {
		if i == (noRecords - 1) {
			listText += val
		} else {
			listText += val + string(definitions.ListItemDelim)
		}
	}
	return listText
}
