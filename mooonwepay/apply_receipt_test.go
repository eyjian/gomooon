// Package mooonwepay
// Wrote by yijian on 2025/03/27
package mooonwepay

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"os"
	"testing"
)

// go test -v -run="TestApplyReceipt" -args mchid serial_no mch_api_v3_key private_key wepay_bill_no
func TestApplyReceipt(t *testing.T) {
	args := os.Args[5:]
	if len(args) != 5 {
		t.Fatal("args num error")
	}
	mchID := args[0]
	mchCertificateSerialNumber := args[1]
	mchAPIv3Key := args[2]
	privateKey := args[3]
	wepayBillNo := args[4]

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(privateKey)
	if err != nil {
		t.Errorf("Load merchant private key error: %+v", err)
		return
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		t.Errorf("new wechat pay client err:%s", err)
		return
	}

	// 调用申请电子回单的函数
	resp, err := ApplyReceipt(client, &ApplyReceiptRequest{
		ctx:         context.Background(),
		OutBillNo:   "",
		WepayBillNo: wepayBillNo,
	})
	if err != nil {
		t.Errorf("申请电子回单时出错: %+v", err)
		return
	}

	t.Logf("申请电子回单成功，状态: %s, 创建时间: %s\n", resp.State, resp.CreateTime)
}