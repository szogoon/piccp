package logging

import (
	"fmt"
	"os"
)

var Verbose bool

func Info(format string, a ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", a...)
}

func Error(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERROR] "+format+"\n", a...)
}

func Debug(format string, a ...interface{}) {
	if Verbose {
		fmt.Printf("[DEBUG] "+format+"\n", a...)
	}
}

func Fatal(format string, a ...interface{}) {
	Error(format, a...)
	os.Exit(1)
}
