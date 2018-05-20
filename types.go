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
	Balance     Decimal `json:"balance"`
	Name        string  `json:"name"`
	Reserved    Decimal `json:"reserved"`
	Unconfirmed Decimal `json:"unconfirmed"`
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
	Base                Decimal `json:"base"`
	CompletedTimestamp  int64   `json:"completed_timestamp"`
	Counter             Decimal `json:"counter"`
	CreationTimestamp   int64   `json:"creation_timestamp"`
	ExpirationTimestamp int64   `json:"expiration_timestamp"`
	FeeBase             Decimal `json:"fee_base"`
	FeeCounter          Decimal `json:"fee_counter"`
	LimitPrice          Decimal `json:"limit_price"`
	LimitVolume         Decimal `json:"limit_volume"`
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
	Price  Decimal `json:"price"`
	Volume Decimal `json:"volume"`
}

type ReceiveAddress struct {
	AccountId        string  `json:"account_id"`
	Address          string  `json:"address"`
	Asset            string  `json:"asset"`
	AssignedAt       int64   `json:"assigned_at"`
	Name             string  `json:"name"`
	ReceiveFee       Decimal `json:"receive_fee"`
	TotalReceived    Decimal `json:"total_received"`
	TotalUnconfirmed Decimal `json:"total_unconfirmed"`
}

type Ticker struct {
	Ask                 Decimal `json:"ask"`
	Bid                 Decimal `json:"bid"`
	LastTrade           Decimal `json:"last_trade"`
	Pair                string  `json:"pair"`
	Rolling24HourVolume Decimal `json:"rolling_24_hour_volume"`
	Timestamp           int64   `json:"timestamp"`
}

type Trade struct {
	Base       Decimal `json:"base"`
	Counter    Decimal `json:"counter"`
	FeeBase    Decimal `json:"fee_base"`
	FeeCounter Decimal `json:"fee_counter"`
	IsBuy      bool    `json:"is_buy"`
	OrderId    string  `json:"order_id"`
	Pair       string  `json:"pair"`
	Price      Decimal `json:"price"`
	Timestamp  int64   `json:"timestamp"`
	Type       string  `json:"type"`
	Volume     Decimal `json:"volume"`
}

type Transaction struct {
	AccountId      string  `json:"account_id"`
	Available      Decimal `json:"available"`
	AvailableDelta Decimal `json:"available_delta"`
	Balance        Decimal `json:"balance"`

	// Transaction amounts computed for convenience.
	BalanceDelta Decimal `json:"balance_delta"`
	Currency     string  `json:"currency"`

	// Human-readable description of the transaction.
	Description string `json:"description"`
	RowIndex    int64  `json:"row_index"`
	Timestamp   int64  `json:"timestamp"`
}

type Withdrawal struct {
	Amount    Decimal `json:"amount"`
	CreatedAt int64   `json:"created_at"`
	Currency  string  `json:"currency"`
	Fee       Decimal `json:"fee"`
	Id        string  `json:"id"`
	Status    string  `json:"status"`
	Type      string  `json:"type"`
}

// vi: ft=go
