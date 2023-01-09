package loadbalance

import (
	"math/rand"
	"time"
)

type Shuffle struct {
	Mode int
}

func (lb *Shuffle) DoBalance(balance []*Balance, key string) (b *Balance, err error) {
	if len(balance) == 0 {
		return nil, ErrEmptyBalance
	}

	rand.Seed(time.Now().UnixNano())

	switch lb.Mode {
	case 0:
		for i := 1; i <= len(balance); i++ {
			lastIdx := len(balance) - 1
			idx := rand.Intn(i)
			balance[idx], balance[lastIdx] = balance[lastIdx], balance[idx]
		}
	default:
		for i := 0; i < len(balance)/2; i++ {
			x := rand.Intn(len(balance))
			y := rand.Intn(len(balance))
			balance[x], balance[y] = balance[y], balance[x]
		}
	}

	b = balance[0]
	b.callTimes++

	return b, nil
}
