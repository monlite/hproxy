package hproxy

import (
	"log"
	"os"
)

// global logger
var L = log.New(os.Stdout, "hproxy: ", log.Lshortfile|log.LstdFlags)
