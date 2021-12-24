package definitions

type OppCode byte

const (
	InitConn  = '0'
	SendMsg   = '1'
	QueryMsg  = '2'
	ToMeMsg   = '3'
	FromMeMsg = '4'
	Nil       = ' '
)

func GetMsgOppCode(msg string) OppCode {
	if msg[0] == InitConn {
		return InitConn
	} else if msg[0] == SendMsg {
		return SendMsg
	} else if msg[0] == QueryMsg {
		return QueryMsg
	}
	return Nil
}
