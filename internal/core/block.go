package core

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

type TxType int32

const (
	NormalTx           TxType = 0 //普通交易
	GenesisTx          TxType = 1 //创世交易
	GenesisCoinCount          = 100
	GenesisTime               = 1630814880000
	GenesisDiff               = "00000fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" //60f
	GenesisPreHash            = "0000000000000000000000000000000000000000000000000000000000000000" //60f
	DiffTargetSpacing         = 1 * 60                                                             //1min 一个区块
	DiffTargetTimeSpan        = 60 * 60                                                            // 每60分钟调整一次难度
	DiffIntervalBlock         = DiffTargetTimeSpan / DiffTargetSpacing                             //60次以后，调整难度

)

var (
	//仅仅是为了后续方便编码和测试，实际上
	GenesisPrivateKeys = [][]byte{{44, 190, 182, 28, 72, 154, 195, 227, 70, 39, 86, 55, 22, 45, 247, 94, 231, 212, 68, 207, 32, 212, 252, 144, 140, 150, 134, 231, 1, 40, 214, 69},
		{37, 175, 36, 250, 25, 142, 150, 140, 15, 59, 114, 33, 160, 85, 234, 46, 232, 8, 148, 252, 209, 35, 247, 208, 198, 208, 180, 87, 199, 123, 21, 163},
		{124, 193, 148, 216, 238, 84, 77, 65, 123, 33, 174, 115, 84, 138, 92, 104, 208, 203, 126, 6, 46, 101, 141, 154, 10, 90, 248, 108, 65, 53, 156, 45},
		{46, 36, 217, 131, 42, 20, 225, 33, 77, 192, 9, 13, 131, 25, 55, 129, 202, 78, 248, 36, 103, 23, 63, 199, 46, 78, 148, 12, 62, 33, 238, 254},
		{189, 204, 180, 135, 97, 95, 152, 255, 132, 51, 102, 4, 100, 111, 175, 247, 227, 152, 149, 246, 69, 251, 238, 114, 55, 205, 60, 17, 36, 82, 180, 216},
		{115, 182, 146, 98, 119, 63, 178, 120, 29, 60, 255, 102, 176, 176, 15, 40, 130, 12, 249, 89, 30, 102, 236, 163, 27, 251, 175, 89, 243, 36, 252, 203},
		{216, 75, 15, 252, 154, 49, 236, 216, 126, 126, 233, 68, 77, 110, 52, 19, 205, 186, 255, 127, 113, 130, 49, 84, 86, 123, 205, 130, 240, 226, 130, 231},
		{174, 39, 70, 72, 166, 168, 162, 221, 205, 9, 50, 194, 57, 6, 61, 141, 89, 143, 163, 126, 39, 68, 160, 59, 244, 234, 204, 175, 222, 246, 47, 34},
		{144, 210, 192, 20, 2, 137, 110, 100, 71, 14, 196, 100, 97, 190, 61, 110, 207, 240, 60, 0, 9, 157, 164, 111, 176, 14, 251, 28, 27, 142, 27, 54},
		{73, 83, 74, 17, 154, 230, 214, 34, 134, 38, 20, 96, 177, 79, 86, 84, 175, 253, 240, 58, 120, 168, 81, 230, 215, 12, 43, 71, 92, 164, 5, 167}}

	Env = &GlobalEnv{UnixTime: func() int64 {
		return time.Now().Unix()
	}}
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

type GlobalEnv struct {
	UnixTime TimeProvider
}

type TimeProvider func() int64

type BlockChain struct {
	Env *GlobalEnv
	*TxDatabase
	Blocks  map[string]*Block
	Current *Block
}

type TxDatabase struct {
	Tx map[string]*Transaction
}

type Transaction struct {
	Timestamp int64
	//交易类型
	Type TxType
	//输入
	Inputs []*Input
	//输出, 规定：每个Transaction中，一个Address只能有一个Output
	Outputs []*Output
	//额外字段，限制长度为 <=100bytes,可以作为备注等
	Extra []byte
	//Hash 以下为推断字段，仅占位用
	Hash string
}

type Input struct {
	//<sig> <pubKey>
	Script Script
	//之前某个 tx 的 Output
	Output Output
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

type Script [][]byte

//---------------------- func below --------------------------
func (c *BlockChain) Size() int {
	return len(c.Blocks)
}

//创世
func Genesis(env *GlobalEnv) *BlockChain {
	chain := &BlockChain{
		TxDatabase: &TxDatabase{
			Tx: make(map[string]*Transaction),
		},
		Blocks: make(map[string]*Block),
		Env:    env,
	}
	block := genesisBlock()
	chain.Append(block)
	return chain
}

// 区块链添加一个新的已校验的区块
func (c *BlockChain) Append(b *Block) {
	_, e := c.Blocks[b.Hash]
	if e {
		Log.Errorf("Same Block Hash found! %s ", b.Hash)
		panic(b.Hash)
	}
	c.Blocks[b.Hash] = b
	for _, t := range b.Tx {
		_, ok := c.Tx[t.Hash]
		if ok {
			Log.Errorf("Same Transaction Hash found!Block %s ", b.Hash)
			panic(b.Hash)
		}
	}
	c.Current = b
}

func genesisBlock() *Block {
	b := &Block{
		Timestamp:    GenesisTime,
		Tx:           createGenesisTx(),
		Height:       0,
		PreHash:      GenesisPreHash,
		PreTxSum:     0,
		PreOutputSum: 0,
	}
	txIds := make([]string, 0)
	for _, tx := range b.Tx {
		txIds = append(txIds, tx.Hash)
	}
	err := b.updateMerk()
	if err != nil {
		panic(err)
	}
	b.Difficulty = GenesisDiff
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

//新建区块,待添加Tx ，计算Hash等操作
func (c *BlockChain) NewBlock(tx []*Transaction) (*Block, error) {

	pre := c.Current
	preOutputCount := 0
	for _, t := range pre.Tx {
		preOutputCount += len(t.Outputs)
	}
	b := &Block{
		Timestamp:    c.Env.UnixTime(),
		PreHash:      pre.Hash,
		Tx:           tx,
		Height:       pre.Height + 1,
		TxCount:      len(tx),
		PreTxSum:     pre.PreTxSum + int64(pre.TxCount),
		PreOutputSum: pre.PreOutputSum + int64(preOutputCount),
		Difficulty:   c.NextDifficulty(),
	}
	err := b.updateMerk()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func checkTx(tx []*Transaction) error {
	l := len(tx)
	if l < 2 {
		return ErrWrapf("tx length < 2")
	}
	for _, it := range tx {
		if it.Timestamp == 0 {
			return ErrWrapf("tx time is 0")
		}
		if it.Type != NormalTx {
			return ErrWrapf("tx should be normal")
		}
		if len(it.Extra) > 100 {
			return ErrWrapf("tx extra len >100")
		}
	}
	return nil

}

func createGenesisTx() []*Transaction {
	txs := make([]*Transaction, 0)
	for i, priv := range GenesisPrivateKeys {
		outs := make([]*Output, 0)
		account := RestoreWallet(priv)
		out := &Output{
			Fee:     GenesisCoinCount,
			Script:  buildP2PKHOutput(account.PublicKey()),
			TxIndex: i,
			Address: account.Address(),
		}
		outs = append(outs, out)
		tx := &Transaction{
			Timestamp: GenesisTime,
			Type:      GenesisTx,
			Outputs:   outs,
		}
		hash, err := tx.CalHash()
		if err != nil {
			panic(err)
		}
		tx.Hash = hex.EncodeToString(hash)
		for _, o := range tx.Outputs {
			o.TxHash = tx.Hash
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

//
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

//添加一笔已交易
func (b *Block) AppendTx(tx *Transaction) {
	tx.Type = NormalTx
	b.Tx = append(b.Tx, tx)
}

//todo
func (b *Block) CheckWith(c *BlockChain) error {
	return nil
}

//cal transaction hash with all field
func (t *Transaction) CalHash() ([]byte, error) {
	all := make([][]byte, 0)
	all = append(all, Int64ToBytes(t.Timestamp))
	all = append(all, Int64ToBytes(int64(int32(t.Type))))
	for _, in := range t.Inputs {
		if inHash, err := in.CalHash(); err != nil {
			return nil, err
		} else {
			all = append(all, inHash)
		}
	}
	for _, out := range t.Outputs {
		all = append(all, out.CalThisTxHash())
	}
	all = append(all, t.Extra)
	allSha256 := ConcatBytes(all...)
	return Sha256(allSha256), nil
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

// ==================================== Difficulty ====================================
func (c *BlockChain) NextDifficulty() string {
	b := c.Current
	if b.Height == 0 {
		return GenesisDiff
	}
	//Only change once per interval
	if (b.Height+1)%DiffIntervalBlock != 0 {
		return b.Difficulty
	}
	var first = b
	for i := 0; i < DiffIntervalBlock-1; i++ {
		first = c.Blocks[first.PreHash]
	}
	var actualSpan = b.Timestamp - first.Timestamp
	if actualSpan < DiffTargetTimeSpan/4 {
		actualSpan = DiffTargetTimeSpan / 4
	}
	if actualSpan > DiffTargetTimeSpan*4 {
		actualSpan = DiffTargetTimeSpan * 4
	}
	return diff(b.Difficulty, actualSpan, DiffTargetTimeSpan)
}

func diff(curDiff string, actualSpan, targetSpan int64) string {
	////新的难度值 = 旧难度值 * （nActualTimespan/nTargetTimespan）
	oldDiff, ok := new(big.Int).SetString(curDiff, 16)
	if !ok {
		panic(ErrWrapf("invalid hex diff %s", curDiff))
	}
	r := new(big.Int)
	actualSpanB := new(big.Int).SetInt64(actualSpan)
	targetSpanB := new(big.Int).SetInt64(targetSpan)
	r.Mul(oldDiff, actualSpanB)
	r = r.Div(r, targetSpanB)
	return r.Text(16)
}

// ==================================== Block Hash ====================================

type HashResult struct {
	Nonce string
	Hash  string
	Ok    bool
	Err   error
}

func (b *Block) TryHash() *HashResult {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	nonceValue := rd.Int63()
	all := make([][]byte, 0)
	preBytes, err := hex.DecodeString(b.PreHash)
	if err != nil {
		return &HashResult{
			Ok:  false,
			Err: err,
		}
	}
	merk, err := hex.DecodeString(b.MerkleTreeRoot)
	if err != nil {
		return &HashResult{
			Ok:  false,
			Err: err,
		}
	}
	all = append(all, Int64ToBytes(b.Timestamp))
	all = append(all, preBytes)
	all = append(all, merk)
	all = append(all, Int64ToBytes(nonceValue))
	nonce := fmt.Sprintf("%08x", nonceValue)
	if b.Difficulty == "" {
		return &HashResult{
			Ok:  false,
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
