// Package mooonredis
// Wrote by yijian on 2025/03/17
package mooonredis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var rentScript = redis.NewScript(`
local key = KEYS[1]
local new_val = ARGV[1]
local ttl_ms = tonumber(ARGV[2])

local current_val = redis.call("GET", key)
if current_val == new_val then
	local remaining = redis.call("PTTL", key)
	if remaining > 0 then
		return 2 -- key已存在，且未过期
	end
elseif current_val then
	return 0 -- key已存在，但值不匹配
end

local ok = redis.call("SET", key, new_val, "PX", ttl_ms, "NX")
if ok then
	return 1 -- key不存在，设置成功
end
return -1 -- key已存在，且值不匹配
`)

var renewScript = redis.NewScript(`
local key = KEYS[1]
local target_val = ARGV[1]
local ttl_ms = tonumber(ARGV[2])

local current_val = redis.call("GET", key)
if not current_val then
	return 0  -- key不存在
end

if current_val ~= target_val then
	return 0  -- 值不匹配
end

redis.call("PEXPIRE", key, ttl_ms)
return 1 -- key存在，续期成功
`)

var releaseScript = redis.NewScript(`
local key = KEYS[1]
local target_val = ARGV[1]

local current_val = redis.call("GET", key)
if not current_val then
	return 1  -- key不存在，无需操作
end

if current_val ~= target_val then
	return 0  -- 值不匹配，拒绝操作
end

redis.call("DEL", key)
return 1 -- key存在，释放成功
`)

type HoldKey struct {
	RedisClient redis.UniversalClient
	Key         string        // 定时器独一无二的 key
	Value       string        // 建议使用 POD_IP
	Expiration  time.Duration // 设置至少 1 分钟的值
}

// RentKey 用于实现分布式锁
// 返回值：
// 1. bool: 是否成功
// 2. error: 错误信息
// 说明：
// 1. 如果 key 不存在，则设置 key 为 value，并返回 true
// 2. 如果 key 存在，但 value 与当前值相同，则返回 true
// 3. 如果 key 存在，但 value 与当前值不同，则返回 false
// 4. 如果设置 key 时发生错误，则返回 false, err
func RentKey(ctx context.Context, hk *HoldKey) (bool, error) {
	keys := []string{hk.Key}
	args := []interface{}{
		hk.Value,      // value
		hk.Expiration, // ttl
	}

	res, err := rentScript.Run(ctx, hk.RedisClient, keys, args...).Int()
	if err != nil {
		return false, err // 错误时返回 false, err
	}

	switch res {
	case 1:
		return true, nil
	case 2:
		return true, nil
	case 0:
		return false, nil
	case -1: // 明确处理竞态失败
		return false, nil
	default:
		return false, nil
	}
}

// RenewKey 用于续期，如果 key 和 value 还未过期，则续期
// 返回值：
// 1. bool: 续期是否成功
// 2. error: 错误信息
// 说明：
// 1. 续期成功返回 true, nil；续期失败返回 false, nil
// 2. 如果续期时发生错误，则返回 false, err
func RenewKey(ctx context.Context, hk *HoldKey) (bool, error) {
	keys := []string{hk.Key}
	args := []interface{}{
		hk.Value,
		hk.Expiration,
	}

	res, err := renewScript.Run(ctx, hk.RedisClient, keys, args...).Int()
	if err != nil {
		return false, err // 网络错误或脚本执行错误
	}

	return res == 1, nil
}

// ReleaseKey 释放分布式锁（仅当值与持有者匹配时删除）
// 返回值：
// 1. bool: true=成功释放，false=释放失败
// 2. error: 错误信息
// 说明：
// 1. 仅当 key 存在且 value 匹配时才删除锁，避免误删其他客户端的锁
// 2. key 不存在时视为已释放，返回 true
func ReleaseKey(ctx context.Context, hk *HoldKey) (bool, error) {
	keys := []string{hk.Key}
	args := []interface{}{hk.Value}

	res, err := releaseScript.Run(ctx, hk.RedisClient, keys, args...).Int()
	if err != nil {
		return false, err // 网络或脚本错误
	}

	return res >= 1, nil // res=1 或 res=0 但 key 不存在均视为成功
}
