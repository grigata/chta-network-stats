package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/grigata/chta-network-stats/internal/rpc"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) GetBlockCount(ctx context.Context) (int64, error) {
	url := c.baseURL + "/api/getblockcount"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("create getblockcount request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("getblockcount request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read getblockcount response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(
			"getblockcount returned HTTP %s: %s",
			resp.Status,
			strings.TrimSpace(string(body)),
		)
	}

	var height int64
	if err := json.Unmarshal(body, &height); err != nil {
		return 0, fmt.Errorf(
			"parse getblockcount response %q: %w",
			strings.TrimSpace(string(body)),
			err,
		)
	}

	return height, nil
}
func (c *Client) GetBlockHash(ctx context.Context, height int64) (string, error) {
	url := fmt.Sprintf(
		"%s/api/getblockhash?index=%d",
		c.baseURL,
		height,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create getblockhash request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("getblockhash request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read getblockhash response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"getblockhash returned HTTP %s: %s",
			resp.Status,
			strings.TrimSpace(string(body)),
		)
	}

	hash := strings.TrimSpace(string(body))

	if hash == "" {
		return "", fmt.Errorf("getblockhash returned an empty response")
	}

	return hash, nil

}
func (c *Client) GetBlock(
	ctx context.Context,
	hash string,
) (*rpc.Block, error) {
	url := fmt.Sprintf(
		"%s/api/getblock?hash=%s",
		c.baseURL,
		hash,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create getblock request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getblock request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read getblock response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"getblock returned HTTP %s: %s",
			resp.Status,
			strings.TrimSpace(string(body)),
		)
	}

	var block rpc.Block
	if err := json.Unmarshal(body, &block); err != nil {
		return nil, fmt.Errorf(
			"parse getblock response: %w",
			err,
		)
	}

	return &block, nil
}
func (c *Client) GetRawTransaction(
	ctx context.Context,
	txid string,
) (*rpc.RawTransaction, error) {
	url := fmt.Sprintf(
		"%s/api/getrawtransaction?txid=%s&decrypt=1",
		c.baseURL,
		txid,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create getrawtransaction request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getrawtransaction request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read getrawtransaction response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"getrawtransaction returned HTTP %s: %s",
			resp.Status,
			strings.TrimSpace(string(body)),
		)
	}

	var tx rpc.RawTransaction
	if err := json.Unmarshal(body, &tx); err != nil {
		return nil, fmt.Errorf(
			"parse getrawtransaction response: %w",
			err,
		)
	}

	return &tx, nil
}
