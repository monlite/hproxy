package hproxy

import (
	"strconv"
	"time"
)

// Duration to e.g. 432ms or 12s, human readable translation
func BeautifyDuration(d time.Duration) string {
	u, ms, s := uint64(d), uint64(time.Millisecond), uint64(time.Second)
	if d < 0 {
		u = -u
	}
	switch {
	case u < ms:
		return "0"
	case u < s:
		return strconv.FormatUint(u/ms, 10) + "ms"
	default:
		return strconv.FormatUint(u/s, 10) + "s"
	}
}

func BeautifySize(s int64) string {
	switch {
	case s < 1<<10:
		return strconv.FormatInt(s, 10) + "B"
	case s < 1<<20:
		return strconv.FormatInt(s/(1<<10), 10) + "KB"
	default:
		return strconv.FormatInt(s/(1<<20), 10) + "MB"
	}
}
