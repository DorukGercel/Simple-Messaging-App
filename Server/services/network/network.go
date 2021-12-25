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

// AcceptConnection starts the server and starts new routine for incoming requests
func AcceptConnection(connPortNo string) {
	// Listen for incoming connections
	l, err := net.Listen(definitions.ConnType, definitions.ConnHost+":"+connPortNo)
	error_handler.CheckFatalError(err)
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + definitions.ConnHost + ":" + connPortNo)
	// Init connections map
	connMap := map[string]net.Conn{}
	for {
		// Listen for an incoming connection
		conn, err := l.Accept()
		error_handler.CheckFatalError(err)
		go handleRequest(conn, connMap)
	}
}

// Handles incoming requests
func handleRequest(currentConn net.Conn, connMap map[string]net.Conn) {
	// User nick of the connection
	var currentUserNick string
	for {
		// Wait for request
		req, err := bufio.NewReader(currentConn).ReadString(definitions.MsgEndChar)
		if err != nil {
			// If client closes connection shut down
			handleClientShutdown(currentUserNick, connMap)
			break
		}
		// Format request
		req = formatUserReq(req)
		// Get the opp code and payload of request
		oppCode, msg := definitions.GetMsgOppCodeAndMessage(req)
		if oppCode == definitions.InitConn {
			// Deal init conn req
			if currentUserNick, err = sendInitConn(msg, currentConn, connMap); err != nil {
				// If error occurred break connection
				break
			}
		} else if oppCode == definitions.SendMsg {
			// Deal send message to other user request
			sendMessage(msg, currentUserNick, currentConn, connMap)
		} else if oppCode == definitions.QueryMsg {
			// Deal query req
			sendQueryResponse(msg, currentUserNick, currentConn)
		}
	}
	currentConn.Close()
}

func sendInitConn(msg string, currentConn net.Conn, connMap map[string]net.Conn) (currentUserNick string, err error) {
	currentUserNick, err = initConn(msg, currentConn, connMap)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// Send users list
	fmt.Fprintf(currentConn, formMsg(server, formTextForList(getUsersList(connMap))))
	return currentUserNick, nil
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

func sendQueryResponse(msg string, currentUserNick string, currentConn net.Conn) {
	records := matchExecuteQuery(currentUserNick, msg)
	fmt.Fprintf(currentConn, formMsg(server, formTextForList(records)))
}

func initConn(msg string, currentConn net.Conn, connMap map[string]net.Conn) (string, error) {
	fmt.Println("Name: " + msg + " joined!")
	if _, ok := connMap[msg]; !ok {
		currentUserNick := msg
		connMap[currentUserNick] = currentConn
		return currentUserNick, nil
	}
	fmt.Println("User already exists")
	return "", errors.New(definitions.ErrUexists)
}

func getUsersList(connMap map[string]net.Conn) []string {
	var list []string
	for key := range connMap {
		list = append(list, key)
	}
	return list
}

func matchExecuteQuery(currentUserNick string, msg string) []string {
	var records []string
	queryOppCode, queryLimit := definitions.GetQueryOppCode(msg)
	if queryOppCode == definitions.ToMeMsg {
		// Check direction of query
		records = db_handler.ExecQuery(currentUserNick, db_handler.TO, queryLimit)
	} else {
		records = db_handler.ExecQuery(currentUserNick, db_handler.FROM, queryLimit)
	}
	return records
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

func formatUserReq(req string) string {
	return req[:len(req)-1]
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
