package streaming

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"sort"
	"sync"
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
	return ol
}

type UpdateCallback func(Update)

type messageProcessor struct {
	seq  int64
	bids map[string]order
	asks map[string]order

	mu sync.Mutex

	updateCallback UpdateCallback
}

func (m *messageProcessor) Reset() {
	m.mu.Lock()
	m.seq = 0
	m.bids = nil
	m.asks = nil
	m.mu.Unlock()
}

func (m *messageProcessor) HandleMessage(message []byte) error {
	if string(message) == "\"\"" {
		return nil
	}

	var ob orderBook
	if err := json.Unmarshal(message, &ob); err != nil {
		return err
	}
	if ob.Asks != nil || ob.Bids != nil {
		if err := m.receivedOrderBook(ob); err != nil {
			return err
		}
		return nil
	}

	var u Update
	if err := json.Unmarshal(message, &u); err != nil {
		return err
	}
	if err := m.receivedUpdate(u); err != nil {
		return err
	}

	return nil
}

func (m *messageProcessor) receivedOrderBook(ob orderBook) error {
	bids, err := convertOrders(ob.Bids)
	if err != nil {
		return err
	}

	asks, err := convertOrders(ob.Asks)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.seq = ob.Sequence
	m.bids = bids
	m.asks = asks
	return nil
}

func (m *messageProcessor) receivedUpdate(u Update) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.seq == 0 {
		// State not initialized so we can't update it.
		return nil
	}

	if u.Sequence <= m.seq {
		// Old update. We can just discard it.
		return nil
	}
	if u.Sequence != m.seq+1 {
		return errors.New("streaming: update received out of sequence")
	}

	// Process trades
	for _, t := range u.TradeUpdates {
		if err := m.processTrade(*t); err != nil {
			return err
		}
	}

	// Process create
	if u.CreateUpdate != nil {
		if err := m.processCreate(*u.CreateUpdate); err != nil {
			return err
		}
	}

	// Process delete
	if u.DeleteUpdate != nil {
		if err := m.processDelete(*u.DeleteUpdate); err != nil {
			return err
		}
	}

	m.seq = u.Sequence

	if m.updateCallback != nil {
		m.updateCallback(u)
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

func (m *messageProcessor) processTrade(t TradeUpdate) error {
	if t.Base.Sign() <= 0 {
		return errors.New("streaming: nonpositive trade")
	}

	ok, err := decTrade(m.bids, t.OrderID, t.Base)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	ok, err = decTrade(m.asks, t.OrderID, t.Base)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	return errors.New("streaming: trade for unknown order")
}

func (m *messageProcessor) processCreate(u CreateUpdate) error {
	o := order{
		ID:     u.OrderID,
		Price:  u.Price,
		Volume: u.Volume,
	}

	if u.Type == string(luno.OrderTypeBid) {
		m.bids[o.ID] = o
	} else if u.Type == string(luno.OrderTypeAsk) {
		m.asks[o.ID] = o
	} else {
		return errors.New("streaming: unknown order type")
	}

	return nil
}

func (m *messageProcessor) processDelete(u DeleteUpdate) error {
	delete(m.bids, u.OrderID)
	delete(m.asks, u.OrderID)
	return nil
}

// OrderBookSnapshot returns the latest order book.
func (m *messageProcessor) OrderBookSnapshot() (
	int64, []luno.OrderBookEntry, []luno.OrderBookEntry) {

	m.mu.Lock()
	defer m.mu.Unlock()

	bids := flatten(m.bids, true)
	asks := flatten(m.asks, false)
	return m.seq, bids, asks
}
