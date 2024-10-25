// Package mooonutils
// Wrote by yijian on 2024/10/25
package mooonutils

import (
	"os"
	"path/filepath"
)

// GetProgramDir 取得程序文件所在目录
func GetProgramDir() (string, error) {
	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// 获取可执行文件所在的目录
	dir := filepath.Dir(exePath)

	return dir, nil
}
