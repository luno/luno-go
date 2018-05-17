package luno

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

type AccountBalance struct {
	AccountId   string  `json:"account_id"`
	Asset       string  `json:"asset"`
	Balance     float64 `json:"balance,string"`
	Name        string  `json:"name"`
	Reserved    float64 `json:"reserved,string"`
	Unconfirmed float64 `json:"unconfirmed,string"`
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
	Base                float64 `json:"base,string"`
	CompletedTimestamp  int64   `json:"completed_timestamp"`
	Counter             float64 `json:"counter,string"`
	CreationTimestamp   int64   `json:"creation_timestamp"`
	ExpirationTimestamp int64   `json:"expiration_timestamp"`
	FeeBase             float64 `json:"fee_base,string"`
	FeeCounter          float64 `json:"fee_counter,string"`
	LimitPrice          float64 `json:"limit_price,string"`
	LimitVolume         float64 `json:"limit_volume,string"`
	OrderId             string  `json:"order_id"`

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
	Price  float64 `json:"price,string"`
	Volume float64 `json:"volume,string"`
}

type ReceiveAddress struct {
	AccountId        string  `json:"account_id"`
	Address          string  `json:"address"`
	Asset            string  `json:"asset"`
	AssignedAt       int64   `json:"assigned_at"`
	Name             string  `json:"name"`
	ReceiveFee       float64 `json:"receive_fee,string"`
	TotalReceived    float64 `json:"total_received,string"`
	TotalUnconfirmed float64 `json:"total_unconfirmed,string"`
}

type Ticker struct {
	Ask                 float64 `json:"ask,string"`
	Bid                 float64 `json:"bid,string"`
	LastTrade           float64 `json:"last_trade,string"`
	Pair                string  `json:"pair"`
	Rolling24HourVolume float64 `json:"rolling_24_hour_volume,string"`
	Timestamp           int64   `json:"timestamp"`
}

type Trade struct {
	Base       float64 `json:"base,string"`
	Counter    float64 `json:"counter,string"`
	FeeBase    float64 `json:"fee_base,string"`
	FeeCounter float64 `json:"fee_counter,string"`
	IsBuy      bool    `json:"is_buy"`
	OrderId    string  `json:"order_id"`
	Pair       string  `json:"pair"`
	Price      float64 `json:"price,string"`
	Timestamp  int64   `json:"timestamp"`
	Type       string  `json:"type"`
	Volume     float64 `json:"volume,string"`
}

type Transaction struct {
	AccountId      string  `json:"account_id"`
	Available      float64 `json:"available,string"`
	AvailableDelta float64 `json:"available_delta,string"`
	Balance        float64 `json:"balance,string"`

	// Transaction amounts computed for convenience.
	BalanceDelta float64 `json:"balance_delta,string"`
	Currency     string  `json:"currency"`

	// Human-readable description of the transaction.
	Description string `json:"description"`
	RowIndex    int64  `json:"row_index"`
	Timestamp   int64  `json:"timestamp"`
}

type Withdrawal struct {
	Amount    float64 `json:"amount,string"`
	CreatedAt int64   `json:"created_at"`
	Currency  string  `json:"currency"`
	Fee       float64 `json:"fee,string"`
	Id        string  `json:"id"`
	Status    string  `json:"status"`
	Type      string  `json:"type"`
}

// vi: ft=go
