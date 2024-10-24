// Package mooonutils
// Wrote by yijian on 2024/10/22
package mooonutils

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractFilenameWithoutExtension 提取不含后缀的文件名
func ExtractFilenameWithoutExtension(path string) string {
	// 提取文件名（包括后缀）
	filenameWithExtension := filepath.Base(path)

	// 提取文件后缀
	extension := filepath.Ext(filenameWithExtension)

	// 去掉文件后缀
	filenameWithoutExtension := strings.TrimSuffix(filenameWithExtension, extension)

	return filenameWithoutExtension
}

// Md5File 计算文件的 md5
// 返回值：文件的 md5 小写值
func Md5File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file://%s error: %s", filePath, err.Error())
	}
	defer file.Close()

	// 创建一个MD5哈希对象
	hash := md5.New()

	// 将文件内容复制到哈希对象中
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("copy file://%s error: %s", filePath, err.Error())
	}

	//获取MD5哈希值的字节表示
	hashBytes := hash.Sum(nil)

	// 将字节表示转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

// Unzip 解压 zip 文件
// 返回值：解压后的文件（含目录部分，如果是当前目录“.”则仅文件名）列表
// options[0]：是否覆盖解压后的同名文件，默认是 true
// options[1]：返回结果是否忽略目录，仅包含文件，默认是 true
// destDir：解压后文件的存放目录，如果不存在会自动创建
func Unzip(zipFile, destDir string, options ...bool) ([]string, error) {
	// 设置默认值
	overwrite := true
	ignoreDir := true
	if len(options) > 0 {
		overwrite = options[0]
	}
	if len(options) > 1 {
		ignoreDir = options[1]
	}

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, fmt.Errorf("open zipfile://%s error: %s", zipFile, err.Error())
	}
	defer reader.Close()

	var paths []string
	var decodeName string
	for _, file := range reader.File {
		if file.Flags == 0 {
			// 解决中文文件名的乱码
			i := bytes.NewReader([]byte(file.Name))
			// GB18030 是 GBK 的扩展集，可以更好地处理中文字符
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := io.ReadAll(decoder)
			decodeName = string(content)
		} else {
			decodeName = file.Name
		}
		path := filepath.Join(destDir, decodeName)

		if file.FileInfo().IsDir() {
			if ignoreDir {
				continue
			}
			paths = append(paths, path)
			os.MkdirAll(path, os.ModePerm)
		} else {
			paths = append(paths, path)
			dirPath := filepath.Dir(path)
			if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return nil, fmt.Errorf("open dir://%s error: %s", dirPath, err.Error())
			}

			flag := os.O_WRONLY | os.O_CREATE
			if overwrite {
				flag = flag | os.O_TRUNC
			}
			outFile, err := os.OpenFile(path, flag, file.Mode())
			if err != nil {
				return nil, fmt.Errorf("open outfile://%s with flag://0x%x error: %s", path, flag, err.Error())
			}

			rc, err := file.Open()
			if err != nil {
				outFile.Close()
				os.Remove(path)
				return nil, fmt.Errorf("open infile://%s error: %s", file.Name, err.Error())
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					os.Remove(path)
				}
				return nil, fmt.Errorf("copy infile://%s to outfile://%s error: %s", file.Name, path, err.Error())
			}
		}
	}

	return paths, nil
}
