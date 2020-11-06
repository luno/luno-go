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

type AddressMeta struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type CryptoDetails struct {
	Address string `json:"address"`
	Txid    string `json:"txid"`
}

type DetailFields struct {
	CryptoDetails CryptoDetails `json:"crypto_details"`
	TradeDetails  TradeDetails  `json:"trade_details"`
}

type Kind string

const (
	KindExchange Kind = "EXCHANGE"
	KindFee      Kind = "FEE"
	KindInterest Kind = "INTEREST"
	KindTransfer Kind = "TRANSFER"
)

type MarketInfo struct {
	// Base currency code
	BaseCurrency string `json:"base_currency"`

	// Counter currency code
	CounterCurrency string `json:"counter_currency"`

	// Fee decimal places
	FeeScale int64 `json:"fee_scale"`

	// Unique identifier for the market
	MarketId string `json:"market_id"`

	// Maximum order price
	MaxPrice decimal.Decimal `json:"max_price"`

	// Maximum order volume
	MaxVolume decimal.Decimal `json:"max_volume"`

	// Minimum order price
	MinPrice decimal.Decimal `json:"min_price"`

	// Minimum order volume
	MinVolume decimal.Decimal `json:"min_volume"`

	// Price decimal places
	PriceScale int64 `json:"price_scale"`

	// Current market trading status:<br>
	// <code>POST_ONLY</code> Trading is indefinitely suspended. This state is
	// commonly used when new markets are being launched to give traders enough
	// time to setup their orders before trading begins. When in this status,
	// orders can only be posted as post-only.<br>
	// <code>ACTIVE</code> Trading is fully enabled.<br>
	// <code>SUSPENDED</code> Trading has been temporarily suspended due to very
	// high volatility. When in this status, orders can only be posted as
	// post-only.<br>
	TradingStatus TradingStatus `json:"trading_status"`

	// Volume decimal places
	VolumeScale int64 `json:"volume_scale"`
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
	// Limit price
	Price decimal.Decimal `json:"price"`

	// Volume available
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

type OrderV2 struct {
	// Amount of base filled
	Base decimal.Decimal `json:"base"`

	// Time of order completion in milliseconds
	CompletedTimestamp Time `json:"completed_timestamp"`

	// Amount of counter filled
	Counter decimal.Decimal `json:"counter"`

	// Time of order creation in milliseconds
	CreationTimestamp Time `json:"creation_timestamp"`

	// Time of order expiration in milliseconds
	ExpirationTimestamp Time `json:"expiration_timestamp"`

	// Base amount of fees to be charged
	FeeBase decimal.Decimal `json:"fee_base"`

	// Counter amount of fees to be charged
	FeeCounter decimal.Decimal `json:"fee_counter"`

	// Limit price to transact
	LimitPrice decimal.Decimal `json:"limit_price"`

	// Limit volume to transact
	LimitVolume decimal.Decimal `json:"limit_volume"`

	// The order reference
	OrderId string `json:"order_id"`

	// Specifies the market
	Pair string `json:"pair"`

	// The order intention
	Side Side `json:"side"`

	// The current state of the order
	//
	// Status meaning:<br>
	// <code>AWAITING</code> The order is awaiting to enter the order book.<br>
	// <code>PENDING</code> The order is in the order book. Some trades may
	// have taken place but the order is not filled yet.<br>
	// <code>COMPLETE</code> The order is no longer in the order book. It has
	// been settled/filled or has been cancelled.
	Status Status `json:"status"`

	// Direction to trigger the order
	StopDirection StopDirection `json:"stop_direction"`

	// Price to trigger the order
	StopPrice decimal.Decimal `json:"stop_price"`

	// The order type
	Type Type `json:"type"`
}

type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusAwaiting Status = "AWAITING"
	StatusComplete Status = "COMPLETE"
	StatusDisabled Status = "DISABLED"
	StatusPending  Status = "PENDING"
	StatusPostonly Status = "POSTONLY"
)

type StopDirection string

const (
	StopDirectionAbove               StopDirection = "ABOVE"
	StopDirectionBelow               StopDirection = "BELOW"
	StopDirectionRelative_last_trade StopDirection = "RELATIVE_LAST_TRADE"
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

type TradeDetails struct {
	// Pair of the market
	Pair string `json:"pair"`

	// Price at which the volume traded for
	Price decimal.Decimal `json:"price"`

	// Sequence identifies the trade within a market
	Sequence int64 `json:"sequence"`

	// Volume is the amount of base traded
	Volume decimal.Decimal `json:"volume"`
}

type TradingStatus string

const (
	TradingStatusPost_only TradingStatus = "POST_ONLY"
	TradingStatusActive    TradingStatus = "ACTIVE"
	TradingStatusSuspended TradingStatus = "SUSPENDED"
)

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
	Details map[string]string `json:"details"`

	// The kind of the transaction indicates the transaction flow
	//
	// Kinds explained:<br>
	// <code>FEE</code> when transaction is towards Luno fees<br>
	// <code>TRANSFER</code> when the transaction is a one way flow of funds, e.g. a deposit or crypto send<br>
	// <code>EXCHANGE</code> when the transaction is part of a two way exchange, e.g. a trade or instant buy
	Kind      Kind  `json:"kind"`
	RowIndex  int64 `json:"row_index"`
	Timestamp Time  `json:"timestamp"`
}

type Type string

const (
	TypeLimit      Type = "LIMIT"
	TypeMarket     Type = "MARKET"
	TypeStop_limit Type = "STOP_LIMIT"
)

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
