// Package mooonwepay
// Wrote by yijian on 2024/08/23
package mooonwepay

import (
	"context"
	"fmt"
	"github.com/eyjian/gomooon/moooncrypto"
	"github.com/eyjian/gomooon/mooonutils"
	"net/http"
	"os"
	"testing"
	"time"
)

// go test -v -run="TestDownloadBill$" -args private_key.pem mchid serial_no bill_type compression_type date
func TestDownloadBill(t *testing.T) {
	numArgs := len(os.Args)
	args := os.Args[5:] // 从左到右到 args 共 5 个
	t.Logf("args num: %d (%s)\n", numArgs, args[0])
	privateKeyFilepath := args[0]
	mchid := args[1]
	serialNo := args[2]
	billType := args[3]
	compressionType := args[4]
	date := args[5]

	privateKey, err := moooncrypto.Filepath2PrivateKey(privateKeyFilepath)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("privateKey ok\n")

	timestamp := time.Now().Unix()
	nonceStr := mooonutils.GetNonceStr(32)
	resp, err := DownloadBill(
		&DownloadBillReq{
			Ctx:        context.Background(),
			HttpClient: &http.Client{},
			PrivateKey: privateKey,

			Host:      "https://api.mch.weixin.qq.com",
			NonceStr:  nonceStr,
			Timestamp: timestamp,
			Mchid:     mchid,
			SerialNo:  serialNo,

			BillType:        billType,
			CompressionType: compressionType,
			Date:            date,

			Filepath: fmt.Sprintf("bill-%s.tar.gz", mooonutils.ConvertDateFormat(date)),
		})
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		t.Logf("%+v\n", *resp)
	}
}
