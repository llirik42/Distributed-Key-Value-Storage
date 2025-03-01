package loops

import "time"

func getDurationMs(ms int) time.Duration {
	return time.Duration(ms) * time.Millisecond
}
