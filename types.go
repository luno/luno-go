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

type CryptoDetails struct {
	Address string `json:"address"`
	Txid    string `json:"txid"`
}

type DetailFields struct {
	CryptoDetails CryptoDetails `json:"crypto_details"`
}

type Order struct {
	Base                decimal.Decimal `json:"base"`
	CompletedTimestamp  Time            `json:"completed_timestamp"`
	Counter             decimal.Decimal `json:"counter"`
	CreationTimestamp   Time            `json:"creation_timestamp"`
	ExpirationTimestamp Time            `json:"expiration_timestamp"`
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
	State OrderState `json:"state"`

	// <code>BID</code> bid (buy) limit order.<br>
	// <code>ASK</code> ask (sell) limit order.
	Type OrderType `json:"type"`
}

type OrderBookEntry struct {
	Price  decimal.Decimal `json:"price"`
	Volume decimal.Decimal `json:"volume"`
}

type OrderState string

const (
	OrderStateComplete OrderState = "COMPLETE"
	OrderStatePending  OrderState = "PENDING"
)

type OrderType string

const (
	OrderTypeAsk  OrderType = "ASK"
	OrderTypeBid  OrderType = "BID"
	OrderTypeBuy  OrderType = "BUY"
	OrderTypeSell OrderType = "SELL"
)

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusDisabled Status = "DISABLED"
	StatusPostonly Status = "POSTONLY"
)

type Ticker struct {
	Ask                 decimal.Decimal `json:"ask"`
	Bid                 decimal.Decimal `json:"bid"`
	LastTrade           decimal.Decimal `json:"last_trade"`
	Pair                string          `json:"pair"`
	Rolling24HourVolume decimal.Decimal `json:"rolling_24_hour_volume"`

	// <code>ACTIVE</code> when the market is trading normally
	//
	// <code>POSTONLY</code> when the market has been suspended and only post-only orders will be accepted
	//
	// <code>DISABLED</code> when the market is shutdown and no orders can be accepted
	Status    Status `json:"status"`
	Timestamp Time   `json:"timestamp"`
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
	Sequence   int64           `json:"sequence"`
	Timestamp  Time            `json:"timestamp"`
	Type       OrderType       `json:"type"`
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
	Description  string       `json:"description"`
	DetailFields DetailFields `json:"detail_fields"`

	// Human-readable label-value attributes.
	Details   map[string]string `json:"details"`
	RowIndex  int64             `json:"row_index"`
	Timestamp Time              `json:"timestamp"`
}

type Withdrawal struct {
	Amount     decimal.Decimal `json:"amount"`
	CreatedAt  Time            `json:"created_at"`
	Currency   string          `json:"currency"`
	ExternalId string          `json:"external_id"`
	Fee        decimal.Decimal `json:"fee"`
	Id         string          `json:"id"`
	Status     string          `json:"status"`
	Type       string          `json:"type"`
}

type beneficiary struct {
	BankAccountBranch string `json:"bank_account_branch"`
	BankAccountNumber string `json:"bank_account_number"`
	BankAccountType   string `json:"bank_account_type"`
	BankCountry       string `json:"bank_country"`
	BankName          string `json:"bank_name"`
	BankRecipient     string `json:"bank_recipient"`
	CreatedAt         Time   `json:"created_at"`
	Id                string `json:"id"`
}

// vi: ft=go
