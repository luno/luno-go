package streaming

type DialOption func(*Connection)

// WithUpdateCallback returns an options which sets a callback function for
// streaming updates. Each update will first be applied to the order book, and
// then passed to the callback function.
func WithUpdateCallback(fn UpdateCallback) DialOption {
	return func(c *Connection) {
		c.MessageProcessor.updateCallback = fn
	}
}
