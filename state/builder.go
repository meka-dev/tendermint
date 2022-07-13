package state

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tendermint/tendermint/types"
)

type BlockBuilder interface {
	BuildBlock(ctx context.Context, req *BuildBlockRequest) (*BuildBlockResponse, error)
}

type BuildBlockRequest struct {
	ProposerAddress string    `json:"proposer_address"`
	ChainID         string    `json:"chain_id"`
	Height          int64     `json:"height"`
	Txs             types.Txs `json:"txs"`
	MaxBytes        int64     `json:"max_bytes"`
	MaxGas          int64     `json:"max_gas"`
}

type BuildBlockResponse struct {
	Txs types.Txs `json:"txs"`
}

//
//
//

type HTTPBlockBuilder struct {
	baseurl string
	client  *http.Client
}

func NewHTTPBlockBuilder(baseurl string, timeout time.Duration) *HTTPBlockBuilder {
	return &HTTPBlockBuilder{
		baseurl: baseurl,
		client:  &http.Client{Timeout: timeout},
	}
}

func (b *HTTPBlockBuilder) BuildBlock(ctx context.Context, req *BuildBlockRequest) (*BuildBlockResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	r, err := http.NewRequestWithContext(ctx, "POST", b.baseurl, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	r.Header.Set("content-type", "application/json")

	res, err := b.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("execute: %w", err)
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code %d (%s)", res.StatusCode, body)
	}

	var resp BuildBlockResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}
