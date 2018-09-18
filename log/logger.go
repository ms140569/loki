package log

import (
	"loki/log/level"
	"fmt"
	"os"
	"strings"
)

// SystemLogLevel is the systemwide loglevel at any given time
var SystemLogLevel level.Level

// SetSystemLogLevel adjusts the systemwide loglevel used by all
// output.
func SetSystemLogLevel(l level.Level) {
	SystemLogLevel = l
}

// SetSystemLogLevelFromString sets the systemwide loglevel to the
// one given in this string
func SetSystemLogLevelFromString(s string) {
	SystemLogLevel = level.ToLoglevel(strings.ToUpper(s))
}

// LeveledLogger is the dispatcher to process all output depending on the given level
func LeveledLogger(l level.Level, format string, v ...interface{}) {
	if l >= SystemLogLevel {
		if l == level.Error || l == level.Fatal {
			fmt.Fprintf(os.Stderr, format+"\n", v...)
		} else {
			fmt.Printf(format+"\n", v...)
		}
	}
}

// Trace logs the given string with the corresponding level
func Trace(format string, v ...interface{}) {
	LeveledLogger(level.Trace, format, v...)
}

// Debug logs the given string with the corresponding level
func Debug(format string, v ...interface{}) {
	LeveledLogger(level.Debug, format, v...)
}

// Info logs the given string with the corresponding level
func Info(format string, v ...interface{}) {
	LeveledLogger(level.Info, format, v...)
}

// Warn logs the given string with the corresponding level
func Warn(format string, v ...interface{}) {
	LeveledLogger(level.Warn, format, v...)
}

// Error logs the given string with the corresponding level
func Error(format string, v ...interface{}) {
	LeveledLogger(level.Error, format, v...)
}

// Fatal logs the given string with the corresponding level
func Fatal(format string, v ...interface{}) {
	LeveledLogger(level.Fatal, format, v...)
}

// All logs the given string with the corresponding level
func All(format string, v ...interface{}) {
	LeveledLogger(level.All, format, v...)
}
