package timex

import "time"

func GetCurrentUTCUnixNano() int64 {
	return time.Now().UTC().UnixNano()
}

func ThirtyMinutes() time.Duration {
	return time.Duration(30) * time.Minute
}

func GetFutureUTCUnixNano(duration time.Duration) int64 {
	return time.Now().UTC().Add(duration).UnixNano()
}
