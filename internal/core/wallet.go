package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

const (
	Version      byte = 0x0
	PubKeyLen         = 64
	LenVersion        = 1
	LenCheckSum       = 4
	LenRipemd160      = 20
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

//Version = 1 byte of 0 (zero); on the test network, this is 1 byte of 111
//Key hash = Version concatenated with RIPEMD-160(SHA-256(public key))
//Checksum = 1st 4 bytes of SHA-256(SHA-256(Key hash))
//Bitcoin Address = Base58Encode(Key hash concatenated with Checksum)
func (a *Wallet) Address() string {
	pub := a.PublicKey()
	mid := Sha160(Sha256(pub))
	checkSum := Sha256(Sha256(ConcatBytes([]byte{Version}, mid)))[:LenCheckSum]
	// 1byte version + 20 byte sha160 + 4 byte checksum
	return Base58(ConcatBytes([]byte{Version}, mid, checkSum))
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

func AddressToRipemd160PubKey(add string) ([]byte, error) {
	bs, err := Base58Decode(add)
	if err != nil {
		return nil, ErrWrap("address convert failed", err)
	}
	l := len(bs)
	//
	if l != LenVersion+LenRipemd160+LenCheckSum {
		return nil, ErrWrapf("invalid address,size %d", l)
	}
	if bs[0] != Version {
		return nil, ErrWrapf("invalid version,%d", bs[0])
	}
	result := bs[LenVersion : LenVersion+LenRipemd160]
	doubleSha := Sha256(Sha256(bs[:LenVersion+LenRipemd160]))
	if !bytes.Equal(bs[LenVersion+LenRipemd160:], doubleSha[:LenCheckSum]) {
		return nil, ErrWrapf("invalid address,checksum failed")
	}
	return result, nil
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
