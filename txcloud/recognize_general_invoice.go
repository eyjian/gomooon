// Package txcloud
// Wrote by yijian on 2026/05/29
//
// 基于腾讯云 OCR `RecognizeGeneralInvoice` 的通用混贴发票识别实现。
// 该实现满足 InvoiceRecognizer 接口，调用方可通过接口持有该实现，
// 后续若新增其它实现（例如新版接口或其它厂商），调用方代码无需改动。
package txcloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

// 默认 region，用于未显式指定时的兜底
const defaultGeneralInvoiceRegion = "ap-guangzhou"

// GeneralInvoiceRecognizer 基于腾讯云 RecognizeGeneralInvoice 的发票识别实现
type GeneralInvoiceRecognizer struct {
	*TxCloud // 继承 Region 等字段
}

// 编译期断言：确保实现满足 InvoiceRecognizer 接口
var _ InvoiceRecognizer = (*GeneralInvoiceRecognizer)(nil)

// NewGeneralInvoiceRecognizer 实例化通用发票识别器
//
// 参数：
//
//	secretId / secretKey 腾讯云访问密钥
//	region              腾讯云地域，传空字符串将使用默认值 "ap-guangzhou"
func NewGeneralInvoiceRecognizer(secretId, secretKey, region string) *GeneralInvoiceRecognizer {
	if region == "" {
		region = defaultGeneralInvoiceRegion
	}
	tx := NewTxCloud(secretId, secretKey, "ocr.tencentcloudapi.com")
	tx.Region = region
	return &GeneralInvoiceRecognizer{
		TxCloud: tx,
	}
}

// RecognizeInvoice 识别发票（实现 InvoiceRecognizer 接口）
func (g *GeneralInvoiceRecognizer) RecognizeInvoice(ctx context.Context, req *InvoiceRequest) (*InvoiceResult, error) {
	if req == nil {
		return nil, fmt.Errorf("invoice request is nil")
	}

	imageURL, imageBase64, isPdf, err := resolveSource(req)
	if err != nil {
		return nil, err
	}

	client, err := ocr.NewClient(g.credential, g.TxCloud.Region, g.clientProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create ocr client: %w", err)
	}

	request := ocr.NewRecognizeGeneralInvoiceRequest()
	if ctx != nil {
		request.SetContext(ctx)
	}
	if imageURL != "" {
		request.ImageUrl = common.StringPtr(imageURL)
	}
	if imageBase64 != "" {
		request.ImageBase64 = common.StringPtr(imageBase64)
	}
	if isPdf {
		request.EnablePdf = common.BoolPtr(true)
		page := req.PdfPageNumber
		if page <= 0 {
			page = 1
		}
		request.PdfPageNumber = common.Int64Ptr(int64(page))
	}
	if len(req.Types) > 0 {
		ts := make([]*int64, 0, len(req.Types))
		for _, t := range req.Types {
			v := t
			ts = append(ts, &v)
		}
		request.Types = ts
	}

	response, err := client.RecognizeGeneralInvoice(request)
	if err != nil {
		// 腾讯云 SDK 错误（*errors.TencentCloudSDKError）原样向上透传
		return nil, err
	}

	return buildInvoiceResult(response), nil
}

// resolveSource 解析 InvoiceRequest 中的输入源，返回最终用于调用腾讯云 API 的
// imageURL、imageBase64、isPdf 三元组。
//
// 处理优先级：
//  1. 高级字段（ImageURL/ImageBase64/ImageFile）优先生效（按上述顺序兜底取值）
//  2. 高级字段全空时使用 Source 字段，自动判别 URL/Base64/本地文件
//
// 若所有输入源全空，返回英文参数错误。
func resolveSource(req *InvoiceRequest) (imageURL, imageBase64 string, isPdf bool, err error) {
	isPdf = req.IsPdf

	// 优先使用高级字段
	if req.ImageURL != "" {
		imageURL = req.ImageURL
		if !isPdf && looksLikePdfPath(imageURL) {
			isPdf = true
		}
		return
	}
	if req.ImageBase64 != "" {
		imageBase64 = stripBase64Prefix(req.ImageBase64)
		if !isPdf && looksLikePdfBase64(imageBase64) {
			isPdf = true
		}
		return
	}
	if req.ImageFile != "" {
		imageBase64, err = readFileAsBase64(req.ImageFile)
		if err != nil {
			return
		}
		if !isPdf && (looksLikePdfPath(req.ImageFile) || looksLikePdfBase64(imageBase64)) {
			isPdf = true
		}
		return
	}

	// 全部高级字段为空时，使用 Source 自动判别
	src := strings.TrimSpace(req.Source)
	if src == "" {
		err = fmt.Errorf("invoice request: at least one of Source/ImageURL/ImageBase64/ImageFile must be provided")
		return
	}

	switch {
	case strings.HasPrefix(src, "http://"), strings.HasPrefix(src, "https://"):
		imageURL = src
		if !isPdf && looksLikePdfPath(src) {
			isPdf = true
		}
	case strings.HasPrefix(src, "data:") && strings.Contains(src, ";base64,"):
		imageBase64 = stripBase64Prefix(src)
		if !isPdf && (strings.Contains(src, "application/pdf") || looksLikePdfBase64(imageBase64)) {
			isPdf = true
		}
	default:
		// 先按本地文件判定：路径存在则当文件读
		if fi, statErr := os.Stat(src); statErr == nil && !fi.IsDir() {
			imageBase64, err = readFileAsBase64(src)
			if err != nil {
				return
			}
			if !isPdf && (looksLikePdfPath(src) || looksLikePdfBase64(imageBase64)) {
				isPdf = true
			}
			return
		}
		// 否则尝试当作 Base64：长度足够且字符集合法
		if len(src) > 256 && isLikelyBase64(src) {
			imageBase64 = src
			if !isPdf && looksLikePdfBase64(imageBase64) {
				isPdf = true
			}
			return
		}
		err = fmt.Errorf("invoice request: cannot determine source type, file not found and not a valid base64 string")
	}
	return
}

// readFileAsBase64 读取本地文件并转为 Base64 字符串
func readFileAsBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", path, err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// stripBase64Prefix 去除 "data:xxx;base64," 前缀
func stripBase64Prefix(s string) string {
	if !strings.HasPrefix(s, "data:") {
		return s
	}
	_, after, ok := strings.Cut(s, ";base64,")
	if !ok {
		return s
	}
	return after
}

// looksLikePdfPath 根据路径或 URL 后缀判断是否为 PDF
func looksLikePdfPath(p string) bool {
	// URL 可能含 query/fragment，先取出纯路径部分
	if u, err := url.Parse(p); err == nil && u.Path != "" {
		p = u.Path
	}
	return strings.EqualFold(filepath.Ext(p), ".pdf")
}

// looksLikePdfBase64 根据 Base64 解码后的前 4 字节是否为 "%PDF" 判断
func looksLikePdfBase64(b64 string) bool {
	// 只解码足够前缀的字节即可
	if len(b64) < 8 {
		return false
	}
	// "%PDF" 编码后为 "JVBE..."
	return strings.HasPrefix(b64, "JVBE")
}

// isLikelyBase64 简单判断字符串是否像 Base64（仅含合法字符集，长度为 4 的倍数附近）
func isLikelyBase64(s string) bool {
	// 去除可能的换行
	s = strings.NewReplacer("\n", "", "\r", "", " ", "").Replace(s)
	if len(s) == 0 {
		return false
	}
	// 字符集校验
	for _, c := range s {
		switch {
		case c >= 'A' && c <= 'Z':
		case c >= 'a' && c <= 'z':
		case c >= '0' && c <= '9':
		case c == '+' || c == '/' || c == '=' || c == '-' || c == '_':
		default:
			return false
		}
	}
	return true
}

// buildInvoiceResult 把 SDK 返回的响应转为统一结构 InvoiceResult
func buildInvoiceResult(response *ocr.RecognizeGeneralInvoiceResponse) *InvoiceResult {
	result := &InvoiceResult{
		RawJSON: response.ToJsonString(),
	}
	if response.Response == nil {
		return result
	}
	if response.Response.RequestId != nil {
		result.RequestId = *response.Response.RequestId
	}
	items := response.Response.MixedInvoiceItems
	result.Invoices = make([]Invoice, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		result.Invoices = append(result.Invoices, parseInvoiceItem(item))
	}
	result.Count = len(result.Invoices)
	return result
}

// commonFieldAliases 通用字段抽取的别名表
//
// key 为 Invoice 的目标字段（如 "Buyer"），value 为子类型 JSON 中可能出现的字段名集合。
// 不同票种字段命名可能不同，按顺序取第一个非空值。
var commonFieldAliases = map[string][]string{
	"Code":        {"Code", "InvoiceCode", "TicketCode", "BillCode", "PaymentCode"},
	"Number":      {"Number", "InvoiceNumber", "TicketNumber", "BillNumber", "ID", "TripNum", "FlightNumber"},
	"Date":        {"Date", "InvoiceDate", "BillingDate", "TransferDate", "TripDate", "ConsumeDate"},
	"Buyer":       {"Buyer", "BuyerName", "Payer", "PassengerName", "Name"},
	"BuyerTaxID":  {"BuyerTaxID", "BuyerTaxNumber", "BuyerTaxId", "PayerTaxID"},
	"Seller":      {"Seller", "SellerName", "Receiver", "MerchantName"},
	"SellerTaxID": {"SellerTaxID", "SellerTaxNumber", "SellerTaxId", "ReceiverTaxID"},
	"Total":       {"Total", "TotalAmount", "Amount", "Fare", "Money"},
	"TotalCn":     {"TotalCn", "TotalCN", "TotalChinese", "AmountInWords"},
	"CheckCode":   {"CheckCode", "CheckNumber"},
	"Remark":      {"Remark", "Remarks", "Note"},
	"Title":       {"Title", "InvoiceTitle", "Name"},
}

// parseInvoiceItem 解析单张发票
func parseInvoiceItem(item *ocr.InvoiceItem) Invoice {
	inv := Invoice{
		OK:   strPtr(item.Code) == "OK",
		Kind: KindUnknown,
	}
	if !inv.OK {
		inv.ErrMsg = strPtr(item.Code)
	}
	if item.Page != nil {
		inv.PageIndex = int(*item.Page)
	}
	if item.Type != nil {
		inv.Kind = InvoiceKind(*item.Type)
	}
	if item.SubType != nil {
		inv.SubType = *item.SubType
	}
	if item.SubTypeDescription != nil {
		inv.SubTypeDesc = *item.SubTypeDescription
	}
	if item.TypeDescription != nil {
		inv.KindDescription = *item.TypeDescription
	} else if d, ok := invoiceKindDescriptions[inv.Kind]; ok {
		inv.KindDescription = d
	}
	if item.QRCode != nil {
		inv.QRCode = *item.QRCode
	}

	// 通过 JSON 反序列化 SingleInvoiceInfos 找到唯一非 nil 的子字段
	if item.SingleInvoiceInfos != nil {
		raw, _ := json.Marshal(item.SingleInvoiceInfos)
		var holder map[string]any
		if err := json.Unmarshal(raw, &holder); err == nil {
			if subKey, subVal := pickSubInvoice(holder); subKey != "" {
				if inv.SubType == "" {
					inv.SubType = subKey
				}
				if m, ok := subVal.(map[string]any); ok {
					inv.Raw = m
					fillCommonFields(&inv, m)
				}
			}
		}
	}

	// 解析金额为 float64
	if inv.Total != "" {
		s := strings.ReplaceAll(inv.Total, ",", "")
		s = strings.TrimSpace(s)
		s = strings.TrimPrefix(s, "￥")
		s = strings.TrimPrefix(s, "¥")
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			inv.TotalAmount = v
		}
	}

	return inv
}

// pickSubInvoice 在 SingleInvoiceInfos 反序列化结果中找到唯一非 nil 的子字段
//
// 腾讯云 SDK 中所有子字段都是指针（omitempty），未识别的子类型为 nil；
// JSON 反序列化后表现为该 key 不存在或值为 nil。
func pickSubInvoice(holder map[string]any) (string, any) {
	for k, v := range holder {
		if v == nil {
			continue
		}
		// map（对象）才视为有效子发票
		if _, ok := v.(map[string]any); ok {
			return k, v
		}
	}
	return "", nil
}

// fillCommonFields 把子发票 map 中的通用字段抽取到 Invoice 扁平字段
func fillCommonFields(inv *Invoice, m map[string]any) {
	get := func(keys []string) string {
		for _, k := range keys {
			if v, ok := m[k]; ok && v != nil {
				if s, ok := v.(string); ok && s != "" {
					return s
				}
				// 数字也允许
				if f, ok := v.(float64); ok {
					return strconv.FormatFloat(f, 'f', -1, 64)
				}
			}
		}
		return ""
	}
	if inv.Title == "" {
		inv.Title = get(commonFieldAliases["Title"])
	}
	inv.Code = get(commonFieldAliases["Code"])
	inv.Number = get(commonFieldAliases["Number"])
	inv.Date = get(commonFieldAliases["Date"])
	inv.Buyer = get(commonFieldAliases["Buyer"])
	inv.BuyerTaxID = get(commonFieldAliases["BuyerTaxID"])
	inv.Seller = get(commonFieldAliases["Seller"])
	inv.SellerTaxID = get(commonFieldAliases["SellerTaxID"])
	inv.Total = get(commonFieldAliases["Total"])
	inv.TotalCn = get(commonFieldAliases["TotalCn"])
	inv.CheckCode = get(commonFieldAliases["CheckCode"])
	inv.Remark = get(commonFieldAliases["Remark"])
}

// strPtr 安全解引用 *string
func strPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
