package main

import (
	"fmt"
	"sync"
	"time"
)

// refer to https://github.com/GUAIK-ORG/go-snowflake

const (
	epoch             = int64(1577808000000)                           // 设置起始时间(时间戳/毫秒):2020-01-01 00:00:00,有效期69年
	timestampBits     = uint(41)                                       // 时间戳占用位数
	dataCenterIdBits  = uint(3)                                        // 数据中心id所占位数
	workerIdBits      = uint(7)                                        // 机器id所占位数
	sequenceBits      = uint(12)                                       // 序列所占的位数
	timestampMax      = int64(-1 ^ (-1 << timestampBits))              // 时间戳最大值
	dataCenterIdMax   = int64(-1 ^ (-1 << dataCenterIdBits))           // 支持的最大数据中心id数量
	workerIdMax       = int64(-1 ^ (-1 << workerIdBits))               // 支持的最大机器id数量
	sequenceMask      = int64(-1 ^ (-1 << sequenceBits))               // 支持的最大序列id数量
	workerIdShift     = sequenceBits                                   // 机器id左移位数
	dataCenterIdShift = sequenceBits + workerIdBits                    // 数据中心id左移位数
	timestampShift    = sequenceBits + workerIdBits + dataCenterIdBits // 时间戳左移位数
)

type Snowflake struct {
	mu           sync.Mutex
	timestamp    int64
	workerId     int64
	dataCenterId int64
	sequence     int64
}

func NewSnowflake(dataCenterId, workerId int64) (sf *Snowflake, err error) {
	if dataCenterId < 0 || dataCenterId > dataCenterIdMax {
		return nil, fmt.Errorf("dataCenterId must be between 0 and %d", dataCenterIdMax-1)
	}

	if workerId < 0 || workerId > workerIdMax {
		return nil, fmt.Errorf("workerId must be between 0 and %d", workerIdMax-1)
	}

	return &Snowflake{
		timestamp:    0,
		dataCenterId: dataCenterId,
		workerId:     workerId,
		sequence:     0,
	}, nil
}

func (s *Snowflake) GenerateId() (id int64, err error) {
	// 获取id最关键的一点 加锁 加锁 加锁
	s.mu.Lock()
	defer s.mu.Unlock() // 生成完成后记得 解锁 解锁 解锁
	now := time.Now().UnixMilli()
	if s.timestamp == now {
		// 当同一时间戳(精度:毫秒)下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度,则需要等待下一毫秒,下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 不同时间戳(精度:毫秒)下直接使用序列号:0
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		return 0, fmt.Errorf("epoch must be between 0 and %d", timestampMax-1)
	}
	s.timestamp = now
	return int64((t)<<timestampShift | (s.dataCenterId << dataCenterIdShift) | (s.workerId << workerIdShift) | (s.sequence)), nil
}
