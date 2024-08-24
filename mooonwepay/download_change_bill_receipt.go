// Package mooonwepay
// Wrote by yijian on 2024/08/23
package mooonwepay

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/eyjian/gomooon/mooonutils"
	"io"
	"net/http"
	"os"
)
import (
	"github.com/eyjian/gomooon/moooncrypto"
)

var (
	downloadChangeBillReceiptErrTag = "download change bill receipt error"
)

type DownloadChangeBillReceiptReq struct {
	Ctx        context.Context
	HttpClient *http.Client
	PrivateKey *rsa.PrivateKey

	Host      string // 主域名：https://api.mch.weixin.qq.com，备域名：https://api2.mch.weixin.qq.com
	NonceStr  string
	Timestamp int64
	Mchid     string
	SerialNo  string

	OutBatchNo  string // 商家转账批次单号
	OutDetailNo string // 商家转账明细单号（如指定表示单笔转账电子回单）
	AcceptType  string // 电子回单受理类型：BATCH_TRANSFER：批量转账明细电子回单 TRANSFER_TO_POCKET：企业付款至零钱电子回单 TRANSFER_TO_BANK：企业付款至银行卡电子回单

	Filepath string // 下载后存放的文件路径
}

type DownloadChangeBillReceiptResp struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`

	HttpStatusCode int `json:"http_status_code,omitempty"`
}

// DownloadChangeBillReceipt 下载转账电子回单
func DownloadChangeBillReceipt(req *DownloadChangeBillReceiptReq) (*DownloadChangeBillReceiptResp, error) {
	ctx := req.Ctx

	// 通过调用 QueryBill，取得下载 url
	queryBillResp, err := getChangeBillReceiptDownloadUrl(req)
	if err != nil {
		if queryBillResp != nil {
			return &DownloadChangeBillReceiptResp{
				Code:           queryBillResp.Code,
				Message:        queryBillResp.Message,
				HttpStatusCode: queryBillResp.HttpStatusCode,
			}, err
		}
		return nil, err
	}

	// 计算签名
	downloadPath := mooonutils.ExtractUrlPath(queryBillResp.DownloadUrl)
	signatureString := makeDownloadChangeBillReceiptSignatureString(req.NonceStr, downloadPath, req.Timestamp)
	signature, err := moooncrypto.RsaSha256SignWithPrivateKey(req.PrivateKey, []byte(signatureString))
	if err != nil {
		return nil, fmt.Errorf("%s: RSA SHA256 sign error: %s", downloadChangeBillReceiptErrTag, err.Error())
	}

	// 生成 Authorization
	authorization := makeChangeBillAuthorization(req.Mchid, req.SerialNo, req.NonceStr, signature, req.Timestamp)

	// 构建请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", queryBillResp.DownloadUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: new http request error: %s", downloadChangeBillReceiptErrTag, err.Error())
	}

	// 设置请求头
	httpReq.Header.Set("Authorization", authorization)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	httpResp, err := req.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s: do http request error: %s", downloadChangeBillReceiptErrTag, err.Error())
	}
	defer httpResp.Body.Close()

	resp := &DownloadChangeBillReceiptResp{
		HttpStatusCode: httpResp.StatusCode,
	}
	if httpResp.StatusCode != http.StatusOK {
		respBodyBytes, err := io.ReadAll(httpResp.Body)
		if err == nil {
			json.Unmarshal(respBodyBytes, resp)
		}
		return resp, fmt.Errorf("%s: http get %s status code error: %d", downloadChangeBillReceiptErrTag, queryBillResp.DownloadUrl, httpResp.StatusCode)
	}

	// 创建文件
	file, err := os.Create(req.Filepath)
	if err != nil {
		return resp, fmt.Errorf("%s: create %s error: %s", downloadChangeBillReceiptErrTag, req.Filepath, err.Error())
	}
	defer file.Close()

	// 写入文件
	_, err = io.Copy(file, httpResp.Body)
	if err != nil {
		return resp, fmt.Errorf("%s: write response to %s error: %s", downloadChangeBillReceiptErrTag, req.Filepath, err.Error())
	}

	return resp, nil
}

// makeDownloadChangeBillReceiptSignatureString 生成签名串
//HTTP请求方法\n
//URL\n
//请求时间戳\n
//请求随机串\n
//请求报文主体\n
func makeDownloadChangeBillReceiptSignatureString(nonceStr, downloadUrl string, timestamp int64) string {
	return fmt.Sprintf("GET\n%s\n%d\n%s\n\n", downloadUrl, timestamp, nonceStr)
}

func getChangeBillReceiptDownloadUrl(req *DownloadChangeBillReceiptReq) (*QueryChangeBillReceiptResp, error) {
	// 调用 ApplyBillReceipt，取得下载 url
	getBillReceiptResp, err := ApplyChangeBillReceipt(
		&ApplyChangeBillReceiptReq{
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
	// ALREADY_EXISTS 转账电子回单申请单数据已存在
	// RESOURCE_ALREADY_EXISTS 该批次回单已申请，您可在通过查询电子回单接口来获取单据信息
	if err == nil || (getBillReceiptResp != nil && getBillReceiptResp.Code != "ALREADY_EXISTS" && getBillReceiptResp.Code != "RESOURCE_ALREADY_EXISTS") {
		return &QueryChangeBillReceiptResp{
			AcceptType:      getBillReceiptResp.AcceptType,
			OutBatchNo:      getBillReceiptResp.OutBatchNo,
			OutDetailNo:     getBillReceiptResp.OutDetailNo,
			SignatureNo:     getBillReceiptResp.SignatureNo,
			SignatureStatus: getBillReceiptResp.SignatureStatus,
			HashType:        getBillReceiptResp.HashType,
			HashValue:       getBillReceiptResp.HashValue,
			DownloadUrl:     getBillReceiptResp.DownloadUrl,
			CreateTime:      getBillReceiptResp.CreateTime,
			UpdateTime:      getBillReceiptResp.UpdateTime,

			Code:    getBillReceiptResp.Code,
			Message: getBillReceiptResp.Message,

			HttpStatusCode: getBillReceiptResp.HttpStatusCode,
		}, err
	}
	if err != nil {
		if getBillReceiptResp == nil || (getBillReceiptResp.Code != "ALREADY_EXISTS" && getBillReceiptResp.Code != "RESOURCE_ALREADY_EXISTS") {
			return nil, err
		}
	}

	// 调用 QueryBillReceipt，取得下载 url
	queryBillReceiptResp, err := QueryChangeBillReceipt(
		&QueryChangeBillReceiptReq{
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
	return queryBillReceiptResp, err
}
