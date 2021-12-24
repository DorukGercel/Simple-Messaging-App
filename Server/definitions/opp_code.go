package definitions

import (
	"strconv"
)

type OppCode byte

const (
	InitConn  = '0'
	SendMsg   = '1'
	QueryMsg  = '2'
	ToMeMsg   = '3'
	FromMeMsg = '4'
	Nil       = ' '
)

func GetMsgOppCodeAndMessage(req string) (OppCode, string) {
	if req[0] == InitConn {
		return InitConn, req[1:]
	} else if req[0] == SendMsg {
		return SendMsg, req[1:]
	} else if req[0] == QueryMsg {
		return QueryMsg, req[1:]
	}
	return Nil, ""
}

func GetQueryOppCode(msg string) (OppCode, int) {
	var queryOpp OppCode
	if len(msg) >= 2 {
		queryOpp = OppCode(msg[0])
		limit, err := strconv.Atoi(msg[1:])
		if err == nil {
			return queryOpp, limit
		} else {
			return queryOpp, -1
		}
	}
	return Nil, -1
}
