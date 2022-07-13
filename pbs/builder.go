package pbs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tendermint/tendermint/types"
)

type BlockBuilder interface {
	RegisterProposer(context.Context, *RegisterProposerRequest) (*RegisterProposerResponse, error)
	BuildBlock(context.Context, *BuildBlockRequest) (*BuildBlockResponse, error)
}

type RegisterProposerRequest struct {
	PaymentAddress string `json:"payment_address"`
	PubKey         []byte `json:"pub_key"`
	PubKeyType     string `json:"pub_key_type"`
	ChainID        string `json:"chain_id"`
}

type RegisterProposerResponse struct {
	Result string `json:"result"`
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

func NewHTTPBlockBuilder(baseurl string, timeout time.Duration) (*HTTPBlockBuilder, error) {
	if !strings.HasPrefix(baseurl, "https://") {
		return nil, fmt.Errorf("HTTPBlockBuilder needs an `https://` baseurl")
	}

	if timeout == 0 {
		timeout = 500 * time.Millisecond
	}

	return &HTTPBlockBuilder{
		baseurl: baseurl,
		client:  &http.Client{Timeout: timeout},
	}, nil
}

func (b *HTTPBlockBuilder) RegisterProposer(ctx context.Context, req *RegisterProposerRequest) (*RegisterProposerResponse, error) {
	var resp RegisterProposerResponse
	return &resp, b.do(ctx, req, &resp)
}

func (b *HTTPBlockBuilder) BuildBlock(ctx context.Context, req *BuildBlockRequest) (*BuildBlockResponse, error) {
	var resp BuildBlockResponse
	return &resp, b.do(ctx, req, &resp)
}

func (b *HTTPBlockBuilder) do(ctx context.Context, req, resp any) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	r, err := http.NewRequestWithContext(ctx, "POST", b.baseurl, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	r.Header.Set("content-type", "application/json")

	res, err := b.client.Do(r)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response code %d (%s)", res.StatusCode, body)
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}
