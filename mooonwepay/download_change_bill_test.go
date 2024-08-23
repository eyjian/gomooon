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

// go test -v -run="TestDownloadBill$" -args private_key.pem serial_no mchid out_batch_no <out_detail_no>
func TestDownloadBill(t *testing.T) {
	numArgs := len(os.Args)
	t.Log("args num:", numArgs)
	if numArgs != 9 && numArgs != 10 {
		t.Error("args error")
		return
	}

	var (
		privateKeyFilepath string
		mchid              string
		serialNo           string
		outBatchNo         string
		outDetailNo        string
	)
	if numArgs == 9 {
		privateKeyFilepath = os.Args[len(os.Args)-4]
		mchid = os.Args[len(os.Args)-3]
		serialNo = os.Args[len(os.Args)-2]
		outBatchNo = os.Args[len(os.Args)-1]
	} else {
		privateKeyFilepath = os.Args[len(os.Args)-5]
		mchid = os.Args[len(os.Args)-4]
		serialNo = os.Args[len(os.Args)-3]
		outBatchNo = os.Args[len(os.Args)-2]
		outDetailNo = os.Args[len(os.Args)-1]
	}

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

			Mchid:       mchid,
			SerialNo:    serialNo,
			OutBatchNo:  outBatchNo,
			OutDetailNo: outDetailNo,
			AcceptType:  "BATCH_TRANSFER",

			Filepath: fmt.Sprintf("change_bill-%s.pdf", time.Now().Format("20060102150405")),
		})
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		t.Logf("%+v\n", *resp)
	}
}