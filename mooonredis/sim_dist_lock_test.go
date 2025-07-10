// Package mooonredis
// Wrote by yijian on 2025/07/10
package mooonredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 生成测试用的Redis客户端（单机模式，方便本地测试）
func getTestRedisClient(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // 本地Redis地址
	})

	// 测试连接
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Fatalf("无法连接Redis: %v", err)
	}
	return client
}

// 生成测试用的Redis集群客户端（需提前部署集群）
func getTestRedisClusterClient(t *testing.T) *redis.ClusterClient {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:7000", // 集群节点1
			"localhost:7001", // 集群节点2
			"localhost:7002", // 集群节点3
		},
	})

	// 测试连接
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Fatalf("无法连接Redis集群: %v", err)
	}
	return client
}

// 生成测试用的锁键
func getTestKeys() []string {
	return []string{
		"test-lock-1",
		"test-lock-2",
		"test-lock-3",
	}
}

// TestNewSimDistLock 测试锁初始化
// go test -v -run="TestNewSimDistLock"
func TestNewSimDistLock(t *testing.T) {
	ctx := context.Background()
	client := getTestRedisClient(t)
	expiration := 2 * time.Minute
	keys := getTestKeys()

	// 测试自动生成value
	lock := NewSimDistLock(ctx, client, expiration, "", keys)
	assert.NotEmpty(t, lock.value, "自动生成的value不应为空")
	assert.Equal(t, 3, len(lock.keys), "锁键数量应为3个")
	assert.Equal(t, expiration, lock.expiration, "过期时间应正确设置")

	// 测试自定义value
	customValue := "test-value-123"
	lock = NewSimDistLock(ctx, client, expiration, customValue, keys)
	assert.Equal(t, customValue, lock.value, "自定义value应正确设置")
}

// TestTryLock_SingleClient 单客户端测试TryLock
// go test -v -run="TestTryLock_SingleClient"
func TestTryLock_SingleClient(t *testing.T) {
	ctx := context.Background()
	client := getTestRedisClient(t)
	keys := getTestKeys()
	lock := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)
	defer lock.Unlock()

	// 第一次获取锁应成功
	acquired, err := lock.TryLock()
	assert.NoError(t, err, "TryLock不应返回错误")
	assert.True(t, acquired, "第一次获取锁应成功")

	// 第二次获取锁应失败（同一客户端）
	acquired, err = lock.TryLock()
	assert.NoError(t, err, "TryLock不应返回错误")
	assert.False(t, acquired, "同一客户端第二次获取锁应失败")
}

// TestTryLock_MultiClient 多客户端测试TryLock
// go test -v -run="TestTryLock_MultiClient"
func TestTryLock_MultiClient(t *testing.T) {
	ctx := context.Background()
	client := getTestRedisClient(t)
	keys := getTestKeys()
	lock1 := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)
	defer lock1.Unlock()

	// 客户端1获取锁成功
	acquired, err := lock1.TryLock()
	assert.NoError(t, err)
	assert.True(t, acquired)

	// 客户端2获取锁失败
	lock2 := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)
	defer lock2.Unlock()
	acquired, err = lock2.TryLock()
	assert.NoError(t, err)
	assert.False(t, acquired)
}

// TestUnlock 测试释放锁
// go test -v -run="TestUnlock"
func TestUnlock(t *testing.T) {
	ctx := context.Background()
	client := getTestRedisClient(t)
	keys := getTestKeys()
	lock := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)

	// 先获取锁
	acquired, err := lock.TryLock()
	assert.NoError(t, err)
	assert.True(t, acquired)

	// 释放锁
	err = lock.Unlock()
	assert.NoError(t, err, "Unlock不应返回错误")

	// 释放后应能再次获取锁
	acquired, err = lock.TryLock()
	assert.NoError(t, err)
	assert.True(t, acquired, "释放锁后应能再次获取")
}

// TestUnlock_NotHeld 测试释放未持有的锁
// go test -v -run="TestUnlock_NotHeld"
func TestUnlock_NotHeld(t *testing.T) {
	ctx := context.Background()
	client := getTestRedisClient(t)
	keys := getTestKeys()
	lock := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)

	// 未获取锁时释放
	err := lock.Unlock()
	assert.Error(t, err, "释放未持有的锁应返回错误")
}

// TestTimedLock 测试TimedLock
// go test -v -run="TestTimedLock"
func TestTimedLock(t *testing.T) {
	ctx := context.Background()
	client := getTestRedisClient(t)
	keys := getTestKeys()
	lock1 := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)
	defer lock1.Unlock()

	// 客户端1先获取锁
	acquired, err := lock1.TryLock()
	assert.NoError(t, err)
	assert.True(t, acquired)

	// 客户端2尝试TimedLock，设置超时时间1秒
	lock2 := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)
	defer lock2.Unlock()

	startTime := time.Now()
	acquired, err = lock2.TimedLock(1 * time.Second)
	assert.NoError(t, err)
	assert.False(t, acquired, "TimedLock在超时前应未获取到锁")
	assert.GreaterOrEqual(t, time.Since(startTime), 1*time.Second, "TimedLock应等待至少超时时间")
}

// TestTimedLock_Success 测试TimedLock成功获取锁的情况
// go test -v -run="TestTimedLock_Success"
func TestTimedLock_Success(t *testing.T) {
	ctx := context.Background()
	client := getTestRedisClient(t)
	keys := getTestKeys()
	lock1 := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)

	// 客户端1获取锁后立即释放
	acquired, err := lock1.TryLock()
	assert.NoError(t, err)
	assert.True(t, acquired)
	lock1.Unlock()

	// 客户端2尝试TimedLock应成功
	lock2 := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)
	defer lock2.Unlock()
	acquired, err = lock2.TimedLock(1 * time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired, "TimedLock应成功获取锁")
}

// TestConcurrency 并发测试（验证分布式锁的互斥性）
// go test -v -run="TestConcurrency"
func TestConcurrency(t *testing.T) {
	ctx := context.Background()
	client := getTestRedisClient(t)
	keys := getTestKeys()
	counter := 0
	var wg sync.WaitGroup
	workerCount := 10
	mutex := sync.Mutex{} // 用于保护counter的本地锁

	// 启动10个并发工作协程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			lock := NewSimDistLock(ctx, client, 2*time.Minute, "", keys)

			// 获取锁
			acquired, err := lock.TimedLock(5 * time.Second)
			if !assert.NoError(t, err) || !assert.True(t, acquired) {
				return
			}

			// 临界区：对counter进行累加
			mutex.Lock()
			counter++
			mutex.Unlock()

			// 模拟业务操作
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

			// 释放锁
			assert.NoError(t, lock.Unlock())
		}(i)
	}

	wg.Wait()
	assert.Equal(t, workerCount, counter, "并发累加结果应等于工作协程数量（验证互斥性）")
}

// TestClusterSupport 测试Redis集群支持（需提前部署集群）
// go test -v -run="TestClusterSupport"
func TestClusterSupport(t *testing.T) {
	// 跳过集群测试（如需测试请注释此行）
	t.Skip("如需测试集群支持，请注释此行并确保集群已部署")

	ctx := context.Background()
	clusterClient := getTestRedisClusterClient(t)
	keys := getTestKeys()
	lock := NewSimDistLock(ctx, clusterClient, 2*time.Minute, "", keys)
	defer lock.Unlock()

	// 测试TryLock
	acquired, err := lock.TryLock()
	assert.NoError(t, err)
	assert.True(t, acquired)

	// 测试再次获取失败
	acquired, err = lock.TryLock()
	assert.NoError(t, err)
	assert.False(t, acquired)

	// 测试释放锁
	assert.NoError(t, lock.Unlock())

	// 测试释放后可再次获取
	acquired, err = lock.TryLock()
	assert.NoError(t, err)
	assert.True(t, acquired)
}