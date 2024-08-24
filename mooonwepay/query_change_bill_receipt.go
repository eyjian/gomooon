// Package mooonwepay
// Wrote by yijian on 2024/08/23
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
	QueryChangeBillBatchReceiptPath  = "/v3/transfer/bill-receipt"
	QueryChangeBillDetailReceiptPath = "/v3/transfer-detail/electronic-receipts"
	queryChangeBillReceiptErrTag     = "query change bill receipt error"
)

type QueryChangeBillReceiptReq struct {
	Ctx        context.Context
	HttpClient *http.Client
	PrivateKey *rsa.PrivateKey

	Host      string // 主域名：https://api.mch.weixin.qq.com，备域名：https://api2.mch.weixin.qq.com
	NonceStr  string
	Timestamp int64

	Mchid       string
	SerialNo    string
	OutBatchNo  string // 商家转账批次单号
	OutDetailNo string // 商家转账明细单号（如指定表示单笔转账电子回单）
	AcceptType  string // 电子回单受理类型：BATCH_TRANSFER：批量转账明细电子回单 TRANSFER_TO_POCKET：企业付款至零钱电子回单 TRANSFER_TO_BANK：企业付款至银行卡电子回单
}

type QueryChangeBillReceiptResp struct {
	AcceptType      string `json:"accept_type,omitempty"`
	OutBatchNo      string `json:"out_batch_no,omitempty"`
	OutDetailNo     string `json:"out_detail_no,omitempty"`
	SignatureNo     string `json:"signature_no,omitempty"`     // 电子回单申请单号，申请单据的唯一标识
	SignatureStatus string `json:"signature_status,omitempty"` // ACCEPTED:已受理，电子签章已受理成功 FINISHED:已完成。电子签章已处理完成
	HashType        string `json:"hash_type,omitempty"`
	HashValue       string `json:"hash_value,omitempty"`
	DownloadUrl     string `json:"download_url,omitempty"` // 电子回单文件的下载地址
	CreateTime      string `json:"create_time,omitempty"`
	UpdateTime      string `json:"update_time,omitempty"`

	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`

	HttpStatusCode int `json:"http_status_code,omitempty"`
}

// QueryChangeBillReceipt 查询转账电子回单
func QueryChangeBillReceipt(req *QueryChangeBillReceiptReq) (*QueryChangeBillReceiptResp, error) {
	ctx := req.Ctx
	url := getQueryChangeBillReceiptUrl(req)

	// 计算签名
	signatureString := makeQueryChangeBillReceiptSignatureString(req)
	signature, err := moooncrypto.RsaSha256SignWithPrivateKey(req.PrivateKey, []byte(signatureString))
	if err != nil {
		return nil, fmt.Errorf("%s: rsa sha256 sign error: %s", queryChangeBillReceiptErrTag, err.Error())
	}

	// 生成 Authorization
	authorization := makeChangeBillAuthorization(req.Mchid, req.SerialNo, req.NonceStr, signature, req.Timestamp)

	// 构建请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: new http request error: %s", queryChangeBillReceiptErrTag, err.Error())
	}

	// 设置请求头
	httpReq.Header.Set("Authorization", authorization)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := req.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s: do http request error: %s", queryChangeBillReceiptErrTag, err.Error())
	}
	defer httpResp.Body.Close()

	// 读取响应
	resp := &QueryChangeBillReceiptResp{
		HttpStatusCode: httpResp.StatusCode,
	}
	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: read http body error: %s", queryChangeBillReceiptErrTag, err.Error())
	}

	// 解析响应
	err = json.Unmarshal(respBodyBytes, resp)
	if httpResp.StatusCode != http.StatusOK {
		// {"code":"RESOURCE_ALREADY_EXISTS","message":"该批次回单已申请，您可在通过查询电子回单接口来获取单据信息"}
		if httpResp.StatusCode == http.StatusUnauthorized {
			return resp, fmt.Errorf("%s: unauthorized, possible authorization incorrect or out_batch_no error", queryChangeBillReceiptErrTag)
		} else {
			return resp, fmt.Errorf("%s: http response %d", queryChangeBillReceiptErrTag, httpResp.StatusCode)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("%s: json unmarshal http response error: %s\n", queryChangeBillReceiptErrTag, err.Error())
	}

	return resp, nil
}

// makeQueryBillReceiptSignatureString 生成签名串
//HTTP请求方法\n
//URL\n
//请求时间戳\n
//请求随机串\n
//请求报文主体\n
func makeQueryChangeBillReceiptSignatureString(req *QueryChangeBillReceiptReq) string {
	if req.AcceptType == "" || req.OutDetailNo == "" {
		return fmt.Sprintf("GET\n%s/%s\n%d\n%s\n\n",
			QueryChangeBillBatchReceiptPath, req.OutBatchNo, req.Timestamp, req.NonceStr)
	} else {
		return fmt.Sprintf("GET\n%s?out_batch_no=%s&out_detail_no=%s&accept_type=%s\n%d\n%s\n\n",
			QueryChangeBillDetailReceiptPath, req.OutBatchNo, req.OutDetailNo, req.AcceptType,
			req.Timestamp, req.NonceStr)
	}
}

func getQueryChangeBillReceiptUrl(req *QueryChangeBillReceiptReq) string {
	if req.AcceptType == "" || req.OutDetailNo == "" {
		return fmt.Sprintf("%s%s/%s",
			req.Host, QueryChangeBillBatchReceiptPath, req.OutBatchNo)
	} else {
		return fmt.Sprintf("%s%s?out_batch_no=%s&out_detail_no=%s&accept_type=%s",
			req.Host, QueryChangeBillDetailReceiptPath, req.OutBatchNo, req.OutDetailNo, req.AcceptType)
	}
}
