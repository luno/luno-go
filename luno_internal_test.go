package luno

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
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

func TestUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		suffix   string
		expected string
	}{
		{
			name:     "default user agent without suffix",
			suffix:   "",
			expected: fmt.Sprintf("LunoGoSDK/%s", Version),
		},
		{
			name:     "user agent with custom suffix",
			suffix:   "luno-mcp/0.1.0",
			expected: fmt.Sprintf("LunoGoSDK/%s", Version),
		},
		{
			name:     "user agent with complex suffix",
			suffix:   "myapp/1.2.3-beta",
			expected: fmt.Sprintf("LunoGoSDK/%s", Version),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedUserAgent string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedUserAgent = r.Header.Get("User-Agent")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{}"))
			}))
			defer srv.Close()

			cl := NewClient()
			cl.SetBaseURL(srv.URL)
			if tt.suffix != "" {
				cl.SetUserAgentSuffix(tt.suffix)
			}

			var res interface{}
			err := cl.do(context.Background(), "GET", "/test", nil, &res, false)
			if err != nil {
				t.Errorf("Expected success, got %v", err)
			}

			// Check that the user agent starts with the expected prefix
			if !strings.HasPrefix(capturedUserAgent, tt.expected) {
				t.Errorf("Expected User-Agent to start with %q, got %q", tt.expected, capturedUserAgent)
			}

			// Check suffix behavior
			if tt.suffix == "" {
				// Should not contain parentheses when no suffix
				if strings.Contains(capturedUserAgent, "(") || strings.Contains(capturedUserAgent, ")") {
					t.Errorf("Expected no parentheses in User-Agent without suffix, got %q", capturedUserAgent)
				}
			} else {
				// Should contain the suffix in parentheses
				expectedSuffix := "(" + tt.suffix + ")"
				if !strings.Contains(capturedUserAgent, expectedSuffix) {
					t.Errorf("Expected User-Agent to contain %q, got %q", expectedSuffix, capturedUserAgent)
				}
			}
		})
	}
}

func TestMakeUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		suffix   string
		contains []string
		notContains []string
	}{
		{
			name:   "default user agent",
			suffix: "",
			contains: []string{
				"LunoGoSDK/" + Version,
				runtime.Version(),
				runtime.GOOS,
				runtime.GOARCH,
			},
			notContains: []string{"(", ")"},
		},
		{
			name:   "user agent with suffix",
			suffix: "test-app/1.0",
			contains: []string{
				"LunoGoSDK/" + Version,
				runtime.Version(),
				runtime.GOOS,
				runtime.GOARCH,
				"(test-app/1.0)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := NewClient()
			if tt.suffix != "" {
				cl.SetUserAgentSuffix(tt.suffix)
			}

			userAgent := cl.makeUserAgent()
			t.Logf("Generated User-Agent: %s", userAgent)

			for _, expected := range tt.contains {
				if !strings.Contains(userAgent, expected) {
					t.Errorf("Expected User-Agent to contain %q, got %q", expected, userAgent)
				}
			}

			for _, notExpected := range tt.notContains {
				if strings.Contains(userAgent, notExpected) {
					t.Errorf("Expected User-Agent not to contain %q, got %q", notExpected, userAgent)
				}
			}
		})
	}
}
