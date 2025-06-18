package streaming

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

type backoffParams struct {
	attempts    int
	lastAttempt time.Time
}

func defaultBackoffHandler(attempts int) time.Duration {
	randInt, err := rand.Int(rand.Reader, big.NewInt(200))
	if err != nil {
		panic(err)
	}
	jitter := time.Duration(randInt.Int64()-100) * time.Millisecond                      // ±100ms
	backoff := time.Duration(math.Min(math.Pow(2, float64(attempts)), 60)) * time.Second // Exponential backoff up to 60s
	return backoff + jitter
}
