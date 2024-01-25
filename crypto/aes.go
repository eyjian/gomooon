// Package crypto
// Wrote by yijian on 2024/01/25
package crypto

import (
	"bytes"
	"fmt"
	"strings"
)

// 在 AES 加密时，如果明文的长度不是 16 字节的整数倍，就需要使用 PKCS7 填充算法对明文进行填充，
// 使其长度满足 16 字节的整数倍。在解密时，也需要使用相同的 PKCS7 填充算法对密文进行解密，以还原原始明文。

// pkcs7Padding 填充明文
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7UnPadding 去除填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:(length - unpadding)], nil
}

func padToLength(str string) string {
	length := len(str)
	if length < 16 {
		return str + strings.Repeat("0", 16-length)
	} else if length < 24 {
		return str + strings.Repeat("0", 24-length)
	} else if length < 32 {
		return str + strings.Repeat("0", 32-length)
	} else {
		return str
	}
}
