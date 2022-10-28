package snowflake

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestSnowFlake(t *testing.T) {
	s, err := NewSnowflake(int64(0), int64(0))
	assert.NoError(t, err)
	t.Log(s.GenerateId())
	t.Log("epoch:", epoch)
	t.Log("timestampBits:", timestampBits)
	t.Log("dataCenterIdBits:", dataCenterIdBits)
	t.Log("workerIdBits:", workerIdBits)
	t.Log("sequenceBits:", sequenceBits)
	t.Log("timestampMax:", timestampMax)
	t.Log("dataCenterIdMax:", dataCenterIdMax)
	t.Log("workerIdMax:", workerIdMax)
	t.Log("sequenceMask:", sequenceMask)
	t.Log("workerIdShift:", workerIdShift)
	t.Log("dataCenterIdShift:", dataCenterIdShift)
	t.Log("timestampShift:", timestampShift)
}
