package errorlogs

import (
	"github.com/Siroshun09/logs/v2"
)

func NewNilLogger() logs.Logger {
	var l *logger
	return l
}
