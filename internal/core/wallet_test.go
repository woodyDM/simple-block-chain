package core

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestNewAccount(t *testing.T) {
	account, _ := NewWallet()
	ac2 := RestoreWallet(account.PrivateKey())
	if !bytes.Equal(ac2.PublicKey(), account.PublicKey()) {
		t.Fatal("account fail")
	}
	if !bytes.Equal(ac2.PrivateKey(), account.PrivateKey()) {
		t.Fatal("account priv fail")
	}
	address := account.Address()
	t.Log(address)
}

func TestWallet_Sign_Verify(t *testing.T) {
	msg := "Key 对 txCopy.ID 进行签名。一个 ECDSA 签名就是一对数字，我们对这对数字连接起来，并存储在输入的 Signature 字段"

	wallet := getTestWallet()
	sign, err := wallet.Sign([]byte(msg))
	if err != nil {
		t.Fatal(err)
	}
	ok := Verify([]byte(msg), sign, wallet.PublicKey())
	if !ok {
		t.Fatal("verify should ok")
	}
}

func TestGenAccount(t *testing.T) {
	all := new(bytes.Buffer)

	for i := 0; i < 10; i++ {
		wallet, _ := NewWallet()
		buf := new(bytes.Buffer)
		buf.WriteString("{")
		pv := wallet.PrivateKey()
		for i, it := range pv {
			buf.WriteString(strconv.Itoa(int(it)))
			if i < len(pv)-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("},")
		all.WriteString(buf.String())
	}
}

func TestAddressToRipemd160PubKey(t *testing.T) {
	w := getTestWallet()
	pubkey, err := AddressToRipemd160PubKey(w.Address())
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := Sha160(Sha256(w.PublicKey()))
	if !bytes.Equal(pubkey, pubKey2) {
		t.Fatal("should same")
	}
}

func TestAddressToRipemd160PubKey_WithError(t *testing.T) {
	w := getTestWallet()
	b := []byte(w.Address())
	b[0] = 1
	_, err := AddressToRipemd160PubKey(string(b))
	if !strings.Contains(err.Error(), "invalid base58 string") {
		t.Fatal(err)
	}
}

func getTestWallet() *Wallet {
	return RestoreWallet(GenesisPrivateKeys[0])
}

func getTestWallet2() *Wallet {
	return RestoreWallet(GenesisPrivateKeys[1])
}
