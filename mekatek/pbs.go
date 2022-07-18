package mekatek

import (
	"fmt"
	"github.com/tendermint/tendermint/types"
)

type Proposer struct {
	v types.PrivValidator
}

func NewProposer(v types.PrivValidator) *Proposer {
	return &Proposer{v: v}
}

func (p Proposer) PubKey() (bytes []byte, typ, addr string, err error) {
	pubKey, err := p.v.GetPubKey()
	if err != nil {
		return nil, "", "", fmt.Errorf("mekatek.Proposer PubKey error: %w", err)
	}

	return pubKey.Bytes(), pubKey.Type(), pubKey.Address().String(), nil
}

func (p Proposer) Sign(b []byte) ([]byte, error) {
	return p.v.SignBytes(b)
}
