// Package mooonutils
// Wrote by yijian on 2025/06/14
package mooonutils

import (
	"errors"
	"fmt"
	"os"
)

// IsDir 判断路径是否为目录
func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 目录不存在
		}
		return false, err // 其他错误
	}
	return info.IsDir(), nil // 判断是否为目录
}

// PathExists 判断路径是否存在（文件或目录）
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // 路径存在
	}
	if os.IsNotExist(err) {
		return false, nil // 路径不存在
	}
	return false, err // 其他错误（如权限不足）
}

// DirExists 判断目录是否存在
func DirExists(path string) (bool, error) {
	isDir, err := IsDir(path)
	if err != nil {
		return false, err
	}
	if !isDir {
		return false, errors.New(fmt.Sprintf("`%s` is not a directory", path))
	}
	return PathExists(path)
}
