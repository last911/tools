package tests

import (
	"github.com/last911/tools"
	"hash/crc32"
	"strconv"
	"testing"
	"time"
)

func TestAbsolutePath(t *testing.T) {
	path, err := tools.AbsolutePath()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("path:", path)
}

func TestRand(t *testing.T) {
	r := tools.NewRand()
	t.Log(r.Int())
	t.Log(tools.RandRangeInt(0, 1000))
}

func BenchmarkConsistentHash(b *testing.B) {
	ch := tools.NewConsistentHash(32, crc32.ChecksumIEEE)
	n := 1000
	for i := 0; i < n; i++ {
		ch.Add(strconv.Itoa(i))
	}

	for i := 0; i < n; i++ {
		ch.Get(strconv.Itoa(tools.NewRand().Int()))
	}
}

func TestRangeArray(t *testing.T) {
	m := tools.RangeArray(0, 256)
	for _, v := range m {
		t.Log(v)
	}
}
func TestAuthcode(t *testing.T) {
	now := time.Now().UnixNano()
	nn := string(now)
	n, err := tools.Authcode(nn, true, "abc")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("raw:", nn)
	// n := utils.Authcode(nn, true, "abc4sd")
	t.Log("authcode:", n)
	x, err := tools.Authcode(n, false, "abc")
	if err != nil {
		t.Fatal(err)
	}

	// x := utils.Authcode(n, false, "abc4sd")
	if x != nn {
		t.Fatal("authcode decode error")
	}
	t.Log("authcode decode:", x)
}
