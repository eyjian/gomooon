// Package mooonwepay
// Wrote by yijian on 2024/08/24
package mooonwepay

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/eyjian/gomooon/moooncrypto"
	"io"
	"net/http"
)

var (
	ApplyTradeBillPath = "/v3/bill/tradebill"    // 交易订单
	ApplyFundBillPath  = "/v3/bill/fundflowbill" // 资金流订单
	applyBillErrTag    = "apply bill error"
)

type ApplyBillReq struct {
	Ctx        context.Context
	HttpClient *http.Client
	PrivateKey *rsa.PrivateKey

	Host      string // 主域名：https://api.mch.weixin.qq.com，备域名：https://api2.mch.weixin.qq.com
	NonceStr  string
	Timestamp int64
	Mchid     string
	SerialNo  string

	// 交易账单类型：
	// ALL: 返回当日所有订单信息（不含充值退款订单）
	// SUCCESS: 返回当日成功支付的订单（不含充值退款订单）
	// REFUND: 返回当日退款订单（不含充值退款订单）
	// RECHARGE_REFUND: 返回当日充值退款订单
	// ALL_SPECIAL: 返回个性化账单当日所有订单信息
	// SUC_SPECIAL: 返回个性化账单当日成功支付的订单
	// REF_SPECIAL: 返回个性化账单当日退款订单
	//
	// 资金账单类型：
	// BASIC: 基本账户
	// OPERATION: 运营账户
	// FEES: 手续费账户
	BillType string // 账单类型

	CompressionType string // 压缩类型，目前仅支持取值 GZIP
	Date            string // 账单日期（格式：yyyy-MM-DD，仅支持三个月内的账单下载申请）
}

type ApplyBillResp struct {
	HashType    string `json:"hash_type,omitempty"`    //哈希类型，目前仅有 SHA1 一种取值
	HashValue   string `json:"hash_value,omitempty"`   // 用于校验文件的完整性
	DownloadUrl string `json:"download_url,omitempty"` // 5 分钟内有效（下载后为 gzip 压缩过的 csv 文件）

	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`

	HttpStatusCode int `json:"http_status_code,omitempty"`
}

var (
	// 交易账单类型
	tradeBillType = map[string]bool{
		"ALL":             true, // 返回当日所有订单信息（不含充值退款订单）
		"SUCCESS":         true, // 返回当日成功支付的订单（不含充值退款订单）
		"REFUND":          true, // 返回当日退款订单（不含充值退款订单）
		"RECHARGE_REFUND": true, // 返回当日充值退款订单
		"ALL_SPECIAL":     true, // 返回个性化账单当日所有订单信息
		"SUC_SPECIAL":     true, // 返回个性化账单当日成功支付的订单
		"REF_SPECIAL":     true, // 返回个性化账单当日退款订单
	}

	// 资金账单类型
	fundBillType = map[string]bool{
		"BASIC":     true, // 基本账户
		"OPERATION": true, // 运营账户
		"FEES":      true, // 手续费账户
	}
)

// IsTradeBill 是否为交易账单类型
func IsTradeBill(billType string) bool {
	_, ok := tradeBillType[billType]
	return ok
}

// IsFundBill 是否为资金账单类型
func IsFundBill(billType string) bool {
	_, ok := fundBillType[billType]
	return ok
}

// ApplyBill 申请交易账单
func ApplyBill(req *ApplyBillReq) (*ApplyBillResp, error) {
	ctx := req.Ctx
	url := getApplyBillUrl(req)

	// 检查传入的账单类型
	if !IsTradeBill(req.BillType) && !IsFundBill(req.BillType) {
		return nil, fmt.Errorf("%s: bill type unsupported: %s", applyBillErrTag, req.BillType)
	}

	// 计算签名
	signatureString := makeApplyBillSignatureString(req)
	signature, err := moooncrypto.RsaSha256SignWithPrivateKey(req.PrivateKey, []byte(signatureString))
	if err != nil {
		return nil, fmt.Errorf("%s: rsa sha256 sign error: %s", applyBillErrTag, err.Error())
	}

	// 生成 Authorization
	authorization := makeChangeBillAuthorization(req.Mchid, req.SerialNo, req.NonceStr, signature, req.Timestamp)

	// 构建请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: new http request error: %s", applyBillErrTag, err.Error())
	}

	// 设置请求头
	httpReq.Header.Set("Authorization", authorization)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := req.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s: do http request error: %s", applyBillErrTag, err.Error())
	}
	defer httpResp.Body.Close()

	// 读取响应
	resp := &ApplyBillResp{
		HttpStatusCode: httpResp.StatusCode,
	}
	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: read http body error: %s", applyBillErrTag, err.Error())
	}

	// 解析响应
	err = json.Unmarshal(respBodyBytes, resp)
	if httpResp.StatusCode != http.StatusOK {
		if httpResp.StatusCode == http.StatusUnauthorized {
			return resp, fmt.Errorf("%s: unauthorized, possible authorization incorrect", applyBillErrTag)
		} else {
			return resp, fmt.Errorf("%s: http response %d", applyBillErrTag, httpResp.StatusCode)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("%s: json unmarshal http response error: %s\n", applyBillErrTag, err.Error())
	}

	return resp, nil
}

// makeApplyBillSignatureString 生成签名串
//HTTP请求方法\n
//URL\n
//请求时间戳\n
//请求随机串\n
//请求报文主体\n
func makeApplyBillSignatureString(req *ApplyBillReq) string {
	if IsTradeBill(req.BillType) {
		return fmt.Sprintf("GET\n%s?bill_date=%s&bill_type=%s&tar_type=%s\n%d\n%s\n\n",
			ApplyTradeBillPath, req.Date, req.BillType, req.CompressionType, req.Timestamp, req.NonceStr)
	} else {
		return fmt.Sprintf("GET\n%s?bill_date=%s&bill_type=%s&tar_type=%s\n%d\n%s\n\n",
			ApplyFundBillPath, req.Date, req.BillType, req.CompressionType, req.Timestamp, req.NonceStr)
	}
}

func getApplyBillUrl(req *ApplyBillReq) string {
	if IsTradeBill(req.BillType) {
		return fmt.Sprintf("%s%s?bill_date=%s&bill_type=%s&tar_type=%s",
			req.Host, ApplyTradeBillPath, req.Date, req.BillType, req.CompressionType)
	} else {
		return fmt.Sprintf("%s%s?bill_date=%s&bill_type=%s&tar_type=%s",
			req.Host, ApplyFundBillPath, req.Date, req.BillType, req.CompressionType)
	}
}
