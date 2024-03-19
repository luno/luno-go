package streaming

import (
	"testing"
	"time"

	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

type expected struct {
	wantErr   bool
	asks      map[string]order
	bids      map[string]order
	lastTrade TradeUpdate
	seq       int64
	status    luno.Status
}

func TestOrderBook(t *testing.T) {
	type args struct {
		orders orderBook
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "success orderBook",
			args: args{
				orders: book(),
			},
			expected: expected{
				asks:   asksMap(),
				bids:   bidsMap(),
				seq:    1,
				status: luno.StatusActive,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{
				seq:  0,
				bids: nil,
				asks: nil,
			}
			err := c.receivedOrderBook(tt.args.orders)
			validateResult(err, t, tt.expected, c)
		})
	}
}

func TestReceivedUpdate(t *testing.T) {
	type args struct {
		u Update
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "error received out of sync",
			args: args{
				u: Update{Sequence: 3},
			},
			expected: expected{
				wantErr: true,
			},
		},
		{
			name: "success skip old seq",
			args: args{
				u: Update{Sequence: 1},
			},
			expected: expected{
				wantErr: false,
				asks:    asksMap(),
				bids:    bidsMap(),
				seq:     1,
				status:  luno.StatusActive,
			},
		},
		{
			name: "success process twice",
			args: args{
				u: Update{
					Sequence: 2,
					TradeUpdates: []*TradeUpdate{
						{
							Sequence:     1,
							Base:         decimal.NewFromFloat64(0.02, 2),
							Counter:      decimal.NewFromFloat64(0.002, 2),
							MakerOrderID: "1",
							TakerOrderID: "32",
						},
						{
							Sequence:     2,
							Base:         decimal.NewFromFloat64(0.01, 2),
							Counter:      decimal.NewFromFloat64(0.001, 2),
							MakerOrderID: "1",
							TakerOrderID: "34",
						},
					},
				},
			},
			expected: expected{
				wantErr: false,
				asks:    asksMap(),
				bids: bidsMap(order{
					ID: "1", Price: decimal.NewFromFloat64(120.0, 1),
					Volume: decimal.NewFromFloat64(0.07, 2),
				}),
				lastTrade: TradeUpdate{
					Sequence:     2,
					Base:         decimal.NewFromFloat64(0.01, 2),
					Counter:      decimal.NewFromFloat64(0.001, 2),
					MakerOrderID: "1",
					TakerOrderID: "34",
				},
				seq:    2,
				status: luno.StatusActive,
			},
		},
		{
			name: "success ask/bid",
			args: args{
				u: Update{
					Sequence: 2,
					TradeUpdates: []*TradeUpdate{
						{
							Sequence:     1,
							Base:         decimal.NewFromFloat64(0.01, 2),
							Counter:      decimal.NewFromFloat64(0.001, 2),
							MakerOrderID: "4",
							TakerOrderID: "32",
						},
						{
							Sequence:     2,
							Base:         decimal.NewFromFloat64(0.01, 2),
							Counter:      decimal.NewFromFloat64(0.001, 2),
							MakerOrderID: "3",
							TakerOrderID: "34",
						},
					},
				},
			},
			expected: expected{
				wantErr: false,
				asks: asksMap(order{
					ID: "4", Price: decimal.NewFromFloat64(180.0, 1),
					Volume: decimal.NewFromFloat64(0.99, 2),
				}),
				bids: bidsMap(order{
					ID: "3", Price: decimal.NewFromFloat64(100.0, 1),
					Volume: decimal.NewFromFloat64(0.99, 2),
				}),
				lastTrade: TradeUpdate{
					Sequence:     2,
					Base:         decimal.NewFromFloat64(0.01, 2),
					Counter:      decimal.NewFromFloat64(0.001, 2),
					MakerOrderID: "3",
					TakerOrderID: "34",
				},
				seq:    2,
				status: luno.StatusActive,
			},
		},
		{
			name: "success delete from ask",
			args: args{
				u: Update{
					Sequence: 2,
					TradeUpdates: []*TradeUpdate{
						{
							Sequence:     1,
							Base:         decimal.NewFromFloat64(0.5, 1),
							Counter:      decimal.NewFromFloat64(0.1, 2),
							MakerOrderID: "2",
							TakerOrderID: "32",
						},
						{
							Sequence:     2,
							Base:         decimal.NewFromFloat64(1, 1),
							Counter:      decimal.NewFromFloat64(1, 1),
							MakerOrderID: "4",
							TakerOrderID: "34",
						},
						{
							Sequence:     3,
							Base:         decimal.NewFromFloat64(0.1, 1),
							Counter:      decimal.NewFromFloat64(1, 1),
							MakerOrderID: "6",
							TakerOrderID: "36",
						},
					},
				},
			},
			expected: expected{
				wantErr: false,
				bids:    bidsMap(),
				lastTrade: TradeUpdate{
					Sequence:     3,
					Base:         decimal.NewFromFloat64(0.1, 1),
					Counter:      decimal.NewFromFloat64(1, 1),
					MakerOrderID: "6",
					TakerOrderID: "36",
				},
				seq:    2,
				status: luno.StatusActive,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{
				asks:   asksMap(),
				bids:   bidsMap(),
				seq:    1,
				status: luno.StatusActive,
			}
			err := c.receivedUpdate(tt.args.u)
			if tt.expected.wantErr {
				if err == nil {
					t.Error("error expected but nil received")
				}
				return
			}
			validateResult(err, t, tt.expected, c)
		})
	}
}

func TestReceivedCreate(t *testing.T) {
	type args struct {
		u Update
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "fail invalid type",
			args: args{
				u: Update{
					Sequence: 2,
					CreateUpdate: &CreateUpdate{
						OrderID: "8",
						Price:   decimal.NewFromFloat64(6700.56, 2),
						Type:    "ASKBID",
						Volume:  decimal.NewFromFloat64(0.01, 2),
					},
				},
			},
			expected: expected{
				wantErr: true,
			},
		},
		{
			name: "create ask",
			args: args{
				u: Update{
					Sequence: 2,
					CreateUpdate: &CreateUpdate{
						OrderID: "8",
						Price:   decimal.NewFromFloat64(6700.56, 2),
						Type:    "ASK",
						Volume:  decimal.NewFromFloat64(0.01, 2),
					},
				},
			},
			expected: expected{
				wantErr: false,
				asks: asksMap(order{
					ID:     "8",
					Price:  decimal.NewFromFloat64(6700.56, 2),
					Volume: decimal.NewFromFloat64(0.01, 2),
				}),
				bids:   bidsMap(),
				seq:    2,
				status: luno.StatusActive,
			},
		},
		{
			name: "create bid",
			args: args{
				u: Update{
					Sequence: 2,
					CreateUpdate: &CreateUpdate{
						OrderID: "7",
						Price:   decimal.NewFromFloat64(6700.54, 2),
						Type:    "BID",
						Volume:  decimal.NewFromFloat64(0.01, 2),
					},
				},
			},
			expected: expected{
				wantErr: false,
				asks:    asksMap(),
				bids: bidsMap(order{
					ID:     "7",
					Price:  decimal.NewFromFloat64(6700.54, 2),
					Volume: decimal.NewFromFloat64(0.01, 2),
				}),
				seq:    2,
				status: luno.StatusActive,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{
				asks:   asksMap(),
				bids:   bidsMap(),
				seq:    1,
				status: luno.StatusActive,
			}
			err := c.receivedUpdate(tt.args.u)
			if tt.expected.wantErr {
				if err == nil {
					t.Error("error expected but nil received")
				}
				return
			}
			validateResult(err, t, tt.expected, c)
		})
	}
}

func TestReceivedDelete(t *testing.T) {
	type args struct {
		u Update
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "success not in asks/bids",
			args: args{
				u: Update{
					Sequence:     2,
					DeleteUpdate: &DeleteUpdate{OrderID: "8"},
				},
			},
			expected: expected{
				wantErr: false,
				asks:    asksMap(),
				bids:    bidsMap(),
				seq:     2,
				status:  luno.StatusActive,
			},
		},
		{
			name: "delete ask",
			args: args{
				u: Update{
					Sequence:     2,
					DeleteUpdate: &DeleteUpdate{OrderID: "4"},
				},
			},
			expected: expected{
				wantErr: false,
				asks: map[string]order{
					"2": {
						ID: "2", Price: decimal.NewFromFloat64(150.0, 1),
						Volume: decimal.NewFromFloat64(0.5, 1),
					},
					"6": {
						ID: "6", Price: decimal.NewFromFloat64(200.0, 1),
						Volume: decimal.NewFromFloat64(0.1, 1),
					},
				},
				bids:   bidsMap(),
				seq:    2,
				status: luno.StatusActive,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{
				asks:   asksMap(),
				bids:   bidsMap(),
				seq:    1,
				status: luno.StatusActive,
			}
			err := c.receivedUpdate(tt.args.u)
			if tt.expected.wantErr {
				if err == nil {
					t.Error("error expected but nil received")
				}
				return
			}
			validateResult(err, t, tt.expected, c)
		})
	}
}

func bidsMap(o ...order) map[string]order {
	res := map[string]order{
		"1": {
			ID: "1", Price: decimal.NewFromFloat64(120.0, 1),
			Volume: decimal.NewFromFloat64(0.1, 1),
		},
		"3": {
			ID: "3", Price: decimal.NewFromFloat64(100.0, 1),
			Volume: decimal.NewFromFloat64(1.0, 1),
		},
		"5": {
			ID: "5", Price: decimal.NewFromFloat64(110.0, 1),
			Volume: decimal.NewFromFloat64(0.5, 1),
		},
	}
	for _, v := range o {
		res[v.ID] = v
	}
	return res
}

func asksMap(o ...order) map[string]order {
	res := map[string]order{
		"2": {
			ID: "2", Price: decimal.NewFromFloat64(150.0, 1),
			Volume: decimal.NewFromFloat64(0.5, 1),
		},
		"4": {
			ID: "4", Price: decimal.NewFromFloat64(180.0, 1),
			Volume: decimal.NewFromFloat64(1.0, 1),
		},
		"6": {
			ID: "6", Price: decimal.NewFromFloat64(200.0, 1),
			Volume: decimal.NewFromFloat64(0.1, 1),
		},
	}
	for _, v := range o {
		res[v.ID] = v
	}
	return res
}

func compareOrderMaps(a, b map[string]order) bool {
	if &a == &b {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	return existsInOtherMap(b, a)
}

func existsInOtherMap(a, b map[string]order) bool {
	for kb, vb := range b {
		va, ok := a[kb]
		if !ok {
			return false
		}
		if va.ID != vb.ID {
			return false
		}
		if va.Price.Cmp(vb.Price) != 0 {
			return false
		}
		if va.Volume.Cmp(vb.Volume) != 0 {
			return false
		}
	}
	return true
}

func compareLastTrade(a, b TradeUpdate) bool {
	if a.Base.Cmp(b.Base) != 0 {
		return false
	}
	a.Base = b.Base
	if a.Counter.Cmp(b.Counter) != 0 {
		return false
	}
	a.Counter = b.Counter

	return a == b
}

func validateResult(err error, t *testing.T, exp expected, c *Conn) {
	if err != nil {
		t.Errorf("Expected success got: %v", err)
	}
	if !compareOrderMaps(exp.asks, c.asks) {
		t.Errorf("Invalid asks. Expected:%v, got:%v", exp.asks, c.asks)
	}
	if !compareOrderMaps(exp.bids, c.bids) {
		t.Errorf("Invalid bids. Expected:%v, got:%v", exp.bids, c.bids)
	}
	if !compareLastTrade(exp.lastTrade, c.lastTrade) {
		t.Errorf("Invalid lastTrade. Expected:%v, got:%v", exp.lastTrade, c.lastTrade)
	}
	if exp.seq != c.seq {
		t.Errorf("Invalid seq. Expected:%v, got:%v", exp.seq, c.seq)
	}
	if exp.status != c.status {
		t.Errorf("Invalid status. Expected:%v, got:%v", exp.status, c.status)
	}
}

func book() orderBook {
	return orderBook{
		Bids: []*order{
			{
				ID: "1", Price: decimal.NewFromFloat64(120.0, 1),
				Volume: decimal.NewFromFloat64(0.1, 1),
			},
			{
				ID: "5", Price: decimal.NewFromFloat64(110.0, 1),
				Volume: decimal.NewFromFloat64(0.5, 1),
			},
			{
				ID: "3", Price: decimal.NewFromFloat64(100.0, 1),
				Volume: decimal.NewFromFloat64(1.0, 1),
			},
		},
		Asks: []*order{
			{
				ID: "2", Price: decimal.NewFromFloat64(150.0, 1),
				Volume: decimal.NewFromFloat64(0.5, 1),
			},
			{
				ID: "4", Price: decimal.NewFromFloat64(180.0, 1),
				Volume: decimal.NewFromFloat64(1.0, 1),
			},
			{
				ID: "6", Price: decimal.NewFromFloat64(200.0, 1),
				Volume: decimal.NewFromFloat64(0.1, 1),
			},
		},
		Sequence: 1,
		Status:   luno.StatusActive,
	}
}

func Test_flatten(t *testing.T) {
	orders := map[string]order{
		"1": {ID: "1", Price: decimal.NewFromInt64(1), Volume: decimal.NewFromInt64(1)},
		"2": {ID: "2", Price: decimal.NewFromInt64(2), Volume: decimal.NewFromInt64(2)},
		"3": {ID: "3", Price: decimal.NewFromInt64(3), Volume: decimal.NewFromInt64(3)},
	}
	expForward := []luno.OrderBookEntry{
		{Price: decimal.NewFromInt64(1), Volume: decimal.NewFromInt64(1)},
		{Price: decimal.NewFromInt64(2), Volume: decimal.NewFromInt64(2)},
		{Price: decimal.NewFromInt64(3), Volume: decimal.NewFromInt64(3)},
	}
	expReverse := []luno.OrderBookEntry{
		{Price: decimal.NewFromInt64(3), Volume: decimal.NewFromInt64(3)},
		{Price: decimal.NewFromInt64(2), Volume: decimal.NewFromInt64(2)},
		{Price: decimal.NewFromInt64(1), Volume: decimal.NewFromInt64(1)},
	}

	forward := flatten(orders, false)
	for i := 0; i < len(orders); i++ {
		compareOrderBookEntry(t, expForward[i], forward[i])
	}

	reverse := flatten(orders, true)
	for i := 0; i < len(orders); i++ {
		compareOrderBookEntry(t, expReverse[i], reverse[i])
	}
}

func compareOrderBookEntry(t *testing.T, want, got luno.OrderBookEntry) {
	if got.Price.Cmp(want.Price) != 0 {
		t.Errorf("Price = %v, want %v", got.Price, want.Price)
	}
	got.Price = want.Price
	if got.Volume.Cmp(want.Volume) != 0 {
		t.Errorf("Volume = %v, want %v", got.Volume, want.Volume)
	}
	got.Volume = want.Volume

	if got != want {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestReceiveUpdateSnapshot(t *testing.T) {
	chDone := make(chan struct{})

	c := &Conn{
		asks:   asksMap(),
		bids:   bidsMap(),
		seq:    1,
		status: luno.StatusActive,
	}

	onUpdate := func(up Update) {
		// Get snapshot to confirm mutex does not create deadlock
		_ = c.Snapshot()
		chDone <- struct{}{}
	}

	c.updateCallback = onUpdate

	tu := []*TradeUpdate{
		{
			Sequence:     2,
			Base:         decimal.NewFromFloat64(0.02, 2),
			Counter:      decimal.NewFromFloat64(0.002, 2),
			MakerOrderID: "1",
			TakerOrderID: "32",
		},
		{
			Sequence:     3,
			Base:         decimal.NewFromFloat64(0.01, 2),
			Counter:      decimal.NewFromFloat64(0.001, 2),
			MakerOrderID: "1",
			TakerOrderID: "34",
		},
	}

	go func() {
		err := c.receivedUpdate(Update{Sequence: 2, TradeUpdates: tu})
		if err != nil {
			t.Errorf("Expected success got: %v", err)
		}
	}()

	select {
	case <-chDone:
	case <-time.After(time.Second):
		t.Errorf("timeout trying to retrieve snapshot on update")
	}
}
