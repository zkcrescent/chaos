package utils

import "time"

func ExitOnErr(err error, msg string, args ...interface{}) {
	args = append(args, err)
	if err != nil {
		log.Fatalf(msg+" , error: %v", args...)
	}
}

// Duration for duration in second
func Duration(d float64) time.Duration {
	return time.Duration(d*1000.0) * time.Millisecond
}
