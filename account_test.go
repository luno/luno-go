package luno

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luno/luno-go/decimal"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountWithAccountType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/1/accounts", r.URL.Path)

		require.NoError(t, r.ParseForm())
		require.Equal(t, "XBT", r.Form.Get("currency"))
		require.Equal(t, "Trading ACC", r.Form.Get("name"))
		require.Equal(t, "SPOT", r.Form.Get("account_type"))

		_ = json.NewEncoder(w).Encode(CreateAccountResponse{
			Currency: "XBT",
			Id:       "12345",
			Name:     "Trading ACC",
		})
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)
	require.NoError(t, cl.SetAuth("key-id", "key-secret"))

	res, err := cl.CreateAccount(context.Background(), &CreateAccountRequest{
		Currency:    "XBT",
		Name:        "Trading ACC",
		AccountType: WalletAccountTypeSpot,
	})
	require.NoError(t, err)
	require.Equal(t, "12345", res.Id)
}

func TestGetBalancesWithAccountTypeFilter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/1/balance", r.URL.Path)

		require.NoError(t, r.ParseForm())
		require.Equal(t, []string{"SPOT", "SPOT_MARGIN"}, r.Form["account_type"])

		_ = json.NewEncoder(w).Encode(GetBalancesResponse{
			Balance: []AccountBalance{
				{
					AccountId:   "12345",
					AccountType: WalletAccountTypeSpot,
					Asset:       "XBT",
					Balance:     decimal.NewFromInt64(1),
				},
			},
		})
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)
	require.NoError(t, cl.SetAuth("key-id", "key-secret"))

	res, err := cl.GetBalances(context.Background(), &GetBalancesRequest{
		AccountType: []WalletAccountType{WalletAccountTypeSpot, WalletAccountTypeSpotMargin},
	})
	require.NoError(t, err)
	require.Len(t, res.Balance, 1)
	require.Equal(t, WalletAccountTypeSpot, res.Balance[0].AccountType)
}
