// Package moooncrypto
// Wrote by yijian on 2024/08/22
package moooncrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/eyjian/gomooon/mooonstr"
	"golang.org/x/crypto/ssh"
)

type KeyType int
type keySize int

// 私钥类型
const (
	RSAPrivateKey KeyType = iota
	PKCS8PrivateKey
	ECPrivateKey
	OpenSSHPrivateKey
)

// 私钥长度，
// 可取值 1024/2048/3072/4096，
// 其中 1024 因安全性低不推荐使用，
// 而 4086 因性能过低也不推荐使用
const (
	RSAKey1024 keySize = 1024
	RSAKey2048         = 2048
	RSAKey3072         = 3072
	RSAKey4096         = 4096
)

func GeneratePrivateKeyFile(keyType KeyType, keySize keySize, filepath string) error {
	privateKeyString, err := GeneratePrivateKeyString(keyType, keySize)
	if err != nil {
		return err
	}
	return mooonstr.WriteString2File(filepath, privateKeyString)
}

// GeneratePrivateKeyString 生成私钥字符串
// keyType: 私钥类型
// keySize: 私钥长度，可取值 1024/2048/3072/4096，其中 1024 因安全性低不推荐使用，而 4086 因性能过低也不推荐使用
// 错误“too few primes of given length to generate an RSA key”表示 keySize 参数太小，导致无法生成有效的 RSA 密钥，
// RSA 密钥生成需要至少两个大素数，这些素数的位数之和应等于 keySize，当 keySize 太小时可能无法找到足够的大素数来生成有效的密钥。
func GeneratePrivateKeyString(keyType KeyType, keySize keySize) (string, error) {
	var privateKey interface{}
	var err error

	switch keyType {
	case RSAPrivateKey, PKCS8PrivateKey:
		// 生成 RSA 私钥
		privateKey, err = rsa.GenerateKey(rand.Reader, int(keySize))
	case ECPrivateKey, OpenSSHPrivateKey:
		// 生成 ECDSA 私钥
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	default:
		return "", errors.New("unsupported key type")
	}
	if err != nil {
		return "", fmt.Errorf("generate key error: %s", err.Error())
	}

	// 将私钥编码为相应的格式
	var privateKeyBytes []byte
	var pemType string

	switch keyType {
	case RSAPrivateKey:
		pemType = "RSA PRIVATE KEY"
		privateKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey))
	case PKCS8PrivateKey:
		pemType = "PRIVATE KEY"
		privateKeyBytes, err = x509.MarshalPKCS8PrivateKey(privateKey)
	case ECPrivateKey:
		pemType = "EC PRIVATE KEY"
		privateKeyBytes, err = x509.MarshalECPrivateKey(privateKey.(*ecdsa.PrivateKey))
	case OpenSSHPrivateKey:
		pemType = "OPENSSH PRIVATE KEY"
		signer, err := ssh.NewSignerFromKey(privateKey)
		if err != nil {
			return "", err
		}
		privateKeyBytes = ssh.MarshalAuthorizedKey(signer.PublicKey())
	}
	if err != nil {
		return "", fmt.Errorf("marshal key error: %s", err.Error())
	}

	// 创建PEM数据结构
	privateKeyBlock := &pem.Block{
		Type:  pemType,
		Bytes: privateKeyBytes,
	}
	return string(pem.EncodeToMemory(privateKeyBlock)), nil
}