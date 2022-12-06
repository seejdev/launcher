package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type transport struct {
	authToken string
	base      http.Transport
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.authToken))
	return t.base.RoundTrip(req)
}

type client struct {
	baseURL *url.URL
	base    http.Client
}

type status struct {
	Status string `json:"status"`
}

func New(authToken, socketPath string) client {
	transport := &transport{
		authToken: authToken,
		base: http.Transport{
			DialContext: dialContext(socketPath),
		},
	}

	client := client{
		base: http.Client{
			Transport: transport,
			Timeout:   1 * time.Second,
		},
	}

	return client
}

func (c *client) Shutdown() error {
	resp, err := c.base.Get("http://unix/shutdown")
	if err != nil {
		return err
	}

	if resp.Body != nil {
		resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *client) Ping() error {
	resp, err := c.base.Get("http://unix/ping")
	if err != nil {
		return err
	}

	if resp.Body != nil {
		resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

type controlRequest struct {
	// TODO: This is a temporary and simple data format for phase 1
	Message string `json:"message"`
}

func (c *client) SetStatus(st string) error {
	params := &status{
		Status: st,
	}

	response, err := c.do("POST", "http://unix/status", params)
	if err != nil {
		// level.Error(c.logger).Log(
		// 	"msg", "error making request to control server endpoint",
		// 	"err", err,
		// )
		return err
	}
	defer response.Body.Close()

	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	// level.Error(c.logger).Log(
	// 	// 	"msg", "error reading response body from control server",
	// 	// 	"err", err,
	// 	// )
	// 	return err
	// }

	// bodyBytes, err := json.Marshal(params)
	// if err != nil {
	// 	return fmt.Errorf("marshaling json: %w", err)
	// }

	// resp, err := c.base.Post("http://unix/status", "application/json", bytes.NewBuffer(bodyBytes))
	// if err != nil {
	// 	return err
	// }

	// if resp.Body != nil {
	// 	resp.Body.Close()
	// }

	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }

	return nil
}

func (c *client) do(verb string, url string, params interface{}) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	return c.doWithHeaders(verb, url, params, headers)
}

func (c *client) doWithHeaders(verb, url string, params interface{}, headers map[string]string) (*http.Response, error) {
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
		url,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("creating request object: %w", err)
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return c.base.Do(request)
}

func (c *client) url(path string) *url.URL {
	u := *c.baseURL
	u.Path = path
	return &u
}
