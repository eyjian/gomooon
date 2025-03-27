// Package mooonwepay
// Wrote by yijian on 2025/03/27
package mooonwepay

import (
	"context"
	"os"
	"testing"

	"github.com/tjfoc/gmsm/sm3"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// go test -v -run="TestDownloadReceipt" -args mchid serial_no mch_api_v3_key private_key download_url local_file_path
func TestDownloadReceipt(t *testing.T) {
	args := os.Args[5:]
	if len(args) != 6 {
		t.Fatal("args num error")
	}
	mchID := args[0]
	mchCertificateSerialNumber := args[1]
	mchAPIv3Key := args[2]
	privateKey := args[3]
	downloadUrl := args[4]
	localFilePath := args[5]

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(privateKey)
	if err != nil {
		t.Errorf("Load merchant private key error: %+v", err)
		return
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	// 参考 https://github.com/wechatpay-apiv3/wechatpay-go/blob/main/FAQ.md
	// 其中第二步的应答中不包含应答数字签名，无法验签，应使用 WithoutValidator() 跳过应答签名的校验
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
		option.WithoutValidator(),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		t.Errorf("new wechat pay client err:%s", err)
		return
	}

	// 构建下载电子回单的请求
	req := &DownloadReceiptRequest{
		ctx:           context.Background(),
		DownloadUrl:   downloadUrl,
		LocalFilePath: localFilePath,
		HashType:      "SM3", // 或 "SHA256"
		HashValue:     "40F77A7B8B121B3E4841D9363C7F6293476973709D4B9C6552FD9CA9996D7CBB",
	}

	t.Logf("下载电子回单请求: %+v\n", *req)

	// 调用下载电子回单的函数
	resp, err := DownloadReceipt(client, req)
	if err != nil {
		t.Errorf("下载电子回单时出错: %+v", err)
		return
	}
	if resp.Code != "SUCCESS" {
		t.Errorf("下载电子回单时出错: %+v", *resp)
	} else {
		t.Logf("下载电子回单成功: %+v\n", *resp)

		if req.HashType == "" {
			fileBytes, err := os.ReadFile(localFilePath)
			if err != nil {
				t.Errorf("读取本地文件失败: %+v", err)
			} else {
				// 使用 SHA256 还是 SM3 视 QueryReceipt 的返回而定
				//hashValue := moooncrypto.Sha256Sign(string(fileBytes), "")
				hashValue := sm3.Sm3Sum(fileBytes)
				t.Logf("hashValue: %X\n", hashValue)
			}
		}
	}
}
