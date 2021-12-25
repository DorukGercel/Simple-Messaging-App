package error_handler

import (
	"fmt"
	"os"
	"strings"
)

// CheckFatalError closes process with fatal error
func CheckFatalError(err error) {
	if err != nil {
		fmt.Println(strings.Title(err.Error()))
		os.Exit(1)
	}
}
