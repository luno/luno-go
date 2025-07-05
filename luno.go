// Package luno is a wrapper for the Luno API.
package luno

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type Limiter interface {
	Wait(context.Context) error
}

// ClientInterface defines all Luno API operations.
// This interface allows for easy mocking of the Luno client in tests.
type ClientInterface interface {
	CancelWithdrawal(ctx context.Context, req *CancelWithdrawalRequest) (*CancelWithdrawalResponse, error)
	CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error)
	CreateBeneficiary(ctx context.Context, req *CreateBeneficiaryRequest) (*CreateBeneficiaryResponse, error)
	CreateFundingAddress(ctx context.Context, req *CreateFundingAddressRequest) (*CreateFundingAddressResponse, error)
	CreateWithdrawal(ctx context.Context, req *CreateWithdrawalRequest) (*CreateWithdrawalResponse, error)
	DeleteBeneficiary(ctx context.Context, req *DeleteBeneficiaryRequest) error
	GetBalances(ctx context.Context, req *GetBalancesRequest) (*GetBalancesResponse, error)
	GetCandles(ctx context.Context, req *GetCandlesRequest) (*GetCandlesResponse, error)
	GetFeeInfo(ctx context.Context, req *GetFeeInfoRequest) (*GetFeeInfoResponse, error)
	GetFundingAddress(ctx context.Context, req *GetFundingAddressRequest) (*GetFundingAddressResponse, error)
	GetMove(ctx context.Context, req *GetMoveRequest) (*GetMoveResponse, error)
	GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error)
	GetOrderBook(ctx context.Context, req *GetOrderBookRequest) (*GetOrderBookResponse, error)
	GetOrderBookFull(ctx context.Context, req *GetOrderBookFullRequest) (*GetOrderBookFullResponse, error)
	GetOrderV2(ctx context.Context, req *GetOrderV2Request) (*GetOrderV2Response, error)
	GetOrderV3(ctx context.Context, req *GetOrderV3Request) (*GetOrderV3Response, error)
	GetTicker(ctx context.Context, req *GetTickerRequest) (*GetTickerResponse, error)
	GetTickers(ctx context.Context, req *GetTickersRequest) (*GetTickersResponse, error)
	GetWithdrawal(ctx context.Context, req *GetWithdrawalRequest) (*GetWithdrawalResponse, error)
	ListBeneficiaries(ctx context.Context, req *ListBeneficiariesRequest) (*ListBeneficiariesResponse, error)
	ListMoves(ctx context.Context, req *ListMovesRequest) (*ListMovesResponse, error)
	ListOrders(ctx context.Context, req *ListOrdersRequest) (*ListOrdersResponse, error)
	ListOrdersV2(ctx context.Context, req *ListOrdersV2Request) (*ListOrdersV2Response, error)
	ListPendingTransactions(ctx context.Context, req *ListPendingTransactionsRequest) (*ListPendingTransactionsResponse, error)
	ListTrades(ctx context.Context, req *ListTradesRequest) (*ListTradesResponse, error)
	ListTransactions(ctx context.Context, req *ListTransactionsRequest) (*ListTransactionsResponse, error)
	ListTransfers(ctx context.Context, req *ListTransfersRequest) (*ListTransfersResponse, error)
	ListUserTrades(ctx context.Context, req *ListUserTradesRequest) (*ListUserTradesResponse, error)
	ListWithdrawals(ctx context.Context, req *ListWithdrawalsRequest) (*ListWithdrawalsResponse, error)
	Markets(ctx context.Context, req *MarketsRequest) (*MarketsResponse, error)
	Move(ctx context.Context, req *MoveRequest) (*MoveResponse, error)
	PostLimitOrder(ctx context.Context, req *PostLimitOrderRequest) (*PostLimitOrderResponse, error)
	PostMarketOrder(ctx context.Context, req *PostMarketOrderRequest) (*PostMarketOrderResponse, error)
	Send(ctx context.Context, req *SendRequest) (*SendResponse, error)
	SendFee(ctx context.Context, req *SendFeeRequest) (*SendFeeResponse, error)
	StopOrder(ctx context.Context, req *StopOrderRequest) (*StopOrderResponse, error)
	UpdateAccountName(ctx context.Context, req *UpdateAccountNameRequest) (*UpdateAccountNameResponse, error)
	Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error)
}

// Client is a Luno API client.
type Client struct {
	httpClient   *http.Client
	rateLimiter  Limiter
	baseURL      string
	apiKeyID     string
	apiKeySecret string
	debug        bool
}

const (
	defaultBaseURL = "https://api.luno.com"
	defaultTimeout = 10 * time.Second
	// Rate limiting parameters:
	// Refer to https://www.luno.com/en/developers/api#tag/Rate-Limiting
	// The default configuration allows for a burst of 1 request every 200ms
	// which aggregates to 300 requests per minute.
	//
	// defaultRate specifies the rate at which requests are allowed.
	defaultRate = time.Minute / 300
	// defaultBurst specifies the maximum number of requests allowed in a burst.
	defaultBurst = 1
)

// NewClient creates a new Luno API client with the default base URL.
func NewClient() *Client {
	return &Client{
		httpClient:  &http.Client{Timeout: defaultTimeout},
		baseURL:     defaultBaseURL,
		rateLimiter: rate.NewLimiter(rate.Every(defaultRate), defaultBurst),
	}
}

// SetAuth provides the client with an API key and secret.
func (cl *Client) SetAuth(apiKeyID, apiKeySecret string) error {
	if apiKeyID == "" || apiKeySecret == "" {
		return errors.New("luno: no credentials provided")
	}
	cl.apiKeyID = apiKeyID
	cl.apiKeySecret = apiKeySecret
	return nil
}

// SetHTTPClient sets the HTTP client that will be used for API calls.
func (cl *Client) SetHTTPClient(httpClient *http.Client) {
	cl.httpClient = httpClient
}

// SetRateLimiter sets the rate limiter that will be used to throttle calls
// made through the client.
func (cl *Client) SetRateLimiter(rateLimiter Limiter) {
	cl.rateLimiter = rateLimiter
}

// SetTimeout sets the timeout for requests made by this client. Note: if you
// set a timeout and then call .SetHTTPClient(), the timeout in the new HTTP
// client will be used.
func (cl *Client) SetTimeout(timeout time.Duration) {
	cl.httpClient.Timeout = timeout
}

// SetBaseURL overrides the default base URL. For internal use.
func (cl *Client) SetBaseURL(baseURL string) {
	cl.baseURL = strings.TrimRight(baseURL, "/")
}

// SetDebug enables or disables debug mode. In debug mode, HTTP requests and
// responses will be logged.
func (cl *Client) SetDebug(debug bool) {
	cl.debug = debug
}

func (cl *Client) do(ctx context.Context, method, path string,
	req, res interface{}, auth bool,
) error {
	err := cl.rateLimiter.Wait(ctx)
	if err != nil {
		return err
	}

	url := cl.baseURL + "/" + strings.TrimLeft(path, "/")

	if cl.debug {
		log.Printf("luno: Call: %s %s", method, path)
		log.Printf("luno: Request: %#v", req)
	}

	var contentType string
	var body io.Reader
	if req != nil {
		values := makeURLValues(req)
		if strings.Contains(path, "{id}") {
			url = strings.ReplaceAll(url, "{id}", values.Get("id"))
			values.Del("id")
		}
		if method == http.MethodGet {
			url = url + "?" + values.Encode()
		} else {
			body = strings.NewReader(values.Encode())
			contentType = "application/x-www-form-urlencoded"
		}
	}

	httpReq, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	httpReq = httpReq.WithContext(ctx)
	httpReq.Header.Set("User-Agent", makeUserAgent())
	if contentType != "" {
		httpReq.Header.Set("Content-Type", contentType)
	}

	if auth {
		httpReq.SetBasicAuth(cl.apiKeyID, cl.apiKeySecret)
	}

	if method != http.MethodGet {
		httpReq.Header.Set("content-type", "application/x-www-form-urlencoded")
	}

	httpRes, err := cl.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	body = httpRes.Body
	if cl.debug {
		b, err := io.ReadAll(body)
		if err != nil {
			log.Printf("luno: Error reading response body: %v", err)
		} else {
			log.Printf("Response: %s", string(b))
		}
		body = bytes.NewReader(b)
	}

	if httpRes.StatusCode == http.StatusTooManyRequests {
		return errors.New("luno: too many requests")
	}

	if httpRes.StatusCode != http.StatusOK {
		var e Error
		err := json.NewDecoder(body).Decode(&e)
		if err != nil {
			return fmt.Errorf("luno: error decoding response (%d %s)",
				httpRes.StatusCode, http.StatusText(httpRes.StatusCode))
		}
		return e
	}

	return json.NewDecoder(body).Decode(res)
}

func makeUserAgent() string {
	return fmt.Sprintf("LunoGoSDK/%s %s %s %s",
		Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// Ensure Client implements ClientInterface at compile time.
var _ ClientInterface = (*Client)(nil)
