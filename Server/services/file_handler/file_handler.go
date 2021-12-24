package file_handler

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// Name of records file
const fileName = "./chat_records.txt"

// Index of field
type fieldIndex int

const (
	RecordId = iota
	FromLen
	FromNick
	ToLen
	ToNick
	MsgLen
	MsgAll
)

// Index of meta data
type metaIndex int

const (
	NextIdIndex = iota
	NoRecordsIndex
)

// File op types
type fileOp int

const (
	read = iota
	write
)

// Meta op types
type metaOp int

const (
	noChange = iota
	increment
	decrement
)

// Special characters in editing file
const (
	MetaSpr       = "\n"
	NoMeta        = 2
	FieldSpr      = "$"
	LineEnd       = "\n"
	DateLayout    = "2006-01-02"
	SinglePadding = " "
	MultiPadding  = "          "
)

// Index of record fields
const (
	IdIndex = iota
	DateIndex
	StrLenIndex
	StrIndex
)

func GetNoRecords() (int, error) {
	metas, err := getTableMetas()
	if err != nil {
		return -1, err
	}
	noRecords, err := strconv.Atoi(metas[1])
	if err != nil {
		return -1, err
	}
	return noRecords, nil
}

func ReadIdRecord(id int) (record string, err error) {
	var f *os.File
	if f, err = openFile(read); err != nil {
		return "", err
	}
	defer f.Close()
	for i := 0; i < NoMeta; i++ {
		getLine(f)
	}
	for i := 0; i < id; i++ {
		text, _ := getLine(f)
		if text == "" {
			break
		}
	}
	record, _ = getLine(f)
	return record, nil
}

// ReadRecord performs fetching both metadata and records
func ReadRecord() (metas []string, records []string, err error) {
	var f *os.File
	if f, err = openFile(read); err != nil {
		return nil, nil, err
	}
	defer f.Close()
	if metas, err = getTableMetas(); err != nil {
		return nil, nil, err
	}
	records = make([]string, 0)
	for i := 0; i < NoMeta; i++ {
		getLine(f)
	}
	for {
		text, _ := getLine(f)
		if text == "" {
			break
		}
		records = append(records, text)
	}
	return metas, records, nil
}

// WriteRecord adds new record
func WriteRecord(msgInfo [3]string) (err error) {
	var f *os.File
	if f, err = openFile(write); err != nil {
		return err
	}
	defer f.Close()
	if err = findNewRecordPlace(f); err != nil {
		return err
	}
	if _, err = f.Write([]byte(generateRecord(msgInfo))); err != nil {
		return err
	}
	if err = setMetasNewRecord(); err != nil {
		return err
	}
	return nil
}

// Opens file according to read or write type
func openFile(op fileOp) (*os.File, error) {
	if _, err := os.Stat(fileName); err == nil {
		var f *os.File
		if op == read {
			if f, err = os.Open(fileName); err != nil {
				return nil, err
			}
		} else {
			if f, err = os.OpenFile(fileName, os.O_RDWR, 0777); err != nil {
				return nil, err
			}
		}
		return f, nil
	} else if errors.Is(err, os.ErrNotExist) {
		createFile()
		return openFile(op)
	} else {
		return nil, err
	}
}

// Creates records file with default metadata
func createFile() {
	err := os.WriteFile(fileName, []byte(initMetas()), 0777)
	if err != nil {
		return
	}
}

// Returns init meta data values as string
func initMetas() string {
	return "0" + MultiPadding + MetaSpr + "0" + MultiPadding + LineEnd
}

// Create record string
func generateRecord(msgInfo [3]string) string {
	id, err := getNextRecordId()
	if err != nil {
		return ""
	}
	return strconv.Itoa(id) +
		FieldSpr +
		strconv.Itoa(len(msgInfo[0])) +
		FieldSpr +
		msgInfo[0] +
		FieldSpr +
		strconv.Itoa(len(msgInfo[1])) +
		FieldSpr +
		msgInfo[1] +
		FieldSpr +
		strconv.Itoa(len(msgInfo[2])) +
		FieldSpr +
		msgInfo[2] +
		LineEnd
}

// Fetch meta data
func getTableMetas() ([]string, error) {
	var f *os.File
	var err error
	if f, err = openFile(read); err != nil {
		return nil, err
	}
	metas := make([]string, NoMeta)
	for i := 0; i < NoMeta; i++ {
		text, _ := getLine(f)
		metas[i] = text[0:strings.Index(text, SinglePadding)]
	}
	return metas, nil
}

// Updates given metadata value
func setTableMetas(mIndex metaIndex, mOp metaOp) error {
	var f *os.File
	var err error
	if f, err = openFile(write); err != nil {
		return err
	}
	f.Seek(0, 0)
	var text string
	var oldPos int64
	for i := 0; i < int(mIndex)+1; i++ {
		text, oldPos = getLine(f)
	}
	oldValStr := text
	var oldVal int
	if oldVal, err = strconv.Atoi(oldValStr[0:strings.Index(oldValStr, SinglePadding)]); err != nil {
		return err
	}
	var newVal int
	if newVal = oldVal + 1; mOp == decrement {
		newVal = oldVal - 1
	}
	f.Seek(oldPos, 0)
	newValStr := strconv.Itoa(newVal)
	if len(oldValStr[0:strings.Index(oldValStr, SinglePadding)]) > len(newValStr) {
		newValStr = newValStr + SinglePadding
	}
	if _, err := f.Write([]byte(newValStr)); err != nil {
		return err
	}
	return nil
}

// Get next record id
func getNextRecordId() (int, error) {
	if metas, err := getTableMetas(); err != nil {
		return -1, err
	} else {
		nextId, err := strconv.Atoi(metas[NextIdIndex])
		if err != nil {
			return -1, nil
		}
		return nextId, nil
	}
}

// Set next record id
func setNextRecordId() error {
	return setTableMetas(NextIdIndex, increment)
}

// Set no record
func setNoRecord(op metaOp) error {
	return setTableMetas(NoRecordsIndex, op)
}

// Set metas for new record op
func setMetasNewRecord() error {
	if err := setNoRecord(increment); err != nil {
		return err
	}
	if err := setNextRecordId(); err != nil {
		return err
	}
	return nil
}

// Find appropriate place for new record
func findNewRecordPlace(f *os.File) error {
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}
	for {
		text, _ := getLine(f)
		if text == "" {
			break
		}
	}
	return nil
}

// Get single line from file
func getLine(f *os.File) (text string, oldPos int64) {
	text = ""
	buf := make([]byte, 10)
	oldPos, _ = f.Seek(0, 1)
	for {
		if _, err := f.Read(buf); err != nil {
			return
		}
		text = text + string(buf)
		if i := strings.Index(text, LineEnd); i > -1 {
			text = text[0:i]
			if _, err := f.Seek(oldPos+int64(len(text))+1, 0); err != nil {
				return
			}
			return
		}
	}
}
