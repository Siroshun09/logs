//go:generate go tool mockgen -destination=mock.gen.go -package=logmock -typed github.com/Siroshun09/logs Logger
package logmock

import (
	"github.com/Siroshun09/logs"
)

var _ logs.Logger = (*MockLogger)(nil)
