package core

type Miner struct {
	p  *TxPool
	tx []*Transaction
	w  *Wallet
}

func NewMiner(p *TxPool, w *Wallet) *Miner {
	m := &Miner{
		p:  p,
		tx: make([]*Transaction, 0),
		w:  w,
	}
	go m.Start()
	return m
}

func (m *Miner) Start() {
	Log.Info("miner started!")
	for {
		select {
		case tx := <-m.p.txCh:
			m.handleNewTransaction(tx)
		case <-EndCh:
			Log.Info("Miner stop when ch end")
			return
		}
	}
}

func (m *Miner) handleNewTransaction(tx *Transaction) {
	m.tx = append(m.tx, tx)
	l := len(m.tx)
	if l < TxPerBlock {
		Log.Debug("miner hold tx Len  ", l, " new tx:", tx)
		return
	}
	//create new Block
	toTx := m.tx
	m.tx = make([]*Transaction, 0)
	//to create coinbase tx and bonus
	txAll := m.createNewBlockTx(toTx)
	newBlock, err := m.p.chain.NewBlock(txAll)
	if err != nil {
		Log.Info("Error when create new block!", err)
		return
	}
	var hash *HashResult
	for {
		hash = newBlock.TryHash()
		if hash.Ok {
			break
		}
	}
	newBlock.UpdateHash(hash)
	Log.Info("============ >>  New  block [", newBlock.Height, "] with hash "+newBlock.Hash+" << ==========")
	err = m.p.chain.Append(newBlock)
	if err != nil {
		Log.Error("Error when append to chain ", err)
		return
	}
	m.p.txBlockCh <- newBlock
	//todo sinal tx to clear
}

func (m *Miner) createNewBlockTx(tx []*Transaction) []*Transaction {
	coinbase := &Transaction{
		Timestamp: m.p.chain.Env.UnixTime(),
		Type:      NormalTx,
		Inputs:    make([]*Input, 0),
		Outputs: []*Output{
			{
				Fee:     CoinBaseCount,
				Script:  buildP2PKHOutput(m.w.PublicKey()),
				TxIndex: 0,
				Address: m.w.Address(),
			},
		},
		Extra: []byte("coinbase"),
	}
	err := coinbase.UpdateHash()
	if err != nil {
		panic(err)
	}
	r := make([]*Transaction, 0)
	r = append(r, coinbase)
	r = append(r, tx...)
	return r
}
