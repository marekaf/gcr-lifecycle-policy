package utils

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// GetLastBillingPeriodStart return the first day of last month in format YYYY-DD-MM
func GetLastBillingPeriodStart() string {
	t := time.Now()
	log.WithFields(log.Fields{
		"time": t,
	}).Debug("current time:")

	// this year, last month, first day of the month, hour, min, sec, nanosec, location
	d := time.Date(t.Year(), t.Month()-1, 1, 12, 30, 0, 0, t.Location())

	log.WithFields(log.Fields{
		"time": d,
	}).Debug("start of last month:")

	// this is the **magical reference date**
	return d.Format("2006-01-02")
}

// GetLastBillingPeriodEnd return the first day of this month (considered the end of billing period) in format YYYY-DD-MM
func GetLastBillingPeriodEnd() string {
	t := time.Now()
	log.WithFields(log.Fields{
		"time": t,
	}).Debug("current time:")

	// this year, last month, first day of the month, hour, min, sec, nanosec, location
	d := time.Date(t.Year(), t.Month(), 1, 12, 30, 0, 0, t.Location())

	log.WithFields(log.Fields{
		"time": d,
	}).Debug("end of last month:")

	// this is the **magical reference date**
	return d.Format("2006-01-02")
}

// GetTimeTwoWeeksAgoStart returns timestamp two weeks ago
func GetTimeTwoWeeksAgoStart() *time.Time {
	t := time.Now()
	log.WithFields(log.Fields{
		"time": t,
	}).Debug("current time:")

	// minus 14 days
	d := t.AddDate(0, 0, -14)

	log.WithFields(log.Fields{
		"time": d,
	}).Debug("start of last month:")

	return &d
}

// GetTimeNow returns timestamp two weeks ago
func GetTimeNow() *time.Time {
	t := time.Now()
	log.WithFields(log.Fields{
		"time": t,
	}).Debug("current time:")

	return &t
}
