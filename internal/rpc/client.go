package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/grigata/chta-network-stats/internal/config"
)

type Client struct {
	url      string
	user     string
	password string
	http     *http.Client
}

type request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type response struct {
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
	ID     int             `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(cfg config.RPCConfig) *Client {
	return &Client{
		url:      fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port),
		user:     cfg.User,
		password: cfg.Password,
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) Call(ctx context.Context, method string, params any, result any) error {
	body, err := json.Marshal(request{
		JSONRPC: "1.0",
		ID:      1,
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return fmt.Errorf("marshal rpc request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create rpc request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.user, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("rpc request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read rpc response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rpc returned HTTP %s: %s", resp.Status, string(responseBody))
	}

	var rpcResp response
	if err := json.Unmarshal(responseBody, &rpcResp); err != nil {
		return fmt.Errorf("parse rpc response: %w", err)
	}

	if rpcResp.Error != nil {
		return fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	if result == nil {
		return nil
	}

	if err := json.Unmarshal(rpcResp.Result, result); err != nil {
		return fmt.Errorf("parse rpc result: %w", err)
	}

	return nil
}
