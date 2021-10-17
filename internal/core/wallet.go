package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

const (
	Version   byte = 0x0
	PubKeyLen      = 64
)

type Wallet struct {
	priv *ecdsa.PrivateKey
}

func (a *Wallet) PublicKey() []byte {
	return append(a.priv.X.Bytes(), a.priv.Y.Bytes()...)
}

func (a *Wallet) PrivateKey() []byte {
	return a.priv.D.Bytes()
}

//使用私钥签名
func (a *Wallet) Sign(msg []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, a.priv, msg)
	if err != nil {
		return nil, ErrWrap("sign error", err)
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

//使用pubKey 校验 msgHash 的签名是否正确
func Verify(msgHash, sign, pubKey []byte) bool {
	pubLen := len(pubKey)
	if pubLen != PubKeyLen {
		Log.Info("Found size invalid pub key")
		return false
	}
	x := new(big.Int)
	y := new(big.Int)
	x.SetBytes(pubKey[:(pubLen / 2)])
	y.SetBytes(pubKey[(pubLen / 2):])

	r := new(big.Int)
	s := new(big.Int)
	sigLen := len(sign)
	r.SetBytes(sign[:(sigLen / 2)])
	s.SetBytes(sign[(sigLen / 2):])
	pub := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return ecdsa.Verify(pub, msgHash, r, s)
}

func (a *Wallet) Address() string {
	pub := a.PublicKey()
	mid := Sha160(Sha256(pub))
	checkSum := Sha256(Sha256(pub))[:4]
	// 1bit version + 20 bit sha160 + 4 bit checksum
	return Base58(ConcatBytes([]byte{Version}, mid, checkSum))
}

//从私钥byte中还原账户
func RestoreWallet(b []byte) *Wallet {
	D := new(big.Int)
	D.SetBytes(b)

	c := elliptic.P256()
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = c
	priv.D = D
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(priv.D.Bytes())
	return &Wallet{priv}
}

//生成一个全新账户
func NewWallet() (*Wallet, error) {
	c := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return nil, ErrWrap("Wallet error", err)
	}
	return &Wallet{priv: privateKey}, nil
}
