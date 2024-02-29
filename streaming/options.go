package streaming

import (
	"time"
)

type DialOption func(*Conn)

// WithUpdateCallback returns an options which sets a callback function for
// streaming updates. Each update will first be applied to the order book, and
// then passed to the callback function.
func WithUpdateCallback(fn UpdateCallback) DialOption {
	return func(c *Conn) {
		c.updateCallback = fn
	}
}

// WithConnectCallback returns an options which sets a callback function for
// when the connection is fully initialised and the orderbook has been set up.
func WithConnectCallback(fn ConnectCallback) DialOption {
	return func(c *Conn) {
		c.connectCallback = fn
	}
}

// WithBackoffHandler specifies a custom handler to calculate backoff duration after each disconnect. Attempt will increment
// with each subsequent call until the attemptReset duration exceeds the duration since the last disconnect, at which point it
// will reset to 0.
func WithBackoffHandler(fn BackoffHandler, attemptReset time.Duration) DialOption {
	return func(c *Conn) {
		c.backoffHandler = fn
		c.attemptReset = attemptReset
	}
}
