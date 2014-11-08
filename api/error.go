package api

import (
	"fmt"
	"os"
)

func HandleError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func HandleErrorWithMsg(err error, msg string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("error: %s, %s", err, msg))
	}
}

func HandleErrorAndExit(err error) {
	if err != nil {
		fmt.Println("exit: ", err)
		panic(err)
		os.Exit(1)
	}
}

func HandleErrorAndExitWithMsg(err error, msg string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("exit: %s, %s", err, msg))
		panic(err)
		os.Exit(1)
	}
}
