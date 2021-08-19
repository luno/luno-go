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
	Amount     decimal.Decimal `json:"amount"`
	CreatedAt  Time            `json:"created_at"`
	Currency   string          `json:"currency"`
	ExternalId string          `json:"external_id"`
	Fee        decimal.Decimal `json:"fee"`
	Id         string          `json:"id"`
	Status     string          `json:"status"`
	Type       string          `json:"type"`
}

// CancelWithdrawal makes a call to DELETE /api/1/withdrawals/{id}.
//
// Cancel a withdrawal request.
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

// CreateQuoteRequest is the request struct for CreateQuote.
type CreateQuoteRequest struct {
	// Amount to buy or sell in the pair base currency.
	//
	// required: true
	BaseAmount decimal.Decimal `json:"base_amount" url:"base_amount"`

	// Currency pair to trade. The pair can also be flipped if you want to buy
	// or sell the counter currency (e.g. ZARXBT).
	//
	// required: true
	Pair string `json:"pair" url:"pair"`

	// <code>BUY</code> or <code>SELL</code>.
	//
	// required: true
	Type string `json:"type" url:"type"`

	// Optional Account for the pair's base currency.
	BaseAccountId int64 `json:"base_account_id" url:"base_account_id"`

	// Optional Account for the pair's counter currency.
	CounterAccountId int64 `json:"counter_account_id" url:"counter_account_id"`
}

// CreateQuoteResponse is the response struct for CreateQuote.
type CreateQuoteResponse struct {
	BaseAmount    decimal.Decimal `json:"base_amount"`
	CounterAmount decimal.Decimal `json:"counter_amount"`
	CreatedAt     Time            `json:"created_at"`
	Discarded     bool            `json:"discarded"`
	Exercised     bool            `json:"exercised"`
	ExpiresAt     Time            `json:"expires_at"`
	Id            string          `json:"id"`
	Pair          string          `json:"pair"`
	Type          string          `json:"type"`
}

// CreateQuote makes a call to POST /api/1/quotes.
//
// Creates a new quote to buy or sell a particular amount of a base currency for a counter currency.
//
// Users can specify either the exact amount to pay or the exact amount to receive.
//
// For example, to buy exactly 0.1 Bitcoin using ZAR, you would create a quote to BUY 0.1 XBTZAR.
// The returned quote includes the appropriate ZAR amount.
// To buy Bitcoin using exactly ZAR 100, create a quote to SELL 100 ZARXBT.
// The returned quote specifies the Bitcoin as the counter amount returned.
//
// An error is returned if the Account is not verified for the currency pair,
// or if the Account would have insufficient balance to ever exercise the quote.
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

	// For internal use.
	Reference string `json:"reference" url:"reference"`
}

// CreateWithdrawalResponse is the response struct for CreateWithdrawal.
type CreateWithdrawalResponse struct {
	Amount     decimal.Decimal `json:"amount"`
	CreatedAt  Time            `json:"created_at"`
	Currency   string          `json:"currency"`
	ExternalId string          `json:"external_id"`
	Fee        decimal.Decimal `json:"fee"`
	Id         string          `json:"id"`
	Status     string          `json:"status"`
	Type       string          `json:"type"`
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

// DiscardQuoteRequest is the request struct for DiscardQuote.
type DiscardQuoteRequest struct {
	// ID of the quote to discard.
	//
	// required: true
	Id int64 `json:"id" url:"id"`
}

// DiscardQuoteResponse is the response struct for DiscardQuote.
type DiscardQuoteResponse struct {
	BaseAmount    decimal.Decimal `json:"base_amount"`
	CounterAmount decimal.Decimal `json:"counter_amount"`
	CreatedAt     Time            `json:"created_at"`
	Discarded     bool            `json:"discarded"`
	Exercised     bool            `json:"exercised"`
	ExpiresAt     Time            `json:"expires_at"`
	Id            string          `json:"id"`
	Pair          string          `json:"pair"`
	Type          string          `json:"type"`
}

// DiscardQuote makes a call to DELETE /api/1/quotes/{id}.
//
// Discard a Quote.
// Once a Quote has been discarded, it cannot be exercised even if it has not expired.
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

// ExerciseQuoteRequest is the request struct for ExerciseQuote.
type ExerciseQuoteRequest struct {
	// ID of the quote to exercise.
	//
	// required: true
	Id int64 `json:"id" url:"id"`
}

// ExerciseQuoteResponse is the response struct for ExerciseQuote.
type ExerciseQuoteResponse struct {
	BaseAmount    decimal.Decimal `json:"base_amount"`
	CounterAmount decimal.Decimal `json:"counter_amount"`
	CreatedAt     Time            `json:"created_at"`
	Discarded     bool            `json:"discarded"`
	Exercised     bool            `json:"exercised"`
	ExpiresAt     Time            `json:"expires_at"`
	Id            string          `json:"id"`
	Pair          string          `json:"pair"`
	Type          string          `json:"type"`
}

// ExerciseQuote makes a call to PUT /api/1/quotes/{id}.
//
// Exercise a quote to perform the Trade.
// If there is sufficient balance available in the Account,
// it will be debited and the counter amount credited.
//
// An error is returned if the quote has expired or if the Account has insufficient available balance.
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

// GetBalancesRequest is the request struct for GetBalances.
type GetBalancesRequest struct {
	// Only return balances for wallets with these currencies (if not provided,
	// all balances will be returned)
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
//
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

// GetOrderRequest is the request struct for GetOrder.
type GetOrderRequest struct {
	// The order ID.
	//
	// required: true
	Id string `json:"id" url:"id"`
}

// GetOrderResponse is the response struct for GetOrder.
type GetOrderResponse struct {
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
	Asks      []OrderBookEntry `json:"asks"`
	Bids      []OrderBookEntry `json:"bids"`
	Timestamp int64            `json:"timestamp"`
}

// GetOrderBook makes a call to GET /api/1/orderbook_top.
//
// Returns a list of the top 100 <em>bids</em> and <em>asks</em> for the currency pair specified in the Order Book.
//
// Ask Orders are sorted by price ascending.
//
// Bid Orders are sorted by price descending.
//
// Orders of the same price are aggregated.
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
	Asks      []OrderBookEntry `json:"asks"`
	Bids      []OrderBookEntry `json:"bids"`
	Timestamp int64            `json:"timestamp"`
}

// GetOrderBookFull makes a call to GET /api/1/orderbook.
//
// This request returns a list of all <em>bids</em> and <em>asks</em> for the currency pair specified in the Order Book.
//
// Ask orders are sorted by price ascending.
//
// Bid orders are sorted by price descending.
//
// Multiple orders at the same price are not aggregated.
//
// <b>Warning:</b> This may return a large amount of data.
// Users are recommended to use the <a href="#operation/getOrderBook">top 100 bids and asks</a>
// or the <a href="#tag/streaming-API-(beta)">Streaming API</a>.
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

// GetOrderV2 makes a call to GET /api/exchange/2/orders/{id}.
//
// Get the details for an order.<br>
// This endpoint is in BETA, behaviour and specification may change without
// any previous notice.
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

// GetQuoteRequest is the request struct for GetQuote.
type GetQuoteRequest struct {
	// ID of the quote to retrieve.
	//
	// required: true
	Id int64 `json:"id" url:"id"`
}

// GetQuoteResponse is the response struct for GetQuote.
type GetQuoteResponse struct {
	BaseAmount    decimal.Decimal `json:"base_amount"`
	CounterAmount decimal.Decimal `json:"counter_amount"`
	CreatedAt     Time            `json:"created_at"`
	Discarded     bool            `json:"discarded"`
	Exercised     bool            `json:"exercised"`
	ExpiresAt     Time            `json:"expires_at"`
	Id            string          `json:"id"`
	Pair          string          `json:"pair"`
	Type          string          `json:"type"`
}

// GetQuote makes a call to GET /api/1/quotes/{id}.
//
// Get the latest status of a quote by its id.
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

// GetTickerRequest is the request struct for GetTicker.
type GetTickerRequest struct {
	// Currency pair
	//
	// required: true
	Pair string `json:"pair" url:"pair"`
}

// GetTickerResponse is the response struct for GetTicker.
type GetTickerResponse struct {
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
	Amount     decimal.Decimal `json:"amount"`
	CreatedAt  Time            `json:"created_at"`
	Currency   string          `json:"currency"`
	ExternalId string          `json:"external_id"`
	Fee        decimal.Decimal `json:"fee"`
	Id         string          `json:"id"`
	Status     string          `json:"status"`
	Type       string          `json:"type"`
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
// The list is truncated after 100 items.
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
// Returns a list of the most recently placed orders. This endpoint will list
// up to 100 open orders by default.<br>
// This endpoint is in BETA, behaviour and specification may change without
// any previous notice.
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
	// Currency pair
	//
	// required: true
	Pair string `json:"pair" url:"pair"`

	// Fetch trades executed after this time, specified as a Unix timestamp in
	// milliseconds.
	Since Time `json:"since" url:"since"`
}

// ListTradesResponse is the response struct for ListTrades.
type ListTradesResponse struct {
	Trades []Trade `json:"trades"`
}

// ListTrades makes a call to GET /api/1/trades.
//
// Returns a list of the most recent Trades for the specified currency pair in the last 24 hours.
// At most 100 results are returned per call.
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
	Currency     string        `json:"currency"`
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	Pending      []Transaction `json:"pending"`
	Transactions []Transaction `json:"transactions"`
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

// ListUserTradesRequest is the request struct for ListUserTrades.
type ListUserTradesRequest struct {
	// Filter to trades of this currency pair.
	//
	// required: true
	Pair string `json:"pair" url:"pair"`

	// Filter to trades from (including) this sequence number.
	// Default behaviour is not to include this filter.
	AfterSeq int64 `json:"after_seq" url:"after_seq"`

	// Filter to trades before this timestamp.
	Before Time `json:"before" url:"before"`

	// Filter to trades before (excluding) this sequence number.
	// Default behaviour is not to include this filter.
	BeforeSeq int64 `json:"before_seq" url:"before_seq"`

	// Limit to this number of trades (default 100).
	Limit int64 `json:"limit" url:"limit"`

	// Filter to trades on or after this timestamp.
	Since Time `json:"since" url:"since"`

	// If set to true, sorts trades in descending order, otherwise ascending
	// order will be assumed.
	SortDesc bool `json:"sort_desc" url:"sort_desc"`
}

// ListUserTradesResponse is the response struct for ListUserTrades.
type ListUserTradesResponse struct {
	Trades []Trade `json:"trades"`
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
}

// MarketsResponse is the response struct for Markets.
type MarketsResponse struct {
	Markets []MarketInfo `json:"markets"`
}

// Markets makes a call to GET /api/exchange/1/markets.
//
// Get all supported markets parameter information like price scale, min and
// max volumes and market ID.
func (cl *Client) Markets(ctx context.Context, req *MarketsRequest) (*MarketsResponse, error) {
	var res MarketsResponse
	err := cl.do(ctx, "GET", "/api/exchange/1/markets", req, &res, false)
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
}

// PostLimitOrderResponse is the response struct for PostLimitOrder.
type PostLimitOrderResponse struct {
	OrderId string `json:"order_id"`
}

// PostLimitOrder makes a call to POST /api/1/postorder.
//
// Create a new Trade Order.
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

	// The counter currency account to use in the trade.
	CounterAccountId int64 `json:"counter_account_id" url:"counter_account_id"`

	// For a <code>BUY</code> order: amount of the counter currency to use (e.g. how much EUR to use to buy BTC in the BTC/EUR market)
	CounterVolume decimal.Decimal `json:"counter_volume" url:"counter_volume"`
}

// PostMarketOrderResponse is the response struct for PostMarketOrder.
type PostMarketOrderResponse struct {
	OrderId string `json:"order_id"`
}

// PostMarketOrder makes a call to POST /api/1/marketorder.
//
// Create a new Market Order.
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
// Permissions required: <code>Perm_W_Send</code>
func (cl *Client) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	var res SendResponse
	err := cl.do(ctx, "POST", "/api/1/send", req, &res, true)
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
// Request to stop an Order.
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
