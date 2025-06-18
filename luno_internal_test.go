package luno

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDoGet(t *testing.T) {
	type testRes struct {
		Value string `json:"value"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(testRes{Value: "test"})
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)

	var res testRes
	err := cl.do(context.Background(), "GET", "/test", nil, &res, false)
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
}

func TestDoGetWithParametes(t *testing.T) {
	type testReq struct {
		Value string `url:"value"`
	}
	type testRes struct {
		Value string `json:"value"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(testRes{Value: r.FormValue("value")})
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)

	now := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	req := testReq{Value: now}

	var res testRes
	err := cl.do(context.Background(), "GET", "/test", &req, &res, false)
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
	if res.Value != now {
		t.Errorf("Expected %q, got %q", now, res.Value)
	}
}

func TestDoPost(t *testing.T) {
	type testReq struct {
		Value string `url:"value"`
	}
	type testRes struct {
		Value string `json:"value"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}
		json.NewEncoder(w).Encode(testRes{Value: r.FormValue("value")})
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)

	now := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	req := testReq{Value: now}

	var res testRes
	err := cl.do(context.Background(), "POST", "/test", &req, &res, false)
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
	if res.Value != now {
		t.Errorf("Expected %q, got %q", now, res.Value)
	}
}

func TestDoAuth(t *testing.T) {
	type testRes struct {
		Username string `json:"string"`
		Password string `json:"password"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()
		json.NewEncoder(w).Encode(testRes{
			Username: username,
			Password: password,
		})
	}))
	defer srv.Close()

	now := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	username := now + "_user"
	password := now + "_password"

	cl := NewClient()
	cl.SetBaseURL(srv.URL)
	err := cl.SetAuth(username, password)
	require.Nil(t, err)
	cl.SetDebug(true)
	cl.SetTimeout(10 * time.Second)
	var res testRes

	// No auth provided:
	err = cl.do(context.Background(), "POST", "/test", nil, &res, false)
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
	if res.Username != "" {
		t.Errorf("Expected empty username, got %q", res.Username)
	}

	// Auth provided:
	err = cl.do(context.Background(), "POST", "/test", nil, &res, true)
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
	if res.Username != username || res.Password != password {
		t.Errorf("Expected %s:%s, got %s:%s", username, password,
			res.Username, res.Password)
	}
}

func TestDoURLReplacement(t *testing.T) {
	type testReq struct {
		ID    int64  `url:"id"`
		Value string `url:"value"`
	}
	type testRes struct {
		Path  string `json:"path"`
		Value string `json:"value"`
		ID    string `json:"id"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(testRes{
			Path:  r.URL.Path,
			Value: r.FormValue("value"),
			ID:    r.FormValue("id"), // should be blank
		})
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)

	r := time.Now().UnixNano() / 1e6
	value := fmt.Sprintf("%d_value", r)
	req := testReq{ID: r, Value: value}

	var res testRes
	err := cl.do(context.Background(), "GET", "/test/{id}", &req, &res, false)
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
	expPath := fmt.Sprintf("/test/%d", req.ID)
	if res.Path != expPath {
		t.Errorf("Expected %q, got %q", expPath, res.Path)
	}
	if res.ID != "" {
		t.Errorf("Expected blank string, got %q", res.ID)
	}
}

func TestDoJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("this is invalid JSON"))
	}))
	defer srv.Close()

	cl := NewClient()
	cl.SetBaseURL(srv.URL)

	var res interface{}
	err := cl.do(context.Background(), "GET", "/", nil, &res, false)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	str400 := "400 Bad Request"
	if !strings.Contains(err.Error(), str400) {
		t.Errorf("Expected error string to contain %q, got %q",
			str400, err.Error())
	}
}
