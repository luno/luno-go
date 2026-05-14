package luno

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinked(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/1/users/linked", r.URL.Path)

		_, _, ok := r.BasicAuth()
		require.True(t, ok, "Linked must send credentials")

		_ = json.NewEncoder(w).Encode(LinkedResponse{
			Users: []LinkedUser{
				{
					CreatedAt:   1700000000000,
					Label:       "alice",
					Permissions: []string{"Perm_R_Balance", "Perm_R_Orders"},
					UserId:      "u-1",
				},
				{
					CreatedAt:   1700000001000,
					Label:       "bob",
					Permissions: []string{"Perm_R_Balance"},
					UserId:      "u-2",
				},
			},
		})
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)
	require.NoError(t, cl.SetAuth("key-id", "key-secret"))

	res, err := cl.Linked(context.Background(), &LinkedRequest{})
	require.NoError(t, err)
	require.Len(t, res.Users, 2)

	require.Equal(t, "u-1", res.Users[0].UserId)
	require.Equal(t, "alice", res.Users[0].Label)
	require.Equal(t, int64(1700000000000), res.Users[0].CreatedAt)
	require.Equal(t, []string{"Perm_R_Balance", "Perm_R_Orders"}, res.Users[0].Permissions)

	require.Equal(t, "u-2", res.Users[1].UserId)
	require.Equal(t, []string{"Perm_R_Balance"}, res.Users[1].Permissions)
}
