// Package mooonwepay
// Wrote by yijian on 2025/03/27
package mooonwepay

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eyjian/gomooon/moooncrypto"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
)

// DownloadReceiptRequest 电子回单下载请求结构体
type DownloadReceiptRequest struct {
	ctx           context.Context
	HashType      string `json:"hash_type"` // 哈希类型（SM3、SHA256），取值需同 QueryReceipt 返回的一致，如果未指定则不校验文件的哈希值
	HashValue     string `json:"hash_value"`
	DownloadUrl   string `json:"download_url"`    // 下载地址
	LocalFilePath string `json:"local_file_path"` // 本地文件路径（指定了 HashType 才会核验）
}

// DownloadReceiptResponse 电子回单下载响应结构体
type DownloadReceiptResponse struct {
	Code    string `json:"code"` // 成功值为 SUCCESS
	Message string `json:"message"`
}

// DownloadReceipt 下载电子回单
// 官方文档：https://pay.weixin.qq.com/doc/v3/merchant/4013866774
func DownloadReceipt(client *core.Client, req *DownloadReceiptRequest) (*DownloadReceiptResponse, error) {
	// 发送 GET 请求
	apiResult, err := client.Get(req.ctx, req.DownloadUrl)
	if err != nil {
		return nil, fmt.Errorf("DownloadReceipt failed to download electronic receipt: %w", err)
	}
	defer apiResult.Response.Body.Close()

	// 检查响应状态码
	if apiResult.Response.StatusCode != 200 {
		return nil, fmt.Errorf("DownloadReceipt failed with status code: %d", apiResult.Response.StatusCode)
	}

	// 读取响应体内容
	respBody, err := io.ReadAll(apiResult.Response.Body)
	if err != nil {
		return nil, fmt.Errorf("DownloadReceipt failed to read response body: %w", err)
	}

	// 检查响应体是否为空
	if len(respBody) == 0 {
		return nil, fmt.Errorf("DownloadReceipt response body is empty")
	}

	resp := &DownloadReceiptResponse{}
	err = json.Unmarshal(respBody, resp)
	if err == nil {
		return resp, nil
	}

	// 验证哈希值
	hashType := strings.ToUpper(req.HashType)
	if hashType == "SM3" {
		hashValue := fmt.Sprintf("%X", sm3.Sm3Sum(respBody))
		if req.HashValue != hashValue {
			return nil, fmt.Errorf("DownloadReceipt failed to verify SM3 hash value: %s", hashValue)
		}
	} else if hashType == "SHA256" {
		hashValue := moooncrypto.Sha256Sign(string(respBody), "")
		if req.HashValue != hashValue {
			return nil, fmt.Errorf("DownloadReceipt failed to verify SHA256 hash value: %s", hashValue)
		}
	}

	// 写入本地文件
	file, err := os.Create(req.LocalFilePath)
	if err != nil {
		return nil, fmt.Errorf("DownloadReceipt failed to create file: %w", err)
	}
	defer file.Close()
	_, err = file.Write(respBody)
	if err != nil {
		os.Remove(req.LocalFilePath)
		return nil, fmt.Errorf("DownloadReceipt failed to write file: %w", err)
	}

	// 返回响应
	return &DownloadReceiptResponse{
		Code:    "SUCCESS",
		Message: "SUCCESS",
	}, nil
}
