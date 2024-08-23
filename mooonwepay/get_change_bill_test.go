// Package mooonwepay
// Wrote by yijian on 2024/08/23
package mooonwepay

import (
	"context"
	"github.com/eyjian/gomooon/moooncrypto"
	"github.com/eyjian/gomooon/mooonutils"
	"net/http"
	"os"
	"testing"
	"time"
)

// go test -v -run="TestGetBill$" -args private_key.pem mchid serial_no out_batch_no <out_detail_no>
func TestGetBill(t *testing.T) {
	privateKeyFilepath := os.Args[len(os.Args)-4]
	mchid := os.Args[len(os.Args)-3]
	serialNo := os.Args[len(os.Args)-2]
	outBatchNo := os.Args[len(os.Args)-1]

	privateKey, err := moooncrypto.Filepath2PrivateKey(privateKeyFilepath)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("privateKey ok\n")

	timestamp := time.Now().Unix()
	nonceStr := mooonutils.GetNonceStr(32)
	req := &GetBillReq{
		Ctx:        context.Background(),
		HttpClient: &http.Client{},
		PrivateKey: privateKey,

		Host:      "https://api.mch.weixin.qq.com",
		NonceStr:  nonceStr,
		Timestamp: timestamp,

		Mchid:      mchid,
		SerialNo:   serialNo,
		OutBatchNo: outBatchNo,
	}
	t.Log(*req)
	resp, err := GetBill(req)
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		t.Logf("%+v\n", *resp)
	}
}