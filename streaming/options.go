package streaming

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
