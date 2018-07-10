package streaming

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/luno/luno-go"
)

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

type UpdateCallback func(UpdateMessage)

type messageProcessor struct {
	orderbook      orderbookState
	updateCallback UpdateCallback
}

func (m *messageProcessor) Reset() {
	m.orderbook.Reset()
}

func (m *messageProcessor) HandleMessage(message []byte) error {
	if string(message) == "\"\"" {
		return nil
	}

	var ob orderbookMessage
	if err := json.Unmarshal(message, &ob); err != nil {
		return err
	}
	if ob.Asks != nil || ob.Bids != nil {
		m.orderbook.Set(ob.Sequence, ob.Bids, ob.Asks)
		return nil
	}

	var u UpdateMessage
	if err := json.Unmarshal(message, &u); err != nil {
		return err
	}
	if err := m.receivedUpdate(u); err != nil {
		return err
	}

	return nil
}

func (m *messageProcessor) receivedUpdate(u UpdateMessage) error {
	m.orderbook.Lock()
	defer m.orderbook.Unlock()

	if m.orderbook.GetStateId() == 0 {
		return nil
	}

	if u.Sequence <= m.orderbook.GetStateId() {
		return nil
	}

	if u.Sequence != m.orderbook.GetStateId()+1 {
		return errors.New("streaming: update received out of sequence")
	}

	for _, t := range u.TradeUpdates {
		if err := m.processTrade(*t); err != nil {
			return err
		}
	}

	if u.CreateUpdate != nil {
		if err := m.processCreate(*u.CreateUpdate); err != nil {
			return err
		}
	}

	if u.DeleteUpdate != nil {
		m.orderbook.RemoveOrder(u.DeleteUpdate.OrderID)
	}

	m.orderbook.SetStateId(u.Sequence)

	if m.updateCallback != nil {
		go m.updateCallback(u)
	}

	return nil
}

func (m *messageProcessor) processTrade(t TradeUpdateMessage) error {
	if t.Base.Sign() <= 0 {
		return errors.New("streaming: nonpositive trade")
	}

	return m.orderbook.DecrementOrder(t.OrderID, t.Base)
}

func (m *messageProcessor) processCreate(u CreateUpdateMessage) error {
	o := order{
		ID:     u.OrderID,
		Price:  u.Price,
		Volume: u.Volume,
	}

	if u.Type != string(luno.OrderTypeBid) && u.Type != string(luno.OrderTypeAsk) {
		return errors.New("streaming: unknown order type")
	}

	m.orderbook.AddOrder(luno.OrderType(u.Type), o)

	return nil
}
