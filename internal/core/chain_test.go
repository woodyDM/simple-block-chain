package core

import "testing"


func TestGenesis(t *testing.T) {
	c := Genesis(MockGlobalEvn)
	for hash, tx := range c.Tx {
		if len(tx.Outputs)!=1{
			t.Fatal("len")
		}
		o:=tx.Outputs[0]
		utxo:=c.GetUtxo(o.Address)
		if len(utxo)!=1{
			t.Fatal("utxo len 1")
		}
		utx:=utxo[0]
		if o.TxHash!=utx.TxHash {
			t.Fatal("Tx hash")
		}
		if o.TxIndex!=utx.TxOutputIndex {
			t.Fatal("txidx ")
		}
		if o.TxHash!= hash {
			t.Fatal("tx hash")
		}
		if utx.TxHash!=hash{
			t.Fatal("utxo tx hash")
		}
		if o.Fee!=utx.Fee {
			t.Fatal("fee ")
		}
	}
}

func TestDiff(t *testing.T) {
	_testDiff("fff", "fff", 2, 2, t)
	_testDiff("111", "222", 4, 2, t)
	_testDiff("333", "111", 1, 3, t)
	_testDiff("fff", "555", 1, 3, t)
	//10： 4929388
	_testDiff("2341fac", "4b376c", 2, 15, t)
	//10 135558177
	_testDiff("2341fac", "8147421", 11, 3, t)
}

func _testDiff(in, expect string, a, t int64, te *testing.T) {
	result := diff(in, a, t)
	if result != expect {
		te.Fatal("fail")
	}
}

func TestMemUtxoDb(t *testing.T) {
	db := NewInMemUtxoDatabase()
	add := getTestWallet().Address()
	l1 := db.GetUtxo(add)
	if len(l1) != 0 {
		t.Fatal("0")
	}
	u := &Utxo{
		Address:       add,
		TxHash:        "111",
		TxOutputIndex: 0,
		Fee:           100,
	}
	db.AddUtxo(u)

	if len(db.GetUtxo(add)) != 1 {
		t.Fatal()
	}
	err := db.RemoveUtxo(u)
	if err != nil {
		t.Fatal()
	}
	if len(db.GetUtxo(add)) != 0 {
		t.Fatal()
	}
}


func TestMemUtxoDb_withNotExistUtxo(t *testing.T) {
	db := NewInMemUtxoDatabase()
	add := getTestWallet().Address()

	u := &Utxo{
		Address:       add,
		TxHash:        "111",
		TxOutputIndex: 0,
		Fee:           100,
	}
	db.AddUtxo(u)

	u2 := &Utxo{
		Address:       add,
		TxHash:        "1111",
		TxOutputIndex: 0,
		Fee:           100,
	}

	err := db.RemoveUtxo(u2)
	if err == nil {
		t.Fatal()
	}

}
