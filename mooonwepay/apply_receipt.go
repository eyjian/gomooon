// Package mooonwepay
// Wrote by yijian on 2025/03/27
package mooonwepay

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
)

// ApplyReceiptRequest 电子回单申请请求结构体
type ApplyReceiptRequest struct {
	ctx         context.Context
	OutBillNo   string `json:"out_bill_no"`   // 商户单号
	WepayBillNo string `json:"wepay_bill_no"` // 微信单号
}

// ApplyReceiptResponse 电子回单申请响应结构体
type ApplyReceiptResponse struct {
	State      string `json:"state"`
	CreateTime string `json:"create_time"`
}

// ApplyReceipt 申请电子回单
func ApplyReceipt(client *core.Client, req *ApplyReceiptRequest) (*ApplyReceiptResponse, error) {
	// 定义请求 URL
	var url string
	if req.OutBillNo != "" {
		url = "https://api.mch.weixin.qq.com/v3/fund-app/mch-transfer/elecsign/out-bill-no"
	} else if req.WepayBillNo != "" {
		url = "https://api.mch.weixin.qq.com/v3/fund-app/mch-transfer/elecsign/transfer-bill-no"
	} else {
		return nil, fmt.Errorf("ApplyReceipt with invalid request: out_bill_no and wepay_bill_no are both empty")
	}

	// 构建请求体
	requestBody := map[string]string{}
	if req.OutBillNo != "" {
		requestBody["out_bill_no"] = req.OutBillNo
	} else if req.WepayBillNo != "" {
		requestBody["transfer_bill_no"] = req.WepayBillNo
	}

	// 将请求体转换为 JSON 格式
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("ApplyReceipt failed to marshal request body: %w", err)
	}

	// 发送 POST 请求
	apiResult, err := client.Post(req.ctx, url, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("ApplyReceipt failed to apply electronic receipt: %w", err)
	}

	// 读取响应体内容
	defer apiResult.Response.Body.Close()
	respBody, err := io.ReadAll(apiResult.Response.Body)
	if err != nil {
		return nil, fmt.Errorf("ApplyReceipt failed to read response body: %w", err)
	}

	// 解析响应
	var resp ApplyReceiptResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("ApplyReceipt failed to unmarshal response: %w", err)
	}

	// 返回响应
	return &resp, nil
}