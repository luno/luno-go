/*
Package streaming implements a client for the Luno Streaming API.

Example:

	c, err := streaming.Dial(keyID, keySecret, "XBTZAR")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for {
		seq, bids, asks := c.OrderBookSnapshot()
		log.Printf("%d: %v %v\n", seq, bids[0], asks[0])
		time.Sleep(time.Minute)
	}
*/
package streaming

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

const (
	readTimeout         = time.Minute
	writeTimeout        = 30 * time.Second
	pingInterval        = 20 * time.Second
	defaultAttemptReset = time.Minute * 30
)

func convertOrders(ol []*order) (map[string]order, error) {
	r := make(map[string]order)
	for _, o := range ol {
		r[o.ID] = *o
	}
	return r, nil
}

type orderList []luno.OrderBookEntry

func (ol orderList) Less(i, j int) bool {
	return ol[i].Price.Cmp(ol[j].Price) < 0
}

func (ol orderList) Swap(i, j int) {
	ol[i], ol[j] = ol[j], ol[i]
}

func (ol orderList) Len() int {
	return len(ol)
}

func flatten(m map[string]order, reverse bool) []luno.OrderBookEntry {
	var ol []luno.OrderBookEntry
	for _, o := range m {
		ol = append(ol, luno.OrderBookEntry{
			Price:  o.Price,
			Volume: o.Volume,
		})
	}
	if reverse {
		sort.Sort(sort.Reverse(orderList(ol)))
	} else {
		sort.Sort(orderList(ol))
	}
	ret := make([]luno.OrderBookEntry, 0, len(ol))
	for _, o := range ol {
		if len(ret) == 0 {
			ret = append(ret, o)
			continue
		}
		lastIdx := len(ret) - 1
		if o.Price.Cmp(ret[lastIdx].Price) == 0 {
			ret[lastIdx].Volume = ret[lastIdx].Volume.Add(o.Volume)
			continue
		}
		ret = append(ret, o)
	}
	return ret
}

type (
	ConnectCallback func(*Conn)
	UpdateCallback  func(Update)
	BackoffHandler  func(attempt int) time.Duration
)

type Conn struct {
	keyID, keySecret string
	pair             string
	connectCallback  ConnectCallback
	updateCallback   UpdateCallback

	backoffHandler BackoffHandler
	attemptReset   time.Duration

	closed bool

	seq  int64
	bids map[string]order
	asks map[string]order

	status luno.Status

	lastTrade TradeUpdate

	mu sync.RWMutex
}

// Dial initiates a connection to the streaming service and starts processing
// data for the given market pair.
// The connection will automatically reconnect on error.
func Dial(keyID, keySecret, pair string, opts ...DialOption) (*Conn, error) {
	if keyID == "" || keySecret == "" {
		return nil, errors.New("streaming: streaming API requires credentials")
	}

	c := &Conn{
		keyID:        keyID,
		keySecret:    keySecret,
		pair:         pair,
		attemptReset: defaultAttemptReset,
	}
	for _, opt := range opts {
		opt(c)
	}

	go c.manageForever()
	return c, nil
}

var wsHost = flag.String(
	"luno_websocket_host", "wss://ws.luno.com", "Luno API websocket host")

func (c *Conn) manageForever() {
	p := &backoffParams{}

	for {
		if err := c.connect(); err != nil {
			log.Printf("luno/streaming: Connection error key=%s pair=%s: %v",
				c.keyID, c.pair, err)
		}
		if c.IsClosed() {
			return
		}

		dt := c.calculateBackoff(p, time.Now())

		log.Printf("luno/streaming: Waiting %s before reconnecting", dt)
		time.Sleep(dt)
	}
}

func (c *Conn) calculateBackoff(p *backoffParams, ts time.Time) time.Duration {
	if ts.Sub(p.lastAttempt) >= c.attemptReset {
		p.attempts = 0
	}

	p.attempts++

	backoff := defaultBackoffHandler
	if c.backoffHandler != nil {
		backoff = c.backoffHandler
	}

	p.lastAttempt = ts

	return backoff(p.attempts)
}

func (c *Conn) connect() error {
	url := *wsHost + "/api/1/stream/" + c.pair
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("unable to dial server: %w", err)
	}
	defer func() {
		_ = ws.Close()
		c.reset()
	}()

	cred := credentials{c.keyID, c.keySecret}
	_ = ws.SetWriteDeadline(time.Now().Add(writeTimeout)) // Safe to ignore error as it's always nil.
	if err := ws.WriteJSON(cred); err != nil {
		return fmt.Errorf("failed to send credentials: %w", err)
	}

	log.Printf("luno/streaming: Connection established key=%s pair=%s",
		c.keyID, c.pair)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go c.sendPings(ctx, ws)

	for {
		if c.IsClosed() {
			return nil
		}

		_, data, err := ws.ReadMessage()
		if err != nil {
			return err
		}
		if errors.Is(err, io.EOF) {
			// Server closed the connection. Return gracefully.
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to receive message: %w", err)
		}

		if string(data) == "\"\"" {
			// Ignore server keep alive messages
			continue
		}

		var ob orderBook
		if err := json.Unmarshal(data, &ob); err != nil {
			return fmt.Errorf("failed to unmarshal order book: %w", err)
		}
		if ob.Asks != nil || ob.Bids != nil {
			// Received an order book.
			if err := c.receivedOrderBook(ob); err != nil {
				return fmt.Errorf("failed to process order book: %w", err)
			}
			if c.connectCallback != nil {
				c.connectCallback(c)
			}
			continue
		}

		var u Update
		if err := json.Unmarshal(data, &u); err != nil {
			return fmt.Errorf("failed to unmarshal update: %w", err)
		}
		if err := c.receivedUpdate(u); err != nil {
			return fmt.Errorf("failed to process update: %w", err)
		}
	}
}

func (c *Conn) sendPings(ctx context.Context, ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	// Set initial read deadline
	_ = ws.SetReadDeadline(time.Now().Add(readTimeout))

	ws.SetPongHandler(func(data string) error {
		// Connection is alive, extend read deadline
		return ws.SetReadDeadline(time.Now().Add(readTimeout))
	})

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			if c.IsClosed() || ws == nil {
				return
			}

			_ = ws.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("luno/streaming: Failed to ping server: %v", err)
			}
		}
	}
}

func (c *Conn) receivedOrderBook(ob orderBook) error {
	bids, err := convertOrders(ob.Bids)
	if err != nil {
		return err
	}

	asks, err := convertOrders(ob.Asks)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.seq = ob.Sequence
	c.status = ob.Status
	c.bids = bids
	c.asks = asks
	return nil
}

func (c *Conn) receivedUpdate(u Update) error {
	valid, err := c.processUpdate(u)
	if err != nil {
		return err
	}

	// If update is not valid, ignore
	if !valid {
		return nil
	}

	if c.updateCallback != nil {
		c.updateCallback(u)
	}

	return nil
}

// Validate and process update into orderbook.
// Return bool indicating if update is valid.
func (c *Conn) processUpdate(u Update) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.seq == 0 {
		// State not initialized so we can't update it.
		return false, nil
	}

	if u.Sequence <= c.seq {
		// Old update. We can just discard it.
		return false, nil
	}

	if u.Sequence != c.seq+1 {
		return false, errors.New("streaming: update received out of sequence")
	}

	// Process trades
	for _, t := range u.TradeUpdates {
		if err := c.processTrade(*t); err != nil {
			return false, err
		}
	}

	// Process create
	if u.CreateUpdate != nil {
		if err := c.processCreate(*u.CreateUpdate); err != nil {
			return false, err
		}
	}

	// Process delete
	if u.DeleteUpdate != nil {
		if err := c.processDelete(*u.DeleteUpdate); err != nil {
			return false, err
		}
	}

	// Process status
	if u.StatusUpdate != nil {
		if err := c.processStatus(*u.StatusUpdate); err != nil {
			return false, err
		}
	}

	c.seq = u.Sequence

	return true, nil
}

func decTrade(m map[string]order, id string, base decimal.Decimal) (
	bool, error,
) {
	o, ok := m[id]
	if !ok {
		return false, nil
	}

	o.Volume = o.Volume.Add(base.Neg())

	if o.Volume.Sign() < 0 {
		return false, fmt.Errorf("streaming: negative volume: %s", o.Volume)
	}

	if o.Volume.Sign() == 0 {
		delete(m, id)
	} else {
		m[id] = o
	}
	return true, nil
}

func (c *Conn) processTrade(t TradeUpdate) error {
	if t.Base.Sign() <= 0 {
		return errors.New("streaming: nonpositive trade")
	}

	c.lastTrade = TradeUpdate{
		Sequence:     t.Sequence,
		Base:         t.Base,
		Counter:      t.Counter,
		MakerOrderID: t.MakerOrderID,
		TakerOrderID: t.TakerOrderID,
		OrderID:      t.OrderID,
	}

	ok, err := decTrade(c.bids, t.MakerOrderID, t.Base)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	ok, err = decTrade(c.asks, t.MakerOrderID, t.Base)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return errors.New("streaming: trade for unknown order")
}

func (c *Conn) processCreate(u CreateUpdate) error {
	o := order{
		ID:     u.OrderID,
		Price:  u.Price,
		Volume: u.Volume,
	}

	if u.Type == string(luno.OrderTypeBid) {
		c.bids[o.ID] = o
	} else if u.Type == string(luno.OrderTypeAsk) {
		c.asks[o.ID] = o
	} else {
		return errors.New("streaming: unknown order type")
	}

	return nil
}

func (c *Conn) processDelete(u DeleteUpdate) error {
	delete(c.bids, u.OrderID)
	delete(c.asks, u.OrderID)
	return nil
}

func (c *Conn) processStatus(u StatusUpdate) error {
	switch u.Status {
	case string(luno.StatusActive):
		c.status = luno.StatusActive
	case string(luno.StatusDisabled):
		c.status = luno.StatusDisabled
	case string(luno.StatusPostonly):
		c.status = luno.StatusPostonly
	default:
		return errors.New(fmt.Sprintf("Unknown status update: %s", u.Status))
	}
	return nil
}

// OrderBookSnapshot returns the latest order book.
// Deprecated at v0.0.8, use Snapshot().
func (c *Conn) OrderBookSnapshot() (
	int64, []luno.OrderBookEntry, []luno.OrderBookEntry,
) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	bids := flatten(c.bids, true)
	asks := flatten(c.asks, false)
	return c.seq, bids, asks
}

type Snapshot struct {
	Sequence   int64
	Bids, Asks []luno.OrderBookEntry
	Status     luno.Status
	LastTrade  TradeUpdate
}

// Snapshot returns the current state of the streamed data.
func (c *Conn) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Snapshot{
		Sequence:  c.seq,
		Bids:      flatten(c.bids, true),
		Asks:      flatten(c.asks, false),
		Status:    c.status,
		LastTrade: c.lastTrade,
	}
}

// Status returns the current status of the streaming connection.
func (c *Conn) Status() luno.Status {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.status
}

// Close the stream. After calling this the client will stop receiving new updates and the results of querying the Conn
// struct (Snapshot, Status...) will be zeroed values.
func (c *Conn) Close() {
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()

	c.reset()
}

func (c *Conn) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.seq = 0
	c.bids = nil
	c.asks = nil
	c.status = ""
}

// IsClosed returns true if the Conn has been closed.
func (c *Conn) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}
