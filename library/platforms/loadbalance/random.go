package loadbalance

import (
	"math/rand"
)

type Random struct{}

func (lb *Random) DoBalance(balance []*Balance, key string) (b *Balance, err error) {
	if len(balance) == 0 {
		return nil, ErrEmptyBalance
	}

	index := rand.Intn(len(balance))
	b = balance[index]
	b.callTimes++
	return b, nil
}
