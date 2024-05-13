// Package moooncrypto
// Wrote by yijian on 2024/05/09
package moooncrypto

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "golang.org/x/crypto/pkcs12"
    "time"
)

type CertInfo struct {
    Ver                int       `json:"ver"`                  // 证书版本
    No                 string    `json:"cert_no"`              // 证书序列号
    Subject            string    `json:"subject"`              // 证书主题
    StartTime          time.Time `json:"start_time"`           // 证书开始时间
    StopTime           time.Time `json:"stop_time"`            // 证书结束时间
    Issuer             string    `json:"issuer"`               // 证书颁发者
    PublicKeyAlgorithm string    `json:"public_key_algorithm"` // 公钥算法
    SignatureAlgorithm string    `json:"signature_algorithm"`  // 签名算法
    PublicKey          string    `json:"public_key"`           // 公钥
    Signature          string    `json:"signature"`            // 签名
    KeyUsage           int       `json:"key_usage"`            // 密钥用途
    SubjectKeyId       string    `json:"subject_key_id"`       // 主题密钥标识
    AuthorityKeyId     string    `json:"authority_key_id"`     // 颁发者密钥标识
}

func GetCertInfo(certPEM string) (*CertInfo, error) {
    // 解码 PEM 格式的证书
    block, _ := pem.Decode([]byte(certPEM))
    if block == nil {
        return nil, fmt.Errorf("failed to decode PEM block")
    }

    // 解析 X.509 证书
    cert, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse X.509 certificate: %s", err.Error())
    }

    return &CertInfo{
        Ver:                cert.Version,
        No:                 cert.SerialNumber.String(),
        Subject:            cert.Subject.CommonName,
        StartTime:          cert.NotBefore,
        StopTime:           cert.NotAfter,
        Issuer:             cert.Issuer.CommonName,
        PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
        SignatureAlgorithm: cert.SignatureAlgorithm.String(),
        PublicKey:          string(cert.PublicKey.(*rsa.PublicKey).N.Bytes()),
        Signature:          string(cert.Signature),
        KeyUsage:           int(cert.KeyUsage),
        SubjectKeyId:       string(cert.SubjectKeyId),
        AuthorityKeyId:     string(cert.AuthorityKeyId),
    }, nil
}

// ExtractCertAndKeyFromP12 P12 文件是一种用于存储和传输用户或服务器私钥、公钥和证书的二进制格式文件，也称为 PFX 文件。
// 它遵循 Public Key Cryptography Standards #12（PKCS#12）标准，该标准为这些密钥和证书提供了一个可移植的格式。
// P12 文件通常包含开发者的公钥和私钥，以及一个证书链，用于验证开发者的身份。
//
// PEM（Privacy-Enhanced Mail）证书是一种使用 Base64 ASCII 编码的纯文本格式，通常具有 .crt 或 .pem 扩展名。
// PEM 证书包含证书主体的公开信息、公钥信息以及签署证书的证书颁发机构（CA）的信息。
// PEM 证书主要用于存储和传输证书，例如在 HTTPS 服务器上使用的 SSL/TLS 证书。
//
// 可从 P12 文件提取出证书和私钥，也可将证书和私钥打包为 P12 文件：
// 1）openssl genrsa -out private_key.pem 2048 # 生成一个 PKCS#1 格式的 2048 位 RSA 私钥
// 2）openssl pkcs8 -topk8 -inform PEM -outform PEM -in private_key.pem -out private_key_pkcs8.pem -nocrypt # 将 PKCS#1 格式的私钥转换为 PKCS#8 格式
// 3）openssl req -new -key private_key.pem -out cert_request.csr # 生成证书签名请求（CSR）
// 4）openssl x509 -req -days 365 -in cert_request.csr -signkey private_key.pem -out self_signed_cert.pem # 生成自签名证书
// 5）openssl pkcs12 -export -in self_signed_cert.pem -inkey private_key.pem -out certificate.p12 # 生成 P12 文件
//
// 返回值分别为：PEM 证书、PEM 私钥和 error
// 参数 password 为解密 P12 数据 p12Data 的密码
func ExtractCertAndKeyFromP12(p12Data []byte, password string) (string, string, error) {
    // 解码 p12 文件
    privateKeyInf, cert, err := pkcs12.Decode(p12Data, password)
    if err != nil {
        return "", "", fmt.Errorf("failed to decode p12: %s", err.Error())
    }

    // 将私钥从 interface{} 类型转换为 *rsa.PrivateKey 类型
    privateKey, ok := privateKeyInf.(*rsa.PrivateKey)
    if !ok {
        return "", "", fmt.Errorf("private key type assertion failed")
    }

    // 将证书转换为 PEM 格式（含公钥、证书序列号等）
    certPemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})

    // 将私钥转换为 PEM 格式
    keyPemBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

    return string(certPemBytes), string(keyPemBytes), nil
}

// IsPemCertificate 判断字符串是否为 PEM 格式的 X.509 证书
func IsPemCertificate(s string) bool {
    block, _ := pem.Decode([]byte(s))
    if block == nil {
        return false
    }
    return block.Type == "CERTIFICATE"
}

// IsP7PemCertificate 判断字符串是否为 PKCS#7 格式的 PEM 格式的证书
func IsP7PemCertificate(s string) bool {
    block, _ := pem.Decode([]byte(s))
    if block == nil {
        return false
    }
    return block.Type == "PKCS7"
}

// IsP8PemPrivateKey 判断字符串是否为 PKCS#8 格式的 PEM 格式的私钥
func IsP8PemPrivateKey(s string) bool {
    block, _ := pem.Decode([]byte(s))
    if block == nil {
        return false
    }
    return block.Type == "PRIVATE KEY"
}

// IsP1PemPrivateKey 判断字符串是否为 PKCS#1 格式的 PEM 格式的私钥
func IsP1PemPrivateKey(s string) bool {
    block, _ := pem.Decode([]byte(s))
    if block == nil {
        return false
    }
    return block.Type == "RSA PRIVATE KEY"
}

// IsEcPemPrivateKey 判断字符串是否为 ECDSA 格式的 PEM 格式的私钥
func IsEcPemPrivateKey(s string) bool {
    block, _ := pem.Decode([]byte(s))
    if block == nil {
        return false
    }
    return block.Type == "EC PRIVATE KEY"
}

// IsOpenSslPemPrivateKey 判断字符串是否为 OpenSSL 格式的 PEM 格式的私钥
func IsOpenSslPemPrivateKey(s string) bool {
    block, _ := pem.Decode([]byte(s))
    if block == nil {
        return false
    }
    return block.Type == "OPENSSL PRIVATE KEY"
}
