package v1

import (
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/tendermint/tendermint/types"
)

var (
	txLimit   = flag.Int("num-txn", 1, "Number of transactions")
	ptrLog    = flag.String("logfile", "pointer.log", "Pointer log")
	doEvict   = flag.Bool("evict", false, "Simulate eviction")
	doRecheck = flag.Bool("recheck", false, "Simulate re-CheckTx")
)

func TestMempoolAddRemove(t *testing.T) {
	txmp := setup(t, 10000)
	txch := make(chan *WrappedTx, 10)

	f, err := os.Create(*ptrLog)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("Close pointer log: %v", err)
		}
	}()

	m := make([]byte, 1<<20)
	fmt.Fprintf(f, "* %p\n", &m)
	runtime.SetFinalizer(&m, func(m *[]byte) {
		fmt.Fprintf(f, "x %p\n", m)
	})

	go func() {
		defer close(txch)

		b := make([]byte, 10000)
		for i := 0; i < *txLimit; i++ {
			rand.Read(b)
			tx := types.Tx(b)
			wtx := &WrappedTx{
				tx:        tx,
				hash:      tx.Key(),
				timestamp: time.Now().UTC(),
				height:    txmp.height,
			}
			fmt.Fprintf(f, "+ %p\n", wtx)
			runtime.SetFinalizer(wtx, func(w *WrappedTx) {
				fmt.Fprintf(f, "- %p\n", w)
			})

			if *doEvict && txmp.canAddTx(wtx) != nil {
				ev := txmp.priorityIndex.GetEvictableTxs(1000, int64(wtx.Size()),
					txmp.SizeBytes(), txmp.config.MaxTxsBytes)
				for _, v := range ev {
					txmp.removeTx(v, true)
				}
			}
			txmp.insertTx(wtx)
			txch <- wtx
		}
	}()
	var txn []*WrappedTx
	for tx := range txch {
		txn = append(txn, tx)
	}
	for _, tx := range txn {
		if !tx.removed && tx.heapIndex >= 0 {
			txmp.removeTx(tx, true)
		}
	}
	if *doRecheck {
		t.Log("Simulating re-CheckTx")
		txmp.updateReCheckTxs()
		runtime.GC()
	}
	txmp.priorityIndex = NewTxPriorityQueue()
	txmp.Flush()
	runtime.GC()
}
