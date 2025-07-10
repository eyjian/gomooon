// Package mooonredis
// Wrote by yijian on 2025/07/10
package mooonredis

import (
	"context"
	"fmt"
	"github.com/eyjian/gomooon/mooonutils"
	"github.com/go-redis/redis/v8"
	"math"
	"sync"
	"time"
)

// SimDistLock 简单的分布式锁
type SimDistLock struct {
	ctx        context.Context
	client     redis.UniversalClient
	expiration time.Duration // 设置至少 1 分钟的值
	value      string        // 生成唯一值，用于标识锁的持有者
	keys       []string      // 独一无二的 Key 数组（固定为 3 个，对 3 个 keys 都成功才算成功，同 redis 集群节点数无关）
}

// NewSimDistLock 生成分布式锁实例
// 参数 value 如果为空，则自动生成生成唯一值用于标识锁的持有者
func NewSimDistLock(ctx context.Context, client redis.UniversalClient, expiration time.Duration, value string, keys []string) *SimDistLock {
	if len(keys) != 3 {
		panic("keys must be exactly 3 for this distributed lock implementation")
	}
	val := value
	if val == "" {
		val = mooonutils.GetNonceStr(32)
	}
	return &SimDistLock{
		ctx:        ctx,
		client:     client,
		expiration: expiration,
		value:      value,
		keys:       keys,
	}
}

// TryLock 尝试获取锁（非阻塞）
// 成功返回 true, nil；超时返回 false, nil；出错返回 false, err
func (dl *SimDistLock) TryLock() (bool, error) {
	var (
		acquiredKeys []string
		mu           sync.Mutex
	)

	var wg sync.WaitGroup
	errs := make(chan error, len(dl.keys))

	// 并行尝试在每个key上获取锁
	for _, key := range dl.keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			// 使用SET命令的NX和PX选项原子性地获取锁
			set, err := dl.client.SetNX(dl.ctx, k, dl.value, dl.expiration).Result()
			if err != nil {
				errs <- err
				return
			}

			if set {
				mu.Lock()
				acquiredKeys = append(acquiredKeys, k)
				mu.Unlock()
			}
		}(key)
	}

	// 等待所有尝试完成
	wg.Wait()
	close(errs)

	// 检查是否有错误
	var errorList []error
	for err := range errs {
		errorList = append(errorList, err)
	}

	if len(errorList) > 0 {
		// 释放已获取的锁
		dl.unlockKeys(acquiredKeys)
		return false, fmt.Errorf("errors while acquiring lock: %v", errorList)
	}

	// 检查是否成功获取了所有3个key的锁
	if len(acquiredKeys) == len(dl.keys) {
		return true, nil
	}

	// 获取锁失败，释放已获取的部分锁
	if len(acquiredKeys) > 0 {
		dl.unlockKeys(acquiredKeys)
	}

	return false, nil
}

// TimedLock 获取锁（阻塞）
// 成功返回 true, nil；超时返回 false, nil；出错返回 false, err
func (dl *SimDistLock) TimedLock(timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(dl.ctx, timeout)
	defer cancel()

	// 初始重试间隔
	retryInterval := 50 * time.Millisecond
	// 最大重试间隔
	maxRetryInterval := 500 * time.Millisecond

	for {
		// 尝试获取锁
		acquired, err := dl.TryLock()
		if err != nil {
			return false, err
		}

		if acquired {
			return true, nil
		}

		// 检查是否超时
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			// 指数退避：每次重试间隔翻倍，但不超过最大值
			time.Sleep(retryInterval)
			retryInterval = time.Duration(math.Min(float64(retryInterval*2), float64(maxRetryInterval)))
		}
	}
}

// Unlock 释放锁（原子操作）
// 成功返回 nil，出错返回 err
func (dl *SimDistLock) Unlock() error {
	var wg sync.WaitGroup
	errs := make(chan error, len(dl.keys))

	// 并行释放所有key上的锁
	for _, key := range dl.keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			// 使用Lua脚本确保原子性释放锁
			script := `
				if redis.call("GET", KEYS[1]) == ARGV[1] then
					return redis.call("DEL", KEYS[1])
				else
					return 0
				end
			`

			_, err := dl.client.Eval(dl.ctx, script, []string{k}, dl.value).Result()
			if err != nil {
				errs <- err
			}
		}(key)
	}

	// 等待所有释放操作完成
	wg.Wait()
	close(errs)

	// 收集并返回所有错误
	var errorList []error
	for err := range errs {
		errorList = append(errorList, err)
	}

	if len(errorList) > 0 {
		return fmt.Errorf("failed to unlock on %d keys: %v", len(errorList), errorList)
	}

	return nil
}

// 释放指定keys上的锁
func (dl *SimDistLock) unlockKeys(keys []string) {
	var wg sync.WaitGroup

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			script := `
				if redis.call("GET", KEYS[1]) == ARGV[1] then
					return redis.call("DEL", KEYS[1])
				else
					return 0
				end
			`

			dl.client.Eval(dl.ctx, script, []string{k}, dl.value).Result()
		}(key)
	}

	wg.Wait()
}