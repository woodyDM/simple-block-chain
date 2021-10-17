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

//计算Merkle根
// https://blog.csdn.net/shangsongwww/article/details/85339243
// https://en.bitcoin.it/wiki/Protocol_documentation
// Note:  Hashes in Merkle Tree displayed in the Block Explorer are of little-endian notation.
// For some implementations and calculations,
// the bytes need to be reversed before they are hashed, and again after the hashing operation.
func MerkleRootStr(txIds []string) (string, error) {
	l := len(txIds)
	if l == 0 {
		panic("Merkle array is 0")
	}
	newTxIds := make([][]byte, 0)
	for _, it := range txIds {
		bs, _ := hex.DecodeString(it)
		// TXIDs  反转一下数组,转化到大端
		newTxIds = append(newTxIds, Reverse(bs))
	}
	root := merkleRoot(newTxIds)
	// 结果是大端，转化为小端
	return hex.EncodeToString(Reverse(root)), nil
}

func Reverse(r []byte) []byte {
	runes := CopyBytes(r)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return runes
}

func merkleRoot(bs [][]byte) []byte {
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
