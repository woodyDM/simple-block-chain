package core

import (
	"encoding/hex"
	"testing"
)

func TestSha256(t *testing.T) {
	dig := Sha256Str([]byte("你好"))
	if dig != "670d9743542cae3ea7ebe36af56bd53648b0a1126162e78d81a32934a711302e" {
		t.Fail()
	}
}

func TestSha160(t *testing.T) {
	dig := Sha160([]byte("0a8efab5eba4157330b3113690508ee944e684a3c6949c00d64fcab2d565e5a6"))
	if hex.EncodeToString(dig) != "467e91fb39f73c685378861f49f53fe6225932c7" {
		t.Fail()
	}
}

func TestInt64ToBytes(t *testing.T) {
	var i int64 = 257
	bs := Int64ToBytes(i)
	if bs[7] != 1 {
		t.Fail()
	}
}

func TestConcatBytes(t *testing.T) {
	a1 := []byte{1, 2, 3}
	a2 := []byte{}
	a3 := []byte{5, 6, 7}
	exp := []byte{1, 2, 3, 5, 6, 7}
	r := ConcatBytes(a1, a2, a3)
	for i, it := range r {
		if it != exp[i] {
			t.Fail()
		}
	}
	b := [][]byte{{1, 2, 3}, {5, 6, 7}}
	r2 := ConcatBytes(b...)
	for i, it := range r2 {
		if it != exp[i] {
			t.Fail()
		}
	}
}
