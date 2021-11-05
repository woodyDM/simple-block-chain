package core

type Miner struct {
	p  *TxPool
	tx []*Transaction
}

func NewMiner(p *TxPool) *Miner {
	m := &Miner{p: p, tx: make([]*Transaction, 0)}
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
	if l < 5 {
		Log.Debug("miner hold tx Len  ", l, " new tx:", tx)
		return
	}
	//create new Block
	toTx := m.tx
	m.tx = make([]*Transaction, 0)
	newBlock, err := m.p.chain.NewBlock(toTx)
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
	Log.Info("New hash found for block ", newBlock.Height, " with hash "+newBlock.Hash)
	err = m.p.chain.Append(newBlock)
	if err != nil {
		Log.Error("Error when append to chain ", err)
		return
	}
	m.p.txBlockCh <- newBlock
	//todo sinal tx to clear
}
