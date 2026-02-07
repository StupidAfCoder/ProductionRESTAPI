package utils

import (
	"fmt"
	"log"
	"os"
)

func ErrorHandler(err error, message string) error {
	errorLogger := log.New(os.Stderr, "ERROR :", log.Ltime|log.Ldate|log.Lshortfile)
	errorLogger.Println(message, err)
	return fmt.Errorf(message)
}
