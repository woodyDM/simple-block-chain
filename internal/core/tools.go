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
func MerkleRootStr(bs []string) (string, error) {
	d := make([][]byte, 0)
	for _, it := range bs {
		bs, err := hex.DecodeString(it)
		if err != nil {
			return "", err
		}
		d = append(d, bs)
	}
	return hex.EncodeToString(MerkleRoot(d)), nil
}

func MerkleRoot(bs [][]byte) []byte {
	l := len(bs)
	if l == 0 {
		panic("Merkle array is 0")
	}
	levelOffset := 0
	for levelSize := l; levelSize > 1; levelSize = (levelSize + 1) / 2 {
		for left := 0; left < levelSize; left += 2 {
			right := min(left+1, levelSize-1)
			v := ConcatBytes(bs[left+levelOffset], bs[right+levelOffset])
			bs = append(bs, Sha256(Sha256(v)))
		}
		levelOffset += levelSize
	}
	return bs[len(bs)-1]
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
