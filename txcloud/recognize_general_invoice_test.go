// Package txcloud
// Wrote by yijian on 2026/05/29
package txcloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

// ----- resolveSource 测试 -----

func TestResolveSource_URL(t *testing.T) {
	url, b64, isPdf, err := resolveSource(&InvoiceRequest{Source: "https://example.com/a.pdf"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if url != "https://example.com/a.pdf" || b64 != "" {
		t.Fatalf("unexpected url=%q b64=%q", url, b64)
	}
	if !isPdf {
		t.Fatalf("expected isPdf=true for .pdf URL")
	}
}

func TestResolveSource_Base64WithPrefix(t *testing.T) {
	raw := []byte("%PDF-1.7\nhello")
	b64 := base64.StdEncoding.EncodeToString(raw)
	src := "data:application/pdf;base64," + b64
	url, gotB64, isPdf, err := resolveSource(&InvoiceRequest{Source: src})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if url != "" || gotB64 != b64 {
		t.Fatalf("expected pure b64 returned, got url=%q b64=%q", url, gotB64)
	}
	if !isPdf {
		t.Fatalf("expected isPdf=true")
	}
}

func TestResolveSource_LocalFile(t *testing.T) {
	tmpDir := t.TempDir()
	pdfPath := filepath.Join(tmpDir, "test.pdf")
	if err := os.WriteFile(pdfPath, []byte("%PDF-1.7\nfake"), 0644); err != nil {
		t.Fatalf("write tmp file: %v", err)
	}
	url, b64, isPdf, err := resolveSource(&InvoiceRequest{Source: pdfPath})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if url != "" || b64 == "" {
		t.Fatalf("expected b64 from file, got url=%q b64=%q", url, b64)
	}
	if !isPdf {
		t.Fatalf("expected isPdf=true for .pdf file")
	}
}

func TestResolveSource_AdvancedFieldsPriority(t *testing.T) {
	// 同时提供 Source 与 ImageURL，应优先使用 ImageURL
	url, _, _, err := resolveSource(&InvoiceRequest{
		Source:   "./should-not-be-used.jpg",
		ImageURL: "https://x.com/a.jpg",
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if url != "https://x.com/a.jpg" {
		t.Fatalf("expected ImageURL to win, got %q", url)
	}
}

func TestResolveSource_AllEmpty(t *testing.T) {
	_, _, _, err := resolveSource(&InvoiceRequest{})
	if err == nil {
		t.Fatalf("expected error for empty request")
	}
}

func TestResolveSource_ImageBase64StripPrefix(t *testing.T) {
	raw := []byte("%PDF-1.7\n")
	b64 := base64.StdEncoding.EncodeToString(raw)
	_, gotB64, isPdf, err := resolveSource(&InvoiceRequest{
		ImageBase64: "data:image/jpeg;base64," + b64,
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if gotB64 != b64 {
		t.Fatalf("prefix not stripped: %q", gotB64)
	}
	// 文件头是 PDF，应自动检测为 PDF
	if !isPdf {
		t.Fatalf("expected isPdf=true based on PDF magic bytes")
	}
}

// ----- looksLikePdfPath / looksLikePdfBase64 -----

func TestLooksLikePdfPath(t *testing.T) {
	cases := map[string]bool{
		"a.pdf":                         true,
		"a.PDF":                         true,
		"a.jpg":                         false,
		"https://x.com/foo/a.pdf?x=1":   true,
		"https://x.com/a.jpg":           false,
		"/var/data/2024-08-02-bill.pdf": true,
	}
	for in, want := range cases {
		if got := looksLikePdfPath(in); got != want {
			t.Errorf("looksLikePdfPath(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestLooksLikePdfBase64(t *testing.T) {
	pdf := base64.StdEncoding.EncodeToString([]byte("%PDF-1.7\nx"))
	if !looksLikePdfBase64(pdf) {
		t.Errorf("expected pdf magic recognized")
	}
	jpg := base64.StdEncoding.EncodeToString([]byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00})
	if looksLikePdfBase64(jpg) {
		t.Errorf("jpg should not be recognized as pdf")
	}
}

// ----- parseInvoiceItem 基于真实响应样例 -----

// nonTaxIncomeBillJSON 取自 tmp/response.txt 中 NonTaxIncomeGeneralBill 字段
const nonTaxIncomeBillJSON = `{
    "Buyer": "陈杏红",
    "BuyerTaxID": "******",
    "CheckCode": "bngFFX",
    "Code": "32050124",
    "Date": "2024年08月02日",
    "Number": "0000817610",
    "Remark": "2019年捐赠",
    "Reviewer": "焦菲",
    "Seller": "孙琳",
    "Title": "公益事业捐赠统一票据(电子)",
    "Total": "1.00",
    "TotalCn": "壹元整"
}`

func TestParseInvoiceItem_NonTaxIncomeBill(t *testing.T) {
	// 构造 SingleInvoiceItem，仅 NonTaxIncomeGeneralBill 非 nil
	// 该字段在 SDK 中类型为 *NonTaxIncomeBill
	var bill ocr.NonTaxIncomeBill
	if err := json.Unmarshal([]byte(nonTaxIncomeBillJSON), &bill); err != nil {
		t.Fatalf("unmarshal sample: %v", err)
	}
	item := &ocr.InvoiceItem{
		Code:               common.StringPtr("OK"),
		Type:               common.Int64Ptr(15),
		Page:               common.Int64Ptr(1),
		SubType:            common.StringPtr("NonTaxIncomeGeneralBill"),
		SubTypeDescription: common.StringPtr("非税收入通用票据"),
		TypeDescription:    common.StringPtr("非税收入发票"),
		SingleInvoiceInfos: &ocr.SingleInvoiceItem{
			NonTaxIncomeGeneralBill: &bill,
		},
	}

	inv := parseInvoiceItem(item)
	if !inv.OK {
		t.Fatalf("expected OK=true")
	}
	if inv.Kind != KindNonTaxIncomeBill {
		t.Errorf("Kind=%v, want %v", inv.Kind, KindNonTaxIncomeBill)
	}
	if inv.SubType != "NonTaxIncomeGeneralBill" {
		t.Errorf("SubType=%q", inv.SubType)
	}
	if inv.SubTypeDesc != "非税收入通用票据" {
		t.Errorf("SubTypeDesc=%q", inv.SubTypeDesc)
	}
	if inv.KindDescription != "非税收入发票" {
		t.Errorf("KindDescription=%q", inv.KindDescription)
	}
	if inv.Code != "32050124" {
		t.Errorf("Code=%q", inv.Code)
	}
	if inv.Number != "0000817610" {
		t.Errorf("Number=%q", inv.Number)
	}
	if inv.Date != "2024年08月02日" {
		t.Errorf("Date=%q", inv.Date)
	}
	if inv.Buyer != "陈杏红" {
		t.Errorf("Buyer=%q", inv.Buyer)
	}
	if inv.Seller != "孙琳" {
		t.Errorf("Seller=%q", inv.Seller)
	}
	if inv.Total != "1.00" {
		t.Errorf("Total=%q", inv.Total)
	}
	if inv.TotalAmount != 1.00 {
		t.Errorf("TotalAmount=%v", inv.TotalAmount)
	}
	if inv.TotalCn != "壹元整" {
		t.Errorf("TotalCn=%q", inv.TotalCn)
	}
	if inv.CheckCode != "bngFFX" {
		t.Errorf("CheckCode=%q", inv.CheckCode)
	}
	if inv.Remark != "2019年捐赠" {
		t.Errorf("Remark=%q", inv.Remark)
	}
	if inv.Title != "公益事业捐赠统一票据(电子)" {
		t.Errorf("Title=%q", inv.Title)
	}
	if inv.PageIndex != 1 {
		t.Errorf("PageIndex=%d", inv.PageIndex)
	}
	if inv.Raw == nil {
		t.Errorf("Raw should be non-nil")
	}
	// 通用字段未覆盖的 Reviewer 仍可在 Raw 中访问
	if v, _ := inv.Raw["Reviewer"].(string); v != "焦菲" {
		t.Errorf("Raw[Reviewer]=%v", inv.Raw["Reviewer"])
	}
}

func TestParseInvoiceItem_NotOK(t *testing.T) {
	item := &ocr.InvoiceItem{
		Code: common.StringPtr("FailedOperation.UnsupportedInvoice"),
		Type: common.Int64Ptr(-1),
	}
	inv := parseInvoiceItem(item)
	if inv.OK {
		t.Fatalf("expected OK=false")
	}
	if inv.ErrMsg != "FailedOperation.UnsupportedInvoice" {
		t.Errorf("ErrMsg=%q", inv.ErrMsg)
	}
}

// ----- 辅助方法 -----

func TestInvoiceKindHelpers(t *testing.T) {
	cases := []struct {
		kind  InvoiceKind
		check func(*Invoice) bool
		want  bool
	}{
		{KindVatInvoice, (*Invoice).IsVatInvoice, true},
		{KindVatElectronic, (*Invoice).IsVatInvoice, true},
		{KindVatInvoiceRoll, (*Invoice).IsVatInvoice, true},
		{KindTrainTicket, (*Invoice).IsVatInvoice, false},
		{KindTrainTicket, (*Invoice).IsTrainTicket, true},
		{KindBankSlip, (*Invoice).IsBankSlip, true},
		{KindNonTaxIncomeBill, (*Invoice).IsNonTaxIncomeBill, true},
		{KindAirTransport, (*Invoice).IsAirTransport, true},
		{KindTaxiTicket, (*Invoice).IsTaxiTicket, true},
		{KindOnlineTaxi, (*Invoice).IsOnlineTaxi, true},
		{KindMedicalInvoice, (*Invoice).IsMedicalInvoice, true},
	}
	for i, c := range cases {
		inv := &Invoice{Kind: c.kind}
		if got := c.check(inv); got != c.want {
			t.Errorf("case %d kind=%v got=%v want=%v", i, c.kind, got, c.want)
		}
	}
}

// ----- 端到端测试（默认 skip，仅设置环境变量时运行）-----
//
// 设置以下环境变量启用：
//
//	TENCENTCLOUD_SECRET_ID
//	TENCENTCLOUD_SECRET_KEY
//	TXCLOUD_INVOICE_TEST_FILE  本地图片或 PDF 路径
func TestRecognizeInvoice_E2E(t *testing.T) {
	id := os.Getenv("TENCENTCLOUD_SECRET_ID")
	key := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	file := os.Getenv("TXCLOUD_INVOICE_TEST_FILE")
	t.Log("\nid:", id, "\nkey:", key, "\nfile:", file)
	if id == "" || key == "" || file == "" {
		t.Skip("skip E2E test: env not set")
	}
	r := NewGeneralInvoiceRecognizer(id, key, "")
	result, err := r.RecognizeInvoice(context.Background(), &InvoiceRequest{Source: file})
	if err != nil {
		t.Fatalf("recognize: %v", err)
	}
	t.Logf("RequestId=%s Count=%d", result.RequestId, result.Count)
	for i, inv := range result.Invoices {
		t.Logf("[%d] OK=%v Kind=%v(%s) SubType=%s Number=%s Date=%s Total=%s(%.2f) Buyer=%s Seller=%s",
			i, inv.OK, inv.Kind, inv.KindDescription, inv.SubType,
			inv.Number, inv.Date, inv.Total, inv.TotalAmount, inv.Buyer, inv.Seller)
	}
}
