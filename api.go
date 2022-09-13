package luno

import (
	"context"

	"github.com/luno/luno-go/decimal"
)

// CancelWithdrawalRequest is the request struct for CancelWithdrawal.
type CancelWithdrawalRequest struct {
	// ID of the withdrawal to cancel.
	//
	// required: true
	Id int64 `json:"id" url:"id"`
}

// CancelWithdrawalResponse is the response struct for CancelWithdrawal.
type CancelWithdrawalResponse struct {
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

	// Type distinguishes between different withdrawal methods where more than one is supported
	// for the given currency.
	Type string `json:"type"`
}

// CancelWithdrawal makes a call to DELETE /api/1/withdrawals/{id}.
//
// Cancels a withdrawal request.
// This can only be done if the request is still in state <code>PENDING</code>.
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

// CreateAccountRequest is the request struct for CreateAccount.
type CreateAccountRequest struct {
	// The currency code for the Account you want to create.  Please see the Currency section for a detailed list of currencies supported by the Luno platform.
	//
	// Users must be verified to trade currency in order to be able to create an Account.  For more information on the verification process, please see <a href="/help/en/articles/1000168396">How do I verify my identity?</a>.
	//
	// Users have a limit of 4 accounts per currency.
	//
	// required: true
	Currency string `json:"currency" url:"currency"`

	// The label to use for this account
	//
	// required: true
	Name string `json:"name" url:"name"`
}

// CreateAccountResponse is the response struct for CreateAccount.
type CreateAccountResponse struct {
	Currency     string        `json:"currency"`
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	Pending      []Transaction `json:"pending"`
	Transactions []Transaction `json:"transactions"`
}

// CreateAccount makes a call to POST /api/1/accounts.
//
// This request creates an Account for the specified currency.  Please note that the balances for the Account will be displayed based on the <code>asset</code> value, which is the currency the Account is based on.
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

// CreateFundingAddressRequest is the request struct for CreateFundingAddress.
type CreateFundingAddressRequest struct {
	// Currency code of the asset.
	//
	// required: true
	Asset string `json:"asset" url:"asset"`

	// An optional name for the new Receive Address
	Name string `json:"name" url:"name"`
}

// CreateFundingAddressResponse is the response struct for CreateFundingAddress.
type CreateFundingAddressResponse struct {
	AccountId        string          `json:"account_id"`
	Address          string          `json:"address"`
	AddressMeta      []AddressMeta   `json:"address_meta"`
	Asset            string          `json:"asset"`
	AssignedAt       Time            `json:"assigned_at"`
	Name             string          `json:"name"`
	QrCodeUri        string          `json:"qr_code_uri"`
	ReceiveFee       decimal.Decimal `json:"receive_fee"`
	TotalReceived    decimal.Decimal `json:"total_received"`
	TotalUnconfirmed decimal.Decimal `json:"total_unconfirmed"`
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

// CreateWithdrawalRequest is the request struct for CreateWithdrawal.
type CreateWithdrawalRequest struct {
	// Amount to withdraw. The currency withdrawn depends on the type setting.
	//
	// required: true
	Amount decimal.Decimal `json:"amount" url:"amount"`

	// Withdrawal type.
	//
	// required: true
	Type string `json:"type" url:"type"`

	// The beneficiary ID of the bank account the withdrawal will be paid out to.
	// This parameter is required if the user has set up multiple beneficiaries.
	// The beneficiary ID can be found by selecting on the beneficiary name on the user’s <a href="/wallet/beneficiaries">Beneficiaries</a> page.
	BeneficiaryId int64 `json:"beneficiary_id" url:"beneficiary_id"`

	// Optional unique ID to associate with this withdrawal.
	// Useful to prevent duplicate sends.
	// This field supports all alphanumeric characters including "-" and "_".
	ExternalId string `json:"external_id" url:"external_id"`

	// If true, it will be a fast withdrawal if possible. Fast withdrawals come with a fee.
	// Currently fast withdrawals are only available for `type=ZAR_EFT`; for other types, an error is returned.
	// Fast withdrawals are not possible for Bank of Baroda, Deutsche Bank, Merrill Lynch South Africa, UBS, Postbank and Tyme Bank.
	// The fee to be charged is the same as when withdrawing from the UI.
	Fast bool `json:"fast" url:"fast"`

	// For internal use.
	Reference string `json:"reference" url:"reference"`
}

// CreateWithdrawalResponse is the response struct for CreateWithdrawal.
type CreateWithdrawalResponse struct {
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

	// Type distinguishes between different withdrawal methods where more than one is supported
	// for the given currency.
	Type string `json:"type"`
}

// CreateWithdrawal makes a call to POST /api/1/withdrawals.
//
// Creates a new withdrawal request to the specified beneficiary.
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

// GetBalancesRequest is the request struct for GetBalances.
type GetBalancesRequest struct {
	// Only return balances for wallets with these currencies (if not provided,
	// all balances will be returned). To request balances for multiple currencies,
	// pass the parameter multiple times,
	// e.g. `assets=XBT&assets=ETH`.
	Assets []string `json:"assets" url:"assets"`
}

// GetBalancesResponse is the response struct for GetBalances.
type GetBalancesResponse struct {
	Balance []AccountBalance `json:"balance"`
}

// GetBalances makes a call to GET /api/1/balance.
//
// The list of all Accounts and their respective balances for the requesting user.
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

// GetCandlesRequest is the request struct for GetCandles.
type GetCandlesRequest struct {
	// Candle duration in seconds.
	// For example, 300 corresponds to 5m candles. Currently supported
	// durations are: 60 (1m), 300 (5m), 900 (15m), 1800 (30m), 3600 (1h),
	// 10800 (3h), 14400 (4h), 28800 (8h), 86400 (24h), 259200 (3d), 604800
	// (7d).
	//
	// required: true
	Duration int64 `json:"duration" url:"duration"`

	// Currency pair
	//
	// required: true
	Pair string `json:"pair" url:"pair"`

	// Filter to candles starting on or after this timestamp (Unix milliseconds).
	// Only up to 1000 of the earliest candles are returned.
	//
	// required: true
	Since Time `json:"since" url:"since"`
}

// GetCandlesResponse is the response struct for GetCandles.
type GetCandlesResponse struct {
	Candles []Candle `json:"candles"`

	// Duration in seconds
	Duration int64  `json:"duration"`
	Pair     string `json:"pair"`
}

// GetCandles makes a call to GET /api/exchange/1/candles.
//
// Get candlestick market data from the specified time until now, from the oldest to the most recent.
//
// Permissions required: <code>MP_None</code>
func (cl *Client) GetCandles(ctx context.Context, req *GetCandlesRequest) (*GetCandlesResponse, error) {
	var res GetCandlesResponse
	err := cl.do(ctx, "GET", "/api/exchange/1/candles", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetFeeInfoRequest is the request struct for GetFeeInfo.
type GetFeeInfoRequest struct {
	// Get fee information about this pair.
	//
	// required: true
	Pair string `json:"pair" url:"pair"`
}

// GetFeeInfoResponse is the response struct for GetFeeInfo.
type GetFeeInfoResponse struct {
	MakerFee        string `json:"maker_fee"`
	TakerFee        string `json:"taker_fee"`
	ThirtyDayVolume string `json:"thirty_day_volume"`
}

// GetFeeInfo makes a call to GET /api/1/fee_info.
//
// Returns the fees and 30 day trading volume (as of midnight) for a given currency pair.  For complete details, please see <a href="en/countries">Fees & Features</a>.
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

// GetFundingAddressRequest is the request struct for GetFundingAddress.
type GetFundingAddressRequest struct {
	// Currency code of the asset.
	//
	// required: true
	Asset string `json:"asset" url:"asset"`

	// Specific cryptocurrency address to retrieve. If not provided, the
	// default address will be used.
	Address string `json:"address" url:"address"`
}

// GetFundingAddressResponse is the response struct for GetFundingAddress.
type GetFundingAddressResponse struct {
	AccountId        string          `json:"account_id"`
	Address          string          `json:"address"`
	AddressMeta      []AddressMeta   `json:"address_meta"`
	Asset            string          `json:"asset"`
	AssignedAt       Time            `json:"assigned_at"`
	Name             string          `json:"name"`
	QrCodeUri        string          `json:"qr_code_uri"`
	ReceiveFee       decimal.Decimal `json:"receive_fee"`
	TotalReceived    decimal.Decimal `json:"total_received"`
	TotalUnconfirmed decimal.Decimal `json:"total_unconfirmed"`
}

// GetFundingAddress makes a call to GET /api/1/funding_address.
//
// Returns the default receive address associated with your account and the
// amount received via the address. Users can specify an optional address parameter to return information for a non-default receive address.
// In the response, <code>total_received</code> is the total confirmed amount received excluding unconfirmed transactions.
// <code>total_unconfirmed</code> is the total sum of unconfirmed receive transactions.
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

// GetMoveRequest is the request struct for GetMove.
type GetMoveRequest struct {
	// Get by the user defined ID. This is mutually exclusive with <code>id</code> and is required if <code>id</code> is
	// not provided.
	ClientMoveId string `json:"client_move_id" url:"client_move_id"`

	// Get by the system ID. This is mutually exclusive with <code>client_move_id</code> and is required if
	// <code>client_move_id</code> is not provided.
	Id string `json:"id" url:"id"`
}

// GetMoveResponse is the response struct for GetMove.
type GetMoveResponse struct {
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

// GetMove makes a call to GET /api/exchange/1/move.
//
// Get a specific move funds instruction by either <code>id</code> or
// <code>client_move_id</code>. If both are provided an API error will be
// returned.
//
// This endpoint is in BETA, behaviour and specification may change without
// any previous notice.
//
// Permissions required: <code>MP_None</code>
func (cl *Client) GetMove(ctx context.Context, req *GetMoveRequest) (*GetMoveResponse, error) {
	var res GetMoveResponse
	err := cl.do(ctx, "GET", "/api/exchange/1/move", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetOrderRequest is the request struct for GetOrder.
type GetOrderRequest struct {
	// Order reference
	//
	// required: true
	Id string `json:"id" url:"id"`
}

// GetOrderResponse is the response struct for GetOrder.
type GetOrderResponse struct {
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

// GetOrder makes a call to GET /api/1/orders/{id}.
//
// Get an Order's details by its ID.
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

// GetOrderBookRequest is the request struct for GetOrderBook.
type GetOrderBookRequest struct {
	// Currency pair of the Orders to retrieve
	//
	// required: true
	Pair string `json:"pair" url:"pair"`
}

// GetOrderBookResponse is the response struct for GetOrderBook.
type GetOrderBookResponse struct {
	// List of asks sorted from lowest to highest price
	Asks []OrderBookEntry `json:"asks"`

	// List of bids sorted from highest to lowest price
	Bids []OrderBookEntry `json:"bids"`

	// Unix timestamp in milliseconds
	Timestamp int64 `json:"timestamp"`
}

// GetOrderBook makes a call to GET /api/1/orderbook_top.
//
// This request returns the best 100 `bids` and `asks`, for the currency pair specified, in the Order Book.
//
// `asks` are sorted by price ascending and `bids` are sorted by price descending.
//
// Multiple orders at the same price are aggregated.
func (cl *Client) GetOrderBook(ctx context.Context, req *GetOrderBookRequest) (*GetOrderBookResponse, error) {
	var res GetOrderBookResponse
	err := cl.do(ctx, "GET", "/api/1/orderbook_top", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetOrderBookFullRequest is the request struct for GetOrderBookFull.
type GetOrderBookFullRequest struct {
	// Currency pair of the Orders to retrieve
	//
	// required: true
	Pair string `json:"pair" url:"pair"`
}

// GetOrderBookFullResponse is the response struct for GetOrderBookFull.
type GetOrderBookFullResponse struct {
	// List of asks sorted from lowest to highest price
	Asks []OrderBookEntry `json:"asks"`

	// List of bids sorted from highest to lowest price
	Bids []OrderBookEntry `json:"bids"`

	// Unix timestamp in milliseconds
	Timestamp int64 `json:"timestamp"`
}

// GetOrderBookFull makes a call to GET /api/1/orderbook.
//
// This request returns all `bids` and `asks`, for the currency pair specified, in the Order Book.
//
// `asks` are sorted by price ascending and `bids` are sorted by price descending.
//
// Multiple orders at the same price are not aggregated.
//
// <b>WARNING:</b> This may return a large amount of data.
// Users are recommended to use the <a href="#operation/getOrderBookTop">top 100 bids and asks</a>
// or the <a href="#tag/Streaming-API">Streaming API</a>.
func (cl *Client) GetOrderBookFull(ctx context.Context, req *GetOrderBookFullRequest) (*GetOrderBookFullResponse, error) {
	var res GetOrderBookFullResponse
	err := cl.do(ctx, "GET", "/api/1/orderbook", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetOrderV2Request is the request struct for GetOrderV2.
type GetOrderV2Request struct {
	// Order reference
	//
	// required: true
	Id string `json:"id" url:"id"`
}

// GetOrderV2Response is the response struct for GetOrderV2.
type GetOrderV2Response struct {
	// Amount of base filled, this value is always positive.
	//
	// Use this field and `side` to determine credit or debit of funds.
	Base decimal.Decimal `json:"base"`

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

// GetOrderV2 makes a call to GET /api/exchange/2/orders/{id}.
//
// Get the details for an order.
//
// Permissions required: <code>Perm_R_Orders</code>
func (cl *Client) GetOrderV2(ctx context.Context, req *GetOrderV2Request) (*GetOrderV2Response, error) {
	var res GetOrderV2Response
	err := cl.do(ctx, "GET", "/api/exchange/2/orders/{id}", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetOrderV3Request is the request struct for GetOrderV3.
type GetOrderV3Request struct {
	// Client Order ID has the value that was passed in when the Order was posted.
	ClientOrderId string `json:"client_order_id" url:"client_order_id"`

	// Order reference
	Id string `json:"id" url:"id"`
}

// GetOrderV3Response is the response struct for GetOrderV3.
type GetOrderV3Response struct {
	// Amount of base filled, this value is always positive.
	//
	// Use this field and `side` to determine credit or debit of funds.
	Base decimal.Decimal `json:"base"`

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

// GetOrderV3 makes a call to GET /api/exchange/3/order.
//
// Get the details for an order by order reference or client order ID.
// Exactly one of the two parameters must be provided, otherwise an error is returned.
// Permissions required: <code>Perm_R_Orders</code>
func (cl *Client) GetOrderV3(ctx context.Context, req *GetOrderV3Request) (*GetOrderV3Response, error) {
	var res GetOrderV3Response
	err := cl.do(ctx, "GET", "/api/exchange/3/order", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTickerRequest is the request struct for GetTicker.
type GetTickerRequest struct {
	// Currency pair
	//
	// required: true
	Pair string `json:"pair" url:"pair"`
}

// GetTickerResponse is the response struct for GetTicker.
type GetTickerResponse struct {
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

// GetTicker makes a call to GET /api/1/ticker.
//
// Returns the latest ticker indicators for the specified currency pair.
//
// Please see the <a href="#tag/currency ">Currency list</a> for the complete list of supported currency pairs.
func (cl *Client) GetTicker(ctx context.Context, req *GetTickerRequest) (*GetTickerResponse, error) {
	var res GetTickerResponse
	err := cl.do(ctx, "GET", "/api/1/ticker", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTickersRequest is the request struct for GetTickers.
type GetTickersRequest struct {
	// Return tickers for multiple markets (if not provided, all tickers will be returned).
	// To request tickers for multiple markets, pass the parameter multiple times,
	// e.g. `pair=XBTZAR&pair=ETHZAR`.
	Pair []string `json:"pair" url:"pair"`
}

// GetTickersResponse is the response struct for GetTickers.
type GetTickersResponse struct {
	Tickers []Ticker `json:"tickers"`
}

// GetTickers makes a call to GET /api/1/tickers.
//
// Returns the latest ticker indicators from all active Luno exchanges.
//
// Please see the <a href="#tag/currency ">Currency list</a> for the complete list of supported currency pairs.
func (cl *Client) GetTickers(ctx context.Context, req *GetTickersRequest) (*GetTickersResponse, error) {
	var res GetTickersResponse
	err := cl.do(ctx, "GET", "/api/1/tickers", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetWithdrawalRequest is the request struct for GetWithdrawal.
type GetWithdrawalRequest struct {
	// Withdrawal ID to retrieve.
	//
	// required: true
	Id int64 `json:"id" url:"id"`
}

// GetWithdrawalResponse is the response struct for GetWithdrawal.
type GetWithdrawalResponse struct {
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

	// Type distinguishes between different withdrawal methods where more than one is supported
	// for the given currency.
	Type string `json:"type"`
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

// ListBeneficiariesResponseRequest is the request struct for ListBeneficiariesResponse.
type ListBeneficiariesResponseRequest struct {
}

// ListBeneficiariesResponseResponse is the response struct for ListBeneficiariesResponse.
type ListBeneficiariesResponseResponse struct {
	Beneficiaries []beneficiary `json:"beneficiaries"`
}

// ListBeneficiariesResponse makes a call to GET /api/1/beneficiaries.
//
// Returns a list of bank beneficiaries.
//
// Permissions required: <code>Perm_R_Beneficiaries</code>
func (cl *Client) ListBeneficiariesResponse(ctx context.Context, req *ListBeneficiariesResponseRequest) (*ListBeneficiariesResponseResponse, error) {
	var res ListBeneficiariesResponseResponse
	err := cl.do(ctx, "GET", "/api/1/beneficiaries", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ListMovesRequest is the request struct for ListMoves.
type ListMovesRequest struct {
	// Filter to moves requested before this timestamp (Unix milliseconds)
	Before int64 `json:"before" url:"before"`

	// Limit to this many moves
	Limit int64 `json:"limit" url:"limit"`
}

// ListMovesResponse is the response struct for ListMoves.
type ListMovesResponse struct {
	Moves []FundsMove `json:"moves"`
}

// ListMoves makes a call to GET /api/exchange/1/move/list_moves.
//
// Returns a list of the most recent moves ordered from newest to oldest.
// This endpoint will list up to 100 most recent moves by default.
//
// This endpoint is in BETA, behaviour and specification may change without
// any previous notice.
//
// Permissions required: <code>MP_None</code>
func (cl *Client) ListMoves(ctx context.Context, req *ListMovesRequest) (*ListMovesResponse, error) {
	var res ListMovesResponse
	err := cl.do(ctx, "GET", "/api/exchange/1/move/list_moves", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ListOrdersRequest is the request struct for ListOrders.
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

// ListOrdersResponse is the response struct for ListOrders.
type ListOrdersResponse struct {
	Orders []Order `json:"orders"`
}

// ListOrders makes a call to GET /api/1/listorders.
//
// Returns a list of the most recently placed Orders.
// Users can specify an optional <code>state=PENDING</code> parameter to restrict the results to only open Orders.
// Users can also specify the market by using the optional currency pair parameter.
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

// ListOrdersV2Request is the request struct for ListOrdersV2.
type ListOrdersV2Request struct {
	// If true, will return closed orders instead of open orders.
	Closed bool `json:"closed" url:"closed"`

	// Filter to orders created before this timestamp (Unix milliseconds)
	CreatedBefore int64 `json:"created_before" url:"created_before"`

	// Limit to this many orders
	Limit int64 `json:"limit" url:"limit"`

	// Filter to only orders of this currency pair.
	Pair string `json:"pair" url:"pair"`
}

// ListOrdersV2Response is the response struct for ListOrdersV2.
type ListOrdersV2Response struct {
	Orders []OrderV2 `json:"orders"`
}

// ListOrdersV2 makes a call to GET /api/exchange/2/listorders.
//
// Returns a list of the most recently placed orders ordered from newest to
// oldest. This endpoint will list up to 100 most recent open orders by
// default.
//
// Permissions required: <Code>Perm_R_Orders</Code>
func (cl *Client) ListOrdersV2(ctx context.Context, req *ListOrdersV2Request) (*ListOrdersV2Response, error) {
	var res ListOrdersV2Response
	err := cl.do(ctx, "GET", "/api/exchange/2/listorders", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ListPendingTransactionsRequest is the request struct for ListPendingTransactions.
type ListPendingTransactionsRequest struct {
	// Account ID
	//
	// required: true
	Id int64 `json:"id" url:"id"`
}

// ListPendingTransactionsResponse is the response struct for ListPendingTransactions.
type ListPendingTransactionsResponse struct {
	Currency     string        `json:"currency"`
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	Pending      []Transaction `json:"pending"`
	Transactions []Transaction `json:"transactions"`
}

// ListPendingTransactions makes a call to GET /api/1/accounts/{id}/pending.
//
// Return a list of all transactions that have not completed for the Account.
//
// Pending transactions are not numbered, and may be reordered, deleted or updated at any time.
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

// ListTradesRequest is the request struct for ListTrades.
type ListTradesRequest struct {
	// Currency pair of the market to list the trades from
	//
	// required: true
	Pair string `json:"pair" url:"pair"`

	// Fetch trades executed after this time, specified as a Unix timestamp in
	// milliseconds. An error will be returned if this is before 24h ago. Use
	// this parameter to either restrict to a shorter window or to iterate over
	// the trades in case you need more than the 100 most recent trades.
	Since Time `json:"since" url:"since"`
}

// ListTradesResponse is the response struct for ListTrades.
type ListTradesResponse struct {
	Trades []PublicTrade `json:"trades"`
}

// ListTrades makes a call to GET /api/1/trades.
//
// Returns a list of recent trades for the specified currency pair. At most
// 100 trades are returned per call and never trades older than 24h. The
// trades are sorted from newest to oldest.
//
// Please see the <a href="#tag/currency ">Currency list</a> for the complete list of supported currency pairs.
func (cl *Client) ListTrades(ctx context.Context, req *ListTradesRequest) (*ListTradesResponse, error) {
	var res ListTradesResponse
	err := cl.do(ctx, "GET", "/api/1/trades", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ListTransactionsRequest is the request struct for ListTransactions.
type ListTransactionsRequest struct {
	// Account ID - the unique identifier for the specific Account.
	//
	// required: true
	Id int64 `json:"id" url:"id"`

	// Maximum of the row range to return (exclusive)
	//
	// required: true
	MaxRow int64 `json:"max_row" url:"max_row"`

	// Minimum of the row range to return (inclusive)
	//
	// required: true
	MinRow int64 `json:"min_row" url:"min_row"`
}

// ListTransactionsResponse is the response struct for ListTransactions.
type ListTransactionsResponse struct {
	Id           string           `json:"id"`
	Transactions []StatementEntry `json:"transactions"`
}

// ListTransactions makes a call to GET /api/1/accounts/{id}/transactions.
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
	err := cl.do(ctx, "GET", "/api/1/accounts/{id}/transactions", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ListTransfersRequest is the request struct for ListTransfers.
type ListTransfersRequest struct {
	// Unique identifier of the account to list the transfers from.
	//
	// required: true
	AccountId int64 `json:"account_id" url:"account_id"`

	// Filter to transfers created before this timestamp (Unix milliseconds).
	// The default value (0) will return the latest transfers on the account.
	Before int64 `json:"before" url:"before"`

	// Limit to this many transfers.
	Limit int64 `json:"limit" url:"limit"`
}

// ListTransfersResponse is the response struct for ListTransfers.
type ListTransfersResponse struct {
	Transfers []Transfer `json:"transfers"`
}

// ListTransfers makes a call to GET /api/exchange/1/transfers.
//
// Returns a list of the most recent confirmed transfers ordered from newest to
// oldest.
// This includes bank transfers, card payments, or on-chain transactions that
// have been reflected on your account available balance.
//
// Note that the Transfer `amount` is always a positive value and you should
// use the `inbound` flag to determine the direction of the transfer.
//
// If you need to paginate the results you can set the `before` parameter to
// the last returned transfer `created_at` field value and repeat the request
// until you have all the transfers you need.
// This endpoint will list up to 100 transfers at a time by default.
//
// This endpoint is in BETA, behaviour and specification may change without
// any previous notice.
//
// Permissions required: <Code>Perm_R_Transfers</Code>
func (cl *Client) ListTransfers(ctx context.Context, req *ListTransfersRequest) (*ListTransfersResponse, error) {
	var res ListTransfersResponse
	err := cl.do(ctx, "GET", "/api/exchange/1/transfers", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ListUserTradesRequest is the request struct for ListUserTrades.
type ListUserTradesRequest struct {
	// Filter to trades of this currency pair.
	//
	// required: true
	Pair string `json:"pair" url:"pair"`

	// Filter to trades from (including) this sequence number.
	// Default behaviour is not to include this filter.
	AfterSeq int64 `json:"after_seq" url:"after_seq"`

	// Filter to trades before this timestamp (Unix milliseconds).
	Before Time `json:"before" url:"before"`

	// Filter to trades before (excluding) this sequence number.
	// Default behaviour is not to include this filter.
	BeforeSeq int64 `json:"before_seq" url:"before_seq"`

	// Limit to this number of trades (default 100).
	Limit int64 `json:"limit" url:"limit"`

	// Filter to trades on or after this timestamp (Unix milliseconds).
	Since Time `json:"since" url:"since"`

	// If set to true, sorts trades in descending order, otherwise ascending
	// order will be assumed.
	SortDesc bool `json:"sort_desc" url:"sort_desc"`
}

// ListUserTradesResponse is the response struct for ListUserTrades.
type ListUserTradesResponse struct {
	Trades []TradeV2 `json:"trades"`
}

// ListUserTrades makes a call to GET /api/1/listtrades.
//
// Returns a list of the recent Trades for a given currency pair for this user, sorted by oldest first.
// If <code>before</code> is specified, then Trades are returned sorted by most-recent first.
//
// <code>type</code> in the response indicates the type of Order that was placed to participate in the trade.
// Possible types: <code>BID</code>, <code>ASK</code>.
//
// If <code>is_buy</code> in the response is true, then the Order which completed the trade (market taker) was a Bid Order.
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

// ListWithdrawalsRequest is the request struct for ListWithdrawals.
type ListWithdrawalsRequest struct {
	// Filter to withdrawals requested on or before the withdrawal with this ID.
	// Can be used for pagination.
	BeforeId int64 `json:"before_id" url:"before_id"`

	// Limit to this many withdrawals
	Limit int64 `json:"limit" url:"limit"`
}

// ListWithdrawalsResponse is the response struct for ListWithdrawals.
type ListWithdrawalsResponse struct {
	Withdrawals []Withdrawal `json:"withdrawals"`
}

// ListWithdrawals makes a call to GET /api/1/withdrawals.
//
// Returns a list of withdrawal requests.
//
// Permissions required: <code>Perm_R_Withdrawals</code>
func (cl *Client) ListWithdrawals(ctx context.Context, req *ListWithdrawalsRequest) (*ListWithdrawalsResponse, error) {
	var res ListWithdrawalsResponse
	err := cl.do(ctx, "GET", "/api/1/withdrawals", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// MarketsRequest is the request struct for Markets.
type MarketsRequest struct {
	// List of market pairs to return. Requesting only the required pairs will improve response times.
	Pair []string `json:"pair" url:"pair"`
}

// MarketsResponse is the response struct for Markets.
type MarketsResponse struct {
	Markets []MarketInfo `json:"markets"`
}

// Markets makes a call to GET /api/exchange/1/markets.
//
// List all supported markets parameter information like price scale, min and
// max order volumes and market ID.
func (cl *Client) Markets(ctx context.Context, req *MarketsRequest) (*MarketsResponse, error) {
	var res MarketsResponse
	err := cl.do(ctx, "GET", "/api/exchange/1/markets", req, &res, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// MoveRequest is the request struct for Move.
type MoveRequest struct {
	// Amount to transfer. Must be positive.
	//
	// required: true
	Amount decimal.Decimal `json:"amount" url:"amount"`

	// The account to credit the funds to.
	//
	// required: true
	CreditAccountId int64 `json:"credit_account_id" url:"credit_account_id"`

	// The account to debit the funds from.
	//
	// required: true
	DebitAccountId int64 `json:"debit_account_id" url:"debit_account_id"`

	// Client move ID.
	// May only contain alphanumeric (0-9, a-z, or A-Z) and special characters (_ ; , . -). Maximum length: 255.
	// It will be available in read endpoints, so you can use it to avoid duplicate moves between the same accounts.
	// Values must be unique across all your successful calls of this endpoint; trying to create a move request
	// with the same `client_move_id` as one of your past move requests will result in a HTTP 409 Conflict response.
	ClientMoveId string `json:"client_move_id" url:"client_move_id"`
}

// MoveResponse is the response struct for Move.
type MoveResponse struct {
	// Move unique identifier
	Id string `json:"id"`

	// The current state of the move.
	//
	// Status meaning:<br>
	// <code>CREATED</code> The move is awaiting execution.<br>
	// <code>MOVING</code> The funds have been reserved and the move is being executed.<br>
	// <code>SUCCESSFUL</code> The move has completed successfully and should be reflected in both accounts available
	// balance.<br>
	// <code>FAILED</code> The move has failed. There could be many reasons for this but the most likely is that the
	// debit account doesn't have enough available funds to move.<br>
	Status Status `json:"status"`
}

// Move makes a call to POST /api/exchange/1/move.
//
// Move funds between two of your accounts with the same currency
// The funds may not be moved by the time the request returns. The GET method
// can be used to poll for the move's status.
//
// Note: moves will show as transactions, but not as transfers.
//
// This endpoint is in BETA, behaviour and specification may change without
// any previous notice.
//
// Permissions required: <code>MP_None_Write</code>
func (cl *Client) Move(ctx context.Context, req *MoveRequest) (*MoveResponse, error) {
	var res MoveResponse
	err := cl.do(ctx, "POST", "/api/exchange/1/move", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// PostLimitOrderRequest is the request struct for PostLimitOrder.
type PostLimitOrderRequest struct {
	// The currency pair to trade.
	//
	// required: true
	Pair string `json:"pair" url:"pair"`

	// Limit price as a decimal string in units of ZAR/BTC.
	//
	// required: true
	Price decimal.Decimal `json:"price" url:"price"`

	// <code>BID</code> for a bid (buy) limit order<br>
	// <code>ASK</code> for an ask (sell) limit order
	//
	// required: true
	Type OrderType `json:"type" url:"type"`

	// Amount of cryptocurrency to buy or sell as a decimal string in units of the currency.
	//
	// required: true
	Volume decimal.Decimal `json:"volume" url:"volume"`

	// The base currency Account to use in the trade.
	BaseAccountId int64 `json:"base_account_id" url:"base_account_id"`

	// Client order ID.
	// May only contain alphanumeric (0-9, a-z, or A-Z) and special characters (_ ; , . -). Maximum length: 255.
	// It will be available in read endpoints, so you can use it to reconcile Luno with your internal system.
	// Values must be unique across all your successful order creation endpoint calls; trying to create an order
	// with the same `client_order_id` as one of your past orders will result in a HTTP 409 Conflict response.
	ClientOrderId string `json:"client_order_id" url:"client_order_id"`

	// The counter currency Account to use in the trade.
	CounterAccountId int64 `json:"counter_account_id" url:"counter_account_id"`

	// Post-only Orders will be cancelled if they would otherwise have traded
	// immediately.
	// For example, if there's a bid at ZAR 100,000 and you place a post-only ask at ZAR 100,000,
	// your order will be cancelled instead of trading.
	// If the best bid is ZAR 100,000 and you place a post-only ask at ZAR 101,000,
	// your order won't trade but will go into the order book.
	PostOnly bool `json:"post_only" url:"post_only"`

	// Side of the trigger price to activate the order. This should be set if `stop_price` is also
	// set.
	//
	// `RELATIVE_LAST_TRADE` will automatically infer the direction based on the last
	// trade price and the stop price. If last trade price is less than stop price then stop
	// direction is ABOVE otherwise is BELOW.
	StopDirection StopDirection `json:"stop_direction" url:"stop_direction"`

	// Trigger trade price to activate this order as a decimal string. If this
	// is set then this is treated as a Stop Limit Order and `stop_direction`
	// is expected to be set too.
	StopPrice decimal.Decimal `json:"stop_price" url:"stop_price"`

	// <code>GTC</code> Good 'Til Cancelled. The order remains open until it is filled or cancelled by the user.</br>
	// <code>IOC</code> Immediate Or Cancel. The part of the order that cannot be filled immediately will be cancelled. Cannot be post-only.</br>
	// <code>FOK</code> Fill Or Kill. If the order cannot be filled immediately and completely it will be cancelled before any trade. Cannot be post-only.
	TimeInForce TimeInForce `json:"time_in_force" url:"time_in_force"`

	// Unix timestamp in milliseconds of when the request was created and sent.
	Timestamp int64 `json:"timestamp" url:"timestamp"`

	// Specifies the number of milliseconds after timestamp the request is valid for.
	// If `timestamp` is not specified, `ttl` will not be used.
	Ttl int64 `json:"ttl" url:"ttl"`
}

// PostLimitOrderResponse is the response struct for PostLimitOrder.
type PostLimitOrderResponse struct {
	// Unique order identifier
	OrderId string `json:"order_id"`
}

// PostLimitOrder makes a call to POST /api/1/postorder.
//
// <b>Warning!</b> Orders cannot be reversed once they have executed.
// Please ensure your program has been thoroughly tested before submitting Orders.
//
// If no <code>base_account_id</code> or <code>counter_account_id</code> are specified,
// your default base currency or counter currency account will be used.
// You can find your Account IDs by calling the <a href="#operation/getBalances">Balances</a> API.
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

// PostMarketOrderRequest is the request struct for PostMarketOrder.
type PostMarketOrderRequest struct {
	// The currency pair to trade.
	//
	// required: true
	Pair string `json:"pair" url:"pair"`

	// <code>BUY</code> to buy an asset<br>
	// <code>SELL</code> to sell an asset
	//
	// required: true
	Type OrderType `json:"type" url:"type"`

	// The base currency account to use in the trade.
	BaseAccountId int64 `json:"base_account_id" url:"base_account_id"`

	// For a <code>SELL</code> order: amount of the base currency to use (e.g. how much BTC to sell for EUR in the BTC/EUR market)
	BaseVolume decimal.Decimal `json:"base_volume" url:"base_volume"`

	// Client order ID.
	// May only contain alphanumeric (0-9, a-z, or A-Z) and special characters (_ ; , . -). Maximum length: 255.
	// It will be available in read endpoints, so you can use it to reconcile Luno with your internal system.
	// Values must be unique across all your successful order creation endpoint calls; trying to create an order
	// with the same `client_order_id` as one of your past orders will result in a HTTP 409 Conflict response.
	ClientOrderId string `json:"client_order_id" url:"client_order_id"`

	// The counter currency account to use in the trade.
	CounterAccountId int64 `json:"counter_account_id" url:"counter_account_id"`

	// For a <code>BUY</code> order: amount of the counter currency to use (e.g. how much EUR to use to buy BTC in the BTC/EUR market)
	CounterVolume decimal.Decimal `json:"counter_volume" url:"counter_volume"`

	// Unix timestamp in milliseconds of when the request was created and sent.
	Timestamp int64 `json:"timestamp" url:"timestamp"`

	// Specifies the number of milliseconds after timestamp the request is valid for.
	// If `timestamp` is not specified, `ttl` will not be used.
	Ttl int64 `json:"ttl" url:"ttl"`
}

// PostMarketOrderResponse is the response struct for PostMarketOrder.
type PostMarketOrderResponse struct {
	// Unique order identifier
	OrderId string `json:"order_id"`
}

// PostMarketOrder makes a call to POST /api/1/marketorder.
//
// A Market Order executes immediately, and either buys as much of the asset that can be bought for a set amount of fiat currency, or sells a set amount of the asset for as much as possible.
//
// <b>Warning!</b> Orders cannot be reversed once they have executed.
// Please ensure your program has been thoroughly tested before submitting Orders.
//
// If no <code>base_account_id</code> or <code>counter_account_id</code> are specified, the default base currency or counter currency account will be used.
// Users can find their account IDs by calling the <a href="#operation/getBalances">Balances</a> request.
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

// SendRequest is the request struct for Send.
type SendRequest struct {
	// Destination address or email address.
	//
	// <b>Note</b>:
	// <ul>
	// <li>Ethereum addresses must be
	// <a href="https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md" target="_blank" rel="nofollow">checksummed</a>.</li>
	// <li>Ethereum sends to email addresses are not supported.</li>
	// </ul>
	//
	// required: true
	Address string `json:"address" url:"address"`

	// Amount to send as a decimal string.
	//
	// required: true
	Amount decimal.Decimal `json:"amount" url:"amount"`

	// Currency to send.
	//
	// required: true
	Currency string `json:"currency" url:"currency"`

	// User description for the transaction to record on the account statement.
	Description string `json:"description" url:"description"`

	// Optional XRP destination tag. Note that HasDestinationTag must be true if this value is provided.
	DestinationTag int64 `json:"destination_tag" url:"destination_tag"`

	// Optional unique ID to associate with this withdrawal.
	// Useful to prevent duplicate sends in case of failure.
	// This supports all alphanumeric characters, as well as "-" and "_".
	ExternalId string `json:"external_id" url:"external_id"`

	// Optional boolean flag indicating that a XRP destination tag is provided (even if zero).
	HasDestinationTag bool `json:"has_destination_tag" url:"has_destination_tag"`

	// Message to send to the recipient.
	// This is only relevant when sending to an email address.
	Message string `json:"message" url:"message"`
}

// SendResponse is the response struct for Send.
type SendResponse struct {
	Success      bool   `json:"success"`
	WithdrawalId string `json:"withdrawal_id"`
}

// Send makes a call to POST /api/1/send.
//
// Send assets from an Account. Please note that the asset type sent must match the receive address of the same cryptocurrency of the same type - Bitcoin to Bitcoin, Ethereum to Ethereum, etc.
//
// Sends can be to a cryptocurrency receive address, or the email address of another Luno platform user.
//
// <b>Note:</b> This is currently unavailable to users who are verified in countries with money travel rules.
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

// SendFeeRequest is the request struct for SendFee.
type SendFeeRequest struct {
	// Destination address or email address.
	//
	// <b>Note</b>:
	// <ul>
	// <li>Ethereum addresses must be
	// <a href="https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md" target="_blank" rel="nofollow">checksummed</a>.</li>
	// <li>Ethereum sends to email addresses are not supported.</li>
	// </ul>
	//
	// required: true
	Address string `json:"address" url:"address"`

	// Amount to send as a decimal string.
	//
	// required: true
	Amount decimal.Decimal `json:"amount" url:"amount"`

	// Currency to send.
	//
	// required: true
	Currency string `json:"currency" url:"currency"`
}

// SendFeeResponse is the response struct for SendFee.
type SendFeeResponse struct {
	Currency string          `json:"currency"`
	Fee      decimal.Decimal `json:"fee"`
}

// SendFee makes a call to GET /api/1/send_fee.
//
// Calculate fees involved with a crypto send request.
//
// Send address can be to a cryptocurrency receive address, or the email address of another Luno platform user.
//
// Permissions required: <code>MP_None</code>
func (cl *Client) SendFee(ctx context.Context, req *SendFeeRequest) (*SendFeeResponse, error) {
	var res SendFeeResponse
	err := cl.do(ctx, "GET", "/api/1/send_fee", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// StopOrderRequest is the request struct for StopOrder.
type StopOrderRequest struct {
	// The Order identifier as a string.
	//
	// required: true
	OrderId string `json:"order_id" url:"order_id"`
}

// StopOrderResponse is the response struct for StopOrder.
type StopOrderResponse struct {
	Success bool `json:"success"`
}

// StopOrder makes a call to POST /api/1/stoporder.
//
// Request to cancel an Order.
//
// <b>Note!</b>: Once an Order has been completed, it can not be reversed.
// The return value from this request will indicate if the Stop request was successful or not.
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

// UpdateAccountNameRequest is the request struct for UpdateAccountName.
type UpdateAccountNameRequest struct {
	// Account ID - the unique identifier for the specific Account.
	//
	// required: true
	Id int64 `json:"id" url:"id"`

	// The label to use for this account
	//
	// required: true
	Name string `json:"name" url:"name"`
}

// UpdateAccountNameResponse is the response struct for UpdateAccountName.
type UpdateAccountNameResponse struct {
	Success bool `json:"success"`
}

// UpdateAccountName makes a call to PUT /api/1/accounts/{id}/name.
//
// Update the name of an account with a given ID.
//
// Permissions required: <code>Perm_W_Addresses</code>
func (cl *Client) UpdateAccountName(ctx context.Context, req *UpdateAccountNameRequest) (*UpdateAccountNameResponse, error) {
	var res UpdateAccountNameResponse
	err := cl.do(ctx, "PUT", "/api/1/accounts/{id}/name", req, &res, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// vi: ft=go
