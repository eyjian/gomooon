// Package txcloud
// Wrote by yijian on 2026/05/29
//
// 本文件定义发票识别能力的统一接口与请求/响应数据结构。
// 当前包内仅提供基于腾讯云 OCR `RecognizeGeneralInvoice` 的实现，
// 后续可在不破坏调用方代码的前提下，平滑新增其它实现（如新版接口或其它厂商）。
package txcloud

import (
	"context"
)

// InvoiceRecognizer 发票识别接口
//
// 调用方应当依赖该接口而非具体实现，以便后续平滑切换到其它识别方案。
// 例如：
//
//	var r txcloud.InvoiceRecognizer = txcloud.NewGeneralInvoiceRecognizer(secretId, secretKey, "")
//	result, err := r.RecognizeInvoice(ctx, &txcloud.InvoiceRequest{Source: "./bill.pdf"})
type InvoiceRecognizer interface {
	// RecognizeInvoice 识别一张图片或一份 PDF 中的发票，返回结构化结果
	//
	// 当腾讯云返回 *errors.TencentCloudSDKError 错误时，原样向上透传，
	// 调用方可使用 GetErrCodeAndErrMsg 提取错误码和错误描述。
	RecognizeInvoice(ctx context.Context, req *InvoiceRequest) (*InvoiceResult, error)
}

// InvoiceRequest 发票识别请求
//
// 推荐用法（极简）：仅设置 Source 字段即可，内部会自动判别输入是 URL、Base64 还是本地文件路径。
//
//	req := &txcloud.InvoiceRequest{Source: "./invoice.pdf"}
//	req := &txcloud.InvoiceRequest{Source: "https://x.com/a.jpg"}
//	req := &txcloud.InvoiceRequest{Source: base64Str}
//
// 高级用法：分别设置 ImageURL/ImageBase64/ImageFile 中之一，与 Source 互斥（高级字段优先生效）。
type InvoiceRequest struct {
	// === 极简字段（推荐使用）===

	// Source 输入源，自动按以下优先级判别：
	//   1) 以 "http://" 或 "https://" 开头     -> 当作 URL
	//   2) 以 "data:xxx;base64," 前缀开头     -> 剥离前缀后当作 Base64
	//   3) 长度较大且为合法 Base64 字符串    -> 当作 Base64
	//   4) 否则当作本地文件路径，自动 ReadFile + Base64 编码
	Source string

	// === 可选字段 ===

	// IsPdf 输入是否为 PDF
	// 不显式指定时，会根据以下规则自动检测：
	//   - 文件路径或 URL 含 ".pdf" 后缀
	//   - Base64 解码前 4 字节为 "%PDF"
	IsPdf bool

	// PdfPageNumber PDF 页码，仅当 IsPdf=true 生效，0 或负值会被自动置为 1
	PdfPageNumber int

	// === 高级字段（与 Source 互斥；若设置则优先使用）===

	// ImageURL 图片或 PDF 的 URL
	ImageURL string
	// ImageBase64 图片或 PDF 的 Base64
	// 含 "data:xxx;base64," 前缀时实现会自动剥离
	ImageBase64 string
	// ImageFile 本地文件路径，实现会自动读取并 Base64 编码
	ImageFile string

	// === 预留扩展 ===

	// Types 限定识别的票种（对应腾讯云 Types 参数），不填表示识别全部
	// 取值参考：
	//   0:出租车票 1:定额发票 2:火车票 3:增值税发票 5:机票行程单
	//   8:通用机打发票 9:汽车票 10:轮船票 11:增值税卷票 12:购车发票
	//   13:过路过桥费 15:非税发票 16:全电发票 17:医疗发票 18:完税凭证
	//   19:海关缴款书 20:银行回单 21:网约车行程单 22:海关报关单
	//   23:海外发票 24:购物小票 25:销货清单 -1:其它发票
	Types []int64
}

// InvoiceResult 发票识别结果
type InvoiceResult struct {
	RequestId string    // 腾讯云请求 ID
	Count     int       // 识别出的发票张数（= len(Invoices)）
	Invoices  []Invoice // 每张发票的结构化信息
	RawJSON   string    // 腾讯云原始响应 JSON（调试用，调用方一般不需要）
}

// InvoiceKind 发票大类（对应腾讯云 InvoiceItem.Type 字段）
//
// 提供枚举常量，方便业务代码进行分支判断，无需记忆数字含义。
type InvoiceKind int

const (
	KindUnknown          InvoiceKind = -1 // 未知类型 / 其它发票
	KindTaxiTicket       InvoiceKind = 0  // 出租车发票
	KindQuotaInvoice     InvoiceKind = 1  // 定额发票
	KindTrainTicket      InvoiceKind = 2  // 火车票
	KindVatInvoice       InvoiceKind = 3  // 增值税发票（专普票）
	KindAirTransport     InvoiceKind = 5  // 机票行程单
	KindMachinePrinted   InvoiceKind = 8  // 通用机打发票
	KindBusInvoice       InvoiceKind = 9  // 汽车票
	KindShippingInvoice  InvoiceKind = 10 // 轮船票
	KindVatInvoiceRoll   InvoiceKind = 11 // 增值税发票（卷票）
	KindMotorVehicle     InvoiceKind = 12 // 购车发票
	KindTollInvoice      InvoiceKind = 13 // 过路过桥费发票
	KindNonTaxIncomeBill InvoiceKind = 15 // 非税收入票据
	KindVatElectronic    InvoiceKind = 16 // 全电发票
	KindMedicalInvoice   InvoiceKind = 17 // 医疗发票
	KindTaxPayment       InvoiceKind = 18 // 完税凭证
	KindCustomsPayment   InvoiceKind = 19 // 海关缴款书
	KindBankSlip         InvoiceKind = 20 // 银行回单
	KindOnlineTaxi       InvoiceKind = 21 // 网约车行程单
	KindCustomsDeclare   InvoiceKind = 22 // 海关进/出口货物报关单
	KindOverseasInvoice  InvoiceKind = 23 // 海外发票
	KindShoppingReceipt  InvoiceKind = 24 // 购物小票
	KindSaleInventory    InvoiceKind = 25 // 销货清单
)

// Invoice 单张发票的结构化信息
//
// 通用字段尽量扁平化为基础类型（string/int/bool/float64），
// 调用方可直接通过字段访问。若需访问通用字段未覆盖的特殊数据，可访问 Raw。
type Invoice struct {
	// === 状态 ===

	OK        bool   // 该发票是否识别成功（腾讯云 Code == "OK" 时为 true）
	ErrMsg    string // 识别失败时的错误说明（来自腾讯云 Code 字段，非 OK 时填充）
	PageIndex int    // 在原文档中的页码（从 1 开始）

	// === 票种元信息 ===

	Kind            InvoiceKind // 票据大类（已转为枚举常量）
	KindDescription string      // 票据大类中文描述，如 "非税收入发票"
	SubType         string      // 票据子类型 KEY，如 "NonTaxIncomeGeneralBill"
	SubTypeDesc     string      // 票据子类型中文描述，如 "非税收入通用票据"
	Title           string      // 票据抬头/名称，如 "公益事业捐赠统一票据(电子)"

	// === 通用要素（扁平化，所有票种尽量填充）===

	Code        string  // 票据代码
	Number      string  // 票据号码 / 发票号
	Date        string  // 开票日期，原始字符串（如 "2024年08月02日"）
	Buyer       string  // 购买方/付款方
	BuyerTaxID  string  // 购买方税号
	Seller      string  // 销售方/收款方
	SellerTaxID string  // 销售方税号
	Total       string  // 金额（原始字符串，如 "1.00"）
	TotalAmount float64 // 金额（解析为 float64，便于直接计算；解析失败为 0）
	TotalCn     string  // 金额大写
	CheckCode   string  // 校验码
	Remark      string  // 备注
	QRCode      string  // 二维码内容（若识别到）

	// === 原始数据（高级用法）===

	// Raw 该张发票 SingleInvoiceInfos 中非 nil 子字段的原始内容（JSON 反序列化为 map）
	// 用于访问通用字段未覆盖的特殊数据，例如：
	//   - 火车票座位号、车次
	//   - 机票登机口、航班号
	//   - 医疗发票就诊科目
	//   - 银行回单收付双方账号
	Raw map[string]any
}

// IsVatInvoice 是否为增值税系列发票（专票/普票/卷票/全电）
func (i *Invoice) IsVatInvoice() bool {
	switch i.Kind {
	case KindVatInvoice, KindVatInvoiceRoll, KindVatElectronic:
		return true
	}
	return false
}

// IsTrainTicket 是否为火车票
func (i *Invoice) IsTrainTicket() bool { return i.Kind == KindTrainTicket }

// IsAirTransport 是否为机票行程单
func (i *Invoice) IsAirTransport() bool { return i.Kind == KindAirTransport }

// IsTaxiTicket 是否为出租车票
func (i *Invoice) IsTaxiTicket() bool { return i.Kind == KindTaxiTicket }

// IsBankSlip 是否为银行回单
func (i *Invoice) IsBankSlip() bool { return i.Kind == KindBankSlip }

// IsNonTaxIncomeBill 是否为非税收入票据
func (i *Invoice) IsNonTaxIncomeBill() bool { return i.Kind == KindNonTaxIncomeBill }

// IsMedicalInvoice 是否为医疗发票
func (i *Invoice) IsMedicalInvoice() bool { return i.Kind == KindMedicalInvoice }

// IsOnlineTaxi 是否为网约车行程单
func (i *Invoice) IsOnlineTaxi() bool { return i.Kind == KindOnlineTaxi }

// invoiceKindDescriptions 票据大类中文描述映射
// 当腾讯云未返回 TypeDescription 时使用本映射兜底
var invoiceKindDescriptions = map[InvoiceKind]string{
	KindUnknown:          "其它发票",
	KindTaxiTicket:       "出租车发票",
	KindQuotaInvoice:     "定额发票",
	KindTrainTicket:      "火车票",
	KindVatInvoice:       "增值税发票",
	KindAirTransport:     "机票行程单",
	KindMachinePrinted:   "通用机打发票",
	KindBusInvoice:       "汽车票",
	KindShippingInvoice:  "轮船票",
	KindVatInvoiceRoll:   "增值税发票（卷票）",
	KindMotorVehicle:     "购车发票",
	KindTollInvoice:      "过路过桥费发票",
	KindNonTaxIncomeBill: "非税收入发票",
	KindVatElectronic:    "全电发票",
	KindMedicalInvoice:   "医疗发票",
	KindTaxPayment:       "完税凭证",
	KindCustomsPayment:   "海关缴款书",
	KindBankSlip:         "银行回单",
	KindOnlineTaxi:       "网约车行程单",
	KindCustomsDeclare:   "海关报关单",
	KindOverseasInvoice:  "海外发票",
	KindShoppingReceipt:  "购物小票",
	KindSaleInventory:    "销货清单",
}
