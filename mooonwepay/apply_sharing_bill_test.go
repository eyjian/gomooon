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

// go test -v -run="TestApplySharingBill$" -args private_key.pem mchid serial_no compression_type date
func TestApplySharingBill(t *testing.T) {
	numArgs := len(os.Args)
	args := os.Args[5:] // 从左到右到 args 共 5 个
	t.Logf("args num: %d (%s)\n", numArgs, args[0])
	privateKeyFilepath := args[0]
	mchid := args[1]
	serialNo := args[2]
	compressionType := args[3]
	date := args[4]

	privateKey, err := moooncrypto.Filepath2PrivateKey(privateKeyFilepath)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("privateKey ok\n")

	timestamp := time.Now().Unix()
	nonceStr := mooonutils.GetNonceStr(32)
	req := &ApplySharingBillReq{
		Ctx:        context.Background(),
		HttpClient: &http.Client{},
		PrivateKey: privateKey,

		Host:      "https://api.mch.weixin.qq.com",
		NonceStr:  nonceStr,
		Timestamp: timestamp,
		Mchid:     mchid,
		SerialNo:  serialNo,

		CompressionType: compressionType,
		Date:            date,
	}
	t.Log(*req)
	resp, err := ApplySharingBill(req)
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		t.Logf("%+v\n", *resp)
	}
}
