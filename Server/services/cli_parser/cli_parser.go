package cli_parser

import (
	"GoChatApp/Server/definitions"
	"errors"
	"os"
	"strconv"
)

// Parse performs the validation and parsing of port no
func Parse() (string, error) {
	portNo, err := getArg()
	if err != nil {
		return "", err
	}
	return portNo, nil
}

// Check arg correctness
func getArg() (string, error) {
	argLen := len(os.Args)
	if argLen < 2 {
		return "", errors.New(definitions.ErrNap)
	} else if argLen > 2 {
		return "", errors.New(definitions.ErrTap)
	} else {
		portArg := os.Args[1]
		if portNo, err := strconv.Atoi(portArg); err == nil {
			if portNo >= definitions.MinPortNo && portNo <= definitions.MaxPortNo {
				return portArg, nil
			}
			return "", errors.New(definitions.ErrBpno)
		}
		return "", errors.New(definitions.ErrBpno)
	}
}
