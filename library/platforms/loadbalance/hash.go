package loadbalance

import (
	"hash/crc32"
)

type Hash struct{}

func (lb *Hash) DoBalance(balance []*Balance, key string) (b *Balance, err error) {
	if len(balance) == 0 {
		return nil, ErrEmptyBalance
	}
	if len(key) == 0 {
		return nil, ErrEmptyKeyForHash
	}

	hash := crc32.Checksum([]byte(key), crc32.MakeTable(crc32.IEEE))
	index := int(hash) % len(balance)
	b = balance[index]
	b.callTimes++
	return b, nil
}
