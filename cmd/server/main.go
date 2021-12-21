package main

import (
	"github.com/woodyDM/simple-block-chain/internal/core"
	"math/rand"
	"time"
)

func main() {
	pool := core.NewTxPool(core.Genesis(core.Env))
	core.NewMiner(pool, core.GetTestWallet(9))
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tick := time.Tick(1 * time.Second)
	for {
		w1 := core.GetTestWallet(rd.Intn(9))
		w2 := core.GetTestWallet(rd.Intn(9))
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
