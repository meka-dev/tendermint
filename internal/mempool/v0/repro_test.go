package v0

import (
	"context"
	"crypto/rand"
	"errors"
	"flag"
	"runtime"
	"testing"

	abciclient "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/abci/example/kvstore"
	"github.com/tendermint/tendermint/internal/mempool"
	"github.com/tendermint/tendermint/types"
)

var (
	txLimit = flag.Int("num-txn", 1, "Number of transactions")
	ptrLog  = flag.String("logfile", "pointer.log", "Pointer log")
)

func TestMempoolAddRemove(t *testing.T) {
	app := kvstore.NewApplication()
	cc := abciclient.NewLocalCreator(app)
	mp, cleanup, err := newMempoolWithApp(cc)
	if err != nil {
		t.Fatalf("Setup: %v", err)
	}
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < *txLimit; i++ {
		info := mempool.TxInfo{SenderID: mempool.UnknownPeerID}
		data := make([]byte, 10000)
		rand.Read(data)
		if err := mp.CheckTx(ctx, data, nil, info); err != nil {
			var mpf types.ErrMempoolIsFull
			if !errors.As(err, &mpf) {
				t.Fatalf("CheckTx: %v", err)
			}
		}
	}

	mp.Flush()
	runtime.GC()
}
