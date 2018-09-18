package level

import "strings"

type Level int

const (
	All Level = iota
	Trace
	Debug
	Info
	Warn
	Error
	Fatal
	Off
)

var levels = []Level{All, Trace, Debug, Info, Warn, Error, Fatal, Off}

func ToLoglevel(s string) Level {
	for _, level := range levels {
		if strings.ToUpper(s) == level.String() {
			return level
		}
	}

	return All // FIXME: this ought not happen, have to come up with some error
}

func (l Level) String() string {
	switch l {
	case Off:
		return "OFF"
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	case All:
		return "ALL"
	}

	return "Level-not-found"
}
