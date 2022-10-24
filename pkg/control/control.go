package control

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/log"

	"go.etcd.io/bbolt"
)

type Client struct {
	addr            string
	baseURL         *url.URL
	cancel          context.CancelFunc
	client          *http.Client
	db              *bbolt.DB
	requestInterval time.Duration
	insecure        bool
	disableTLS      bool
	logger          log.Logger
}

func NewControlClient(db *bbolt.DB, addr string, opts ...Option) (*Client, error) {
	baseURL, err := url.Parse("https://" + addr)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}
	c := &Client{
		logger:          log.NewNopLogger(),
		baseURL:         baseURL,
		client:          http.DefaultClient,
		db:              db,
		addr:            addr,
		requestInterval: 60 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.disableTLS {
		c.baseURL.Scheme = "http"
	}

	return c, nil
}

func (c *Client) Start(ctx context.Context) {
	ctx, c.cancel = context.WithCancel(ctx)
	requestTicker := time.NewTicker(c.requestInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-requestTicker.C:
			c.makeRequest(ctx)
		}
	}
}

func (c *Client) do(verb, path string, params interface{}) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	return c.doWithHeaders(verb, path, params, headers)
}

func (c *Client) doWithHeaders(verb, path string, params interface{}, headers map[string]string) (*http.Response, error) {
	var bodyBytes []byte
	var err error
	if params != nil {
		bodyBytes, err = json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("marshaling json: %w", err)
		}
	}

	request, err := http.NewRequest(
		verb,
		c.url(path).String(),
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("creating request object: %w", err)
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return c.client.Do(request)
}

func (c *Client) url(path string) *url.URL {
	u := *c.baseURL
	u.Path = path
	return &u
}

// Stop stops the client
func (c *Client) Stop() {
	c.cancel()
}
