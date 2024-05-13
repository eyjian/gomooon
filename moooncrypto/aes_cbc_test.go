// Package moooncrypto
// Wrote by yijian on 2024/01/25
package moooncrypto

import "testing"

// go test -v -run="TestAesCBCEncryptText"
func TestAesCBCEncryptText(t *testing.T) {
    key := "0123456789"
    data := "gomooon"

    if testAesCBCEncryptText(t, key, data) {
        key = "0123456789123456"
        if testAesCBCEncryptText(t, key, data) {
            key = "0123456789123456aa"
            data = "hello, welcome"
            if testAesCBCEncryptText(t, key, data) {
                key := "0123456789012345678912345"
                _ = testAesCBCEncryptText(t, key, data)
            }
        }
    }
}
func testAesCBCEncryptText(t *testing.T, key, data string) bool {
    // 加密
    encodedCiphertext, err := AesCBCEncryptText(key, data)
    if err != nil {
        t.Errorf("AesCBCEncryptText error: %s\n", err.Error())
        return false
    } else {
        t.Logf("AesCBCEncryptText ok: %s\n", encodedCiphertext)

        // 解密
        plaintext, err := AesCBCDecryptText(key, encodedCiphertext)
        if err != nil {
            t.Errorf("AesCBCDecryptText error: %s\n", err.Error())
            return false
        } else {
            t.Logf("AesCBCDecryptText ok: %s\n", plaintext)
            return true
        }
    }
}
