// Package mooonstr
// Wrote by yijian on 2024/09/03
package mooonstr

import (
	"strconv"
	"strings"
)

// JoinInt32 使用分隔符将多个 int32 串拼接成一个字符串
func JoinInt32(elems []int32, sep string) string {
	// 将 []int32 转换为 []string
	strElems := make([]string, len(elems))
	for i, num := range elems {
		strElems[i] = strconv.FormatInt(int64(num), 10)
	}

	// 使用 strings.Join 函数拼接字符串
	return strings.Join(strElems, sep)
}

// JoinUint32 使用分隔符将多个 uint32 串拼接成一个字符串
func JoinUint32(elems []uint32, sep string) string {
	// 将 []uint32 转换为 []string
	strElems := make([]string, len(elems))
	for i, num := range elems {
		strElems[i] = strconv.FormatInt(int64(num), 10)
	}

	// 使用 strings.Join 函数拼接字符串
	return strings.Join(strElems, sep)
}

// JoinInt64 使用分隔符将多个 int64 串拼接成一个字符串
func JoinInt64(elems []int64, sep string) string {
	// 将 []int64 转换为 []string
	strElems := make([]string, len(elems))
	for i, num := range elems {
		strElems[i] = strconv.FormatInt(num, 10)
	}

	// 使用 strings.Join 函数拼接字符串
	return strings.Join(strElems, sep)
}

// JoinUint64 使用分隔符将多个 uint64 串拼接成一个字符串
func JoinUint64(elems []uint64, sep string) string {
	// 将 []uint64 转换为 []string
	strElems := make([]string, len(elems))
	for i, num := range elems {
		strElems[i] = strconv.FormatUint(num, 10)
	}

	// 使用 strings.Join 函数拼接字符串
	return strings.Join(strElems, sep)
}