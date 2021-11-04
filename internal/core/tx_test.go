package core

import (
	"strings"
	"testing"
)

func TestFilterUsedUtxo(t *testing.T) {
	u := []*Utxo{uxto("1", "h1", 0, 1),
		uxto("2", "h2", 1, 12),
		uxto("3", "h3", 2, 13),
		uxto("4", "h4", 3, 14)}

	extra := uxto("5", "h2", 0, 1)
	extra2 := uxto("6", "h6", 1, 1)
	r := append(u, extra)
	r = append(r, extra2)

	result := filterUsedUtxo(r, u)
	if len(result) != 2 {
		t.Fatal("len")
	}
	if result[0] != extra {
		t.Fatal("should same 1")
	}
	if result[1] != extra2 {
		t.Fatal("should same 2")
	}

}

func TestPickUtxo(t *testing.T) {
	u := []*Utxo{uxto("1", "h1", 0, 10),
		uxto("2", "h2", 1, 20),
		uxto("3", "h3", 2, 30),
		uxto("4", "h4", 3, 40)}
	r := pickUtxo(u, 55)
	if len(r) != 3 {
		t.Fatal("len 3")
	}
	if r[0].Fee != 10 {
		t.Fatal("fee 10")
	}
	if r[1].Fee != 20 {
		t.Fatal("fee 10")
	}
	if r[2].Fee != 30 {
		t.Fatal("fee 10")
	}
}

func TestPickUtxo_NotEnough(t *testing.T) {
	u := []*Utxo{uxto("1", "h1", 0, 10),
		uxto("2", "h2", 1, 20),
		uxto("3", "h3", 2, 30),
		uxto("4", "h4", 3, 40)}
	r := pickUtxo(u, 105)
	if r != nil {
		t.Fatal("fail")
	}
}

func TestCreateNormalTx(t *testing.T) {
	pool := NewTxPool(Genesis(MockGlobalEvn))
	w1 := getTestWallet()
	w2 := getTestWallet2()
	resp := pool.Transform(&TxRequest{
		From:  w1.Address(),
		To:    w2.Address(),
		Fee:   5,
		Extra: "你好👋",
		w:     w1,
	})
	if resp.err != nil {
		t.Fatal("err")
	}
	if len(resp.tx.Outputs) != 2 {
		t.Fatal("len 2")
	}
	tx := resp.tx
	if tx.Type != NormalTx {
		t.Fatal("type err")
	}
	if string(tx.Extra) != "你好👋" {
		t.Fatal("fail")
	}
	if tx.Outputs[0].Fee != 95 {
		t.Fatal("left 95")
	}
	if tx.Outputs[1].Fee != 5 {
		t.Fatal("trans 5")
	}

	utxo := pool.usedUtxo.GetUtxo(w1.Address())
	if len(utxo) != 1 {
		t.Fatal("should use one")
	}
	if len(pool.tx) != 1 {
		t.Fatal("should create 1 tx")
	}
}

func TestCreateNormalTx2(t *testing.T) {
	pool := NewTxPool(Genesis(MockGlobalEvn))
	w1 := getTestWallet()
	w2 := getTestWallet2()
	resp := pool.Transform(&TxRequest{
		From:  w1.Address(),
		To:    w2.Address(),
		Fee:   100,
		Extra: "你好👋",
		w:     w1,
	})
	if resp.err != nil {
		t.Fatal("err")
	}
	if len(resp.tx.Outputs) != 1 {
		t.Fatal("len 1")
	}
	tx := resp.tx
	if tx.Type != NormalTx {
		t.Fatal("type err")
	}
	if string(tx.Extra) != "你好👋" {
		t.Fatal("fail")
	}
	out := tx.Outputs[0]
	if out.Fee != 100 {
		t.Fatal("left 95")
	}
	if out.Address != w2.Address() {
		t.Fatal("add err")
	}
	if out.TxIndex != 0 {
		t.Fatal("idx 0")
	}
	if out.TxHash != tx.Hash {
		t.Fatal("need tx hash")
	}

}

func TestCreateNormalTxErr(t *testing.T) {
	pool := NewTxPool(Genesis(MockGlobalEvn))
	w1 := getTestWallet()
	w2 := getTestWallet2()
	resp := pool.Transform(&TxRequest{
		From:  w1.Address(),
		To:    w2.Address(),
		Fee:   101,
		Extra: "你好👋",
		w:     w1,
	})
	if resp.err == nil {
		t.Fatal("err")
	}
	if !strings.Contains(resp.err.Error(), "No enough utxo fo") {
		t.Fatal("should not enough")
	}
}

func uxto(add, txHash string, idx int, fee int64) *Utxo {
	return &Utxo{
		Address: add,
		TxHash:  txHash,
		TxIndex: idx,
		Fee:     fee,
	}
}
