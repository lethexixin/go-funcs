package loadbalance

type RoundRobinWeight struct{}

func (lb *RoundRobinWeight) DoBalance(balance []*Balance, key string) (b *Balance, err error) {
	if len(balance) == 0 {
		return nil, ErrEmptyBalance
	}

	totalWeight := int64(0)
	max := int64(0)
	maxBalanceIdx := 0

	for idx, v := range balance {
		totalWeight += v.effectiveWeight
		v.currentWeight += v.effectiveWeight
		if v.currentWeight > max {
			max = v.currentWeight
			maxBalanceIdx = idx
		}
	}

	balance[maxBalanceIdx].currentWeight -= totalWeight

	b = balance[maxBalanceIdx]
	b.callTimes++

	return b, nil
}
