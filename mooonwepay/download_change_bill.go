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
	"net/url"
	"os"
)
import (
	"github.com/eyjian/gomooon/moooncrypto"
)

type DownloadBillReq struct {
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

	Filepath string // 下载后的文件路径
}

type DownloadBillResp struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`

	HttpStatusCode int `json:"http_status_code,omitempty"`
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

// DownloadBill 下载电子回单
func DownloadBill(req *DownloadBillReq) (*DownloadBillResp, error) {
	ctx := req.Ctx

	// 通过调用 QueryBill，取得下载 url
	queryBillResp, err := getDownloadUrl(req)
	if err != nil {
		return &DownloadBillResp{
			Code:           queryBillResp.Code,
			Message:        queryBillResp.Message,
			HttpStatusCode: queryBillResp.HttpStatusCode,
		}, err
	}

	// 计算签名
	downloadPath := extractUrlPath(queryBillResp.DownloadUrl)
	signatureString := makeDownloadBillSignatureString(req.NonceStr, downloadPath, req.Timestamp)
	signature, err := moooncrypto.RsaSha256SignWithPrivateKey(req.PrivateKey, []byte(signatureString))
	if err != nil {
		return nil, fmt.Errorf("download bill error: RSA SHA256 sign error: %s", err.Error())
	}

	// 生成 Authorization
	authorization := makeChangeBillAuthorization(req.Mchid, req.SerialNo, req.NonceStr, signature, req.Timestamp)

	// 构建请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", queryBillResp.DownloadUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("download bill error: new http request error: %s", err.Error())
	}

	// 设置请求头
	httpReq.Header.Set("Authorization", authorization)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	httpResp, err := req.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("download bill error: do http request error: %s", err.Error())
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
		return resp, fmt.Errorf("download bill error: http get %s status code error: %d", queryBillResp.DownloadUrl, httpResp.StatusCode)
	}

	// 创建文件
	file, err := os.Create(req.Filepath)
	if err != nil {
		return resp, fmt.Errorf("download bill error: create %s error: %s", req.Filepath, err.Error())
	}
	defer file.Close()

	// 写入文件
	_, err = io.Copy(file, httpResp.Body)
	if err != nil {
		return resp, fmt.Errorf("download bill error: write response to %s error: %s", req.Filepath, err.Error())
	}

	return resp, nil
}

func getDownloadUrl(req *DownloadBillReq) (*QueryBillResp, error) {
	// 调用 QueryBill，取得下载 url
	queryBillResp, err := QueryBill(
		&QueryBillReq{
			Ctx:        req.Ctx,
			HttpClient: req.HttpClient,
			PrivateKey: req.PrivateKey,

			Host:      req.Host,
			NonceStr:  req.NonceStr,
			Timestamp: req.Timestamp,

			Mchid:       req.Mchid,
			SerialNo:    req.SerialNo,
			OutBatchNo:  req.OutBatchNo,
			OutDetailNo: req.OutDetailNo,
			AcceptType:  req.AcceptType,
		})
	return queryBillResp, err
}

// extractUrlPath 提取 url 路径
func extractUrlPath(urlStr string) string {
	parsedUrl, _ := url.Parse(urlStr)
	parsedUrl.Scheme = ""
	parsedUrl.Host = ""
	return parsedUrl.String()
}