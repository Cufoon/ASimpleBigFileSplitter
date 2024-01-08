package main

import (
	"errors"
	"fmt"
)

var ErrShouldToExit = errors.New("ERR_TO_EXIT")

func handleShouldToExitErr(err error) {
	if errors.Is(err, ErrShouldToExit) {
		return
	}
	fmt.Println(err)
}
