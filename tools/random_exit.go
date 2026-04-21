package tools

import (
	"os"
	"runtime/debug"
	"time"
)

func randomExit() {
	if time.Now().Nanosecond()%2 == 0 {
		debug.PrintStack()
		os.Exit(1)
	}
}