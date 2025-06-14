// Package mooonutils
// Wrote by yijian on 2025/06/14
package mooonutils

import (
	"os"
	"path/filepath"
	"strings"
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

// GetFilesBySuffix 获取指定目录下指定后缀的文件，结构存放到字符串数组中
func GetFilesBySuffix(dirPath string, suffixes []string) ([]string, error) {
	var fileList []string

	// 检查目录是否存在
	exists, isDir, err := PathExists(dirPath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &os.PathError{}
	}
	if !isDir {
		return nil, &os.PathError{}
	}

	// 遍历目录
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查文件后缀
		ext := filepath.Ext(info.Name())
		for _, suffix := range suffixes {
			if strings.EqualFold(ext, suffix) {
				fileList = append(fileList, path)
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}
