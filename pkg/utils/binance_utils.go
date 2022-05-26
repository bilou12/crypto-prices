package utils

import "time"

func ConvertTimestampToDate(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, 0)
}
