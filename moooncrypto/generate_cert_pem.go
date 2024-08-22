// Package moooncrypto
// Wrote by yijian on 2024/08/22
package moooncrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

// CertTemplate 证书模板
type CertTemplate struct {
	SerialNumber          *big.Int      // 证书的唯一标识符，通常是一个大整数，如：big.NewInt(1234567890)
	Subject               CertSubject   // 证书主题
	NotBefore             time.Time     // time.Now()
	NotAfter              time.Time     // time.Now().Add(time.Hour * 24 * 365)
	KeyUsage              x509.KeyUsage // x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	BasicConstraintsValid bool          // 是否启用基本约束
	IsCA                  bool          // 是否为 CA 证书

	// 在生成 PEM 格式的证书时，Type 字段通常设置为 CERTIFICATE。
	// 这是因为 X.509 证书是用于公钥加密的最常见和广泛使用的证书格式。
	// 在大多数情况下，您应该使用 CERTIFICATE 类型来表示 PEM 编码的 X.509 证书。
	Type string // 证书类型，如："CERTIFICATE"
}

// CERTIFICATE：用于表示 X.509 证书
// X509 CRL：用于表示X.509证书吊销列表（CRL），CRL 是一个由证书颁发机构（CA）维护的列表，用于记录已经吊销的证书
// ATTRIBUTE CERTIFICATE：用于表示属性证书。属性证书是一种扩展 X.509 证书的方法，用于包含与证书持有者相关的其他属性信息
// CERTIFICATE PAIR：用于表示证书对，证书对通常包含一个签名证书和一个加密证书，分别用于签名和加密操作。
//
// CERTIFICATE REQUEST：用于表示证书签名请求（CSR）
// RSA PRIVATE KEY：用于表示 PKCS#1 格式的 RSA 私钥
// PRIVATE KEY：用于表示 PKCS#8 格式的私钥，可以是 RSA、ECDSA 或其他类型的私钥
// EC PRIVATE KEY：用于表示 ECDSA 私钥
// PUBLIC KEY：用于表示公钥，可以是 RSA、ECDSA 或其他类型的公钥
// PKCS7：用于表示 PKCS#7 格式的数据，通常用于封装签名或加密的消息
// SSH2 PUBLIC KEY：用于表示 OpenSSH 格式的公钥

// CertSubject 证书主题
type CertSubject struct {
	Organization []string // 示例："Example Corp."
	CommonName   string   // 示例："example.com"
}

// GenerateCertPemStringFromPrivateKeyFilepath 从私钥文件生成证书 PEM 字符串
func GenerateCertPemStringFromPrivateKeyFilepath(privateKeyFilepath string, ct *CertTemplate) (string, error) {
	privateKey, err := Filepath2PrivateKey(privateKeyFilepath)
	if err != nil {
		return "", err
	}
	return GenerateCertPemStringFromPrivateKey(privateKey, ct)
}

// GenerateCertPemStringFromPrivateKeyFile 从私钥文件生成证书 PEM 字符串
func GenerateCertPemStringFromPrivateKeyFile(privateKeyFile *os.File, ct *CertTemplate) (string, error) {
	privateKey, err := File2PrivateKey(privateKeyFile)
	if err != nil {
		return "", err
	}
	return GenerateCertPemStringFromPrivateKey(privateKey, ct)
}

// GenerateCertPemStringFromPrivateKeyString 从私钥字符串生成证书 PEM 字符串
func GenerateCertPemStringFromPrivateKeyString(privateKeyString string, ct *CertTemplate) (string, error) {
	privateKey, err := String2PrivateKey(privateKeyString)
	if err != nil {
		return "", err
	}
	return GenerateCertPemStringFromPrivateKey(privateKey, ct)
}

// GenerateCertPemStringFromPrivateKey 从私钥生成证书 PEM 字符串
func GenerateCertPemStringFromPrivateKey(privateKey *rsa.PrivateKey, ct *CertTemplate) (string, error) {
	// 创建证书模板
	certTemplate := &x509.Certificate{
		SerialNumber: ct.SerialNumber,
		Subject: pkix.Name{
			Organization: ct.Subject.Organization,
			CommonName:   ct.Subject.CommonName,
		},
		NotBefore:             ct.NotBefore,
		NotAfter:              ct.NotAfter,
		KeyUsage:              ct.KeyUsage,
		BasicConstraintsValid: ct.BasicConstraintsValid,
		IsCA:                  ct.IsCA,
	}

	// 生成证书
	certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", err
	}

	// 将证书编码为 PEM 格式
	certBlock := &pem.Block{
		Type:  ct.Type,
		Bytes: certBytes,
	}
	certPem := pem.EncodeToMemory(certBlock)
	return string(certPem), nil
}