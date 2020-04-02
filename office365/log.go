package office365

import (
	"log"
	"os"
)

// Event represents a logging event.
type Event struct {
	id         int
	messageFmt string
}

var (
	debugMessage = Event{id: 1, messageFmt: "[DEBUG] %s"}
)

// GoLogger wraps a standard go logger.
type GoLogger struct {
	*log.Logger
}

// NewLogger returns an instance of the GoLogger.
func NewLogger(l *log.Logger) *GoLogger {
	if l == nil {
		l = log.New(os.Stdout, "office365", log.Flags())
	}
	return &GoLogger{l}
}

// Debug output a formatted message with a debug header.
func (l *GoLogger) Debug(message string) {
	l.Printf(debugMessage.messageFmt, message)
}
