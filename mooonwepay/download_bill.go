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
	"os"
)
import (
	"github.com/eyjian/gomooon/moooncrypto"
	"github.com/eyjian/gomooon/mooonutils"
)

var (
	downloadBillErrTag = "download bill error"
)

type DownloadBillReq struct {
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

	Filepath string // 下载后存放的文件路径（下载后为 gzip 压缩过的 csv 文件）
}

type DownloadBillResp struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`

	HttpStatusCode int `json:"http_status_code,omitempty"`
}

// DownloadBill 下载转账电子回单
// 账单文件格式，参加微信支付官方文档：https://pay.weixin.qq.com/docs/merchant/apis/bill-download/download-bill.html
func DownloadBill(req *DownloadBillReq) (*DownloadBillResp, error) {
	ctx := req.Ctx

	// 通过调用 ApplyBill，取得下载 url
	applyBillResp, err := getBillDownloadUrl(req)
	if err != nil {
		if applyBillResp != nil {
			return &DownloadBillResp{
				Code:           applyBillResp.Code,
				Message:        applyBillResp.Message,
				HttpStatusCode: applyBillResp.HttpStatusCode,
			}, err
		}
		return nil, err
	}

	// 计算签名
	downloadPath := mooonutils.ExtractUrlPath(applyBillResp.DownloadUrl)
	signatureString := makeDownloadBillSignatureString(req.NonceStr, downloadPath, req.Timestamp)
	signature, err := moooncrypto.RsaSha256SignWithPrivateKey(req.PrivateKey, []byte(signatureString))
	if err != nil {
		return nil, fmt.Errorf("%s: RSA SHA256 sign error: %s", downloadBillErrTag, err.Error())
	}

	// 生成 Authorization
	authorization := makeChangeBillAuthorization(req.Mchid, req.SerialNo, req.NonceStr, signature, req.Timestamp)

	// 构建请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", applyBillResp.DownloadUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: new http request error: %s", downloadBillErrTag, err.Error())
	}

	// 设置请求头
	httpReq.Header.Set("Authorization", authorization)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	httpResp, err := req.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s: do http request error: %s", downloadBillErrTag, err.Error())
	}
	defer httpResp.Body.Close()

	resp := &DownloadBillResp{
		HttpStatusCode: httpResp.StatusCode,
	}
	if httpResp.StatusCode != http.StatusOK {
		respBodyBytes, err := io.ReadAll(httpResp.Body)
		if err == nil {
			json.Unmarshal(respBodyBytes, resp)
		}
		return resp, fmt.Errorf("%s: http get %s status code error: %d", downloadBillErrTag, applyBillResp.DownloadUrl, httpResp.StatusCode)
	}

	// 创建文件
	file, err := os.Create(req.Filepath)
	if err != nil {
		return resp, fmt.Errorf("%s: create %s error: %s", downloadBillErrTag, req.Filepath, err.Error())
	}
	defer file.Close()

	// 写入文件
	_, err = io.Copy(file, httpResp.Body)
	if err != nil {
		return resp, fmt.Errorf("%s: write response to %s error: %s", downloadBillErrTag, req.Filepath, err.Error())
	}

	return resp, nil
}

// makeDownloadBillSignatureString 生成签名串
//HTTP请求方法\n
//URL\n
//请求时间戳\n
//请求随机串\n
//请求报文主体\n
func makeDownloadBillSignatureString(nonceStr, downloadUrl string, timestamp int64) string {
	return fmt.Sprintf("GET\n%s\n%d\n%s\n\n", downloadUrl, timestamp, nonceStr)
}

func getBillDownloadUrl(req *DownloadBillReq) (*ApplyBillResp, error) {
	// 调用 ApplyBill，取得下载 url
	applyBillResp, err := ApplyBill(
		&ApplyBillReq{
			Ctx:        req.Ctx,
			HttpClient: req.HttpClient,
			PrivateKey: req.PrivateKey,

			Host:      req.Host,
			NonceStr:  req.NonceStr,
			Timestamp: req.Timestamp,
			Mchid:     req.Mchid,
			SerialNo:  req.SerialNo,

			BillType:        req.BillType,
			CompressionType: req.CompressionType,
			Date:            req.Date,
		})
	if err != nil {
		return nil, err
	}
	return applyBillResp, nil
}
