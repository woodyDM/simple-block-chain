package core

type TxPool struct {
	Chain     *BlockChain
	usedUtxo  UtxoDatabase
	txReqCh   chan *TxRequest
	txRespCh  chan *TxResponse
	txCh      chan *Transaction
	txBlockCh chan *Block
	endl      chan bool
}

type TxRequest struct {
	From  string
	To    string
	Fee   int64
	Extra string
	w     *Wallet
}

type TxResponse struct {
	tx  *Transaction
	err error
}

func NewErrTxResponse(err error) *TxResponse {
	return &TxResponse{
		err: err,
	}
}

func NewTxPool(c *BlockChain) *TxPool {
	pool := TxPool{
		Chain:     c,
		usedUtxo:  NewInMemUtxoDatabase(),
		txReqCh:   make(chan *TxRequest),
		txRespCh:  make(chan *TxResponse),
		txCh:      make(chan *Transaction, 100),
		txBlockCh: make(chan *Block, 0),
		endl:      make(chan bool),
	}
	go pool.start()
	return &pool
}

func (p *TxPool) Stop() {
	close(p.endl)
}

func (p *TxPool) start() {
	for {
		select {
		case <-p.endl:
			Log.Info("Tx pool stop")
			return
		case req := <-p.txReqCh:
			resp := p.transform0(req)
			if resp.err == nil {
				p.txCh <- resp.tx
			}
			p.txRespCh <- resp
		case block := <-p.txBlockCh:
			p.receiveBlock(block)
		}
	}
}

//todo transform when stop
func (p *TxPool) Transform(tx *TxRequest) *TxResponse {
	extraB := []byte(tx.Extra)
	if len(extraB) > ExtraLen {
		return NewErrTxResponse(ErrWrapf("Extra len exceed max len"))
	}
	_, e := AddressToRipemd160PubKey(tx.w.Address())
	if e != nil {
		return NewErrTxResponse(ErrWrap("Invalid address", e))
	}
	fee := tx.Fee
	if fee <= 0 {
		return NewErrTxResponse(ErrWrapf("Invalid fee %d", fee))
	}
	p.txReqCh <- tx
	return <-p.txRespCh
}

func (p *TxPool) transform0(tx *TxRequest) *TxResponse {
	valid := p.Chain.GetUtxo(tx.From)
	used := p.usedUtxo.GetUtxo(tx.From)
	unused := filterUsedUtxo(valid, used)
	thisUtxo := pickUtxo(unused, tx.Fee)
	if thisUtxo == nil {
		Log.Debug("Not enough utxo for ", tx)
		return NewErrTxResponse(ErrWrapf("No enough utxo for %s", tx.From))
	} else {
		transaction, err := p.createNormalTx(thisUtxo, tx)
		if err != nil {
			return NewErrTxResponse(err)
		}
		err = transaction.UpdateHash()
		if err != nil {
			return NewErrTxResponse(err)
		}
		for _, it := range unused {
			p.usedUtxo.AddUtxo(it)
		}
		Log.Info("TxPool put transaction ", transaction.Hash, " to pool. Request is ", tx)
		return &TxResponse{
			tx:  transaction,
			err: nil,
		}
	}
}

func pickUtxo(uxto []*Utxo, fee int64) []*Utxo {
	used := make([]*Utxo, 0)
	var total int64 = 0
	for _, it := range uxto {
		total += it.Fee
		used = append(used, it)
		if total >= fee {
			return used
		}
	}
	return nil
}

//??????utxo ?????? ??????
func (p *TxPool) createNormalTx(used []*Utxo, tx *TxRequest) (*Transaction, error) {
	trans := &Transaction{
		Timestamp: p.Chain.Env.UnixTime(),
		Type:      NormalTx,
		Extra:     []byte(tx.Extra),
	}
	w := tx.w
	//build input
	inputs := make([]*Input, 0)
	var total int64 = 0
	for _, it := range used {
		if inTx, exist := p.Chain.Tx[it.TxHash]; !exist {
			return nil, ErrWrapf("Transaction %s not found !", it.TxHash)
		} else {
			i := len(inTx.Outputs)
			if i <= it.TxOutputIndex {
				return nil, ErrWrapf("Transaction %s out of index [%d] of total [%d]", it.TxHash, it.TxOutputIndex, i)
			}
			//create input
			output := inTx.Outputs[it.TxOutputIndex]
			script, err := buildP2PKHInput([]byte(inTx.Hash), tx.w)
			if err != nil {
				return nil, ErrWrap("can't create tx", err)
			}
			err = VerifyScript(inTx.Hash, script, output.Script)
			if err != nil {
				return nil, ErrWrap("script verify fail", err)
			}
			in := &Input{
				Script: script,
				Output: output,
			}
			inputs = append(inputs, in)
			total += output.Fee
		}
	}
	trans.Inputs = inputs
	//build output
	left := total - tx.Fee
	if left < 0 {
		panic("expect bonus >=0 ")
	}
	outputs := make([]*Output, 0)
	if left > 0 {
		//create left output
		o1 := &Output{
			Fee:     left,
			Script:  buildP2PKHOutput(w.PublicKey()),
			Address: w.Address(),
		}
		outputs = append(outputs, o1)
	}
	sc, err := buildP2PKHOutputWithAddress(tx.To)
	if err != nil {
		return nil, ErrWrap("can't build output script", err)
	}
	o2 := &Output{
		Fee:     tx.Fee,
		Script:  sc,
		Address: tx.To,
	}
	outputs = append(outputs, o2)
	for i, it := range outputs {
		it.TxIndex = i
	}
	trans.Outputs = outputs
	return trans, nil
}

func (p *TxPool) receiveBlock(block *Block) {
	for _, o := range block.Tx {
		for _, i := range o.Inputs {
			out := i.Output
			e := p.usedUtxo.RemoveUtxo(newUtxo(out))
			if e != nil {
				panic(ErrWrapf("Not found used utxo %v", out))
			}
		}
	}
	Log.Info("Remove used utxos ")
}

func filterUsedUtxo(valid []*Utxo, used []*Utxo) []*Utxo {
	r := make([]*Utxo, 0)
	uMap := make(map[Utxo]bool)
	for _, it := range used {
		uMap[*it] = true
	}
	for _, it := range valid {
		if _, exist := uMap[*it]; !exist {
			r = append(r, it)
		}
	}
	return r
}
