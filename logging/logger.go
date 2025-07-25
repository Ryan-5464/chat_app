package logging

import (
	"log"
	"runtime"
	"time"
)

func NewLogger(debug bool) *Logger {
	logger := &Logger{Debug: debug}
	return logger
}

type Logger struct {
	Debug bool
}

func (l *Logger) Log(message string) {
	timestamp := time.Now().UTC().Format(time.RFC3339)

	log.Printf("[MESSAGE] %s : %s", timestamp, message)
}

func (l *Logger) LogError(err error) {
	timestamp := time.Now().UTC().Format(time.RFC3339)

	log.Printf("[ERROR] %s : %s", timestamp, err)
}

func (l *Logger) LogFunctionInfo() {

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
