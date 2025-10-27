package luno

import (
	"context"
	"testing"

	"github.com/luno/luno-go/decimal"
)

// MockClient is a simple mock implementation of ClientInterface for testing.
type MockClient struct {
	GetBalancesFunc func(ctx context.Context, req *GetBalancesRequest) (*GetBalancesResponse, error)
	GetTickerFunc   func(ctx context.Context, req *GetTickerRequest) (*GetTickerResponse, error)
}

func (m *MockClient) CancelWithdrawal(ctx context.Context, req *CancelWithdrawalRequest) (*CancelWithdrawalResponse, error) {
	panic("not implemented")
}

func (m *MockClient) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	panic("not implemented")
}

func (m *MockClient) CreateBeneficiary(ctx context.Context, req *CreateBeneficiaryRequest) (*CreateBeneficiaryResponse, error) {
	panic("not implemented")
}

func (m *MockClient) CreateFundingAddress(ctx context.Context, req *CreateFundingAddressRequest) (*CreateFundingAddressResponse, error) {
	panic("not implemented")
}

func (m *MockClient) CreateWithdrawal(ctx context.Context, req *CreateWithdrawalRequest) (*CreateWithdrawalResponse, error) {
	panic("not implemented")
}

func (m *MockClient) DeleteBeneficiary(ctx context.Context, req *DeleteBeneficiaryRequest) error {
	panic("not implemented")
}

func (m *MockClient) GetBalances(ctx context.Context, req *GetBalancesRequest) (*GetBalancesResponse, error) {
	if m.GetBalancesFunc != nil {
		return m.GetBalancesFunc(ctx, req)
	}
	panic("not implemented")
}

func (m *MockClient) GetCandles(ctx context.Context, req *GetCandlesRequest) (*GetCandlesResponse, error) {
	panic("not implemented")
}

func (m *MockClient) GetFeeInfo(ctx context.Context, req *GetFeeInfoRequest) (*GetFeeInfoResponse, error) {
	panic("not implemented")
}

func (m *MockClient) GetFundingAddress(ctx context.Context, req *GetFundingAddressRequest) (*GetFundingAddressResponse, error) {
	panic("not implemented")
}

func (m *MockClient) GetMove(ctx context.Context, req *GetMoveRequest) (*GetMoveResponse, error) {
	panic("not implemented")
}

func (m *MockClient) GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error) {
	panic("not implemented")
}

func (m *MockClient) GetOrderBook(ctx context.Context, req *GetOrderBookRequest) (*GetOrderBookResponse, error) {
	panic("not implemented")
}

func (m *MockClient) GetOrderBookFull(ctx context.Context, req *GetOrderBookFullRequest) (*GetOrderBookFullResponse, error) {
	panic("not implemented")
}

func (m *MockClient) GetOrderV2(ctx context.Context, req *GetOrderV2Request) (*GetOrderV2Response, error) {
	panic("not implemented")
}

func (m *MockClient) GetOrderV3(ctx context.Context, req *GetOrderV3Request) (*GetOrderV3Response, error) {
	panic("not implemented")
}

func (m *MockClient) GetTicker(ctx context.Context, req *GetTickerRequest) (*GetTickerResponse, error) {
	if m.GetTickerFunc != nil {
		return m.GetTickerFunc(ctx, req)
	}
	panic("not implemented")
}

func (m *MockClient) GetTickers(ctx context.Context, req *GetTickersRequest) (*GetTickersResponse, error) {
	panic("not implemented")
}

func (m *MockClient) GetWithdrawal(ctx context.Context, req *GetWithdrawalRequest) (*GetWithdrawalResponse, error) {
	panic("not implemented")
}

func (m *MockClient) ListBeneficiaries(ctx context.Context, req *ListBeneficiariesRequest) (*ListBeneficiariesResponse, error) {
	panic("not implemented")
}

func (m *MockClient) ListMoves(ctx context.Context, req *ListMovesRequest) (*ListMovesResponse, error) {
	panic("not implemented")
}

func (m *MockClient) ListOrders(ctx context.Context, req *ListOrdersRequest) (*ListOrdersResponse, error) {
	panic("not implemented")
}

func (m *MockClient) ListOrdersV2(ctx context.Context, req *ListOrdersV2Request) (*ListOrdersV2Response, error) {
	panic("not implemented")
}

func (m *MockClient) ListPendingTransactions(ctx context.Context, req *ListPendingTransactionsRequest) (*ListPendingTransactionsResponse, error) {
	panic("not implemented")
}

func (m *MockClient) ListTrades(ctx context.Context, req *ListTradesRequest) (*ListTradesResponse, error) {
	panic("not implemented")
}

func (m *MockClient) ListTransactions(ctx context.Context, req *ListTransactionsRequest) (*ListTransactionsResponse, error) {
	panic("not implemented")
}

func (m *MockClient) ListTransfers(ctx context.Context, req *ListTransfersRequest) (*ListTransfersResponse, error) {
	panic("not implemented")
}

func (m *MockClient) ListUserTrades(ctx context.Context, req *ListUserTradesRequest) (*ListUserTradesResponse, error) {
	panic("not implemented")
}

func (m *MockClient) ListWithdrawals(ctx context.Context, req *ListWithdrawalsRequest) (*ListWithdrawalsResponse, error) {
	panic("not implemented")
}

func (m *MockClient) Markets(ctx context.Context, req *MarketsRequest) (*MarketsResponse, error) {
	panic("not implemented")
}

func (m *MockClient) Move(ctx context.Context, req *MoveRequest) (*MoveResponse, error) {
	panic("not implemented")
}

func (m *MockClient) PostLimitOrder(ctx context.Context, req *PostLimitOrderRequest) (*PostLimitOrderResponse, error) {
	panic("not implemented")
}

func (m *MockClient) PostMarketOrder(ctx context.Context, req *PostMarketOrderRequest) (*PostMarketOrderResponse, error) {
	panic("not implemented")
}

func (m *MockClient) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	panic("not implemented")
}

func (m *MockClient) SendFee(ctx context.Context, req *SendFeeRequest) (*SendFeeResponse, error) {
	panic("not implemented")
}

func (m *MockClient) StopOrder(ctx context.Context, req *StopOrderRequest) (*StopOrderResponse, error) {
	panic("not implemented")
}

func (m *MockClient) UpdateAccountName(ctx context.Context, req *UpdateAccountNameRequest) (*UpdateAccountNameResponse, error) {
	panic("not implemented")
}

func (m *MockClient) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	panic("not implemented")
}

// Ensure MockClient implements ClientInterface at compile time.
var _ ClientInterface = (*MockClient)(nil)

// Example function that uses ClientInterface for testability.
func GetAccountBalance(client ClientInterface, ctx context.Context, asset string) (decimal.Decimal, error) {
	req := &GetBalancesRequest{Assets: []string{asset}}
	res, err := client.GetBalances(ctx, req)
	if err != nil {
		return decimal.Zero(), err
	}
	
	for _, balance := range res.Balance {
		if balance.Asset == asset {
			return balance.Balance, nil
		}
	}
	
	return decimal.Zero(), nil
}

// TestClientInterface demonstrates how to use the interface for mocking.
func TestClientInterface(t *testing.T) {
	// Test with mock client
	mockClient := &MockClient{
		GetBalancesFunc: func(ctx context.Context, req *GetBalancesRequest) (*GetBalancesResponse, error) {
			return &GetBalancesResponse{
				Balance: []AccountBalance{
					{
						Asset:   "XBT",
						Balance: decimal.NewFromFloat64(1.5, 8),
					},
				},
			}, nil
		},
	}

	balance, err := GetAccountBalance(mockClient, context.Background(), "XBT")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := decimal.NewFromFloat64(1.5, 8)
	if balance.String() != expected.String() {
		t.Errorf("Expected balance %s, got %s", expected.String(), balance.String())
	}
}

// TestRealClientImplementsInterface verifies that the real Client implements the interface.
func TestRealClientImplementsInterface(t *testing.T) {
	client := NewClient()
	
	// This will compile only if Client implements ClientInterface
	var _ ClientInterface = client
	
	// Verify it's not nil
	if client == nil {
		t.Error("Expected client to be non-nil")
	}
}