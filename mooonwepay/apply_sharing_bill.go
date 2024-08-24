// Package mooonwepay
// Wrote by yijian on 2024/08/24
package mooonwepay

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)
import (
	"github.com/eyjian/gomooon/moooncrypto"
)

var (
	ApplySharingBillPath   = "/v3/profitsharing/bills" // 分账账单
	applySharingBillErrTag = "apply sharing bill error"
)

type ApplySharingBillReq struct {
	Ctx        context.Context
	HttpClient *http.Client
	PrivateKey *rsa.PrivateKey

	Host      string // 主域名：https://api.mch.weixin.qq.com，备域名：https://api2.mch.weixin.qq.com
	NonceStr  string
	Timestamp int64
	Mchid     string
	SerialNo  string

	CompressionType string // 压缩类型，目前仅支持取值 GZIP
	Date            string // 账单日期（格式：yyyy-MM-DD，仅支持三个月内的账单下载申请）
}

type ApplySharingBillResp struct {
	HashType    string `json:"hash_type,omitempty"`    //哈希类型，目前仅有 SHA1 一种取值
	HashValue   string `json:"hash_value,omitempty"`   // 用于校验文件的完整性
	DownloadUrl string `json:"download_url,omitempty"` // 5 分钟内有效（下载后为 gzip 压缩过的 csv 文件）

	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`

	HttpStatusCode int `json:"http_status_code,omitempty"`
}

// ApplySharingBill 申请分账账单
func ApplySharingBill(req *ApplySharingBillReq) (*ApplySharingBillResp, error) {
	ctx := req.Ctx
	url := getApplySharingBillUrl(req)

	// 计算签名
	signatureString := makeApplySharingBillSignatureString(req)
	signature, err := moooncrypto.RsaSha256SignWithPrivateKey(req.PrivateKey, []byte(signatureString))
	if err != nil {
		return nil, fmt.Errorf("%s: rsa sha256 sign error: %s", applySharingBillErrTag, err.Error())
	}

	// 生成 Authorization
	authorization := makeChangeBillAuthorization(req.Mchid, req.SerialNo, req.NonceStr, signature, req.Timestamp)

	// 构建请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: new http request error: %s", applySharingBillErrTag, err.Error())
	}

	// 设置请求头
	httpReq.Header.Set("Authorization", authorization)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := req.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s: do http request error: %s", applySharingBillErrTag, err.Error())
	}
	defer httpResp.Body.Close()

	// 读取响应
	resp := &ApplySharingBillResp{
		HttpStatusCode: httpResp.StatusCode,
	}
	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: read http body error: %s", applySharingBillErrTag, err.Error())
	}

	// 解析响应
	err = json.Unmarshal(respBodyBytes, resp)
	if httpResp.StatusCode != http.StatusOK {
		if httpResp.StatusCode == http.StatusUnauthorized {
			return resp, fmt.Errorf("%s: unauthorized, possible authorization incorrect", applySharingBillErrTag)
		} else {
			return resp, fmt.Errorf("%s: http response %d", applySharingBillErrTag, httpResp.StatusCode)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("%s: json unmarshal http response error: %s\n", applySharingBillErrTag, err.Error())
	}

	return resp, nil
}

// makeApplySharingBillSignatureString 生成签名串
//HTTP请求方法\n
//URL\n
//请求时间戳\n
//请求随机串\n
//请求报文主体\n
func makeApplySharingBillSignatureString(req *ApplySharingBillReq) string {
	return fmt.Sprintf("GET\n%s?bill_date=%s&tar_type=%s\n%d\n%s\n\n",
		ApplyFundBillPath, req.Date, req.CompressionType, req.Timestamp, req.NonceStr)
}

func getApplySharingBillUrl(req *ApplySharingBillReq) string {
	return fmt.Sprintf("%s%s?bill_date=%s&tar_type=%s",
		req.Host, ApplyFundBillPath, req.Date, req.CompressionType)
}
