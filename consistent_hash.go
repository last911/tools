package tools

import (
	"encoding/hex"
	"github.com/golang/groupcache/consistenthash"
	"strconv"
)

// ConsistentHash struct
type ConsistentHash struct {
	*consistenthash.Map
}

// NewConsistentHash return ConsistentHash
func NewConsistentHash(n int, fn ...consistenthash.Hash) *ConsistentHash {
	var f consistenthash.Hash
	if len(fn) > 0 {
		f = fn[0]
	}
	return &ConsistentHash{consistenthash.New(n, f)}
}

// Md5Hash uint32
func Md5Hash(key []byte) uint32 {
	k := Md5Sum(string(key))
	ks := k[0:8]
	dst := make([]byte, hex.DecodedLen(len(ks)))
	hex.Decode(dst, []byte(ks))

	i, err := strconv.ParseUint(string(dst), 16, 32)
	var ii uint32
	if err != nil {
		ii = NewRand().Uint32()
	}

	if ii != 0 {
		return ii
	}

	return uint32(i)
}
