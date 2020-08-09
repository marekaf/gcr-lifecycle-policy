package utils

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// GetTimeNow returns timestamp two weeks ago
func GetTimeNow() *time.Time {
	t := time.Now()
	log.WithFields(log.Fields{
		"time": t,
	}).Debug("current time:")

	return &t
}
