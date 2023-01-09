package loadbalance

import (
	"fmt"
	"math"
	"testing"
)

import (
	"github.com/lethexixin/go-funcs/utils/strs"
)

func TestHashBalance(t *testing.T) {
	var bs []*Balance
	for i := 0; i < 5; i++ {
		addr := fmt.Sprintf("192.168.1.%d", i)
		b := NewBalance(addr, 8080, 0)
		bs = append(bs, b)
	}

	var lb LoadBalancer = new(Hash)
	for i := 0; i < 15; i++ {
		uid := strs.UUID4()
		b, err := lb.DoBalance(bs, uid)
		if err == nil {
			t.Log(uid, b.String())
		}
	}
}

func TestRandomBalance(t *testing.T) {
	var bs []*Balance
	for i := 0; i < 5; i++ {
		addr := fmt.Sprintf("192.168.1.%d", i)
		b := NewBalance(addr, 8080, 0)
		bs = append(bs, b)
	}

	var lb LoadBalancer = new(Random)
	for i := 0; i < 15; i++ {
		b, err := lb.DoBalance(bs, "")
		if err == nil {
			t.Log(b.String())
		}
	}
}

func TestRoundRobinBalance(t *testing.T) {
	var bs []*Balance
	for i := 0; i < 5; i++ {
		addr := fmt.Sprintf("192.168.1.%d", i)
		b := NewBalance(addr, 8080, 0)
		bs = append(bs, b)
	}

	var lb LoadBalancer = new(RoundRobin)
	for i := 0; i < 15; i++ {
		b, err := lb.DoBalance(bs, "")
		if err == nil {
			t.Log(b.String())
		}
	}
}

func TestRRWBalance(t *testing.T) {
	var bs []*Balance
	for i := 0; i < 3; i++ {
		addr := fmt.Sprintf("192.168.1.%d", i)
		b := NewBalance(addr, 8080, int64(math.Pow(2, float64(i))))
		bs = append(bs, b)
	}

	var lb LoadBalancer = new(RoundRobinWeight)
	for i := 0; i < 7; i++ {
		b, err := lb.DoBalance(bs, "")
		if err == nil {
			t.Log(b.String())
		}
	}
}

func TestShuffleBalance(t *testing.T) {
	var bs []*Balance
	for i := 0; i < 5; i++ {
		addr := fmt.Sprintf("192.168.1.%d", i)
		b := NewBalance(addr, 8080, 0)
		bs = append(bs, b)
	}

	var lb0 LoadBalancer = new(Shuffle)
	for i := 0; i < 100; i++ {
		b, err := lb0.DoBalance(bs, "")
		if err == nil {
			t.Log("mode 0:", b.String())
		}
	}

	var lb1 LoadBalancer = &Shuffle{Mode: 1}
	for i := 0; i < 100; i++ {
		b, err := lb1.DoBalance(bs, "")
		if err == nil {
			t.Log("mode 1:", b.String())
		}
	}
}
