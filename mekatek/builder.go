package mekatek

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tendermint/tendermint/types"
)

var (
	DefaultBlockBuilderAPIRURL = &url.URL{Scheme: "https", Host: "api.mekatek.xyz"}
	DefaultBlockBuilderTimeout = 1 * time.Second
)

type BlockBuilder interface {
	BuildBlock(context.Context, *BuildBlockRequest) (*BuildBlockResponse, error)
}

func NewBuilder(
	chainID string,
	privValidator types.PrivValidator,
	apiURL *url.URL,
	apiTimeout time.Duration,
	paymentAddr string,
) (BlockBuilder, error) {
	pubKey, err := privValidator.GetPubKey()
	if err != nil {
		return nil, fmt.Errorf("get public key from validator: %w", err)
	}

	bb, err := newHTTPBlockBuilder(apiURL, apiTimeout, privValidator)
	if err != nil {
		return nil, fmt.Errorf("create HTTP block builder: %w", err)
	}

	if _, err = bb.RegisterProposer(context.Background(), &registerProposerRequest{
		ChainID:        chainID,
		PaymentAddress: paymentAddr,
		PubKeyBytes:    pubKey.Bytes(),
		PubKeyType:     pubKey.Type(),
	}); err != nil {
		return nil, fmt.Errorf("register proposer: %w", err)
	}

	return bb, nil
}

func GetURIFromEnv(envkey string, def *url.URL) *url.URL {
	s := os.Getenv(envkey)
	if s == "" {
		return def
	}

	if !strings.HasPrefix(s, "http") {
		s = def.Scheme + s
	}

	u, err := url.Parse(s)
	if err != nil {
		return def
	}

	return &url.URL{
		Scheme: def.Scheme, // default URL defines scheme (e.g. HTTPS)
		Host:   u.Host,     // provided URI defines host:port
		Path:   def.Path,   // default URL defines path
	}
}

func GetDurationFromEnv(envkey string, def time.Duration) time.Duration {
	s := os.Getenv(envkey)
	if s == "" {
		return def
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return def
	}

	return d
}

//
//
//

type httpBlockBuilder struct {
	baseurl   *url.URL
	client    *http.Client
	validator types.PrivValidator
}

func newHTTPBlockBuilder(
	baseurl *url.URL,
	timeout time.Duration,
	v types.PrivValidator,
) (*httpBlockBuilder, error) {
	return &httpBlockBuilder{
		baseurl:   baseurl,
		client:    &http.Client{Timeout: timeout},
		validator: v,
	}, nil
}

func (b *httpBlockBuilder) BuildBlock(
	ctx context.Context,
	req *BuildBlockRequest,
) (*BuildBlockResponse, error) {
	var resp BuildBlockResponse
	return &resp, b.do(ctx, "", req, &resp)
}

func (b *httpBlockBuilder) RegisterProposer(
	ctx context.Context,
	req *registerProposerRequest,
) (*registerProposerResponse, error) {
	var resp registerProposerResponse
	return &resp, b.do(ctx, "/proposers/register", req, &resp)
}

func (b *httpBlockBuilder) do(ctx context.Context, path string, req, resp interface{}) error {
	pubKey, err := b.validator.GetPubKey()
	if err != nil {
		return fmt.Errorf("get public key: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	// TODO: SECURITY 🚨 review, do we need to sign other things than the body?
	// What about nonces (e.g. timestamp)? Are replay attacks possible or exploitable here?
	signature, err := b.validator.SignBytes(body)
	if err != nil {
		return fmt.Errorf("signature failed: %w", err)
	}

	u := b.baseurl
	u.Path = path
	uri := u.String()

	r, err := http.NewRequestWithContext(ctx, "POST", uri, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	r.Header.Set("content-type", "application/json")
	r.Header.Set("mekatek-proposer-address", pubKey.Address().String())
	r.Header.Set("mekatek-request-signature", hex.EncodeToString(signature))

	res, err := b.client.Do(r)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response code %d (%s)", res.StatusCode, strings.TrimSpace(string(body)))
	}

	if err = json.Unmarshal(body, resp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

//
//
//

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

type registerProposerRequest struct {
	ChainID        string `json:"chain_id"`
	PaymentAddress string `json:"payment_address"`
	PubKeyBytes    []byte `json:"pub_key_bytes"`
	PubKeyType     string `json:"pub_key_type"`
}

type registerProposerResponse struct {
	Result string `json:"result"`
}
