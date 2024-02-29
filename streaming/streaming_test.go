package streaming

import (
	"math/rand"
	"testing"
	"time"

	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func TestFlatten(t *testing.T) {
	testCases := []struct {
		name string

		orders  map[string]order
		reverse bool

		expOrders []luno.OrderBookEntry
	}{
		{name: "empty orders"},
		{
			name: "single order",
			orders: map[string]order{
				"2": {
					ID:     "2",
					Price:  decimal.NewFromInt64(100),
					Volume: decimal.NewFromInt64(1),
				},
			},
			expOrders: []luno.OrderBookEntry{
				{Price: decimal.NewFromInt64(100), Volume: decimal.NewFromInt64(1)},
			},
		},
		{
			name: "sorted orders",
			orders: map[string]order{
				"1": {
					ID:     "1",
					Price:  decimal.NewFromInt64(20),
					Volume: decimal.NewFromInt64(2),
				},
				"3": {
					ID:     "3",
					Price:  decimal.NewFromInt64(40),
					Volume: decimal.NewFromInt64(4),
				},
				"2": {
					ID:     "2",
					Price:  decimal.NewFromInt64(30),
					Volume: decimal.NewFromInt64(3),
				},
			},
			expOrders: []luno.OrderBookEntry{
				{Price: decimal.NewFromInt64(20), Volume: decimal.NewFromInt64(2)},
				{Price: decimal.NewFromInt64(30), Volume: decimal.NewFromInt64(3)},
				{Price: decimal.NewFromInt64(40), Volume: decimal.NewFromInt64(4)},
			},
		},
		{
			name: "reversed orders",
			orders: map[string]order{
				"1": {
					ID:     "1",
					Price:  decimal.NewFromInt64(20),
					Volume: decimal.NewFromInt64(2),
				},
				"3": {
					ID:     "3",
					Price:  decimal.NewFromInt64(40),
					Volume: decimal.NewFromInt64(4),
				},
				"2": {
					ID:     "2",
					Price:  decimal.NewFromInt64(30),
					Volume: decimal.NewFromInt64(3),
				},
			},
			reverse: true,
			expOrders: []luno.OrderBookEntry{
				{Price: decimal.NewFromInt64(40), Volume: decimal.NewFromInt64(4)},
				{Price: decimal.NewFromInt64(30), Volume: decimal.NewFromInt64(3)},
				{Price: decimal.NewFromInt64(20), Volume: decimal.NewFromInt64(2)},
			},
		},
		{
			name: "condense orders",
			orders: map[string]order{
				"1": {
					ID:     "1",
					Price:  decimal.NewFromInt64(1000),
					Volume: decimal.NewFromFloat64(0.4, 1),
				},
				"3": {
					ID:     "3",
					Price:  decimal.NewFromInt64(1000),
					Volume: decimal.NewFromFloat64(0.4, 1),
				},
				"2": {
					ID:     "2",
					Price:  decimal.NewFromInt64(1000),
					Volume: decimal.NewFromFloat64(0.2, 1),
				},
			},
			expOrders: []luno.OrderBookEntry{
				{Price: decimal.NewFromInt64(1000), Volume: decimal.NewFromInt64(1)},
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

func TestBackoff(t *testing.T) {
	tcs := []struct {
		name         string
		fn           BackoffHandler
		attemptReset time.Duration
		p            *backoffParams
		seed         func()
		reqTS        time.Time
		expBackoff   time.Duration
		expParams    backoffParams
	}{
		{
			name: "default",
			fn:   defaultBackoffHandler,
			p:    &backoffParams{},
			seed: func() {
				// Seed random generator to test default backoff handler
				rand.Seed(123)
			},
			reqTS:      time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			expBackoff: 1935000000,
			expParams: backoffParams{
				attempts:    1,
				lastAttempt: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "second per attempt",
			p:    &backoffParams{},
			fn: func(attempt int) time.Duration {
				return time.Second * time.Duration(attempt)
			},
			reqTS:      time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			expBackoff: time.Second,
			expParams: backoffParams{
				attempts:    1,
				lastAttempt: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "second per attempt 3rd attempt",
			p: &backoffParams{
				attempts:    3,
				lastAttempt: time.Date(2024, 2, 3, 0, 10, 0, 0, time.UTC),
			},
			fn: func(attempt int) time.Duration {
				return time.Second * time.Duration(attempt)
			},
			attemptReset: defaultAttemptReset,
			reqTS:        time.Date(2024, 2, 3, 0, 20, 0, 0, time.UTC),
			expBackoff:   time.Second * 4,
			expParams: backoffParams{
				attempts:    4,
				lastAttempt: time.Date(2024, 2, 3, 0, 20, 0, 0, time.UTC),
			},
		},
		{
			name: "second per attempt 3rd attempt reset",
			p: &backoffParams{
				attempts:    3,
				lastAttempt: time.Date(2024, 2, 3, 0, 10, 0, 0, time.UTC),
			},
			fn: func(attempt int) time.Duration {
				return time.Second * time.Duration(attempt)
			},
			attemptReset: time.Minute * 10,
			reqTS:        time.Date(2024, 2, 3, 0, 20, 0, 0, time.UTC),
			expBackoff:   time.Second,
			expParams: backoffParams{
				attempts:    1,
				lastAttempt: time.Date(2024, 2, 3, 0, 20, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.seed != nil {
				tc.seed()
			}

			c := &Conn{
				backoffHandler: tc.fn,
				attemptReset:   tc.attemptReset,
			}

			dt := c.calculateBackoff(tc.p, tc.reqTS)
			if dt != tc.expBackoff {
				t.Errorf("backoff %d doesn't match expect backoff %d", dt, tc.expBackoff)
			}

			if tc.expParams != *tc.p {
				t.Errorf("params doesn't match expected params")
			}
		})
	}
}
