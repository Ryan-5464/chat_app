package util

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

var Log *Logger = NewLogger(false)

func NewLogger(debug bool) *Logger {
	logger := &Logger{Debug: debug}
	return logger
}

type Logger struct {
	Debug bool
}

func (l *Logger) FunctionInfo() {
	if !l.Debug {
		return
	}

	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		log.Println("[LOGERROR] Unable to get function information.")
		return
	}
	fn := runtime.FuncForPC(pc)

	var fname string
	if fn == nil {
		fname = "unknown"
	} else {
		fname = fn.Name()
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	log.Printf("[DEBUG] %s : %s : %s : (line %d)", timestamp, fname, file, line)

}

func (l *Logger) Error(err error) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	log.Printf("[ERROR] %s : %s", timestamp, err)
}

func (l *Logger) Errorf(message string, values ...any) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	log.Printf("[ERROR] %s : %s", timestamp, fmt.Sprintf(message, values...))
}

func (l *Logger) Dbug(message string) {
	if !l.Debug {
		return
	}
	timestamp := time.Now().UTC().Format(time.RFC3339)
	log.Printf("[DEBUG] %s : %s", timestamp, message)
}

func (l *Logger) Dbugf(message string, values ...any) {
	if !l.Debug {
		return
	}
	timestamp := time.Now().UTC().Format(time.RFC3339)
	log.Printf("[DEBUG] %s : %s", timestamp, fmt.Sprintf(message, values...))
}

func (l *Logger) Info(message string) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	log.Printf("[MESSAGE] %s : %s", timestamp, message)
}

func (l *Logger) Infof(message string, values ...any) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	log.Printf("[MESSAGE] %s : %s", timestamp, fmt.Sprintf(message, values...))
}
