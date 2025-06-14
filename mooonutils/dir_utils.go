// Package mooonutils
// Wrote by yijian on 2025/06/14
package mooonutils

import (
	"os"
)

// PathExists 判断路径是否存在
// 返回：是否存在exists、是否为目录isDir、错误err
func PathExists(path string) (bool, bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil // 目录不存在
		}
		return false, false, err // 其他错误
	}
	return true, info.IsDir(), nil // 判断是否为目录
}
