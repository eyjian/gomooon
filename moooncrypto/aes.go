// Package moooncrypto
// Wrote by yijian on 2024/01/25
package moooncrypto

import (
    "bytes"
    "fmt"
    "strings"
)

// 在 AES 加密时，如果明文的长度不是 16 字节的整数倍，就需要使用 PKCS7 填充算法对明文进行填充，
// 使其长度满足 16 字节的整数倍。在解密时，也需要使用相同的 PKCS7 填充算法对密文进行解密，以还原原始明文。
// 综合考虑，如果对安全性要求较高，建议采用 CTR 模式或 OFB 模式。如果对性能要求较高，可以考虑采用 ECB 模式或 CFB 模式。
// 但是，需要注意的是 ECB 模式和 CFB 模式容易受到密码本重放攻击，因此需要根据具体应用场景进行选择。

/*
IV（Initialization Vector）的作用是为加密过程提供一个随机数，以增加加密数据的随机性和不可预测性。
如果密钥没有泄漏，但 IV 泄漏了，加密数据的安全性会受到一定程度的威胁。
AES 加密算法是一个对称加密算法，它使用相同的密钥进行加密和解密。在加密过程中，IV 可以防止相同的明文块产生相同的密文块，从而增加加密数据的安全性。
IV 的长度应该与AES的块大小相同，通常为 128 位。在实际应用中，IV应该是随机生成的，并且在加密和解密过程中都应该使用相同的IV。如果 IV 不同，加密和解密过程将会产生不同的结果。
因为 IV 的作用是增加加密数据的随机性和不可预测性，如果 IV 泄漏了，攻击者就可以利用这个信息来尝试破解加密数据。
然而，如果密钥是安全的，那么攻击者仍然无法完全破解加密数据。因为 AES 加密算法是一个对称加密算法，加密和解密使用的密钥是相同的。
即使攻击者知道了 IV，他仍然无法破解加密数据，除非他同时知道密钥。
总之，如果密钥没有泄漏，但 IV 泄漏了，加密数据的安全性会受到一定程度的威胁。
为了提高加密数据的安全性和可靠性，建议使用不同的 IV，并确保密钥的安全。
*/

// pkcs7Padding 填充明文
func pkcs7Padding(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padtext...)
}

// pkcs7UnPadding 去除填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
    length := len(data)
    unpadding := int(data[length-1])
    if unpadding > length {
        return nil, fmt.Errorf("invalid pkcs7 padding")
    }
    return data[:(length - unpadding)], nil
}

func padToLength(str string) string {
    length := len(str)
    if length < 16 {
        return str + strings.Repeat("0", 16-length)
    } else if length < 24 {
        return str + strings.Repeat("0", 24-length)
    } else if length < 32 {
        return str + strings.Repeat("0", 32-length)
    } else {
        return str
    }
}
