package luno

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendNetworks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/1/send/networks", r.URL.Path)

		require.NoError(t, r.ParseForm())
		require.Equal(t, "ETH", r.Form.Get("currency"))

		_, _, ok := r.BasicAuth()
		require.True(t, ok, "SendNetworks must send credentials")

		_ = json.NewEncoder(w).Encode(SendNetworksResponse{
			Networks: []SendNetwork{
				{Id: 0, Name: "Ethereum", NativeCurrency: "ETH"},
				{Id: 1, Name: "Polygon", NativeCurrency: "MATIC"},
			},
		})
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)
	require.NoError(t, cl.SetAuth("key-id", "key-secret"))

	res, err := cl.SendNetworks(context.Background(), &SendNetworksRequest{
		Currency: "ETH",
	})
	require.NoError(t, err)
	require.Len(t, res.Networks, 2)

	require.Equal(t, int64(0), res.Networks[0].Id)
	require.Equal(t, "Ethereum", res.Networks[0].Name)
	require.Equal(t, "ETH", res.Networks[0].NativeCurrency)

	require.Equal(t, int64(1), res.Networks[1].Id)
	require.Equal(t, "Polygon", res.Networks[1].Name)
	require.Equal(t, "MATIC", res.Networks[1].NativeCurrency)
}
