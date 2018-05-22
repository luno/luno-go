package luno

import "github.com/luno/luno-go/decimal"

type AccountBalance struct {
	AccountId   string          `json:"account_id"`
	Asset       string          `json:"asset"`
	Balance     decimal.Decimal `json:"balance"`
	Name        string          `json:"name"`
	Reserved    decimal.Decimal `json:"reserved"`
	Unconfirmed decimal.Decimal `json:"unconfirmed"`
}

type AccountCapabilities struct {
	CanBuy        bool `json:"can_buy"`
	CanDeposit    bool `json:"can_deposit"`
	CanReceive    bool `json:"can_receive"`
	CanSell       bool `json:"can_sell"`
	CanSend       bool `json:"can_send"`
	CanSocialSend bool `json:"can_social_send"`
	CanWithdraw   bool `json:"can_withdraw"`
}

type Order struct {
	Base                decimal.Decimal `json:"base"`
	CompletedTimestamp  int64           `json:"completed_timestamp"`
	Counter             decimal.Decimal `json:"counter"`
	CreationTimestamp   int64           `json:"creation_timestamp"`
	ExpirationTimestamp int64           `json:"expiration_timestamp"`
	FeeBase             decimal.Decimal `json:"fee_base"`
	FeeCounter          decimal.Decimal `json:"fee_counter"`
	LimitPrice          decimal.Decimal `json:"limit_price"`
	LimitVolume         decimal.Decimal `json:"limit_volume"`
	OrderId             string          `json:"order_id"`

	// Specifies the market.
	Pair string `json:"pair"`

	// <code>PENDING</code> The order has been placed. Some trades may have
	// taken place but the order is not filled yet.<br>
	// <code>COMPLETE</code> The order is no longer active. It has been settled
	// or has been cancelled.
	State string `json:"state"`

	// <code>BID</code> bid (buy) limit order.<br>
	// <code>ASK</code> ask (sell) limit order.
	Type string `json:"type"`
}

type OrderBookEntry struct {
	Price  decimal.Decimal `json:"price"`
	Volume decimal.Decimal `json:"volume"`
}

type OrderState string

const (
	OrderStatePending  OrderState = "PENDING"
	OrderStateComplete OrderState = "COMPLETE"
)

type OrderType string

const (
	OrderTypeAsk  OrderType = "ASK"
	OrderTypeBid  OrderType = "BID"
	OrderTypeBuy  OrderType = "BUY"
	OrderTypeSell OrderType = "SELL"
)

type ReceiveAddress struct {
	AccountId        string          `json:"account_id"`
	Address          string          `json:"address"`
	Asset            string          `json:"asset"`
	AssignedAt       int64           `json:"assigned_at"`
	Name             string          `json:"name"`
	ReceiveFee       decimal.Decimal `json:"receive_fee"`
	TotalReceived    decimal.Decimal `json:"total_received"`
	TotalUnconfirmed decimal.Decimal `json:"total_unconfirmed"`
}

type Ticker struct {
	Ask                 decimal.Decimal `json:"ask"`
	Bid                 decimal.Decimal `json:"bid"`
	LastTrade           decimal.Decimal `json:"last_trade"`
	Pair                string          `json:"pair"`
	Rolling24HourVolume decimal.Decimal `json:"rolling_24_hour_volume"`
	Timestamp           int64           `json:"timestamp"`
}

type Trade struct {
	Base       decimal.Decimal `json:"base"`
	Counter    decimal.Decimal `json:"counter"`
	FeeBase    decimal.Decimal `json:"fee_base"`
	FeeCounter decimal.Decimal `json:"fee_counter"`
	IsBuy      bool            `json:"is_buy"`
	OrderId    string          `json:"order_id"`
	Pair       string          `json:"pair"`
	Price      decimal.Decimal `json:"price"`
	Timestamp  int64           `json:"timestamp"`
	Type       string          `json:"type"`
	Volume     decimal.Decimal `json:"volume"`
}

type Transaction struct {
	AccountId      string          `json:"account_id"`
	Available      decimal.Decimal `json:"available"`
	AvailableDelta decimal.Decimal `json:"available_delta"`
	Balance        decimal.Decimal `json:"balance"`

	// Transaction amounts computed for convenience.
	BalanceDelta decimal.Decimal `json:"balance_delta"`
	Currency     string          `json:"currency"`

	// Human-readable description of the transaction.
	Description string `json:"description"`
	RowIndex    int64  `json:"row_index"`
	Timestamp   int64  `json:"timestamp"`
}

type Withdrawal struct {
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt int64           `json:"created_at"`
	Currency  string          `json:"currency"`
	Fee       decimal.Decimal `json:"fee"`
	Id        string          `json:"id"`
	Status    string          `json:"status"`
	Type      string          `json:"type"`
}

// vi: ft=go
