// The log module provides some slightly more advanced logging in a simple function
package prm

import (
	"log"
)

const (
	LOG_DEBUG = iota
	LOG_INFO  = iota
	LOG_WARN  = iota
	LOG_ERROR = iota
)

func LogLevelToString(loglevel int) string {

	if loglevel == LOG_DEBUG {
		return "debug"
	}

	if loglevel == LOG_INFO {
		return "info"
	}

	if loglevel == LOG_WARN {
		return "warn"
	}

	return "error"

}

// LogPRM is a helper function for logging messages from prm / prmserver
func (prm *PRM) LogPRM(msg string, loglevel int) {
	if loglevel == LOG_ERROR || loglevel >= prm.Config.LogLevel {
		log.Println("[prm:"+LogLevelToString(loglevel)+"]", msg)
	}
}
