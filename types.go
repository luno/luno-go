package luno

import "github.com/luno/luno-go/decimal"

type AccountBalance struct {
	// ID of the account.
	AccountId string `json:"account_id"`

	// Currency code for the asset held in this account.
	Asset string `json:"asset"`

	// The amount available to send or trade.
	Balance decimal.Decimal `json:"balance"`

	// The name set by the user upon creating the account.
	Name string `json:"name"`

	// Amount locked by Luno and cannot be sent or traded. This could be due to
	// open orders.
	Reserved decimal.Decimal `json:"reserved"`

	// Amount that is awaiting some sort of verification to be credited to this
	// account. This could be an on-chain transaction that Luno is waiting for
	// further block verifications to happen.
	Unconfirmed decimal.Decimal `json:"unconfirmed"`
}

type AccountType string

const (
	AccountTypeCurrentcheque AccountType = "Current/Cheque"
	AccountTypeSavings       AccountType = "Savings"
	AccountTypeTransmission  AccountType = "Transmission"
)

type AddressMeta struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type BankAccountType string

const (
	BankAccountTypeCurrentcheque BankAccountType = "Current/Cheque"
	BankAccountTypeSavings       BankAccountType = "Savings"
	BankAccountTypeTransmission  BankAccountType = "Transmission"
)

type Candle struct {
	// Closing price
	Close decimal.Decimal `json:"close"`

	// High price
	High decimal.Decimal `json:"high"`

	// Low price
	Low decimal.Decimal `json:"low"`

	// Opening price
	Open decimal.Decimal `json:"open"`

	// Unix timestamp in milliseconds
	Timestamp Time `json:"timestamp"`

	// Volume traded
	Volume decimal.Decimal `json:"volume"`
}

type CryptoDetails struct {
	Address string `json:"address"`
	Txid    string `json:"txid"`
}

type DetailFields struct {
	CryptoDetails CryptoDetails `json:"crypto_details"`
	TradeDetails  TradeDetails  `json:"trade_details"`
}

type FundsMove struct {
	// The assets quantity to move from the debit account to credit account. This is always a positive value.
	Amount decimal.Decimal `json:"amount"`

	// User defined unique ID
	ClientMoveId string `json:"client_move_id"`

	// Unix time the move was initiated, in milliseconds
	CreatedAt Time `json:"created_at"`

	// The account to credit the funds to.
	CreditAccountId string `json:"credit_account_id"`

	// The account to debit the funds from.
	DebitAccountId string `json:"debit_account_id"`

	// Unique ID, defined by Luno
	Id string `json:"id"`

	// Current status of the move.
	//
	// Status meaning:<br>
	// <code>CREATED</code> The move is awaiting execution.<br>
	// <code>MOVING</code> The funds have been reserved and the move is being executed.<br>
	// <code>SUCCESSFUL</code> The move has completed successfully and should be reflected in both accounts available
	// balance.<br>
	// <code>FAILED</code> The move has failed. There could be many reasons for this but the most likely is that the
	// debit account doesn't have enough available funds to move.<br>
	Status Status `json:"status"`

	// Unix time the move was last updated, in milliseconds
	UpdatedAt Time `json:"updated_at"`
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
	// Amount of base filled, this value is always positive.
	Base decimal.Decimal `json:"base"`

	// Time of order completion (Unix milliseconds)
	//
	// This value is set at the time of this order leaving the order book,
	// either immediately upon posting or later on due to a trade or cancellation.
	// Whilst the order is still pending/live it will be 0.
	CompletedTimestamp Time `json:"completed_timestamp"`

	// Amount of counter filled, this value is always positive.
	Counter decimal.Decimal `json:"counter"`

	// Time of order creation (Unix milliseconds)
	CreationTimestamp Time `json:"creation_timestamp"`

	// Time of order expiration (Unix milliseconds)
	//
	// This value is set at the time of processing a request from you to cancel the order, otherwise it will be 0.
	ExpirationTimestamp Time `json:"expiration_timestamp"`

	// Base amount of fees to be charged
	FeeBase decimal.Decimal `json:"fee_base"`

	// Counter amount of fees to be charged
	FeeCounter decimal.Decimal `json:"fee_counter"`

	// Limit price to transact
	LimitPrice decimal.Decimal `json:"limit_price"`

	// Limit volume to transact
	LimitVolume decimal.Decimal `json:"limit_volume"`
	OrderId     string          `json:"order_id"`

	// Specifies the market.
	Pair string `json:"pair"`

	// <code>PENDING</code> The order has been placed. Some trades may have
	// taken place but the order is not filled yet.<br>
	// <code>COMPLETE</code> The order is no longer active. It has been settled
	// or has been cancelled.
	State OrderState `json:"state"`

	// The Time in force option used when the LimitOrder was posted.
	//
	// Only returned on limit orders.<br>
	// <code>GTC</code> Good 'Til Cancelled. The order remains open until it is filled or cancelled by the user. (default)</br>
	// <code>IOC</code> Immediate Or Cancel. The part of the order that cannot be filled immediately will be cancelled. Cannot be post-only.</br>
	// <code>FOK</code> Fill Or Kill. If the order cannot be filled immediately and completely it will be cancelled before any trade. Cannot be post-only.
	TimeInForce string `json:"time_in_force"`

	// <code>BUY</code> buy market order.<br>
	// <code>SELL</code> sell market order.<br>
	// <code>BID</code> bid (buy) limit order.<br>
	// <code>ASK</code> ask (sell) limit order.
	Type OrderType `json:"type"`
}

type OrderBookEntry struct {
	// Limit price at which orders are trading at
	Price decimal.Decimal `json:"price"`

	// The volume available at the limit price
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
	// Amount of base filled, this value is always positive.
	//
	// Use this field and `side` to determine credit or debit of funds.
	Base decimal.Decimal `json:"base"`

	// The base currency account
	BaseAccountId int64 `json:"base_account_id"`

	// Client Order ID has the value that was passed in when the Order was posted.
	ClientOrderId string `json:"client_order_id"`

	// Time of order completion (Unix milliseconds)
	//
	// This value is set at the time of this order leaving the order book,
	// either immediately upon posting or later on due to a trade or cancellation.
	// Whilst the order is still pending/live it will be 0.
	CompletedTimestamp Time `json:"completed_timestamp"`

	// Amount of counter filled, this value is always positive.
	//
	// Use this field and `side` to determine credit or debit of funds.
	Counter decimal.Decimal `json:"counter"`

	// The counter currency account
	CounterAccountId int64 `json:"counter_account_id"`

	// Time of order creation (Unix milliseconds)
	CreationTimestamp Time `json:"creation_timestamp"`

	// Time of order expiration (Unix milliseconds)
	//
	// This value is set at the time of processing a request from you to cancel the order, otherwise it will be 0.
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

	// The intention of the order, whether to buy or sell funds in the market.
	//
	// You can use this to determine the flow of funds in the order.
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

	// The Time in force option used when the LimitOrder was posted.
	//
	// Only returned on limit orders.<br>
	// <code>GTC</code> Good 'Til Cancelled. The order remains open until it is filled or cancelled by the user. (default)</br>
	// <code>IOC</code> Immediate Or Cancel. The part of the order that cannot be filled immediately will be cancelled. Cannot be post-only.</br>
	// <code>FOK</code> Fill Or Kill. If the order cannot be filled immediately and completely it will be cancelled before any trade. Cannot be post-only.
	TimeInForce string `json:"time_in_force"`

	// The order type
	Type Type `json:"type"`
}

type PublicTrade struct {
	// Whether the taker was buying or not.
	IsBuy bool `json:"is_buy"`

	// Price at which the asset traded at
	Price decimal.Decimal `json:"price"`

	// The ever incrementing trade identifier within a market
	Sequence int64 `json:"sequence"`

	// Unix timestamp in milliseconds
	Timestamp Time `json:"timestamp"`

	// Amount of assets traded
	Volume decimal.Decimal `json:"volume"`
}

type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

type Status string

const (
	StatusActive     Status = "ACTIVE"
	StatusAwaiting   Status = "AWAITING"
	StatusCancelled  Status = "CANCELLED"
	StatusCancelling Status = "CANCELLING"
	StatusComplete   Status = "COMPLETE"
	StatusCompleted  Status = "COMPLETED"
	StatusCreated    Status = "CREATED"
	StatusDisabled   Status = "DISABLED"
	StatusFailed     Status = "FAILED"
	StatusMoving     Status = "MOVING"
	StatusPending    Status = "PENDING"
	StatusPostonly   Status = "POSTONLY"
	StatusProcessing Status = "PROCESSING"
	StatusSuccessful Status = "SUCCESSFUL"
	StatusUnknown    Status = "UNKNOWN"
	StatusWaiting    Status = "WAITING"
)

type StopDirection string

const (
	StopDirectionAbove               StopDirection = "ABOVE"
	StopDirectionBelow               StopDirection = "BELOW"
	StopDirectionRelative_last_trade StopDirection = "RELATIVE_LAST_TRADE"
)

type Ticker struct {
	// The lowest ask price
	Ask decimal.Decimal `json:"ask"`

	// The highest bid price
	Bid decimal.Decimal `json:"bid"`

	// Last trade price
	LastTrade decimal.Decimal `json:"last_trade"`
	Pair      string          `json:"pair"`

	// 24h rolling trade volume
	Rolling24HourVolume decimal.Decimal `json:"rolling_24_hour_volume"`

	// Market current status
	//
	// <code>ACTIVE</code> when the market is trading normally
	//
	// <code>POSTONLY</code> when the market has been suspended and only post-only orders will be accepted
	//
	// <code>DISABLED</code> when the market is shutdown and no orders can be accepted
	Status Status `json:"status"`

	// Unix timestamp in milliseconds of the tick
	Timestamp Time `json:"timestamp"`
}

type TimeInForce string

const (
	TimeInForceGtc TimeInForce = "GTC"
	TimeInForceIoc TimeInForce = "IOC"
	TimeInForceFok TimeInForce = "FOK"
)

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

type TradeV2 struct {
	// Amount of base filled
	Base decimal.Decimal `json:"base"`

	// Client Order ID has the value that was passed in when the Order was posted.
	ClientOrderId string `json:"client_order_id"`

	// Amount of counter filled
	Counter decimal.Decimal `json:"counter"`

	// Base amount of fees charged
	FeeBase decimal.Decimal `json:"fee_base"`

	// Counter amount of fees charged
	FeeCounter decimal.Decimal `json:"fee_counter"`
	IsBuy      bool            `json:"is_buy"`

	// Unique order identifier
	OrderId string `json:"order_id"`

	// Currency pair
	Pair string `json:"pair"`

	// Order price
	Price    decimal.Decimal `json:"price"`
	Sequence int64           `json:"sequence"`

	// Unix timestamp in milliseconds
	Timestamp Time `json:"timestamp"`

	// Order type
	Type OrderType `json:"type"`

	// Order volume
	Volume decimal.Decimal `json:"volume"`
}

type TradingStatus string

const (
	TradingStatusPost_only TradingStatus = "POST_ONLY"
	TradingStatusActive    TradingStatus = "ACTIVE"
	TradingStatusSuspended TradingStatus = "SUSPENDED"
)

type Transaction struct {
	AccountId string `json:"account_id"`

	// Amount available
	Available decimal.Decimal `json:"available"`

	// Change in amount available
	AvailableDelta decimal.Decimal `json:"available_delta"`

	// Account balance
	Balance decimal.Decimal `json:"balance"`

	// Change in balance
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
	Kind Kind `json:"kind"`

	// A unique reference for the transaction this statement entry relates to.
	// There may be multiple statement entries related to the same transaction.
	// E.g. a withdrawal and the withdrawal fee are two separate statement entries with the same reference.
	Reference string `json:"reference"`
	RowIndex  int64  `json:"row_index"`

	// Unix timestamp, in milliseconds
	Timestamp Time `json:"timestamp"`
}

type Transfer struct {
	// Amount that has been credited or debited on the account. This is always a
	// positive value regardless of the transfer direction.
	Amount decimal.Decimal `json:"amount"`

	// Unix timestamp the transfer was initiated, in milliseconds
	CreatedAt Time `json:"created_at"`

	// Fee that has been charged by Luno with regards to this transfer.
	// This is not included in the `amount`.
	// For example, if you receive a transaction with the raw amount of 1 BTC
	// and we charge a `fee` of 0.003 BTC on this transaction you will be
	// credited the `amount` of 0.997 BTC.
	Fee decimal.Decimal `json:"fee"`

	// Transfer unique identifier
	Id string `json:"id"`

	// True for credit transfers, false for debits.
	Inbound bool `json:"inbound"`

	// When the transfer reflects an on-chain transaction this field will have
	// the transaction ID.
	TransactionId string `json:"transaction_id"`
}

type Type string

const (
	TypeLimit      Type = "LIMIT"
	TypeMarket     Type = "MARKET"
	TypeStop_limit Type = "STOP_LIMIT"
)

type Withdrawal struct {
	// Amount to withdraw
	Amount decimal.Decimal `json:"amount"`

	// Unix time the withdrawal was initiated, in milliseconds
	CreatedAt Time `json:"created_at"`

	// Withdrawal currency.
	Currency string `json:"currency"`

	// External ID has the value that was passed in when the Withdrawal request was posted.
	ExternalId string `json:"external_id"`

	// Withdrawal fee
	Fee decimal.Decimal `json:"fee"`
	Id  string          `json:"id"`

	// Status
	Status Status `json:"status"`

	// Transfer ID is the identifier of the Withdrawal's transfer once it completes.
	TransferId string `json:"transfer_id"`

	// Type distinguishes between different withdrawal methods where more than one is supported
	// for the given currency.
	Type string `json:"type"`
}

type beneficiary struct {
	// Bank branch code
	BankAccountBranch string `json:"bank_account_branch"`

	// Beneficiary bank account number
	BankAccountNumber string `json:"bank_account_number"`

	// Bank account type
	BankAccountType BankAccountType `json:"bank_account_type"`

	// Bank country of origin
	BankCountry string `json:"bank_country"`

	// Bank SWIFT code
	BankName string `json:"bank_name"`

	// The owner of the recipient account
	BankRecipient string `json:"bank_recipient"`

	// Time of beneficiary creation
	CreatedAt int64 `json:"created_at"`

	// Unique id referencing beneficiary
	Id string `json:"id"`

	// If the bank account supports fast withdrawals
	SupportsFastWithdrawals bool `json:"supports_fast_withdrawals"`
}

// vi: ft=go
