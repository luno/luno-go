package streaming

import (
	"testing"

	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func TestFlatten(t *testing.T) {
	testCases := []struct {
		name string

		orders map[string]order
		reverse bool

		expOrders []luno.OrderBookEntry
	}{
		{name: "empty orders"},
		{name: "single order",
			orders: map[string]order{
				"2": {
					ID:     "2",
					Price:  decimal.NewFromInt64(100),
					Volume: decimal.NewFromInt64( 1),
				},
			},
			expOrders: []luno.OrderBookEntry{
				{Price:decimal.NewFromInt64(100), Volume: decimal.NewFromInt64(1)},
			},
		},
		{name: "sorted orders",
			orders: map[string]order{
				"1": {
					ID:     "1",
					Price:  decimal.NewFromInt64(20),
					Volume: decimal.NewFromInt64( 2),
				},
				"3": {
					ID:     "3",
					Price:  decimal.NewFromInt64(40),
					Volume: decimal.NewFromInt64( 4),
				},
				"2": {
					ID:     "2",
					Price:  decimal.NewFromInt64(30),
					Volume: decimal.NewFromInt64( 3),
				},
			},
			expOrders: []luno.OrderBookEntry{
				{Price:decimal.NewFromInt64(20), Volume: decimal.NewFromInt64(2)},
				{Price:decimal.NewFromInt64(30), Volume: decimal.NewFromInt64(3)},
				{Price:decimal.NewFromInt64(40), Volume: decimal.NewFromInt64(4)},
			},
		},
		{name: "reversed orders",
			orders: map[string]order{
				"1": {
					ID:     "1",
					Price:  decimal.NewFromInt64(20),
					Volume: decimal.NewFromInt64( 2),
				},
				"3": {
					ID:     "3",
					Price:  decimal.NewFromInt64(40),
					Volume: decimal.NewFromInt64( 4),
				},
				"2": {
					ID:     "2",
					Price:  decimal.NewFromInt64(30),
					Volume: decimal.NewFromInt64( 3),
				},
			},
			reverse: true,
			expOrders: []luno.OrderBookEntry{
				{Price:decimal.NewFromInt64(40), Volume: decimal.NewFromInt64(4)},
				{Price:decimal.NewFromInt64(30), Volume: decimal.NewFromInt64(3)},
				{Price:decimal.NewFromInt64(20), Volume: decimal.NewFromInt64(2)},
			},
		},
		{name: "condense orders",
			orders: map[string]order{
				"1": {
					ID:     "1",
					Price:  decimal.NewFromInt64(1000),
					Volume: decimal.NewFromFloat64( 0.4, 1),
				},
				"3": {
					ID:     "3",
					Price:  decimal.NewFromInt64(1000),
					Volume: decimal.NewFromFloat64( 0.4, 1),
				},
				"2": {
					ID:     "2",
					Price:  decimal.NewFromInt64(1000),
					Volume: decimal.NewFromFloat64( 0.2, 1),
				},
			},
			expOrders: []luno.OrderBookEntry{
				{Price:decimal.NewFromInt64(1000), Volume: decimal.NewFromInt64(1)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orders := flatten(tc.orders, tc.reverse)

			if len(orders) != len(tc.expOrders) {
				t.Errorf("length of orders doesn't match %d, expected %d", len(orders), len(tc.expOrders))
			}

			for i, o := range orders {
				expO := tc.expOrders[i]
				if expO.Price.Cmp(o.Price) != 0 {
					t.Errorf("order %d price doesn't match %s, expected %s", i, o.Price, expO.Price)
				}
				if expO.Volume.Cmp(o.Volume) != 0 {
					t.Errorf("order %d volume doesn't match %s, expected %s", i, o.Volume, expO.Volume)
				}
			}
		})
	}
}