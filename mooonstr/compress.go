// Package mooonstr
// Wrote by yijian on 2025/06/20
package mooonstr

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
)

func CompressChinese(s string) (string, error) {
	// UTF-8 转 GBK（减少中文字符体积）
	gbkBytes, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(s))
	if err != nil {
		return "", fmt.Errorf("GBK编码失败: %w", err)
	}

	// GZIP 压缩
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(gbkBytes); err != nil {
		return "", fmt.Errorf("GZIP写入失败: %w", err)
	}
	if err := gz.Close(); err != nil { // 必须关闭写入器确保数据刷新
		return "", fmt.Errorf("GZIP关闭失败: %w", err)
	}

	// Base64 编码（避免二进制乱码）
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecompressChinese(compressed string) (string, error) {
	// Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(compressed)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %w", err)
	}

	// GZIP 解压
	gzReader, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return "", fmt.Errorf("GZIP读取器创建失败: %w", err)
	}
	defer gzReader.Close()

	var gbkBuf bytes.Buffer
	if _, err := io.Copy(&gbkBuf, gzReader); err != nil {
		return "", fmt.Errorf("GZIP解压失败: %w", err)
	}

	// GBK 转 UTF-8
	utf8Bytes, err := simplifiedchinese.GBK.NewDecoder().Bytes(gbkBuf.Bytes())
	if err != nil {
		return "", fmt.Errorf("GBK转UTF-8失败: %w", err)
	}

	return string(utf8Bytes), nil
}
