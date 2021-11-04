package core

import (
	"testing"
)

func TestFilterUsedUtxo(t *testing.T) {
	u:=[]*Utxo{uxto("1","h1",0,1),
		uxto("2","h2",1,12),
		uxto("3","h3",2,13),
		uxto("4","h4",3,14)}

	extra := uxto("5", "h2", 0, 1)
	extra2 := uxto("6", "h6", 1, 1)
	r:= append(u, extra)
	r=append(r,extra2)

	result:= filterUsedUtxo(r, u)
	if len(result)!=2{
		t.Fatal("len")
	}
	if result[0]!= extra {
		t.Fatal("should same 1")
	}
	if result[1]!= extra2 {
		t.Fatal("should same 2")
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
