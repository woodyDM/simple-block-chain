package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/shengdoushi/base58"
	"golang.org/x/crypto/ripemd160"
)

func ErrWrap(msg string, err error) error {
	return fmt.Errorf("%s: %v ", msg, err)
}

func ErrWrapf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

func Base58(b []byte) string {
	return base58.Encode(b, base58.BitcoinAlphabet)
}

func Base58Decode(str string) ([]byte, error) {
	return base58.Decode(str, base58.BitcoinAlphabet)
}

func Sha256Str(b []byte) string {
	return hex.EncodeToString(Sha256(b))
}

func Sha256(b []byte) []byte {
	hash := sha256.New()
	hash.Write(b)
	return hash.Sum(nil)
}

func Sha160(b []byte) []byte {
	hash := ripemd160.New()
	hash.Write(b)
	return hash.Sum(nil)
}

func Int64ToBytes(i int64) []byte {
	buffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(buffer, binary.BigEndian, i)
	return buffer.Bytes()
}

func CopyBytes(b []byte) []byte {
	des := make([]byte, len(b))
	copy(des, b)
	return des
}

func ConcatBytes(bs ...[]byte) []byte {
	r := make([]byte, 0)
	for _, it := range bs {
		r = append(r, it...)
	}
	return r
}

