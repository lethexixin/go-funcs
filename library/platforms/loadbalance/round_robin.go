package loadbalance

type RoundRobin struct {
	nextIndex int
}

func (lb *RoundRobin) DoBalance(balance []*Balance, key string) (b *Balance, err error) {
	if len(balance) == 0 {
		return nil, ErrEmptyBalance
	}

	if lb.nextIndex >= len(balance) {
		lb.nextIndex = 0
	}

	b = balance[lb.nextIndex]
	b.callTimes++
	lb.nextIndex = (lb.nextIndex + 1) % len(balance)

	return b, nil

}
