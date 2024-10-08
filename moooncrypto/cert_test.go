// Package moooncrypto
// Wrote by yijian on 2024/05/09
package moooncrypto

import (
    "crypto/x509"
    "math/big"
    "os"
    "testing"
    "time"
)

// go test -v -run="TestGetCertInfo$"
func TestGetCertInfo(t *testing.T) {
    // 读取证书文件
    certPEM, err := os.ReadFile("certificate.pem")
    if err != nil {
        t.Errorf("Failed to read certificate file: %v", err)
    } else {
        // 获取证书信息
        certInfo, err := GetCertInfo(string(certPEM))
        if err != nil {
            t.Errorf("Failed to get certificate info: %v", err)
        } else {
            t.Logf("Subject: %s, No: %s, Start: %s, Stop: %s\n", certInfo.Subject, certInfo.No16, certInfo.StartTime, certInfo.StopTime)
        }
    }
}

// go test -v -run="TestExtractCertAndKeyFromP12$"
func TestExtractCertAndKeyFromP12(t *testing.T) {
    // 读取 P12 二进制文件
    p12Data, err := os.ReadFile("certificate.p12")
    if err != nil {
        t.Errorf("Failed to read p12 file: %s\n", err.Error())
    } else {
        password := "" // 测试的没有设置密码
        cert, key, err := ExtractCertAndKeyFromP12(p12Data, password)
        if err != nil {
            t.Errorf("Read p12 file error: %s\n", err.Error())
        } else {
            t.Logf("Cert:\n%s\n", cert)
            t.Logf("Key:\n%s\n", key)
        }
    }
}

// go test -v -run="TestIsPemCertificate$"
func TestIsPemCertificate(t *testing.T) {
    cert, err := os.ReadFile("self_signed_cert.pem")
    if err != nil {
        t.Errorf("Failed to read pem file: %s\n", err.Error())
    } else {
        if !IsPemCertificate(string(cert)) {
            t.Errorf("is not cert file\n")
        } else {
            t.Logf("is cert\n")
        }
    }
}

// go test -v -run="TestIsP8PemPrivateKey$"
func TestIsP8PemPrivateKey(t *testing.T) {
    cert, err := os.ReadFile("private_key_pkcs8.pem")
    if err != nil {
        t.Errorf("Failed to read private key file: %s\n", err.Error())
    } else {
        if !IsP8PemPrivateKey(string(cert)) {
            t.Errorf("is not PKCS#8 private key\n")
        } else {
            t.Logf("is cert\n")
        }
    }
}

// go test -v -run="TestIsP1PemPrivateKey$"
func TestIsP1PemPrivateKey(t *testing.T) {
    cert, err := os.ReadFile("private_key.pem")
    if err != nil {
        t.Errorf("Failed to read private key file: %s\n", err.Error())
    } else {
        if !IsP1PemPrivateKey(string(cert)) {
            t.Errorf("is not PKCS#1 private key\n")
        } else {
            t.Logf("is cert\n")
        }
    }
}

// go test -v -run="TestString2PrivateKey$"
func TestString2PrivateKey(t *testing.T) {
    privateKeyStr, err := GeneratePrivateKeyString(RSAPrivateKey, RSAKey2048)
    if err != nil {
        t.Errorf("Failed to generate private key: %s\n", err.Error())
        return
    } else {
        t.Logf("private key: %s\n", privateKeyStr)
    }

    privateKey, err := String2PrivateKey(privateKeyStr)
    if err != nil {
        t.Errorf("Failed to read private key file: %s\n", err.Error())
    } else {
        t.Logf("private key: ok\n")
    }

    certPemString, err := GenerateCertPemStringFromPrivateKey(privateKey,
        &CertTemplate{
            SerialNumber: big.NewInt(1234567890),
            Subject: CertSubject{
                CommonName:   "localhost",
                Organization: []string{"localhost"},
            },
            NotBefore:             time.Now(),
            NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10),
            KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
            IsCA:                  false,
            BasicConstraintsValid: true,
        })
    if err != nil {
        t.Errorf("Failed to generate cert: %s\n", err.Error())
        return
    } else {
        t.Logf("cert: %s\n", certPemString)
    }
}