// Package mooonredis
// Wrote by yijian on 2025/03/17
// 安装依赖库：
// go get github.com/stretchr/testify
// go get github.com/stretchr/testify/assert
// 运行测试：
// go test -v -cover
package mooonredis

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
)

// 重构客户端初始化逻辑，每个测试用例独立创建客户端
func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		MinIdleConns: 2,                // 保持最小空闲连接
		IdleTimeout:  30 * time.Second, // 调大空闲超时
	})
}

// TestRentKey_Basic 基础场景测试
// go test -v -run="TestRentKey_Basic"
func TestRentKey_Basic(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	hk := &HoldKey{
		RedisClient: rdb,
		Key:         "test_key",
		Value:       "test_value",
		Expiration:  time.Minute,
	}

	// 用例1：首次获取锁
	ok, err := RentKey(ctx, hk)
	assert.True(t, ok)
	assert.Nil(t, err)

	// 用例2：重复获取相同值的锁
	ok, err = RentKey(ctx, hk)
	assert.True(t, ok)
	assert.Nil(t, err)

	// 用例3：值冲突的获取
	// 这里需要修改 hk 的值以模拟值冲突
	hk.Value = "different_value"
	ok, err = RentKey(ctx, hk)
	assert.False(t, ok)
	assert.Nil(t, err)
}

// TestRentKey_Expiration 过期场景测试
// go test -v -run="TestRentKey_Expiration"
func TestRentKey_Expiration(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	hk := &HoldKey{
		RedisClient: rdb,
		Key:         "test_key",
		Value:       "test_value",
		Expiration:  time.Minute,
	}

	key := "expired_lock"
	// 设置键并主动触发过期检查
	rdb.Set(ctx, key, "v1", 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)

	// 正确判断键不存在的方式（需引入 errors 包）
	_, err := rdb.Get(ctx, key).Result()
	assert.True(t, errors.Is(err, redis.Nil)) // 使用 errors.Is 判断错误类型

	// 修改 hk 的 key 为 expired_lock
	hk.Key = key
	ok, err := RentKey(ctx, hk)
	assert.NoError(t, err)
	assert.True(t, ok)
}

// TestRenewKey_Failure 失败场景测试
func TestRenewKey_Failure(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	hk := &HoldKey{
		RedisClient: rdb,
		Key:         "test_key",
		Value:       "test_value",
		Expiration:  time.Minute,
	}

	// 用例1：值不匹配
	RentKey(ctx, hk)
	// 修改 hk 的值以模拟值不匹配
	hk.Value = "different_value"
	ok, err := RenewKey(ctx, hk)
	assert.False(t, ok)
	assert.NoError(t, err)

	// 用例2：key不存在
	hk.Key = "non_existent_key"
	ok, err = RenewKey(ctx, hk)
	assert.False(t, ok)
	assert.NoError(t, err)
}

// TestReleaseKey_SafeDelete 安全删除测试
func TestReleaseKey_SafeDelete(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	hk := &HoldKey{
		RedisClient: rdb,
		Key:         "test_key",
		Value:       "test_value",
		Expiration:  time.Minute,
	}

	// 初始设置
	RentKey(ctx, hk)

	// 用例1：正确释放
	ok, err := ReleaseKey(ctx, hk)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, int64(0), rdb.Exists(ctx, hk.Key).Val())

	// 用例2：值不匹配时拒绝删除
	RentKey(ctx, hk)
	// 修改 hk 的值以模拟值不匹配
	hk.Value = "different_value"
	ok, err = ReleaseKey(ctx, hk)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, int64(1), rdb.Exists(ctx, hk.Key).Val())
}

// TestReleaseKey_Idempotent 幂等性测试
func TestReleaseKey_Idempotent(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	hk := &HoldKey{
		RedisClient: rdb,
		Key:         "test_key",
		Value:       "test_value",
		Expiration:  time.Minute,
	}

	// 释放不存在的key
	ok, err := ReleaseKey(ctx, hk)
	assert.NoError(t, err)
	assert.True(t, ok) // 符合幂等性设计[8](@ref)
}

// TestConcurrentOperations 并发操作测试
func TestConcurrentOperations(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	hk := &HoldKey{
		RedisClient: rdb,
		Key:         "test_key",
		Value:       "test_value",
		Expiration:  time.Minute,
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 并发获取和释放
			ok, err := RentKey(ctx, hk)
			if err != nil {
				t.Errorf("RentKey error: %v", err)
			}
			if ok {
				defer func() {
					_, err := ReleaseKey(ctx, hk)
					if err != nil {
						t.Errorf("ReleaseKey error: %v", err)
					}
				}()
			}
		}()
	}
	wg.Wait()

	// 最终状态验证
	assert.Equal(t, int64(0), rdb.Exists(ctx, hk.Key).Val())
}

// TestErrorHandling 错误处理测试
func TestErrorHandling(t *testing.T) {
	// 模拟网络错误（关闭客户端）
	closedClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	closedClient.Close()

	hk := &HoldKey{
		RedisClient: closedClient,
		Key:         "test_key",
		Value:       "test_value",
		Expiration:  time.Minute,
	}

	// RentKey错误
	ok, err := RentKey(ctx, hk)
	assert.False(t, ok)
	assert.Error(t, err)

	// ReleaseKey错误
	ok, err = ReleaseKey(ctx, hk)
	assert.False(t, ok)
	assert.Error(t, err)
}
