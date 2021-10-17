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

// see : https://btc.com/btc/block/123456
func TestMerkleRoot(t *testing.T) {
	tx := []string{
		"5b75086dafeede555fc8f9a810d8b10df57c46f9f176ccc3dd8d2fa20edd685b",
		"e3d0425ab346dd5b76f44c222a4bb5d16640a4247050ef82462ab17e229c83b4",
		"137d247eca8b99dee58e1e9232014183a5c5a9e338001a0109df32794cdcc92e",
		"5fd167f7b8c417e59106ef5acfe181b09d71b8353a61a55a2f01aa266af5412d",
		"60925f1948b71f429d514ead7ae7391e0edf965bf5a60331398dae24c6964774",
		"d4d5fc1529487527e9873256934dfb1e4cdcb39f4c0509577ca19bfad6c5d28f",
		"7b29d65e5018c56a33652085dbb13f2df39a1a9942bfe1f7e78e97919a6bdea2",
		"0b89e120efd0a4674c127a76ff5f7590ca304e6a064fbc51adffbd7ce3a3deef",
		"603f2044da9656084174cfb5812feaf510f862d3addcf70cacce3dc55dab446e",
		"9a4ed892b43a4df916a7a1213b78e83cd83f5695f635d535c94b2b65ffb144d3",
		"dda726e3dad9504dce5098dfab5064ecd4a7650bfe854bb2606da3152b60e427",
		"e46ea8b4d68719b65ead930f07f1f3804cb3701014f8e6d76c4bdbc390893b94",
		"864a102aeedf53dd9b2baab4eeb898c5083fde6141113e0606b664c41fe15e1f",
	}

	result, err := MerkleRootStr(tx)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	if result != "0e60651a9934e8f0decd1c5fde39309e48fca0cd1c84a21ddfde95033762d86c" {
		t.Fatal("merklet fail")
	}
}
