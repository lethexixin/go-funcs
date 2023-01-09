package loadbalance

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyBalance    = errors.New("balance data is empty")
	ErrEmptyKeyForHash = errors.New("hash balance key is empty")
)

// LoadBalancer refer to https://github.com/henkgo/blancer
type LoadBalancer interface {
	DoBalance(balance []*Balance, key string) (b *Balance, err error)
}

type Balance struct {
	addr      string
	port      int64
	callTimes int64

	weight          int64
	currentWeight   int64
	effectiveWeight int64
}

func NewBalance(addr string, port, weight int64) *Balance {
	return &Balance{
		addr:      addr,
		port:      port,
		callTimes: 0,

		weight:          weight,
		currentWeight:   weight,
		effectiveWeight: weight,
	}
}

func (lb *Balance) Addr() string {
	return lb.addr
}

func (lb *Balance) Port() int64 {
	return lb.port
}

func (lb *Balance) Weight() int64 {
	return lb.weight
}

func (lb *Balance) CallTimes() int64 {
	return lb.callTimes
}

func (lb *Balance) String() string {
	return fmt.Sprintf("%+v", *lb)
}
