// Package moooncrypto
// Wrote by yijian on 2024/01/25
package moooncrypto

import "testing"

// go test -v -run="TestAesOFBEncryptText"
func TestAesOFBEncryptText(t *testing.T) {
    key := "0123456789"
    data := "gomooon"

    if testAesOFBEncryptText(t, key, data) {
        key = "0123456789123456"
        if testAesOFBEncryptText(t, key, data) {
            key = "0123456789123456aa"
            data = "hello, welcome"
            if testAesOFBEncryptText(t, key, data) {
                key := "0123456789012345678912345"
                _ = testAesOFBEncryptText(t, key, data)
            }
        }
    }
}
func testAesOFBEncryptText(t *testing.T, key, data string) bool {
    // 加密
    encodedCiphertext, err := AesOFBEncryptText(key, data)
    if err != nil {
        t.Errorf("AesOFBEncryptText error: %s\n", err.Error())
        return false
    } else {
        t.Logf("AesOFBEncryptText ok: %s\n", encodedCiphertext)

        // 解密
        plaintext, err := AesOFBDecryptText(key, encodedCiphertext)
        if err != nil {
            t.Errorf("AesOFBDecryptText error: %s\n", err.Error())
            return false
        } else {
            t.Logf("AesOFBDecryptText ok: %s\n", plaintext)
            return true
        }
    }
}
