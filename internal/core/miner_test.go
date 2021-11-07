package core

import (
	"math/rand"
	"testing"
	"time"
)

func __TestNewMiner(t *testing.T) {

	pool := NewTxPool(Genesis(Env))
	NewMiner(pool,getTestWallet_(9))
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tick := time.Tick(1 * time.Second)
	//n := 0
	for {
		w1 := getTestWallet_(rd.Intn(9))
		w2 := getTestWallet_(rd.Intn(9))
		fee := int64(rd.Intn(3)) + 1
		if w1.Address() == w2.Address() {
			continue
		}
		select {
		case <-tick:
			w1.Transform(pool, w2.Address(), fee, time.Now().String())
		}

	}

}
func next(n int) int {
	n = n + 1
	if n == len(GenesisPrivateKeys) {
		return 0
	} else {
		return n
	}

}
