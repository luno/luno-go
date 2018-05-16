package luno

import "context"

type CancelWithdrawalRequest struct {
	// ID of the withdrawal to cancel.
	Id int64 `json:"id" url:"id"`
}

type CancelWithdrawalResponse struct {
	Amount    float64 `json:"amount,string"`
	CreatedAt int64   `json:"created_at"`
	Currency  string  `json:"currency"`
	Fee       float64 `json:"fee,string"`
	Id        string  `json:"id"`
	Status    string  `json:"status"`
	Type      string  `json:"type"`
}

// CancelWithdrawal makes a call to DELETE /api/1/withdrawals/{id}.
//
// Cancel a withdrawal request. This can only be done if the request is still
// in state <code>PENDING</code>.
//
// Permissions required: <code>Perm_W_Withdrawals</code>
func (cl *Client) CancelWithdrawal(ctx context.Context, req *CancelWithdrawalRequest) (*CancelWithdrawalResponse, error) {
	var res CancelWithdrawalResponse
	err := cl.do(ctx, "DELETE", "/api/1/withdrawals/{id}", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type CreateAccountRequest struct {
	// The currency code for the account you want to create
	//
	// You must be verified to trade currency in order to be able to create an
	// account. A user has a limit of 4 accounts per currency.
	Currency string `json:"currency" url:"currency"`

	// The label to use for this account
	Name string `json:"name" url:"name"`
}

type CreateAccountResponse struct {
	Balance          AccountBalance      `json:"balance"`
	Capabilities     AccountCapabilities `json:"capabilities"`
	Currency         string              `json:"currency"`
	Icon             string              `json:"icon"`
	Id               string              `json:"id"`
	IsDefault        bool                `json:"is_default"`
	Name             string              `json:"name"`
	Pending          []Transaction       `json:"pending"`
	ReceiveAddresses []ReceiveAddress    `json:"receive_addresses"`
	Transactions     []Transaction       `json:"transactions"`
}

type Transaction struct {
	AccountId      string  `json:"account_id"`
	AppExt         AppExt  `json:"app_ext"`
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

type AppExt struct {
	ApproxCurrency  string `json:"approx_currency"`
	ApproxValue     string `json:"approx_value"`
	CancelToken     string `json:"cancel_token"`
	PrimaryCategory string `json:"primary_category"`
}

// CreateAccount makes a call to POST /api/1/accounts.
//
// Create an additional account for the specified currency.
//
// Permissions required: <code>Perm_W_Addresses</code>
func (cl *Client) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	var res CreateAccountResponse
	err := cl.do(ctx, "POST", "/api/1/accounts", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type CreateFundingAddressRequest struct {
	// Currency code of the asset.
	Asset string `json:"asset" url:"asset"`

	// An optional name for the new address
	Name string `json:"name" url:"name"`
}

type CreateFundingAddressResponse struct {
	AccountId        string  `json:"account_id"`
	Address          string  `json:"address"`
	Asset            string  `json:"asset"`
	AssignedAt       int64   `json:"assigned_at"`
	Name             string  `json:"name"`
	ReceiveFee       float64 `json:"receive_fee,string"`
	TotalReceived    float64 `json:"total_received,string"`
	TotalUnconfirmed float64 `json:"total_unconfirmed,string"`
}

// CreateFundingAddress makes a call to POST /api/1/funding_address.
//
// Allocates a new receive address to your account. There is a rate limit of 1
// address per hour, but bursts of up to 10 addresses are allowed. Only 1
// Ethereum receive address can be created.
//
// Permissions required: <code>Perm_W_Addresses</code>
func (cl *Client) CreateFundingAddress(ctx context.Context, req *CreateFundingAddressRequest) (*CreateFundingAddressResponse, error) {
	var res CreateFundingAddressResponse
	err := cl.do(ctx, "POST", "/api/1/funding_address", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type CreateQuoteRequest struct {
	// Amount to buy or sell in the pair base currency.
	BaseAmount float64 `json:"base_amount,string" url:"base_amount"`

	// Currency pair to trade. The pair can also be flipped if you want to buy
	// or sell the counter currency (e.g. ZARXBT).
	Pair string `json:"pair" url:"pair"`

	// <code>BUY</code> or <code>SELL</code>.
	Type string `json:"type" url:"type"`

	// Optional account for the pair's base currency.
	BaseAccountId int64 `json:"base_account_id" url:"base_account_id"`

	// Optional account for the pair's counter currency.
	CounterAccountId int64 `json:"counter_account_id" url:"counter_account_id"`
}

type CreateQuoteResponse struct {
	BaseAmount    float64 `json:"base_amount,string"`
	CounterAmount float64 `json:"counter_amount,string"`
	CreatedAt     int64   `json:"created_at"`
	Discarded     bool    `json:"discarded"`
	Exercised     bool    `json:"exercised"`
	ExpiresAt     int64   `json:"expires_at"`
	Id            string  `json:"id"`
	Pair          string  `json:"pair"`
	Type          string  `json:"type"`
}

// CreateQuote makes a call to POST /api/1/quotes.
//
// Creates a new quote to buy or sell a particular amount.
//
// You can specify either the exact amount that you want to pay or the exact
// amount that you want too receive.
//
// For example, to buy exactly 0.1 Bitcoin using ZAR, you would create a quote
// to BUY 0.1 XBTZAR. The returned quote includes the appropriate ZAR amount. To
// buy Bitcoin using exactly ZAR 100, you would create a quote to SELL 100
// ZARXBT. The returned quote specifies the Bitcoin as the counter amount that
// will be returned.
//
// An error is returned if your account is not verified for the currency pair,
// or if your account would have insufficient balance to ever exercise the
// quote.
//
// Permissions required: <code>Perm_W_Orders</code>
func (cl *Client) CreateQuote(ctx context.Context, req *CreateQuoteRequest) (*CreateQuoteResponse, error) {
	var res CreateQuoteResponse
	err := cl.do(ctx, "POST", "/api/1/quotes", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type CreateWithdrawalRequest struct {
	// Amount to withdraw. The currency depends on the type.
	Amount float64 `json:"amount,string" url:"amount"`

	// Withdrawal type.
	Type string `json:"type" url:"type"`

	// The beneficiary ID of the bank account the withdrawal will be paid out
	// to. This parameter is required if you have multiple bank accounts. Your
	// bank account beneficiary ID can be found by clicking on the beneficiary
	// name on the <a href="/wallet/beneficiaries">Beneficiaries</a> page.
	BeneficiaryId int64 `json:"beneficiary_id" url:"beneficiary_id"`

	// For internal use.
	Reference string `json:"reference" url:"reference"`
}

type CreateWithdrawalResponse struct {
	Amount    float64 `json:"amount,string"`
	CreatedAt int64   `json:"created_at"`
	Currency  string  `json:"currency"`
	Fee       float64 `json:"fee,string"`
	Id        string  `json:"id"`
	Status    string  `json:"status"`
	Type      string  `json:"type"`
}

// CreateWithdrawal makes a call to POST /api/1/withdrawals.
//
// Creates a new withdrawal request.
//
// Permissions required: <code>Perm_W_Withdrawals</code>
func (cl *Client) CreateWithdrawal(ctx context.Context, req *CreateWithdrawalRequest) (*CreateWithdrawalResponse, error) {
	var res CreateWithdrawalResponse
	err := cl.do(ctx, "POST", "/api/1/withdrawals", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DiscardQuoteRequest struct {
	// ID of the quote to discard.
	Id int64 `json:"id" url:"id"`
}

type DiscardQuoteResponse struct {
	BaseAmount    float64 `json:"base_amount,string"`
	CounterAmount float64 `json:"counter_amount,string"`
	CreatedAt     int64   `json:"created_at"`
	Discarded     bool    `json:"discarded"`
	Exercised     bool    `json:"exercised"`
	ExpiresAt     int64   `json:"expires_at"`
	Id            string  `json:"id"`
	Pair          string  `json:"pair"`
	Type          string  `json:"type"`
}

// DiscardQuote makes a call to DELETE /api/1/quotes/{id}.
//
// Discard a quote. Once a quote has been discarded, it cannot be exercised even
// if it has not expired yet.
//
// Permissions required: <code>Perm_W_Orders</code>
func (cl *Client) DiscardQuote(ctx context.Context, req *DiscardQuoteRequest) (*DiscardQuoteResponse, error) {
	var res DiscardQuoteResponse
	err := cl.do(ctx, "DELETE", "/api/1/quotes/{id}", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type ExerciseQuoteRequest struct {
	// ID of the quote to exercise.
	Id int64 `json:"id" url:"id"`
}

type ExerciseQuoteResponse struct {
	BaseAmount    float64 `json:"base_amount,string"`
	CounterAmount float64 `json:"counter_amount,string"`
	CreatedAt     int64   `json:"created_at"`
	Discarded     bool    `json:"discarded"`
	Exercised     bool    `json:"exercised"`
	ExpiresAt     int64   `json:"expires_at"`
	Id            string  `json:"id"`
	Pair          string  `json:"pair"`
	Type          string  `json:"type"`
}

// ExerciseQuote makes a call to PUT /api/1/quotes/{id}.
//
// Exercise a quote to perform the trade. If there is sufficient balance
// available in your account, it will be debited and the counter amount
// credited.
//
// An error is returned if the quote has expired or if you have insufficient
// available balance.
//
// Permissions required: <code>Perm_W_Orders</code>
func (cl *Client) ExerciseQuote(ctx context.Context, req *ExerciseQuoteRequest) (*ExerciseQuoteResponse, error) {
	var res ExerciseQuoteResponse
	err := cl.do(ctx, "PUT", "/api/1/quotes/{id}", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetBalancesRequest struct {
	// Only return balances for wallets with these currencies (if not provided,
	// all balances will be returned)
	Assets []string `json:"assets" url:"assets"`
}

type GetBalancesResponse struct {
	Balance []AccountBalance `json:"balance"`
}

// GetBalances makes a call to GET /api/1/balance.
//
// Return the list of all accounts and their respective balances.
//
// Permissions required: <code>Perm_R_Balance</code>
func (cl *Client) GetBalances(ctx context.Context, req *GetBalancesRequest) (*GetBalancesResponse, error) {
	var res GetBalancesResponse
	err := cl.do(ctx, "GET", "/api/1/balance", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetFeeInfoRequest struct {
	// Get fee information about this pair.
	Pair string `json:"pair" url:"pair"`
}

type GetFeeInfoResponse struct {
	MakerFee        string `json:"maker_fee"`
	TakerFee        string `json:"taker_fee"`
	ThirtyDayVolume string `json:"thirty_day_volume"`
}

// GetFeeInfo makes a call to GET /api/1/fee_info.
//
// Returns your fees and 30 day trading volume (as of midnight) for a given
// pair.
//
// Permissions required: <code>Perm_R_Orders</code>
func (cl *Client) GetFeeInfo(ctx context.Context, req *GetFeeInfoRequest) (*GetFeeInfoResponse, error) {
	var res GetFeeInfoResponse
	err := cl.do(ctx, "GET", "/api/1/fee_info", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetFundingAddressRequest struct {
	// Currency code of the asset.
	Asset string `json:"asset" url:"asset"`

	// Specific Bitcoin or Ethereum address to retrieve. If not provided, the
	// default address will be used.
	Address string `json:"address" url:"address"`
}

type GetFundingAddressResponse struct {
	AccountId        string  `json:"account_id"`
	Address          string  `json:"address"`
	Asset            string  `json:"asset"`
	AssignedAt       int64   `json:"assigned_at"`
	Name             string  `json:"name"`
	ReceiveFee       float64 `json:"receive_fee,string"`
	TotalReceived    float64 `json:"total_received,string"`
	TotalUnconfirmed float64 `json:"total_unconfirmed,string"`
}

// GetFundingAddress makes a call to GET /api/1/funding_address.
//
// Returns the default receive address associated with your account and the
// amount received via the address. You can specify an optional address
// parameter to return information for a non-default receive address. In the
// response, total_received is the total confirmed Bitcoin amount received
// excluding unconfirmed transactions. total_unconfirmed is the total sum of
// unconfirmed receive transactions.
//
// Permissions required: <code>Perm_R_Addresses</code>
func (cl *Client) GetFundingAddress(ctx context.Context, req *GetFundingAddressRequest) (*GetFundingAddressResponse, error) {
	var res GetFundingAddressResponse
	err := cl.do(ctx, "GET", "/api/1/funding_address", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetOrderRequest struct {
	// The order ID.
	Id string `json:"id" url:"id"`
}

type GetOrderResponse struct {
	Base                float64 `json:"base,string"`
	Btc                 string  `json:"btc"`
	CompletedTimestamp  int64   `json:"completed_timestamp"`
	Counter             float64 `json:"counter,string"`
	CreationTimestamp   int64   `json:"creation_timestamp"`
	ExpirationTimestamp int64   `json:"expiration_timestamp"`
	FeeBase             float64 `json:"fee_base,string"`
	FeeBtc              string  `json:"fee_btc"`
	FeeCounter          float64 `json:"fee_counter,string"`
	FeeZar              string  `json:"fee_zar"`
	LimitPrice          float64 `json:"limit_price,string"`
	LimitVolume         float64 `json:"limit_volume,string"`
	OrderId             string  `json:"order_id"`

	// Specifies the market.
	Pair string `json:"pair"`

	// <code>PENDING</code> The order has been placed. Some trades may have
	// taken place but the order is not filled yet.<br>
	// <code>COMPLETE</code> The order is no longer active. It has been settled
	// or has been cancelled.
	State  string  `json:"state"`
	Trades []Trade `json:"trades"`

	// <code>BID</code> bid (buy) limit order.<br>
	// <code>ASK</code> ask (sell) limit order.
	Type string `json:"type"`
	Zar  string `json:"zar"`
}

// GetOrder makes a call to GET /api/1/orders/{id}.
//
// Get an order by its ID.
//
// Permissions required: <code>Perm_R_Orders</code>
func (cl *Client) GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error) {
	var res GetOrderResponse
	err := cl.do(ctx, "GET", "/api/1/orders/{id}", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetOrderBookRequest struct {
	// Currency pair
	Pair string `json:"pair" url:"pair"`
}

type GetOrderBookResponse struct {
	Asks      []OrderBookEntry `json:"asks"`
	Bids      []OrderBookEntry `json:"bids"`
	Timestamp int64            `json:"timestamp"`
}

type OrderBookEntry struct {
	Price  float64 `json:"price,string"`
	Volume float64 `json:"volume,string"`
}

// GetOrderBook makes a call to GET /api/1/orderbook.
//
// Returns a list of bids and asks in the order book. Ask orders are sorted by
// price ascending. Bid orders are sorted by price descending. Note that
// multiple orders at the same price are not necessarily conflated.
func (cl *Client) GetOrderBook(ctx context.Context, req *GetOrderBookRequest) (*GetOrderBookResponse, error) {
	var res GetOrderBookResponse
	err := cl.do(ctx, "GET", "/api/1/orderbook", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetQuoteRequest struct {
	// ID of the quote to retrieve.
	Id int64 `json:"id" url:"id"`
}

type GetQuoteResponse struct {
	BaseAmount    float64 `json:"base_amount,string"`
	CounterAmount float64 `json:"counter_amount,string"`
	CreatedAt     int64   `json:"created_at"`
	Discarded     bool    `json:"discarded"`
	Exercised     bool    `json:"exercised"`
	ExpiresAt     int64   `json:"expires_at"`
	Id            string  `json:"id"`
	Pair          string  `json:"pair"`
	Type          string  `json:"type"`
}

// GetQuote makes a call to GET /api/1/quotes/{id}.
//
// Get the latest status of a quote.
//
// Permissions required: <code>Perm_R_Orders</code>
func (cl *Client) GetQuote(ctx context.Context, req *GetQuoteRequest) (*GetQuoteResponse, error) {
	var res GetQuoteResponse
	err := cl.do(ctx, "GET", "/api/1/quotes/{id}", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetTickerRequest struct {
	// Currency pair
	Pair string `json:"pair" url:"pair"`
}

type GetTickerResponse struct {
	Ask                 float64 `json:"ask,string"`
	Bid                 float64 `json:"bid,string"`
	LastTrade           float64 `json:"last_trade,string"`
	Pair                string  `json:"pair"`
	Rolling24HourVolume float64 `json:"rolling_24_hour_volume,string"`
	Timestamp           int64   `json:"timestamp"`
}

// GetTicker makes a call to GET /api/1/ticker.
//
// Returns the latest ticker indicators.
func (cl *Client) GetTicker(ctx context.Context, req *GetTickerRequest) (*GetTickerResponse, error) {
	var res GetTickerResponse
	err := cl.do(ctx, "GET", "/api/1/ticker", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetTickersRequest struct {
}

type GetTickersResponse struct {
	Tickers []Ticker `json:"tickers"`
}

type Ticker struct {
	Ask                 float64 `json:"ask,string"`
	Bid                 float64 `json:"bid,string"`
	LastTrade           float64 `json:"last_trade,string"`
	Pair                string  `json:"pair"`
	Rolling24HourVolume float64 `json:"rolling_24_hour_volume,string"`
	Timestamp           int64   `json:"timestamp"`
}

// GetTickers makes a call to GET /api/1/tickers.
//
// Returns the latest ticker indicators from all active Luno exchanges.
func (cl *Client) GetTickers(ctx context.Context, req *GetTickersRequest) (*GetTickersResponse, error) {
	var res GetTickersResponse
	err := cl.do(ctx, "GET", "/api/1/tickers", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetWithdrawalRequest struct {
	// Withdrawal ID to retrieve.
	Id int64 `json:"id" url:"id"`
}

type GetWithdrawalResponse struct {
	Amount    float64 `json:"amount,string"`
	CreatedAt int64   `json:"created_at"`
	Currency  string  `json:"currency"`
	Fee       float64 `json:"fee,string"`
	Id        string  `json:"id"`
	Status    string  `json:"status"`
	Type      string  `json:"type"`
}

// GetWithdrawal makes a call to GET /api/1/withdrawals/{id}.
//
// Returns the status of a particular withdrawal request.
//
// Permissions required: <code>Perm_R_Withdrawals</code>
func (cl *Client) GetWithdrawal(ctx context.Context, req *GetWithdrawalRequest) (*GetWithdrawalResponse, error) {
	var res GetWithdrawalResponse
	err := cl.do(ctx, "GET", "/api/1/withdrawals/{id}", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetWithdrawalsRequest struct {
}

type GetWithdrawalsResponse struct {
	Withdrawals []Withdrawal `json:"withdrawals"`
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

// GetWithdrawals makes a call to GET /api/1/withdrawals.
//
// Returns a list of withdrawal requests.
//
// Permissions required: <code>Perm_R_Withdrawals</code>
func (cl *Client) GetWithdrawals(ctx context.Context, req *GetWithdrawalsRequest) (*GetWithdrawalsResponse, error) {
	var res GetWithdrawalsResponse
	err := cl.do(ctx, "GET", "/api/1/withdrawals", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type OrderState string

const (
	OrderStatePending  OrderState = "PENDING"
	OrderStateComplete OrderState = "COMPLETE"
)

type ListOrdersRequest struct {
	// Filter to orders created before this timestamp (Unix milliseconds)
	CreatedBefore int64 `json:"created_before" url:"created_before"`

	// Limit to this many orders
	Limit int64 `json:"limit" url:"limit"`

	// Filter to only orders of this currency pair
	Pair string `json:"pair" url:"pair"`

	// Filter to only orders of this state
	State OrderState `json:"state" url:"state"`
}

type ListOrdersResponse struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	Base                float64 `json:"base,string"`
	Btc                 string  `json:"btc"`
	CompletedTimestamp  int64   `json:"completed_timestamp"`
	Counter             float64 `json:"counter,string"`
	CreationTimestamp   int64   `json:"creation_timestamp"`
	ExpirationTimestamp int64   `json:"expiration_timestamp"`
	FeeBase             float64 `json:"fee_base,string"`
	FeeBtc              string  `json:"fee_btc"`
	FeeCounter          float64 `json:"fee_counter,string"`
	FeeZar              string  `json:"fee_zar"`
	LimitPrice          float64 `json:"limit_price,string"`
	LimitVolume         float64 `json:"limit_volume,string"`
	OrderId             string  `json:"order_id"`

	// Specifies the market.
	Pair string `json:"pair"`

	// <code>PENDING</code> The order has been placed. Some trades may have
	// taken place but the order is not filled yet.<br>
	// <code>COMPLETE</code> The order is no longer active. It has been settled
	// or has been cancelled.
	State  string  `json:"state"`
	Trades []Trade `json:"trades"`

	// <code>BID</code> bid (buy) limit order.<br>
	// <code>ASK</code> ask (sell) limit order.
	Type string `json:"type"`
	Zar  string `json:"zar"`
}

// ListOrders makes a call to GET /api/1/listorders.
//
// Returns a list of the most recently placed orders. You can specify an
// optional <code>state=PENDING</code> parameter to restrict the results to only
// open orders. You can also specify the market by using the optional pair
// parameter. The list is truncated after 100 items.
//
// Permissions required: <code>Perm_R_Orders</code>
func (cl *Client) ListOrders(ctx context.Context, req *ListOrdersRequest) (*ListOrdersResponse, error) {
	var res ListOrdersResponse
	err := cl.do(ctx, "GET", "/api/1/listorders", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type ListPendingTransactionsRequest struct {
	// Account ID
	Id int64 `json:"id" url:"id"`
}

type ListPendingTransactionsResponse struct {
	Balance          AccountBalance      `json:"balance"`
	Capabilities     AccountCapabilities `json:"capabilities"`
	Currency         string              `json:"currency"`
	Icon             string              `json:"icon"`
	Id               string              `json:"id"`
	IsDefault        bool                `json:"is_default"`
	Name             string              `json:"name"`
	Pending          []Transaction       `json:"pending"`
	ReceiveAddresses []ReceiveAddress    `json:"receive_addresses"`
	Transactions     []Transaction       `json:"transactions"`
}

// ListPendingTransactions makes a call to GET /api/1/accounts/{id}/pending.
//
// Return a list of all pending transactions related to the account.
//
// Unlike account entries, pending transactions are not numbered, and may be
// reordered, deleted or updated at any time.
//
// Permissions required: <code>Perm_R_Transactions</code>
func (cl *Client) ListPendingTransactions(ctx context.Context, req *ListPendingTransactionsRequest) (*ListPendingTransactionsResponse, error) {
	var res ListPendingTransactionsResponse
	err := cl.do(ctx, "GET", "/api/1/accounts/{id}/pending", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type ListTradesRequest struct {
	// Currency pair
	Pair string `json:"pair" url:"pair"`

	// Fetch trades executed after this time, specified as a Unix timestamp in
	// milliseconds.
	Since int64 `json:"since" url:"since"`
}

type ListTradesResponse struct {
	Currency string  `json:"currency"`
	Trades   []Trade `json:"trades"`
}

// ListTrades makes a call to GET /api/1/trades.
//
// Returns a list of the most recent trades. At most 100 results are returned
// per call.
func (cl *Client) ListTrades(ctx context.Context, req *ListTradesRequest) (*ListTradesResponse, error) {
	var res ListTradesResponse
	err := cl.do(ctx, "GET", "/api/1/trades", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type ListTransactionsRequest struct {
	// Account ID
	Id int64 `json:"id" url:"id"`

	// Maximum of the row range to return (exclusive)
	MaxRow int64 `json:"max_row" url:"max_row"`

	// Minimum of the row range to return (inclusive)
	MinRow int64 `json:"min_row" url:"min_row"`
}

type ListTransactionsResponse struct {
	Balance          AccountBalance      `json:"balance"`
	Capabilities     AccountCapabilities `json:"capabilities"`
	Currency         string              `json:"currency"`
	Icon             string              `json:"icon"`
	Id               string              `json:"id"`
	IsDefault        bool                `json:"is_default"`
	Name             string              `json:"name"`
	Pending          []Transaction       `json:"pending"`
	ReceiveAddresses []ReceiveAddress    `json:"receive_addresses"`
	Transactions     []Transaction       `json:"transactions"`
}

// ListTransactions makes a call to GET /api/1/{id}/transactions.
//
// Return a list of transaction entries from an account.
//
// Transaction entry rows are numbered sequentially starting from 1, where 1 is
// the oldest entry. The range of rows to return are specified with the
// <code>min_row</code> (inclusive) and <code>max_row</code> (exclusive)
// parameters. At most 1000 rows can be requested per call.
//
// If <code>min_row</code> or <code>max_row</code> is non-positive, the range
// wraps around the most recent row. For example, to fetch the 100 most recent
// rows, use <code>min_row=-100</code> and <code>max_row=0</code>.
//
// Permissions required: <code>Perm_R_Transactions</code>
func (cl *Client) ListTransactions(ctx context.Context, req *ListTransactionsRequest) (*ListTransactionsResponse, error) {
	var res ListTransactionsResponse
	err := cl.do(ctx, "GET", "/api/1/{id}/transactions", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type ListUserTradesRequest struct {
	// Filter to trades of this currency pair.
	Pair string `json:"pair" url:"pair"`

	// Limit to this number of trades (default 100).
	Limit int64 `json:"limit" url:"limit"`

	// Filter to trades on or after this timestamp.
	Since int64 `json:"since" url:"since"`
}

type ListUserTradesResponse struct {
	Trades []Trade `json:"trades"`
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

// ListUserTrades makes a call to GET /api/1/listtrades.
//
// Returns a list of your recent trades for a given pair, sorted by oldest
// first.
//
// <code>type</code> in the response indicates the type of order that you placed
// in order to participate in the trade. Possible types: <code>BID</code>,
// <code>ASK</code>.
//
// If <code>is_buy</code> in the response is true, then the order which
// completed the trade (market taker) was a bid order.
//
// Results of this query may lag behind the latest data.
//
// Permissions required: <code>Perm_R_Orders</code>
func (cl *Client) ListUserTrades(ctx context.Context, req *ListUserTradesRequest) (*ListUserTradesResponse, error) {
	var res ListUserTradesResponse
	err := cl.do(ctx, "GET", "/api/1/listtrades", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

const (
	OrderTypeBid OrderType = "BID"
	OrderTypeAsk OrderType = "ASK"
)

type PostLimitOrderRequest struct {
	// The currency pair to trade.
	Pair string `json:"pair" url:"pair"`

	// Limit price as a decimal string in units of ZAR/BTC.
	Price float64 `json:"price,string" url:"price"`

	// <code>BID</code> for a bid (buy) limit order<br>
	// <code>ASK</code> for ab ask (sell) limit order
	Type OrderType `json:"type" url:"type"`

	// Amount of Bitcoin or Ethereum to buy or sell as a decimal string in units
	// of the currency.
	Volume float64 `json:"volume,string" url:"volume"`

	// The base currency account to use in the trade.
	BaseAccountId string `json:"base_account_id" url:"base_account_id"`

	// The counter currency account to use in the trade.
	CounterAccountId string `json:"counter_account_id" url:"counter_account_id"`
}

type PostLimitOrderResponse struct {
	OrderId string `json:"order_id"`
}

// PostLimitOrder makes a call to POST /api/1/postorder.
//
// Create a new trade order.
//
// Warning! Orders cannot be reversed once they have executed. Please ensure
// your program has been thoroughly tested before submitting orders.
//
// If no <code>base_account_id</code> or <code>counter_account_id</code> are
// specified, your default base currency or counter currency account will be
// used. You can find your account IDs by calling the
// <a href="#operation/getBalances">Balances</a> API.
//
// Permissions required: <code>Perm_W_Orders</code>
func (cl *Client) PostLimitOrder(ctx context.Context, req *PostLimitOrderRequest) (*PostLimitOrderResponse, error) {
	var res PostLimitOrderResponse
	err := cl.do(ctx, "POST", "/api/1/postorder", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type OrderType string

const (
	OrderTypeBuy  OrderType = "BUY"
	OrderTypeSell OrderType = "SELL"
)

type PostMarketOrderRequest struct {
	// The currency pair to trade.
	Pair string `json:"pair" url:"pair"`

	// <code>BUY</code> to buy Bitcoin or Ethereum<br>
	// <code>SELL</code> to sell Bitcoin or Ethereum
	Type OrderType `json:"type" url:"type"`

	// The base currency account to use in the trade.
	BaseAccountId string `json:"base_account_id" url:"base_account_id"`

	// For a <code>SELL</code> order: amount of Bitcoin to sell as a decimal
	// string in units of BTC or ETH.
	BaseVolume float64 `json:"base_volume,string" url:"base_volume"`

	// The counter currency account to use in the trade.
	CounterAccountId string `json:"counter_account_id" url:"counter_account_id"`

	// For a <code>BUY</code> order: amount of local currency (e.g. ZAR, MYR) to
	// spend as a decimal string in units of the local currency.
	CounterVolume float64 `json:"counter_volume,string" url:"counter_volume"`
}

type PostMarketOrderResponse struct {
	OrderId string `json:"order_id"`
}

// PostMarketOrder makes a call to POST /api/1/marketorder.
//
// Create a new market order.
//
// A market order executes immediately, and either buys as much Bitcoin or
// Ethereum that can be bought for a set amount of fiat currency, or sells a
// set amount of Bitcoin or Ethereum for as much fiat as possible.
//
// Warning! Orders cannot be reversed once they have executed. Please ensure
// your program has been thoroughly tested before submitting orders.
//
// If no base_account_id or counter_account_id are specified, your default base
// currency or counter currency account will be used. You can find your account
// IDs by calling the <a href="#operation/getBalances">Balances</a> API.
//
// Permissions required: <code>Perm_W_Orders</code>
func (cl *Client) PostMarketOrder(ctx context.Context, req *PostMarketOrderRequest) (*PostMarketOrderResponse, error) {
	var res PostMarketOrderResponse
	err := cl.do(ctx, "POST", "/api/1/marketorder", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type SendRequest struct {
	// Destination Bitcoin address or email address, or Ethereum address to send
	// to.
	//
	// Note:
	// <ul>
	// <li>Ethereum addresses must be
	// <a href="https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md" target="_blank" rel="nofollow">checksummed</a>.</li>
	// <li>Ethereum sends to email addresses are not supported.</li>
	// </ul>
	Address string `json:"address" url:"address"`

	// Amount to send as a decimal string.
	Amount float64 `json:"amount,string" url:"amount"`

	// Currency to send.
	Currency string `json:"currency" url:"currency"`

	// Description for the transaction to record on the account statement.
	Description string `json:"description" url:"description"`

	// Message to send to the recipient. This is only relevant when sending to
	// an email address.
	Message string `json:"message" url:"message"`
}

type SendResponse struct {
	Success      bool   `json:"success"`
	WithdrawalId string `json:"withdrawal_id"`
}

// Send makes a call to POST /api/1/send.
//
// Send Bitcoin from your account to a Bitcoin address or email address. Send
// Ethereum from your account to an Ethereum address.
//
// If the email address is not associated with an existing Luno account, an
// invitation to create an account and claim the funds will be sent.
//
// Warning! Digital currency transactions are irreversible. Please ensure your
// program has been thoroughly tested before using this call.
//
// Permissions required: <code>Perm_W_Send</code>
func (cl *Client) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	var res SendResponse
	err := cl.do(ctx, "POST", "/api/1/send", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type StopOrderRequest struct {
	// The order reference as a string.
	OrderId string `json:"order_id" url:"order_id"`
}

type StopOrderResponse struct {
	Success bool `json:"success"`
}

// StopOrder makes a call to POST /api/1/stoporder.
//
// Request to stop an order.
//
// Permissions required: <code>Perm_W_Orders</code>
func (cl *Client) StopOrder(ctx context.Context, req *StopOrderRequest) (*StopOrderResponse, error) {
	var res StopOrderResponse
	err := cl.do(ctx, "POST", "/api/1/stoporder", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// vi: ft=go
