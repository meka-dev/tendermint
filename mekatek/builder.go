package mekatek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
)

type BlockBuilder interface {
	BuildBlock(context.Context, *BuildBlockRequest) (*BuildBlockResponse, error)
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

func BuilderFromEnv(ctx context.Context, pubKey crypto.PubKey, chainID string) (BlockBuilder, error) {
	apiURL := os.Getenv("MEKATEK_BLOCK_BUILDER_API_URL")
	timeout, _ := time.ParseDuration(os.Getenv("MEKATEK_BLOCK_BUILDER_TIMEOUT"))
	paymentAddress := os.Getenv("MEKATEK_BLOCK_BUILDER_PAYMENT_ADDRESS")

	b, err := newHTTPBlockBuilder(apiURL, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create block builder: %w", err)
	}

	if _, err = b.RegisterProposer(ctx, &registerProposerRequest{
		PaymentAddress: paymentAddress,
		PubKey:         pubKey.Bytes(),
		PubKeyType:     pubKey.Type(),
		ChainID:        chainID,
	}); err != nil {
		return nil, fmt.Errorf("register proposer: %w", err)
	}

	return b, nil
}

//
//
//

type httpBlockBuilder struct {
	baseurl string
	client  *http.Client
}

func newHTTPBlockBuilder(baseurl string, timeout time.Duration) (*httpBlockBuilder, error) {
	if !strings.HasPrefix(baseurl, "http") {
		baseurl = "https://" + baseurl
	}

	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	u.Scheme = "https"
	u.Path = ""
	baseurl = u.String()

	if timeout == 0 {
		timeout = 500 * time.Millisecond
	}

	return &httpBlockBuilder{
		baseurl: baseurl,
		client:  &http.Client{Timeout: timeout},
	}, nil
}

func (b *httpBlockBuilder) BuildBlock(ctx context.Context, req *BuildBlockRequest) (*BuildBlockResponse, error) {
	var resp BuildBlockResponse
	return &resp, b.do(ctx, req, &resp)
}

type registerProposerRequest struct {
	PaymentAddress string `json:"payment_address"`
	PubKey         []byte `json:"pub_key"`
	PubKeyType     string `json:"pub_key_type"`
	ChainID        string `json:"chain_id"`
}

type registerProposerResponse struct {
	Result string `json:"result"`
}

func (b *httpBlockBuilder) RegisterProposer(
	ctx context.Context,
	req *registerProposerRequest,
) (*registerProposerResponse, error) {
	var resp registerProposerResponse
	return &resp, b.do(ctx, req, &resp)
}

func (b *httpBlockBuilder) do(ctx context.Context, req, resp interface{}) error {
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

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response code %d (%s)", res.StatusCode, body)
	}

	if err = json.Unmarshal(body, resp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}
