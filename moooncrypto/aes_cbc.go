// Package moooncrypto
// Wrote by yijian on 2024/01/25
package moooncrypto

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "errors"
    "fmt"
)

// 在 CBC 模式下，每一个明文分组先与前一个密文分组进行异或操作，再进行加密。
// 这种模式的加密速度适中，安全性较好，但需要额外的存储空间来保存初始向量。
// AES 要求 key 长度只能为 128 或 192 或 256 比特中的一种，即 16 字节或 24 字节或 32 字节中的一种。

// AesCBCEncryptText 加密文本
// key 加密密钥，长度不能超过 32 个字符，否则返回 error。如果不是 16、24 和 32 整的倍数，则会自动补 0 填充到最近的长度。
func AesCBCEncryptText(key, data string) (string, error) {
    keyLen := len(key)
    if keyLen > 256 {
        return "", fmt.Errorf("length of CBC encrypt key exceeds 256")
    }

    // 创建 AES 分组密码的实例
    keyBytes := []byte(padToLength(key))
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", fmt.Errorf("new CBC encrypt cipher error: %s", err.Error())
    }

    // 创建分组模式（这里使用 CBC 模式）
    iv := keyBytes[:aes.BlockSize] // 使用密钥的前 16 个字节作为初始向量
    mode := cipher.NewCBCEncrypter(block, iv)

    // 对明文进行填充
    paddedPlaintext := pkcs7Padding([]byte(data), aes.BlockSize)

    // 加密
    paddedPlaintextLen := len(paddedPlaintext)
    if paddedPlaintextLen%aes.BlockSize != 0 {
        return "", errors.New("plaintext format data")
    }
    ciphertext := make([]byte, paddedPlaintextLen)
    mode.CryptBlocks(ciphertext, paddedPlaintext)

    // 将加密后的结果转换为 Base64 编码
    encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
    return encodedCiphertext, nil
}

// AesCBCDecryptText 解密文本
// key 解密密钥，值同加密密钥，长度不能超过 32 个字符，否则返回 error。如果不是 16、24 和 32 整的倍数，则会自动补 0 填充到最近的长度。
func AesCBCDecryptText(key, data string) (string, error) {
    keyLen := len(key)
    if keyLen > 256 {
        return "", fmt.Errorf("length of CBC decrypt key exceeds 256")
    }

    // 创建 AES 分组密码的实例
    keyBytes := []byte(padToLength(key))
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", fmt.Errorf("new CBC decrypt cipher error: %s", err.Error())
    }

    // 创建分组模式（这里使用 CBC 模式）
    iv := keyBytes[:aes.BlockSize] // 使用密钥的前 16 个字节作为初始向量
    mode := cipher.NewCBCDecrypter(block, iv)

    // 将加密后的密文转换为字节数组
    ciphertext, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        return "", fmt.Errorf("CBC decrypt base64 decode error: %s", err.Error())
    }

    // 解密
    ciphertextLen := len(ciphertext)
    if ciphertextLen%aes.BlockSize != 0 {
        return "", errors.New("ciphertext format data")
    }
    paddedPlaintext := make([]byte, ciphertextLen)
    mode.CryptBlocks(paddedPlaintext, ciphertext)

    // 去除填充
    plaintext, err := pkcs7UnPadding(paddedPlaintext)
    if err != nil {
        return "", fmt.Errorf("CBC decrypt %s", err.Error())
    }

    return string(plaintext), nil
}
