// Package gomooon
// Wrote by yijian on 2024/01/02
package gomooon

import (
	"strconv"
	"time"
)

// GetCurrentTimestamp 获取指定时区的当前时间戳
func GetCurrentTimestamp(timezone string) int64 {
	// 加载时区信息（北京时区：Asia/Shanghai）
	location, _ := time.LoadLocation(timezone)

	// 获取当前时间，并转换为北京时区的时间
	now := time.Now().In(location)

	// 将北京时区的时间转换为时间戳（精确到秒）
	timestamp := now.Unix()

	// 返回时间戳
	return timestamp
}

func GetCurrentTimestampString(timezone string) string {
	currentTimestamp := GetCurrentTimestamp(timezone)
	return strconv.FormatInt(currentTimestamp, 10)
}
