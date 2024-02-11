package util

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type LogLevel int32

const (
	LOG_DEBUG LogLevel = iota
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_FATAL
	LOG_NONE
)

var logPrefixe = map[LogLevel]string{
	LOG_DEBUG: `[D] `,
	LOG_INFO:  `[I] `,
	LOG_WARN:  `[W] `,
	LOG_ERROR: `[E] `,
	LOG_FATAL: `/!\ `,
}

var levelNames = map[string]LogLevel{
	"debug": LOG_DEBUG,
	"info":  LOG_INFO,
	"warn":  LOG_WARN,
	"error": LOG_ERROR,
	"fatal": LOG_FATAL,
	"none":  LOG_NONE,
}

var minLevel = LOG_INFO

func InitLogger(writer io.Writer) {
	log.SetOutput(writer)
	log.SetFlags(log.Ltime | log.Lmsgprefix)
}

func SetLogLevel(newLevel LogLevel) {
	minLevel = newLevel
}

func SetLogLevelByName(newLevelName string) {
	if newLevel, ok := levelNames[strings.ToLower(newLevelName)]; ok {
		SetLogLevel(newLevel)
	}
}

func Log(aktLevel LogLevel, format string, v ...any) string {
	if minLevel > aktLevel {
		return ""
	}

	var msg strings.Builder
	msg.WriteString(logPrefixe[aktLevel])
	msg.WriteString(fmt.Sprintf(format, v...))
	log.Println(msg.String())
	return msg.String()
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
	panic(fmt.Errorf("%s", Log(LOG_FATAL, format, v...)))
}
