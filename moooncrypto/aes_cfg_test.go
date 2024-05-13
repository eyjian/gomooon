// Package moooncrypto
// Wrote by yijian on 2024/01/25
package moooncrypto

import "testing"

// go test -v -run="TestAesCFBEncryptText"
func TestAesCFBEncryptText(t *testing.T) {
    key := "0123456789"
    data := "gomooon"

    if testAesCFBEncryptText(t, key, data) {
        key = "0123456789123456"
        if testAesCFBEncryptText(t, key, data) {
            key = "0123456789123456aa"
            data = "hello, welcome"
            if testAesCFBEncryptText(t, key, data) {
                key := "0123456789012345678912345"
                _ = testAesCFBEncryptText(t, key, data)
            }
        }
    }
}
func testAesCFBEncryptText(t *testing.T, key, data string) bool {
    // 加密
    encodedCiphertext, err := AesCFBEncryptText(key, data)
    if err != nil {
        t.Errorf("AesCFBEncryptText error: %s\n", err.Error())
        return false
    } else {
        t.Logf("AesCFBEncryptText ok: %s\n", encodedCiphertext)

        // 解密
        plaintext, err := AesCFBDecryptText(key, encodedCiphertext)
        if err != nil {
            t.Errorf("AesCFBDecryptText error: %s\n", err.Error())
            return false
        } else {
            t.Logf("AesCFBDecryptText ok: %s\n", plaintext)
            return true
        }
    }
}
