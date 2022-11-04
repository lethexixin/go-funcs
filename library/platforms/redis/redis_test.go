package redis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

import (
	"github.com/lethexixin/go-funcs/utils/times"
)

import (
	"github.com/google/uuid"
)

func initRedis() (r *Redis, err error) {
	host := "127.0.0.1"
	port := 6379
	password := ""
	db := 0

	r = new(Redis)
	if err = r.Init(Host(host), Port(port), Password(password), DB(db)); err != nil {
		return nil, err
	}
	return r, nil
}

func TestScanKeys(t *testing.T) {
	r, err := initRedis()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	keyPre := "test_key:"
	// 模拟生成 50 个 keyPre 前缀的数据
	for i := 0; i < 50; i++ {
		if err = r.Client.Set(ctx, keyPre+uuid.New().String(), times.Millisecond(), 0).Err(); err != nil {
			t.Error(err)
		}
	}

	chKey := make(chan string, 100)
	// 扫描符合正则表达式 test_key:* 的所有 key
	r.ScanKeys(ctx, 0, fmt.Sprintf("%s*", keyPre), 10, chKey)
	go func() {
		for key := range chKey {
			t.Log("key:", key)
		}
	}()

	select {}
}

func TestScanHashInfos(t *testing.T) {
	r, err := initRedis()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	hashKey := "test_hash_" + uuid.New().String()
	t.Logf("hashKey:%s", hashKey)
	for i := 0; i < 50; i++ {
		if err = r.Client.HSet(ctx, hashKey, uuid.New().String(), times.Millisecond()).Err(); err != nil {
			t.Error(err)
		}
	}

	chHash := make(chan HashInfo, 100)
	r.ScanHashInfos(ctx, hashKey, 0, "", 10, chHash)
	go func() {
		for info := range chHash {
			t.Logf("hashKey:%s, field:%s, value:%s", hashKey, info.Field, info.Value)
		}
	}()

	select {}
}

func TestDelKeysByBatch(t *testing.T) {
	r, err := initRedis()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	keyPre := "test_key:"

	// 模拟生成 1000 个 keyPre 前缀的数据
	for i := 0; i < 1000; i++ {
		if err = r.Client.Set(ctx, keyPre+uuid.New().String(), times.Millisecond(), 0).Err(); err != nil {
			t.Error(err)
		}
	}

	// 睡眠 10s, 此时你可以观察下该 keyPre 的数据信息
	time.Sleep(time.Second * 10)

	chKey := make(chan string, 100)
	// 扫描符合正则表达式 test_key:* 的所有 key
	r.ScanKeys(ctx, 0, fmt.Sprintf("%s*", keyPre), 10, chKey)

	err = r.DelKeysByBatch(ctx, chKey, 100)
	if err != nil {
		t.Error(err)
	}

	t.Log("over")
}

func TestDelHashInfosByBatch(t *testing.T) {
	r, err := initRedis()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	hashKey := "test_hash_" + uuid.New().String()
	t.Logf("hashKey:%s", hashKey)

	// 对 hashKey 模拟生成 1000 个 field 数据
	for i := 0; i < 1000; i++ {
		if err = r.Client.HSet(ctx, hashKey, uuid.New().String(), times.Millisecond()).Err(); err != nil {
			t.Error(err)
		}
	}

	// 睡眠 10s, 此时你可以观察下该 hashKey 的 field 数据信息
	time.Sleep(time.Second * 10)

	chHash := make(chan HashInfo, 100)
	r.ScanHashInfos(ctx, hashKey, 0, "", 10, chHash)

	err = r.DelHashInfosByBatch(ctx, hashKey, chHash, 100)
	if err != nil {
		t.Error(err)
	}

	t.Log("over")
}

func TestCheckKeysTTL(t *testing.T) {
	r, err := initRedis()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	keys := make([]string, 0)
	for i := 1; i <= 50; i++ {
		key := uuid.New().String()
		keys = append(keys, key)
		if err = r.Client.SetEX(ctx, key, times.Millisecond(), time.Duration(i)*time.Minute).Err(); err != nil {
			t.Error(err)
		}
	}

	kt, err := r.CheckKeysTTL(ctx, keys)
	if err != nil {
		t.Error(err)
	}
	for _, v := range kt {
		t.Logf("key:%s,ttl:%f/sec", v.Key, v.TTL.Seconds())
	}

}
