package redis

import (
	"context"
	"time"
)

import (
	"github.com/go-redis/redis/v8"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

type HashInfo struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type KeyTTLInfo struct {
	Key string        `json:"key"`
	TTL time.Duration `json:"ttl"`
}

// ScanKeys 通过 scan 命令查找与 match 匹配的 keys.
// 此方法已经开启了一个 goroutine 执行.
// 当 match="" 时相当于全库扫描, match 可以是正则表达式, eg: test* / test_* / test:* / *test*  等等.
// 可参考 redis_test.go TestScanKeys() 方法使用.
func (r *Redis) ScanKeys(ctx context.Context, cursor uint64, match string, count int64, chKey chan<- string) {
	go func() {
		defer close(chKey)
		iter := r.Client.Scan(ctx, cursor, match, count).Iterator()
		for iter.Next(ctx) {
			chKey <- iter.Val()
		}
		if err := iter.Err(); err != nil {
			logger.Errorf("iter.Err():%s", err.Error())
		}
	}()
}

// ScanHashInfos 通过 hscan 命令查找 hashKey 中与 match 匹配的 fields.
// 此方法已经开启了一个 goroutine 执行.
// 当 match="" 时相当于全 hashKey 扫描, match 可以是正则表达式, eg: test* / test_* / test:* / *test*  等等.
// 可参考 redis_test.go TestScanHashInfos() 方法使用.
func (r *Redis) ScanHashInfos(ctx context.Context, hashKey string, cursor uint64, match string, count int64, chHash chan<- HashInfo) {
	go func() {
		defer close(chHash)
		idx := 0
		preIterVal := ""

		// HScan().Iterator() 命令迭代器迭代时, 会以 field1,value1,filed2,value2,field3,value3... 的形式返回数据
		iter := r.Client.HScan(ctx, hashKey, cursor, match, count).Iterator()
		for iter.Next(ctx) {
			idx++
			if idx%2 == 0 {
				chHash <- HashInfo{
					Field: preIterVal,
					Value: iter.Val(),
				}
			}
			preIterVal = iter.Val()
		}
		if err := iter.Err(); err != nil {
			logger.Errorf("iter.Err():%s", err.Error())
		}
	}()
}

// DelKeysByBatch 通过传递 chKey <-chan string , 批量删除 key.
// delBatchSize 默认 100, 代表每批次删除 key 的数量大小.
// 可结合 ScanKeys() 方法使用, 具体可参考 redis_test.go TestDelKeysByBatch() 方法使用.
func (r *Redis) DelKeysByBatch(ctx context.Context, chKey <-chan string, delBatchSize uint) (err error) {
	if delBatchSize <= 0 {
		delBatchSize = 100
	}

	pipe := r.Client.Pipeline()
	for key := range chKey {
		pipe.Del(ctx, key)

		if pipe.Len() < int(delBatchSize) {
			continue
		}

		if _, err = pipe.Exec(ctx); err != nil {
			return err
		}
	}

	// 考虑到 range channel 的时候, 可能 pipe 里面还有数据没有达到 delBatchSize 的大小时就 channel closed 了, 所以仍需要再执行一遍 pipe.Exec(ctx)
	if _, err = pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

// DelHashInfosByBatch 通过传递 chHash <-chan HashInfo , 批量删除 hashKey 里的 field 从而达到删除 hashKey 的功能.
// delBatchSize 默认 100, 代表每批次删除 field 的数量大小.
// 可结合 ScanHashInfos() 方法使用, 具体可参考 redis_test.go TestDelHashInfosByBatch() 方法使用.
func (r *Redis) DelHashInfosByBatch(ctx context.Context, hashKey string, chHash <-chan HashInfo, delBatchSize uint) (err error) {
	if delBatchSize <= 0 {
		delBatchSize = 100
	}

	pipe := r.Client.Pipeline()
	for info := range chHash {
		pipe.HDel(ctx, hashKey, info.Field)

		if pipe.Len() < int(delBatchSize) {
			continue
		}

		if _, err = pipe.Exec(ctx); err != nil {
			return err
		}
	}

	// 考虑到 range channel 的时候, 可能 pipe 里面还有数据没有达到 delBatchSize 的大小时就 channel closed 了, 所以仍需要再执行一遍 pipe.Exec(ctx)
	if _, err = pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

// CheckKeysTTL 批量检查 key 的过期时间.
// 可结合 ScanKeys() 方法使用, 具体可参考 redis_test.go TestCheckKeysTTL() 方法使用.
func (r *Redis) CheckKeysTTL(ctx context.Context, keys []string) (ktInfo []KeyTTLInfo, err error) {
	cmder, err := r.Client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.TTL(ctx, key)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(cmder); i++ {
		d, e := cmder[i].(*redis.DurationCmd).Result()
		if e != nil {
			return nil, e
		}
		ktInfo = append(ktInfo, KeyTTLInfo{
			Key: keys[i],
			TTL: d,
		})
	}
	return ktInfo, nil
}
