// Package mooonutils
// Wrote by yijian on 2024/10/25
package mooonutils

import (
    "encoding/binary"
    "io"
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

// IsExecutableFile 判断 filepath 是否为一个可执行的程序文件
func IsExecutableFile(filepath string) (bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// 定义要读取的字节数
	const headerSize = 8

	// 创建一个缓冲区来存储文件头
	header := make([]byte, headerSize)

	// 读取文件头
	_, err = io.ReadFull(file, header)
	if err != nil {
		return false, err
	}

	// 检查ELF文件
	if IsELF(header) {
		return true, nil
	}

	// 检查PE文件（Windows可执行文件）
	if IsPE(header) {
		return true, nil
	}

	// 检查Mach-O文件（macOS可执行文件）
	if IsMachO(header) {
		return true, nil
	}

	return false, nil
}

func IsELF(header []byte) bool {
	return len(header) >= 4 && string(header[:4]) == "\x7fELF"
}

func IsPE(header []byte) bool {
	return len(header) >= 2 && string(header[:2]) == "MZ"
}

func IsMachO(header []byte) bool {
	return len(header) >= 4 && (binary.LittleEndian.Uint32(header[:4]) == 0xfeedface ||
		binary.BigEndian.Uint32(header[:4]) == 0xfeedface)
}