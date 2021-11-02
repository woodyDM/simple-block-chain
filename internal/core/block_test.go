package core

import (
	"encoding/hex"
	"fmt"
	"testing"
)

var (
	Int1000          = []byte{0, 0, 0, 0, 0, 0, 3, 232}
	Int0             = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	Int1             = []byte{0, 0, 0, 0, 0, 0, 0, 1}
	Int2             = []byte{0, 0, 0, 0, 0, 0, 0, 2}
	MockGlobalEvn    = &GlobalEnv{UnixTime: TimeProvider(MockTime())}
	MockTimeInterval = 2000
)

func MockTime() func() int64 {
	var i int64 = GenesisTime
	return func() int64 {
		i += int64(MockTimeInterval)
		return i
	}
}

type Sha256Lib struct {
	ToSha256Factory   map[string]string
	FromSha256Factory map[string]string
}

func NewSha256Lib() *Sha256Lib {
	return &Sha256Lib{
		ToSha256Factory:   make(map[string]string),
		FromSha256Factory: make(map[string]string),
	}
}

func (l *Sha256Lib) PutBytes(b []byte) string {
	s := hex.EncodeToString(b)
	s256 := Sha256Str(b)
	l.ToSha256Factory[s] = s256
	l.FromSha256Factory[s256] = s
	return s256
}

func (l *Sha256Lib) PutBytes2(b []byte) []byte {
	s := hex.EncodeToString(b)
	s256 := Sha256(b)
	str := hex.EncodeToString(s256)
	l.ToSha256Factory[s] = str
	l.FromSha256Factory[str] = s
	return s256
}

func TestScriptHash(t *testing.T) {
	s := &Script{}
	s.append([]byte{OpDuplicate})
	s.append([]byte{1, 2, 3, 4})
	s.append([]byte{1, 2, 3, 4, 5, 6})

	str := Sha256Str([]byte{OpDuplicate, 1, 2, 3, 4, 1, 2, 3, 4, 5, 6})

	if str != hex.EncodeToString(s.CalHash()) {
		t.Fatal("Fail")
	}

	fmt.Println(Int64ToBytes(1000))
	fmt.Println(Int64ToBytes(0))
	fmt.Println(Int64ToBytes(1))

}

func TestOutput_CalThisTxHashHash(t *testing.T) {

	s := &Script{}
	s.append([]byte{OpDuplicate})
	s.append([]byte{1, 2, 3, 4})
	s.append([]byte{1, 2, 3, 4, 5, 6})

	w := getTestWallet()
	out := Output{
		Fee:     1000,
		Script:  s,
		TxIndex: 2,
		Address: w.Address(),
	}

	scriptSha := Sha256([]byte{OpDuplicate, 1, 2, 3, 4, 1, 2, 3, 4, 5, 6})
	str := Sha256Str(ConcatBytes(Int1000, scriptSha, Int2, Sha256([]byte(w.Address()))))

	if str != hex.EncodeToString(out.CalThisTxHash()) {
		t.Fatal("Fail")
	}
}

func TestOutput_CalPreTxHash(t *testing.T) {

	s := &Script{}
	s.append([]byte{OpDuplicate})
	s.append([]byte{1, 2, 3, 4})
	s.append([]byte{1, 2, 3, 4, 5, 6})

	out := Output{
		Fee:     1000,
		Script:  s,
		TxHash:  Sha256Str([]byte("你好")),
		TxIndex: 2,
	}

	scriptSha := Sha256([]byte{OpDuplicate, 1, 2, 3, 4, 1, 2, 3, 4, 5, 6})
	str := Sha256Str(ConcatBytes(Int1000, scriptSha, Sha256([]byte("你好")), Int2))

	hash, err := out.CalPreTxHash()
	if err != nil {
		t.Fatal("err should nil")
	}
	if str != hex.EncodeToString(hash) {
		t.Fatal("Fail")
	}
}

func TestInput_CalHash(t *testing.T) {

	ins := &Script{}
	ins.append([]byte{OpPushData})
	ins.append([]byte{1, 2, 3})
	ins.append([]byte{4, 5, 6})

	s := &Script{}
	s.append([]byte{OpDuplicate})
	s.append([]byte{1, 2, 3, 4})
	s.append([]byte{1, 2, 3, 4, 5, 6})

	out := Output{
		Fee:     1000,
		Script:  s,
		TxHash:  Sha256Str([]byte("你好")),
		TxIndex: 2,
	}

	input := &Input{
		Script: *ins,
		Output: out,
	}

	inScriptSha := Sha256([]byte{OpPushData, 1, 2, 3, 4, 5, 6})
	outScriptSha, _ := out.CalPreTxHash()
	str := Sha256Str(ConcatBytes(inScriptSha, outScriptSha))

	hash, err := input.CalHash()
	if err != nil {
		t.Fatal("err should nil")
	}
	if str != hex.EncodeToString(hash) {
		t.Fatal("Fail")
	}
}

//集成测试 脚本功能 和 vm
func TestScript_VM(t *testing.T) {
	wallet := RestoreWallet(GenesisPrivateKeys[0])
	txHash := Sha256([]byte("Coinbase是每个区块中第一笔交易的特殊名称。也被叫做“创币交易”。\n\n获" +
		"胜的矿工在其区块模版里创建了这个特殊交易。\n\nCoinbase交易与普通交易具有相同的格式，但与普通交易不同的" +
		"是：\n\n只有一个交易输入。\n交易输入的前序输出哈希是0000…0000。"))
	input, err := buildP2PKHInput(txHash, wallet)
	if err != nil {
		t.Fatal(err)
	}
	output := buildP2PKHOutput(wallet.PublicKey())
	allScript := ConcatScript(input, output)
	vm := NewVm(*allScript)
	vm.SetEnv(VMEnvHash, txHash)

	err = vm.Exec()
	if err != nil {
		t.Fatal(err)
	}

}

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
		if o.TxIndex!=utx.TxIndex{
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

func TestGenesisBlock(t *testing.T) {
	block := genesisBlock()
	merk := "ef551a513148cf836ba134ff59492806bdc5d0256816210630694f1634fdfc25"
	if block.MerkleTreeRoot != merk {
		t.Fatal("merk fail")
	}
	if block.Difficulty != GenesisDiff {
		t.Fatal("diff fail")
	}
	r := block.TryHash()
	if r.Ok {
		if r.Hash >= block.Difficulty {
			t.Fatal("hash fail")
		}
	}
	newR:= block.HashWith(GenesisBlockNonce)
	if !newR.Ok{
		t.Fatal("should ok")
	}
	if newR.Hash!=block.Hash {
		t.Fatal("hash check fail")
	}
}

func __TestGenesisHashCal(t *testing.T) {
	block:=genesisBlock()
	var  r *HashResult
	for  {
		r=block.TryHash()
		if r.Ok {
			break
		}
	}
	fmt.Println(r)

}

func TestMemUtxoDb(t *testing.T) {
	db := NewInMemUtxoDatabase()
	add := getTestWallet().Address()
	l1 := db.GetUtxo(add)
	if len(l1) != 0 {
		t.Fatal("0")
	}
	u := &Utxo{
		Address: add,
		TxHash:  "111",
		TxIndex: 0,
		Fee:     100,
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
		Address: add,
		TxHash:  "111",
		TxIndex: 0,
		Fee:     100,
	}
	db.AddUtxo(u)

	u2 := &Utxo{
		Address: add,
		TxHash:  "1111",
		TxIndex: 0,
		Fee:     100,
	}

	err := db.RemoveUtxo(u2)
	if err == nil {
		t.Fatal()
	}

}

func TestNonceGen(t *testing.T) {
	var nonceValue int64=102099534523455
	nonce := fmt.Sprintf("%016x", nonceValue)
	if nonce!="00005cdbe67cac3f"{
		t.Fatal("hex")
	}
	v,_:=hex.DecodeString(nonce)
	b := Int64ToBytes(nonceValue)
	if len(v)!=len(b){
		t.Fatal("f")
	}
	for i:=0;i<len(v);i++{
		if v[i]!=b[i]{
			t.Fatal("idx i fail")
		}
	}
}

