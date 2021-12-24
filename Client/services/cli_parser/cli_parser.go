package cli_parser

import (
	"GoChatApp/Client/definitions"
	"errors"
	"os"
	"strconv"
	"strings"
)

// Parse performs the validation and parsing of port no
func Parse() (string, string, string, error) {
	nickName, ipAddr, portNo, err := getArg()
	if err != nil {
		return "", "", "", err
	}
	return nickName, ipAddr, portNo, nil
}

// Check arg correctness
func getArg() (string, string, string, error) {
	argLen := len(os.Args)
	if argLen < definitions.MinNickLen {
		return "", "", "", errors.New(definitions.ErrNap)
	} else if argLen > 4 {
		return "", "", "", errors.New(definitions.ErrTap)
	} else {
		nickArg := os.Args[1]
		if checkNickArg(nickArg) != nil {
			return "", "", "", errors.New(definitions.ErrSname)
		}

		ipArg := os.Args[2]
		if checkIpArg(ipArg) != nil {
			return "", "", "", errors.New(definitions.ErrWip)
		}

		portArg := os.Args[3]
		portNo, err := strconv.Atoi(portArg)
		if err != nil || checkPortNo(portNo) != nil {
			return "", "", "", errors.New(definitions.ErrBpno)
		}
		return nickArg, ipArg, portArg, nil
	}
}

func checkNickArg(nick string) error {
	if len(nick) < 4 || strings.Contains(nick, "*") || strings.Contains(nick, ":") {
		return errors.New(definitions.ErrSname)
	}
	return nil
}

func checkIpArg(ip string) error {
	// TODO: Write a ip checker
	return nil
}

func checkPortNo(portNo int) error {
	if portNo >= definitions.MinPortNo && portNo <= definitions.MaxPortNo {
		return nil
	}
	return errors.New(definitions.ErrBpno)
}
