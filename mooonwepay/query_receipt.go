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

// QueryReceiptRequest 查询电子回单请求结构体
type QueryReceiptRequest struct {
	ctx         context.Context
	OutBillNo   string `json:"out_bill_no"`   // 商户单号
	WepayBillNo string `json:"wepay_bill_no"` // 微信单号
}

// QueryReceiptResponse 查询电子回单响应结构体
type QueryReceiptResponse struct {
	State       string `json:"state"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	HashType    string `json:"hash_type"`
	HashValue   string `json:"hash_value"`
	DownloadUrl string `json:"download_url"`
}

// QueryReceipt 查询电子回单
func QueryReceipt(client *core.Client, req *QueryReceiptRequest) (*QueryReceiptResponse, error) {
	// 定义请求 URL
	var url string
	if req.OutBillNo != "" {
		url = fmt.Sprintf("https://api.mch.weixin.qq.com/v3/fund-app/mch-transfer/elecsign/out-bill-no/%s", req.OutBillNo)
	} else if req.WepayBillNo != "" {
		url = fmt.Sprintf("https://api.mch.weixin.qq.com/v3/fund-app/mch-transfer/elecsign/transfer-bill-no/%s", req.WepayBillNo)
	} else {
		return nil, fmt.Errorf("QueryReceipt with invalid request: out_bill_no and wepay_bill_no are both empty")
	}

	// 发送 GET 请求
	apiResult, err := client.Get(req.ctx, url)
	if err != nil {
		return nil, fmt.Errorf("QueryReceipt failed to query electronic receipt: %w", err)
	}

	// 读取响应体内容
	defer apiResult.Response.Body.Close()
	respBody, err := io.ReadAll(apiResult.Response.Body)
	if err != nil {
		return nil, fmt.Errorf("QueryReceipt failed to read response body: %w", err)
	}

	// 解析响应
	var resp QueryReceiptResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("QueryReceipt failed to unmarshal response: %w", err)
	}

	// 返回响应
	return &resp, nil
}
