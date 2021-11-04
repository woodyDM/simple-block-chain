package core

import "math/big"

type TimeProvider func() int64

type GlobalEnv struct {
	UnixTime TimeProvider
}

type BlockChain struct {
	Env *GlobalEnv
	//txs
	*TxDatabase
	//utxo
	UtxoDatabase
	//key block hash
	Blocks map[string]*Block
	//
	Current *Block
}

type TxDatabase struct {
	//key
	Tx map[string]*Transaction
}

type Utxo struct {
	Address string
	TxHash  string
	TxIndex int
	Fee     int64
}

type UtxoDatabase interface {
	AddUtxo(u *Utxo)
	GetUtxo(address string) []*Utxo
	RemoveUtxo(u *Utxo) error
	Clear()
}

type InMemUtxoDatabase struct {
	db map[string][]*Utxo
}

// ==================================== func below ====================================

func NewInMemUtxoDatabase() UtxoDatabase {
	return &InMemUtxoDatabase{db: make(map[string][]*Utxo)}
}

func (i *InMemUtxoDatabase) AddUtxo(u *Utxo) {
	address := u.Address
	_, ok := i.db[address]
	if !ok {
		i.db[address] = make([]*Utxo, 0)
	}
	i.db[address] = append(i.db[address], u)
}

func (i *InMemUtxoDatabase) GetUtxo(address string) []*Utxo {
	return i.db[address]
}

func (i *InMemUtxoDatabase) RemoveUtxo(u *Utxo) error {
	add := u.Address
	l := i.db[add]
	m := make([]*Utxo, 0)
	removed := false
	for _, it := range l {
		if it == u {
			if removed {
				return ErrWrapf("Already removed %v", u)
			}
			removed = true
		} else {
			m = append(m, it)
		}
	}
	if !removed {
		return ErrWrapf("not found utxo %v ", u)
	}
	i.db[add] = m
	return nil
}

func (i *InMemUtxoDatabase) Clear() {
	i.db = make(map[string][]*Utxo)
}

func newUtxo(t *Transaction, txIdx int, o *Output) *Utxo {
	return &Utxo{
		Address: o.Address,
		TxHash:  t.Hash,
		TxIndex: txIdx,
		Fee:     o.Fee,
	}
}

func (c *BlockChain) Size() int {
	return len(c.Blocks)
}

//创世
func Genesis(env *GlobalEnv) *BlockChain {
	chain := &BlockChain{
		TxDatabase: &TxDatabase{
			Tx: make(map[string]*Transaction),
		},
		Blocks:       make(map[string]*Block),
		Env:          env,
		UtxoDatabase: NewInMemUtxoDatabase(),
	}
	block := genesisBlock()
	e := chain.Append(block)
	if e != nil {
		panic(e)
	}
	return chain
}

// 区块链添加一个新的区块，并做简单校验
func (c *BlockChain) Append(b *Block) error {
	ec := checkWhenAppend(b)
	if ec != nil {
		return ec
	}
	_, e := c.Blocks[b.Hash]
	if e {
		Log.Errorf("Same Block Hash found! %s ", b.Hash)
		panic(b.Hash)
	}
	c.Blocks[b.Hash] = b
	c.Current = b
	//update Transactions
	for _, t := range b.Tx {
		_, ok := c.Tx[t.Hash]
		if ok {
			Log.Errorf("Same Transaction Hash found!Block %s ", b.Hash)
			panic(b.Hash)
		}
		c.Tx[t.Hash] = t
	}
	//update utxo
	if b.Height != 0 {
		for idx, t := range b.Tx {
			for _, i := range t.Inputs {
				e := c.RemoveUtxo(newUtxo(t, idx, &i.Output))
				if e != nil {
					panic(ErrWrap("utxo not exist", e))
				}
			}
		}
	}
	for i, t := range b.Tx {
		for _, o := range t.Outputs {
			c.AddUtxo(newUtxo(t, i, o))
		}
	}
	return nil

}

//todo utxo checks
func checkWhenAppend(b *Block) error {
	if b.Nonce == "" {
		return ErrWrapf("Empty nonce in block %v", b.Height)
	}
	if len(b.Tx) == 0 {
		return ErrWrapf("Empty tx in block %v", b.Height)
	}
	if b.Hash == "" {
		return ErrWrapf("Empty hash in block %v", b.Height)
	}
	e := b.HashWith(b.Nonce)
	if !e.Ok {
		return ErrWrapf("Invalid hash in block %v", b.Height)
	}
	if e.Hash != b.Hash {
		return ErrWrapf("Illegal hash in block %v", b.Height)
	}
	return nil
}

//新建区块, 只留下Nonce和 Hash待确定
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

func (c *BlockChain) OnTx(tx TxRequest) {
	Log.Info("chain receive ", tx)

}
