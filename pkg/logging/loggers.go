package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/madeinheaven91/black-turtle-go/pkg/config"
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
)

var (
	criticalLogger = log.New(os.Stdout, "\033[1;30;41m[CRITICAL]\033[0m ", log.LstdFlags)
	telegramLogger = log.New(os.Stdout, "\033[1;34m[TELEGRAM]\033[0m ", log.LstdFlags)
	infoLogger     = log.New(os.Stdout, "\033[1;94m[INFO]\033[0m ", log.LstdFlags)
	warningLogger  = log.New(os.Stdout, "\033[1;33m[WARN]\033[0m ", log.LstdFlags)
	errorLogger    = log.New(os.Stdout, "\033[1;31m[ERROR]\033[0m ", log.LstdFlags)
	debugLogger    = log.New(os.Stdout, "\033[1;92m[DEBUG]\033[0m ", log.LstdFlags|log.Lmicroseconds)
	traceLogger    = log.New(os.Stdout, "\033[1;90m[TRACE]\033[0m ", log.LstdFlags|log.Lmicroseconds)

	logLevels = map[string]int{
		"CRITICAL": 0,
		"ERROR":    1,
		"WARNING":  2,
		"TELEGRAM": 4,
		"INFO":     3,
		"DEBUG":    5,
		"TRACE":    6,
	}
	logLevel = logLevels[os.Getenv("LOG_LEVEL")]

	once sync.Once
)

func InitLoggers() {
	once.Do(func() {
		levelStr := config.LogLevel()
		if level, exists := logLevels[levelStr]; exists {
			logLevel = level
		} else {
			logLevel = logLevels["INFO"]
		}
	})
	fmt.Println(lexicon.Get(lexicon.Greeting))

	var lvl string
	for k, v := range logLevels {
		if logLevel == v {
			lvl = k
			break
		}
	}
	Info("Log level: %s\n", lvl)
}

func Critical(input string, args ...any) {
	if logLevel >= logLevels["CRITICAL"] {
		_, file, no, ok := runtime.Caller(1)
		if ok {
			args = append([]any{no}, args...)
			args = append([]any{file}, args...)
			input = "%s:%d " + input
			criticalLogger.Printf(input, args...)
		}
	}
}

func Telegram(input string, args ...any) {
	if logLevel >= logLevels["TELEGRAM"] {
		telegramLogger.Printf(input, args...)
	}
}

func Info(input string, args ...any) {
	if logLevel >= logLevels["INFO"] {
		infoLogger.Printf(input, args...)
	}
}

func Warning(input string, args ...any) {
	if logLevel >= logLevels["WARNING"] {
		warningLogger.Printf(input, args...)
	}
}

func Error(input string, args ...any) {
	if logLevel >= logLevels["ERROR"] {
		_, file, no, ok := runtime.Caller(1)
		if ok {
			args = append([]any{no}, args...)
			args = append([]any{file}, args...)
			input = "%s:%d " + input
			errorLogger.Printf(input, args...)
		}
	}
}

func Debug(input string, args ...any) {
	if logLevel >= logLevels["DEBUG"] {
		_, file, no, ok := runtime.Caller(1)
		if ok {
			args = append([]any{no}, args...)
			args = append([]any{file}, args...)
			input = "%s:%d " + input
			debugLogger.Printf(input, args...)
		}
	}
}

func Trace(input string, args ...any) {
	if logLevel >= logLevels["TRACE"] {
		_, file, no, ok := runtime.Caller(1)
		if ok {
			args = append([]any{no}, args...)
			args = append([]any{file}, args...)
			input = "%s:%d " + input
			traceLogger.Printf(input, args...)
		}
	}
}
