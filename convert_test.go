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

func TestConvert(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/exchange/1/convert", r.URL.Path)

		require.NoError(t, r.ParseForm())
		require.Equal(t, "100", r.Form.Get("amount"))
		require.Equal(t, "12345", r.Form.Get("source_account_id"))
		require.Equal(t, "12346", r.Form.Get("target_account_id"))
		require.Equal(t, "convert-abc123", r.Form.Get("idempotency_key"))

		_, _, ok := r.BasicAuth()
		require.True(t, ok, "Convert must send credentials")

		_ = json.NewEncoder(w).Encode(ConvertResponse{
			Id:                   "conv-1",
			SourceCurrency:       "ZARU",
			TargetCurrency:       "ZAR",
			ConvertedAmount:      decimal.NewFromInt64(100),
			SourceAccountBalance: decimal.NewFromInt64(900),
			TargetAccountBalance: decimal.NewFromInt64(1100),
		})
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)
	require.NoError(t, cl.SetAuth("key-id", "key-secret"))

	res, err := cl.Convert(context.Background(), &ConvertRequest{
		Amount:          decimal.NewFromInt64(100),
		SourceAccountId: 12345,
		TargetAccountId: 12346,
		IdempotencyKey:  "convert-abc123",
	})
	require.NoError(t, err)
	require.Equal(t, "conv-1", res.Id)
	require.Equal(t, "ZARU", res.SourceCurrency)
	require.Equal(t, "ZAR", res.TargetCurrency)
	require.Equal(t, "100", res.ConvertedAmount.String())
	require.Equal(t, "900", res.SourceAccountBalance.String())
	require.Equal(t, "1100", res.TargetAccountBalance.String())
}

func TestConvertError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"missing idempotency key","error_code":"ErrMissingIdempotencyKey"}`))
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)
	require.NoError(t, cl.SetAuth("key-id", "key-secret"))

	_, err := cl.Convert(context.Background(), &ConvertRequest{
		Amount:          decimal.NewFromInt64(100),
		SourceAccountId: 12345,
		TargetAccountId: 12346,
	})
	require.Error(t, err)

	var apiErr Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, "ErrMissingIdempotencyKey", apiErr.ErrCode())
}
