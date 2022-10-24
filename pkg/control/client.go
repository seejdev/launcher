package control

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log/level"
	"github.com/kolide/launcher/pkg/osquery"
)

type controlRequest struct {
	// NodeKey string `json:"node_key"`
}

type controlResponse struct {
	Tasks []map[string]string `json:"tasks"`
	Err   string              `json:"error,omitempty"`
	// NodeInvalid bool                `json:"node_invalid,omitempty"`
}

func (c *Client) makeRequest(ctx context.Context) {

	publicKey, fingerprint, err := osquery.PublicKeyFromDB(c.db)
	if err != nil {
		level.Debug(c.logger).Log(
			"msg", "error pub key",
			"err", err,
			"pubkey", publicKey,
			"finger", fingerprint,
		)
		return
	}

	privKey, err := osquery.PrivateKeyFromDB(c.db)
	if err != nil {
		level.Debug(c.logger).Log(
			"msg", "error pub key",
			"err", err,
			"privKey", privKey,
			"finger", fingerprint,
		)
		return
	}

	// nodeKey, err := osquery.NodeKeyFromDB(c.db)
	// if err != nil {
	// 	level.Debug(c.logger).Log(
	// 		"msg", "error getting node key from db to make request to control server",
	// 		"err", err,
	// 	)
	// 	return
	// }

	verb, path := "POST", "/api/v1/control"
	params := &controlRequest{
		// NodeKey: nodeKey,
	}
	response, err := c.do(verb, path, params)
	if err != nil {
		level.Debug(c.logger).Log(
			"msg", "error making request to control server endpoint",
			"err", err,
		)
		return
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusNotFound:
		level.Debug(c.logger).Log(
			"msg", "got 404 making control server request",
			"err", err,
		)
		return
	}

	if response.StatusCode != http.StatusOK {
		level.Debug(c.logger).Log(
			"msg", "got not-ok status code from control server",
			"response_code", response.StatusCode,
		)
		return
	}

	var responseBody controlResponse
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		level.Debug(c.logger).Log(
			"msg", "error decoding control server json",
			"err", err,
		)
		return
	}

	if responseBody.Err != "" {
		level.Debug(c.logger).Log(
			"msg", "response body contained error",
			"err", responseBody.Err,
		)
		return
	}

	if len(responseBody.Tasks) > 0 {
		level.Debug(c.logger).Log(
			"msg", "found tasks",
			"count", len(responseBody.Tasks),
		)

		// for every shell, handle the shell in a goroutine
		// for _, session := range responseBody.Sessions {
		// 	go c.connectToShell(ctx, path, session)
		// }
	}
}

/*
func (c *Client) connectToShell(ctx context.Context, path string, session map[string]string) {
	room, ok := session["session_id"]
	if !ok {
		level.Debug(c.logger).Log(
			"msg", "session didn't contain id",
		)
		return
	}

	secret, ok := session["secret"]
	if !ok {
		level.Debug(c.logger).Log(
			"msg", "session didn't contain secret",
		)
		return
	}

	wsPath := path + "/" + room
	client, err := wsrelay.NewClient(c.addr, wsPath, c.disableTLS, c.insecure)
	if err != nil {
		level.Debug(c.logger).Log(
			"msg", "error creating client",
			"err", err,
		)
		return
	}
	defer client.Close()

	pty, err := ptycmd.NewCmd("/bin/bash", []string{"--login"})
	if err != nil {
		level.Debug(c.logger).Log(
			"msg", "error creating PTY command",
			"err", err,
		)
		return
	}

	TTY, err := webtty.New(
		client,
		pty,
		secret,
		webtty.WithPermitWrite(),
		webtty.WithLogger(c.logger),
		webtty.WithKeepAliveDeadline(),
	)
	if err != nil {
		level.Debug(c.logger).Log(
			"msg", "error creating TTY",
			"err", err,
		)
	}
	if err := TTY.Run(ctx); err != nil {
		level.Debug(c.logger).Log(
			"msg", "error running TTY",
			"err", err,
		)
		return
	}
}*/
