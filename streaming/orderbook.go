package streaming

import (
	"errors"
	"fmt"
	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"sync"
)

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

type orderbookState struct {
	seq  int64
	bids map[string]order
	asks map[string]order

	mu sync.Mutex
}

func (ob *orderbookState) Reset() {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	ob.seq = 0
	ob.bids = nil
	ob.asks = nil
}

func convertOrders(ol []*order) map[string]order {
	r := make(map[string]order)
	for _, o := range ol {
		r[o.ID] = *o
	}
	return r
}

func (ob *orderbookState) Set(sequence int64, bids []*order, asks []*order) {
	bidsMap := convertOrders(bids)
	asksMap := convertOrders(asks)

	ob.mu.Lock()
	defer ob.mu.Unlock()

	ob.seq = sequence
	ob.bids = bidsMap
	ob.asks = asksMap
}

func (ob *orderbookState) Lock() {
	ob.mu.Lock()
}

func (ob *orderbookState) Unlock() {
	ob.mu.Unlock()
}

func (ob *orderbookState) GetStateId() int64 {
	return ob.seq
}

func (ob *orderbookState) SetStateId(seq int64) {
	ob.seq = seq
}

func (ob *orderbookState) DecrementOrder(id string, base decimal.Decimal) error {
	ok, err := decTrade(ob.bids, id, base)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	ok, err = decTrade(ob.asks, id, base)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	return errors.New("streaming: trade for unknown order")
}

func (ob *orderbookState) AddOrder(ordertype luno.OrderType, order order) {
	if ordertype == luno.OrderTypeBid {
		ob.bids[order.ID] = order
	} else if ordertype == luno.OrderTypeAsk {
		ob.asks[order.ID] = order
	}
}

func (ob *orderbookState) RemoveOrder(id string) {
	delete(ob.bids, id)
	delete(ob.asks, id)
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

// OrderBookSnapshot returns the latest order book.
func (ob *orderbookState) GetSnapshot() (
	int64, []luno.OrderBookEntry, []luno.OrderBookEntry) {

	ob.mu.Lock()
	defer ob.mu.Unlock()

	bids := flatten(ob.bids, true)
	asks := flatten(ob.asks, false)
	return ob.seq, bids, asks
}
