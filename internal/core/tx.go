package core

type TxPool struct {
	chain    *BlockChain
	usedUtxo UtxoDatabase
	tx       []*Transaction
	txReqCh  chan *TxRequest
	txRespCh chan *TxResponse
	endl     chan bool
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

func NewTxPool(c *BlockChain) *TxPool {
	pool := TxPool{
		chain:    c,
		usedUtxo: NewInMemUtxoDatabase(),
		tx:       make([]*Transaction, 0),
		txReqCh:  make(chan *TxRequest),
		txRespCh: make(chan *TxResponse),
		endl:     make(chan bool),
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
			p.txRespCh <- p.transform0(req)
		}
	}
}

//todo transform when stop
func (p *TxPool) Transform(tx *TxRequest) error {
	p.txReqCh <- tx
	resp := <-p.txRespCh
	return resp.err
}

//todo in some go routine
func (p *TxPool) transform0(tx *TxRequest) *TxResponse {
	valid := p.chain.GetUtxo(tx.From)
	used := p.usedUtxo.GetUtxo(tx.From)
	unused := filterUsedUtxo(valid, used)
	thisUtxo := pickUtxo(unused, tx.Fee)
	if thisUtxo == nil {
		Log.Info("Not enough utxo for ", tx)
		return &TxResponse{
			err: ErrWrapf("No enough utxo for %s", tx.From),
		}
	} else {
		transaction, err := p.createNormalTx(thisUtxo, tx)
		if err != nil {
			return &TxResponse{
				err: err,
			}
		}
		//todo update chain and txpool
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

func (p *TxPool) createNormalTx(used []*Utxo, tx *TxRequest) (*Transaction, error) {
	trans := &Transaction{
		Timestamp: p.chain.Env.UnixTime(),
		Type:      NormalTx,
		Inputs:    nil,
		Outputs:   nil,
		Extra:     []byte(tx.Extra),
	}
	w := tx.w
	//build input
	inputs := make([]*Input, 0)
	var total int64 = 0
	for _, it := range used {
		if inTx, exist := p.chain.Tx[it.TxHash]; !exist {
			return nil, ErrWrapf("Transaction %s not found !", it.TxHash)
		} else {
			i := len(inTx.Outputs)
			if i <= it.TxIndex {
				return nil, ErrWrapf("Transaction %s out of index %d of %d", it.TxHash, it.TxIndex, i)
			}
			//create input
			output := inTx.Outputs[it.TxIndex]
			script, err := buildP2PKHInput([]byte(inTx.Hash), tx.w)
			if err != nil {
				return nil, ErrWrap("can't create tx", err)
			}
			in := &Input{
				Script: *script,
				Output: *output,
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
		//create bonus output fo
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
