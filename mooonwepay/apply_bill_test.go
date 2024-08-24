// Package mooonwepay
// Wrote by yijian on 2024/08/24
package mooonwepay

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)
import (
	"github.com/eyjian/gomooon/moooncrypto"
	"github.com/eyjian/gomooon/mooonutils"
)

// go test -v -run="TestApplyBill$" -args private_key.pem mchid serial_no bill_type compression_type date
func TestApplyBill(t *testing.T) {
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
	req := &ApplyBillReq{
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
	}
	t.Log(*req)
	resp, err := ApplyBill(req)
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		t.Logf("%+v\n", *resp)
	}
}
