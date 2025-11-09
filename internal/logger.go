package internal

import (
	"log"
	"os"
)

var LOGGER = newCustomLogger()

func newCustomLogger() *log.Logger {
	return log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
}
