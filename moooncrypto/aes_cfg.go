// Package moooncrypto
// Wrote by yijian on 2024/01/25
package moooncrypto

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/hex"
    "fmt"
)

// 与 ECB 和 CBC 模式只能够加密块数据不同，CFB 能够将块密文（Block Cipher）转换为流密文（Stream Cipher）。
// AES 要求 key 长度只能为 128 或 192 或 256 比特中的一种，即 16 字节或 24 字节或 32 字节中的一种。

// AesCFBEncryptText 加密文本
// key 加密密钥，长度不能超过 32 个字符，否则返回 error。如果不是 16、24 和 32 整的倍数，则会自动补 0 填充到最近的长度。
func AesCFBEncryptText(key, data string) (string, error) {
    keyLen := len(key)
    if keyLen > 256 {
        return "", fmt.Errorf("length of CFB encrypt key exceeds 256")
    }

    // 创建 AES 分组密码的实例
    keyBytes := []byte(padToLength(key))

    // 创建 AES 密钥
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", fmt.Errorf("new CFB encrypt cipher error: %s", err.Error())
    }

    // 初始化向量
    iv := keyBytes[:aes.BlockSize] // 使用密钥的前 16 个字节作为初始向量
    //iv := make([]byte, aes.BlockSize)
    //if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    //	return "", fmt.Errorf("CFB encrypt rand read full error: %s", err.Error())
    //}

    // 创建 CFB 模式的加密器
    stream := cipher.NewCFBEncrypter(block, iv)

    // 加密明文
    plaintext := []byte(data)
    ciphertext := make([]byte, len(plaintext))
    stream.XORKeyStream(ciphertext, plaintext)

    return hex.EncodeToString(ciphertext), nil
}

// AesCFBDecryptText 解密文本
// key 解密密钥，值同加密密钥，长度不能超过 32 个字符，否则返回 error。如果不是 16、24 和 32 整的倍数，则会自动补 0 填充到最近的长度。
func AesCFBDecryptText(key, data string) (string, error) {
    keyLen := len(key)
    if keyLen > 256 {
        return "", fmt.Errorf("length of CFB decrypt key exceeds 256")
    }

    // 创建 AES 分组密码的实例
    keyBytes := []byte(padToLength(key))

    // 创建 AES 密钥
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", fmt.Errorf("new CFB decrypt cipher error: %s", err.Error())
    }

    // 初始化向量
    iv := keyBytes[:aes.BlockSize] // 使用密钥的前 16 个字节作为初始向量
    //iv := make([]byte, aes.BlockSize)
    //if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    //	return "", fmt.Errorf("CFB decrypt rand read full error: %s", err.Error())
    //}

    // 创建 CFB 模式的解密器
    stream := cipher.NewCFBDecrypter(block, iv)

    // 解密密文
    ciphertext, _ := hex.DecodeString(data)
    plaintext := make([]byte, len(ciphertext))
    stream.XORKeyStream(plaintext, ciphertext)

    return string(plaintext), nil
}
