// Package mooonutils
// Wrote by yijian on 2024/10/22
package mooonutils

import (
	"archive/zip"
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path/filepath"
)

// Unzip 解压 zip 文件
// 返回值：解压后的文件（含目录部分，如果是当前目录“.”则仅文件名）列表
// overwrite：是否覆盖解压后的同名文件
// destDir：解压后文件的存放目录，如果不存在会自动创建
func Unzip(zipFile, destDir string, overwrite ...bool) ([]string, error) {
	// 设置默认值
	truncated := true
	if len(overwrite) > 0 {
		truncated = overwrite[0]
	}

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, err
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
		paths = append(paths, path)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return nil, err
			}

			flag := os.O_WRONLY | os.O_CREATE
			if truncated {
				flag = flag | os.O_TRUNC
			}
			outFile, err := os.OpenFile(path, flag, file.Mode())
			if err != nil {
				return nil, err
			}

			rc, err := file.Open()
			if err != nil {
				outFile.Close()
				return nil, err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return nil, err
			}
		}
	}

	return paths, nil
}