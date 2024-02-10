package util

import (
	"io"
	"log"
	"os"
)

type LogLevel int32

const (
	LOG_DEBUG LogLevel = iota
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_NONE
)

var logPrefixe = map[LogLevel]string{
	LOG_DEBUG: "[D]",
	LOG_INFO:  "[I]",
	LOG_WARN:  "[W]",
	LOG_ERROR: "[E]",
}

var minLevel = LOG_INFO

func InitLogger(writer io.Writer) {
	log.SetOutput(writer)
	log.SetFlags(log.Ltime | log.Lmsgprefix)
}

func SetLogLevel(newLevel LogLevel) {
	minLevel = newLevel
}

func Log(aktLevel LogLevel, format string, v ...any) {
	if minLevel <= aktLevel {
		log.Print(logPrefixe[aktLevel])
		log.Printf(format, v...)
		log.Println()
	}
}

func Debug(format string, v ...any) {
	Log(LOG_DEBUG, format, v...)
}

func Info(format string, v ...any) {
	Log(LOG_INFO, format, v...)
}

func Warn(format string, v ...any) {
	Log(LOG_WARN, format, v...)
}

func Error(format string, v ...any) {
	Log(LOG_ERROR, format, v...)
}

func Fatal(format string, v ...any) {
	Error(format, v...)
	os.Exit(-1)
}
