package streaming

import (
	"math"
	"math/rand"
	"time"
)

type backoffParams struct {
	attempts    int
	lastAttempt time.Time
}

func defaultBackoffHandler(attempts int) time.Duration {
	jitter := time.Duration(rand.Intn(200)-100) * time.Millisecond                       // ±100ms
	backoff := time.Duration(math.Min(math.Pow(2, float64(attempts)), 60)) * time.Second // Exponential backoff up to 60s
	return backoff + jitter
}
