// Package moooncrypto
// Wrote by yijian on 2024/01/25
package moooncrypto

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/hex"
    "fmt"
)

// 在 OFB 模式下，先用块加密器生成密钥流（keystream），然后再将密钥流与明文流异或得到密文流，
// 解密是先用块加密器生成密钥流，再将密钥流与密文流异或得到明文，由于异或操作的对称性所以加密和解密的流程是完全一样的。

// AesOFBEncryptText 加密文本
// key 加密密钥，长度不能超过 32 个字符，否则返回 error。如果不是 16、24 和 32 整的倍数，则会自动补 0 填充到最近的长度。
func AesOFBEncryptText(key, data string) (string, error) {
    keyLen := len(key)
    if keyLen > 256 {
        return "", fmt.Errorf("length of OFB encrypt key exceeds 256")
    }

    // 创建 AES 分组密码的实例
    keyBytes := []byte(padToLength(key))

    // 创建 AES 密钥
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", fmt.Errorf("new OFB encrypt cipher error: %s", err.Error())
    }

    // 初始化向量
    iv := keyBytes[:aes.BlockSize] // 使用密钥的前 16 个字节作为初始向量
    //iv := make([]byte, aes.BlockSize)
    //if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    //	return "", fmt.Errorf("OFB encrypt rand read full error: %s", err.Error())
    //}

    // 创建 OFB 模式的加密器
    stream := cipher.NewOFB(block, iv)

    // 加密明文
    plaintext := []byte(data)
    ciphertext := make([]byte, len(plaintext))
    stream.XORKeyStream(ciphertext, plaintext)

    return hex.EncodeToString(ciphertext), nil
}

// AesOFBDecryptText 解密文本
// key 解密密钥，值同加密密钥，长度不能超过 32 个字符，否则返回 error。如果不是 16、24 和 32 整的倍数，则会自动补 0 填充到最近的长度。
func AesOFBDecryptText(key, data string) (string, error) {
    keyLen := len(key)
    if keyLen > 256 {
        return "", fmt.Errorf("length of OFB decrypt key exceeds 256")
    }

    // 创建 AES 分组密码的实例
    keyBytes := []byte(padToLength(key))

    // 创建 AES 密钥
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", fmt.Errorf("new OFB decrypt cipher error: %s", err.Error())
    }

    // 初始化向量
    iv := keyBytes[:aes.BlockSize] // 使用密钥的前 16 个字节作为初始向量
    //iv := make([]byte, aes.BlockSize)
    //if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    //	return "", fmt.Errorf("OFB decrypt rand read full error: %s", err.Error())
    //}

    // 创建 OFB 模式的解密器
    stream := cipher.NewOFB(block, iv)

    // 解密密文
    ciphertext, err := hex.DecodeString(data)
    if err != nil {
        return "", fmt.Errorf("OFB decrypt decode error: %s", err.Error())
    }

    plaintext := make([]byte, len(ciphertext))
    stream.XORKeyStream(plaintext, ciphertext)

    return string(plaintext), nil
}
