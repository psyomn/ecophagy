package common

import (
	"log"
)

func DefaultLogSettings() { log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) }
