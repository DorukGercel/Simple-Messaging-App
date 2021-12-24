package error_handler

import (
	"fmt"
	"os"
	"strings"
)

func CheckFatalError(err error) {
	if err != nil {
		fmt.Println(strings.Title(err.Error()))
		os.Exit(1)
	}
}
