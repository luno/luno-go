// Package luno is a wrapper for the Luno API.
package luno

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	defaultRate    = time.Minute / 300
	defaultBurst   = 1
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
	req, res interface{}, auth bool) error {

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
			url = strings.Replace(url, "{id}", values.Get("id"), -1)
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
		b, err := ioutil.ReadAll(body)
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
