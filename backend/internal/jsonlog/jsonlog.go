// BIOAFF/backend/internal/jsonlog/jsonlog.go

package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// we can have different servirty levels of logging entries
type Level int8

// levels start at zero
const (
	LevelInfo  Level = iota //value is 0
	LevelError              //value is 1
	LevelFatal              //value is 2
	LevelOff                //value is 3
)

// The severity levels as a human readable friendly format
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Defind a custom logger
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// The new() function creates a new instance of logger
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// Helper methods
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	//Ensure severity level is at least minimal
	if level < l.minLevel {
		return 0, nil
	}

	//Create a struct for holding the log entry data
	data := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitemtpy`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	//should we include the stack trace?
	if level >= LevelError {
		data.Trace = string(debug.Stack())
	}

	//Encode the log entry to JSON
	var line []byte
	line, err := json.Marshal(data)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message" + err.Error())
	}

	//Prepare to write the log entry
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out.Write(append(line, '\n'))
}

// Implement the io.Write interface
func (l *Logger) Wrtie(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
