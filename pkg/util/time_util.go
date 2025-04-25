package util

import "time"

func IsWithinRange(checkTime, startTime, endTime time.Time) bool {
	return (checkTime.After(startTime) || checkTime.Equal(startTime)) && (checkTime.Before(endTime) || checkTime.Equal(endTime))
}
