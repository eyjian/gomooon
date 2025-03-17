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
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
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

	// 用例1：首次获取锁
	ok, err := RentKey(rdb, ctx, "lock1", "v1", time.Minute)
	assert.True(t, ok)
	assert.Nil(t, err)

	// 用例2：重复获取相同值的锁
	ok, _ = RentKey(rdb, ctx, "lock1", "v1", time.Minute)
	assert.True(t, ok) // 值匹配直接返回成功

	// 用例3：值冲突的获取
	ok, _ = RentKey(rdb, ctx, "lock1", "v2", time.Minute)
	assert.False(t, ok)
}

// TestRentKey_Expiration 过期场景测试
// go test -v -run="TestRentKey_Expiration"
func TestRentKey_Expiration(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	key := "expired_lock"
	// 设置键并主动触发过期检查
	rdb.Set(ctx, key, "v1", 1*time.Millisecond)
	//rdb.Exists(ctx, key) // 第一次触发惰性删除
	time.Sleep(10 * time.Millisecond)
	//rdb.Exists(ctx, key) // 第二次确认删除

	// 正确判断键不存在的方式（需引入 errors 包）
	_, err := rdb.Get(ctx, key).Result()
	assert.True(t, errors.Is(err, redis.Nil)) // 使用 errors.Is 判断错误类型

	// 获取锁前重置客户端状态（解决连接池残留问题）
	//rdb.Ping(ctx)

	// 重新获取锁
	ok, _ := RentKey(rdb, ctx, key, "v1", time.Minute)
	assert.True(t, ok)
}

// TestRenewKey_Failure 失败场景测试
func TestRenewKey_Failure(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	// 用例1：值不匹配
	RentKey(rdb, ctx, "renew_fail", "v1", time.Minute)
	ok, _ := RenewKey(rdb, ctx, "renew_fail", "v2", time.Minute)
	assert.False(t, ok)

	// 用例2：key不存在
	ok, _ = RenewKey(rdb, ctx, "not_exist", "v1", time.Minute)
	assert.False(t, ok)
}

// TestReleaseKey_SafeDelete 安全删除测试
func TestReleaseKey_SafeDelete(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	// 初始设置
	RentKey(rdb, ctx, "release_lock", "v1", time.Minute)

	// 用例1：正确释放
	ok, _ := ReleaseKey(rdb, ctx, "release_lock", "v1")
	assert.True(t, ok)
	assert.Equal(t, int64(0), rdb.Exists(ctx, "release_lock").Val())

	// 用例2：值不匹配时拒绝删除
	RentKey(rdb, ctx, "release_lock", "v1", time.Minute)
	ok, _ = ReleaseKey(rdb, ctx, "release_lock", "v2")
	assert.False(t, ok)
	assert.Equal(t, int64(1), rdb.Exists(ctx, "release_lock").Val())
}

// TestReleaseKey_Idempotent 幂等性测试
func TestReleaseKey_Idempotent(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	// 释放不存在的key
	ok, _ := ReleaseKey(rdb, ctx, "ghost_key", "v1")
	assert.True(t, ok) // 符合幂等性设计[8](@ref)
}

// TestConcurrentOperations 并发操作测试
func TestConcurrentOperations(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 并发获取和释放
			ok, _ := RentKey(rdb, ctx, "concurrent_lock", "v1", time.Minute)
			if ok {
				defer ReleaseKey(rdb, ctx, "concurrent_lock", "v1")
			}
		}()
	}
	wg.Wait()

	// 最终状态验证
	assert.Equal(t, int64(0), rdb.Exists(ctx, "concurrent_lock").Val())
}

// TestErrorHandling 错误处理测试
func TestErrorHandling(t *testing.T) {
	// 模拟网络错误（关闭客户端）
	closedClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	closedClient.Close()

	// RentKey错误
	ok, err := RentKey(closedClient, ctx, "error_key", "v1", time.Second)
	assert.False(t, ok)
	assert.Error(t, err)

	// ReleaseKey错误
	ok, err = ReleaseKey(closedClient, ctx, "error_key", "v1")
	assert.False(t, ok)
	assert.Error(t, err)
}