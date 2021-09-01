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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

const websocketTimeout = time.Minute

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
	return ol
}

type ConnectCallback func(*Conn)
type UpdateCallback func(Update)

type Conn struct {
	keyID, keySecret string
	pair             string
	connectCallback  ConnectCallback
	updateCallback   UpdateCallback

	ws     *websocket.Conn
	closed bool

	seq  int64
	bids map[string]order
	asks map[string]order

	status luno.Status

	lastMessage time.Time
	lastTrade   TradeUpdate

	mu sync.Mutex
}

// Dial initiates a connection to the streaming service and starts processing
// data for the given market pair.
// The connection will automatically reconnect on error.
func Dial(keyID, keySecret, pair string, opts ...DialOption) (*Conn, error) {
	if keyID == "" || keySecret == "" {
		return nil, errors.New("streaming: streaming API requires credentials")
	}

	c := &Conn{
		keyID:     keyID,
		keySecret: keySecret,
		pair:      pair,
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
	attempts := 0
	var lastAttempt time.Time
	for {
		c.mu.Lock()
		closed := c.closed
		c.mu.Unlock()
		if closed {
			return
		}

		lastAttempt = time.Now()
		attempts++
		if err := c.connect(); err != nil {
			log.Printf("luno/streaming: Connection error key=%s pair=%s: %v",
				c.keyID, c.pair, err)
		}

		if time.Now().Sub(lastAttempt) > time.Hour {
			attempts = 0
		}
		if attempts > 5 {
			attempts = 5
		}
		wait := 5
		for i := 0; i < attempts; i++ {
			wait = 2 * wait
		}
		wait = wait + rand.Intn(wait)
		dt := time.Duration(wait) * time.Second
		log.Printf("luno/streaming: Waiting %s before reconnecting", dt)
		time.Sleep(dt)
	}
}

func (c *Conn) connect() error {
	url := *wsHost + "/api/1/stream/" + c.pair
	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		return fmt.Errorf("unable to dial server: %w", err)
	}
	defer func() {
		ws.Close()
		c.mu.Lock()
		c.ws = nil
		c.seq = 0
		c.bids = nil
		c.asks = nil
		c.status = ""
		c.mu.Unlock()
	}()

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	} else {
		c.ws = ws
		c.mu.Unlock()
	}

	cred := credentials{c.keyID, c.keySecret}
	if err := websocket.JSON.Send(ws, cred); err != nil {
		return fmt.Errorf("failed to send credentials: %w", err)
	}

	log.Printf("luno/streaming: Connection established key=%s pair=%s",
		c.keyID, c.pair)

	go sendPings(ws)

	for {
		var data []byte
		c.ws.SetReadDeadline(time.Now().Add(websocketTimeout))
		err := websocket.Message.Receive(c.ws, &data)
		if errors.Is(err, io.EOF) {
			// Server closed the connection. Return gracefully.
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to receive message: %w", err)
		}

		if string(data) == "\"\"" {
			c.receivedPing()
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

func sendPings(ws *websocket.Conn) {
	defer ws.Close()
	for {
		ws.SetWriteDeadline(time.Now().Add(websocketTimeout))
		if err := websocket.Message.Send(ws, ""); err != nil {
			log.Printf("luno/streaming: Failed to ping server: %v", err)
			return
		}
		time.Sleep(time.Minute)
	}
}

func (c *Conn) receivedPing() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastMessage = time.Now()
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

	c.lastMessage = time.Now()
	c.seq = ob.Sequence
	c.status = ob.Status
	c.bids = bids
	c.asks = asks
	return nil
}

func (c *Conn) receivedUpdate(u Update) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.seq == 0 {
		// State not initialized so we can't update it.
		return nil
	}

	if u.Sequence <= c.seq {
		// Old update. We can just discard it.
		return nil
	}
	if u.Sequence != c.seq+1 {
		return errors.New("streaming: update received out of sequence")
	}

	// Process trades
	for _, t := range u.TradeUpdates {
		if err := c.processTrade(*t); err != nil {
			return err
		}
	}

	// Process create
	if u.CreateUpdate != nil {
		if err := c.processCreate(*u.CreateUpdate); err != nil {
			return err
		}
	}

	// Process delete
	if u.DeleteUpdate != nil {
		if err := c.processDelete(*u.DeleteUpdate); err != nil {
			return err
		}
	}

	// Process status
	if u.StatusUpdate != nil {
		if err := c.processStatus(*u.StatusUpdate); err != nil {
			return err
		}
	}

	c.lastMessage = time.Now()
	c.seq = u.Sequence

	if c.updateCallback != nil {
		c.updateCallback(u)
	}

	return nil
}

func decTrade(m map[string]order, id string, base decimal.Decimal) (
	bool, error) {

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
		Base:    t.Base,
		Counter: t.Counter,
	}

	ok, err := decTrade(c.bids, t.OrderID, t.Base)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	ok, err = decTrade(c.asks, t.OrderID, t.Base)
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
	int64, []luno.OrderBookEntry, []luno.OrderBookEntry) {

	c.mu.Lock()
	defer c.mu.Unlock()

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
	c.mu.Lock()
	defer c.mu.Unlock()

	return Snapshot{
		Sequence:  c.seq,
		Bids:      flatten(c.bids, true),
		Asks:      flatten(c.asks, false),
		Status:    c.status,
		LastTrade: c.lastTrade,
	}
}

// Status returns the currenct status of the streaming connection.
func (c *Conn) Status() luno.Status {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.status
}

// Close the connection.
func (c *Conn) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.closed = true
	if c.ws != nil {
		c.ws.Close()
	}
}
