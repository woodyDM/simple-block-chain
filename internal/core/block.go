package core

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

type TxType int32

const (
	NormalTx  TxType = 0 //普通交易
	GenesisTx TxType = 1 //创世交易
)

/**
区块
第0个为创世区块
*/
type Block struct {
	Timestamp int64          //时间戳
	Hash      string         //本区块hash
	Nonce     string         //随机数  64bit (8byte) field
	PreHash   string         //前一区块hash
	Tx        []*Transaction //size>1  第0个一定是CoinbaseTransaction, CoinbaseTransaction 的Input脚本可以是任何bytes,不会校验
	/**
	以下字段可以推断出
	*/
	Height         uint64 // 区块在区块链中的高度 0开始
	TxCount        int    //本区块交易总数,值为 len(Tx)
	PreTxSum       int64  //之前的所有区块交易总数
	PreOutputSum   int64  //之前的所有区块output总数
	MerkleTreeRoot string
	Difficulty     string
}

type Transaction struct {
	Timestamp int64
	//交易类型
	Type TxType
	//输入
	Inputs []*Input
	//输出, 规定：每个Transaction中，一个Address只能有一个Output
	Outputs []*Output
	//额外字段，限制长度为 <= ExtraLen ,可以作为备注等
	Extra []byte
	//Hash 以下为推断字段，仅占位用
	Hash      string
	BlockHash string
}

type Script [][]byte

type Input struct {
	//<sig> <pubKey>
	Script *Script
	//之前某个 tx 的 Output
	Output *Output
}

type Output struct {
	//Coin count
	Fee int64
	//OP_DUP OP_HASH160 OP_PUSH <pubKey160Hash> OP_EQ_VERIFY OP_CHECK_SIGN
	Script *Script
	//output在它所在交易的下标,可以推断得出, 参与Hash
	TxIndex int
	//输出的地址 可以从脚本反推脚本的hash160， 参与Hash计算
	Address string
	//Output所在的tx的 hash
	//*注意*：此字段不参与本tx的Hash计算
	//但是在Input中引用的时候，值必须存在,且需要被计算到Input的Hash中；因为是来自之前就计算好了的tx
	TxHash string
}

// ==================================== func below  ====================================
func genesisBlock() *Block {
	b := &Block{
		Timestamp:    GenesisTime,
		Tx:           createGenesisTx(),
		Height:       0,
		PreHash:      GenesisPreHash,
		PreTxSum:     0,
		PreOutputSum: 0,
	}
	b.TxCount = len(b.Tx)
	txIds := make([]string, 0)
	for _, tx := range b.Tx {
		txIds = append(txIds, tx.Hash)
	}
	err := b.updateMerk()
	if err != nil {
		panic(err)
	}
	b.Difficulty = GenesisDiff
	b.Nonce = GenesisBlockNonce
	b.Hash = GenesisBlockHash
	return b
}

func (b *Block) updateMerk() error {
	txIds := make([]string, 0)
	for _, tx := range b.Tx {
		txIds = append(txIds, tx.Hash)
	}
	merk, e := MerkleRootStr(txIds)
	if e != nil {
		return e
	}
	b.MerkleTreeRoot = merk
	return nil
}

func createGenesisTx() []*Transaction {
	txs := make([]*Transaction, 0)
	for _, priv := range GenesisPrivateKeys {
		outs := make([]*Output, 0)
		account := RestoreWallet(priv)
		out := &Output{
			Fee:     GenesisCoinCount,
			Script:  buildP2PKHOutput(account.PublicKey()),
			TxIndex: 0,
			Address: account.Address(),
		}
		outs = append(outs, out)
		tx := &Transaction{
			Timestamp: GenesisTime,
			Type:      GenesisTx,
			Outputs:   outs,
		}
		err := tx.UpdateHash()
		if err != nil {
			panic(err)
		}
		txs = append(txs, tx)
	}
	return txs

}

// 为了简单 coinbase也用 P2PKH
//OP_DUP OP_HASH160 OP_PUSH <pubKey160Hash> OP_EQ_VERIFY OP_CHECK_SIGN
func buildP2PKHOutput(pubKey []byte) *Script {
	ripemd160 := Sha160(Sha256(pubKey))
	return _output(ripemd160)
}

func buildP2PKHOutputWithAddress(add string) (*Script, error) {
	key, err := AddressToRipemd160PubKey(add)
	if err != nil {
		return nil, err
	}
	return _output(key), nil
}

func _output(ripemd160 []byte) *Script {
	return &Script{
		OpDuplicateA,
		OpSha160A,
		OpPushDataA,
		ripemd160,
		OpEqVerifyA,
		OpCheckSignA,
	}
}

//OP_PUSH <sig> OP_PUSH <pubKey>
func buildP2PKHInput(txHash []byte, w *Wallet) (*Script, error) {
	sign, err := w.Sign(txHash)
	if err != nil {
		return nil, ErrWrap("Failed build input", err)
	}
	return &Script{
		OpPushDataA,
		sign,
		OpPushDataA,
		w.PublicKey(),
	}, nil
}

//cal this transaction hash and update hexHash into Output
func (t *Transaction) UpdateHash() error {
	all := make([][]byte, 0)
	all = append(all, Int64ToBytes(t.Timestamp))
	all = append(all, Int64ToBytes(int64(t.Type)))
	for _, in := range t.Inputs {
		if inHash, err := in.CalHash(); err != nil {
			return err
		} else {
			all = append(all, inHash)
		}
	}
	for _, out := range t.Outputs {
		all = append(all, out.CalThisTxHash())
	}
	all = append(all, t.Extra)
	allSha256 := ConcatBytes(all...)
	txHash := Sha256(allSha256)
	txHashHex := hex.EncodeToString(txHash)
	for _, o := range t.Outputs {
		o.TxHash = txHashHex
	}
	t.Hash = txHashHex
	return nil
}

//计算本tx时用的Hash
func (o *Output) CalThisTxHash() []byte {
	feeBytes := Int64ToBytes(o.Fee)
	scriptBytes := o.Script.CalHash()
	idxBytes := Int64ToBytes(int64(o.TxIndex))
	addHash := Sha256([]byte(o.Address))
	all := ConcatBytes(feeBytes, scriptBytes, idxBytes, addHash)
	return Sha256(all)
}

//在作为Input中位于之前tx时，计算Hash
func (o *Output) CalPreTxHash() ([]byte, error) {
	if o.TxHash == "" {
		return nil, ErrWrapf("Pre Hash should not be empty")
	}
	hashBytes, err := hex.DecodeString(o.TxHash)
	if err != nil {
		return nil, ErrWrap("Pre Hash not exit", err)
	}
	feeBytes := Int64ToBytes(o.Fee)
	scriptBytes := o.Script.CalHash()
	idxBytes := Int64ToBytes(int64(o.TxIndex))
	all := ConcatBytes(feeBytes, scriptBytes, hashBytes, idxBytes)
	return Sha256(all), nil
}

//Input Hash计算
func (i *Input) CalHash() ([]byte, error) {
	scriptHash := i.Script.CalHash()
	outHash, err := i.Output.CalPreTxHash()
	if err != nil {
		return nil, ErrWrap("Input Hash Cal Error", err)
	}
	all := ConcatBytes(scriptHash, outHash)
	return Sha256(all), nil
}

//合并脚本
func ConcatScript(in, out *Script) *Script {
	s := new(Script)
	for _, it := range *in {
		*s = append(*s, it)
	}
	for _, it := range *out {
		*s = append(*s, it)
	}
	return s
}

//ScriptHash计算
func (s *Script) CalHash() []byte {
	return Sha256(ConcatBytes([][]byte(*s)...))
}

//添加script
func (s *Script) append(b []byte) {
	*s = append(*s, b)
}

// ==================================== Block Hash ====================================
type HashResult struct {
	Nonce string
	Hash  string
	Ok    bool
	Err   error
}

func NextNonce() string {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	nonceValue := rd.Int63()
	nonce := fmt.Sprintf("%016x", nonceValue)
	return nonce
}

func (b *Block) TryHash() *HashResult {
	return b.HashWith(NextNonce())
}

func (b *Block) UpdateHash(r *HashResult) {
	if r.Ok {
		b.Nonce = r.Nonce
		b.Hash = r.Hash
	} else {
		panic("Should not use this hash result ")
	}
}

func (b *Block) HashWith(nonce string) *HashResult {
	nonceValue, err := hex.DecodeString(nonce)
	if err != nil {
		return &HashResult{
			Err: err,
		}
	}
	all := make([][]byte, 0)
	preBytes, err := hex.DecodeString(b.PreHash)
	if err != nil {
		return &HashResult{
			Err: err,
		}
	}
	merk, err := hex.DecodeString(b.MerkleTreeRoot)
	if err != nil {
		return &HashResult{
			Err: err,
		}
	}
	all = append(all, Int64ToBytes(b.Timestamp))
	all = append(all, preBytes)
	all = append(all, merk)
	all = append(all, nonceValue)

	if b.Difficulty == "" {
		return &HashResult{
			Err: ErrWrapf("empty block difficulty"),
		}
	}

	allSha256 := ConcatBytes(all...)
	hash := hex.EncodeToString(Sha256(Sha256(allSha256)))
	return &HashResult{
		Nonce: nonce,
		Hash:  hash,
		Ok:    hash < b.Difficulty,
		Err:   nil,
	}
}
