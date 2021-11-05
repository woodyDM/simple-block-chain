package core

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewMiner(t *testing.T) {

	pool := NewTxPool(Genesis(Env))
	NewMiner(pool)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tick := time.Tick(5 * time.Second)
	for {
		w1 := getTestWallet_(rd.Intn(len(GenesisPrivateKeys)))
		w2 := getTestWallet_(rd.Intn(len(GenesisPrivateKeys)))
		if w1.Address() == w2.Address() {
			continue
		}
		select {
		case <-tick:
			e := w1.Transform(pool, w2.Address(), 1, time.Now().String())
			if e.err != nil {
				t.Logf("Error %v", e.err)
			}
		}
	}

}
