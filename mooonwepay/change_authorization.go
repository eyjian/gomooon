// Package mooonwepay
// Wrote by yijian on 2024/08/23
package mooonwepay

import "fmt"

func makeChangeBillAuthorization(mchid, serialNo, nonceStr, signature string, timestamp int64) string {
	return fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",signature="%s",timestamp="%d",serial_no="%s"`,
		mchid, nonceStr, signature, timestamp, serialNo)
}