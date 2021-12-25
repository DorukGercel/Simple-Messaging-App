package db_handler

import (
	"GoChatApp/Server/services/file_handler"
	"fmt"
	"strings"
)

type MsgDirection int

const (
	FROM = iota
	TO
)

// WriteRecord write record to the record file
func WriteRecord(recordInfo [3]string) {
	if err := file_handler.WriteRecord(recordInfo); err != nil {
		fmt.Println(err)
	}
}

// ExecQuery perform required query operations
func ExecQuery(nickname string, dir MsgDirection, limit int) []string {
	noRecords, err := file_handler.GetNoRecords()
	if err != nil {
		fmt.Println(err)
	}
	var lookId int
	if lookId = file_handler.FromNick; dir == TO {
		lookId = file_handler.ToNick
	}
	// Related records
	var relatedRecords []string
	// Traverse all the records in the file
	for i := 0; i < noRecords; i++ {
		record, _ := file_handler.ReadIdRecord(i)
		recordSlice := strings.Split(record, file_handler.FieldSpr)
		fmt.Println(recordSlice)
		if recordSlice[lookId] == nickname {
			// Related record
			relatedRecords = append(relatedRecords, recordSlice[file_handler.FromNick]+"->"+recordSlice[file_handler.ToNick]+": "+recordSlice[file_handler.MsgAll])
		}
	}
	if limit > len(relatedRecords) || limit == -1 {
		return relatedRecords
	}
	return relatedRecords[len(relatedRecords)-limit:]
}
