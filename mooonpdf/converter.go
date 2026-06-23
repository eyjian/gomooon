// Package mooonpdf
// Wrote by yijian on 2026/05/29
package mooonpdf

import "github.com/eyjian/gomooon/mooonerror"

// ImageFormat 图片输出格式
type ImageFormat string

const (
	ImageFormatPNG ImageFormat = "png"
	ImageFormatJPG ImageFormat = "jpg"
)

// Pdf2ImageOptions PDF 转图片的可选参数
type Pdf2ImageOptions struct {
	// DPI 分辨率（每英寸点数），值越大图片越清晰但文件越大，默认 300
	// 常用值参考：
	//   72  - 屏幕预览，文件最小，清晰度较低
	//   96  - 网页显示，Windows 标准屏幕 DPI
	//   150 - 普通打印和阅读，平衡清晰度与文件大小
	//   200 - 高质量打印，适合文档归档；电子发票归档推荐值
	//   300 - 印刷级精度，适合正式出版物；银行转账回单推荐值；票据/发票推荐值（默认值）
	//   600 - 极高精度，文件很大，通常仅用于专业印刷
	// 常见场景推荐：
	//   电子发票       → 200（满足财税归档清晰度要求，文件适中）
	//   银行转账回单   → 300（含小字号和印章，需较高清晰度确保可辨识）
	//   票据/发票      → 300（含小字号和印章，需较高清晰度确保可辨识）
	//   合同/协议      → 200~300（200 够用，有签章或小字可选 300）
	//   普通文档/报告  → 150~200（日常阅读和打印足够）
	// 取值范围建议：72~600
	DPI int

	// Format 输出格式，默认 PNG（支持透明背景）
	// 可选值：ImageFormatPNG（无损，文件较大）、ImageFormatJPG（有损压缩，文件较小）
	Format ImageFormat

	// UseCropBox 是否使用 CropBox 而非 MediaBox 来生成图片，默认 true
	//
	// PDF 页面有多个 Box 属性，决定了页面的大小和可见区域：
	//   MediaBox - 页面的物理大小（"画布"），类似画纸的整体尺寸
	//   CropBox  - 页面的裁剪区域（"画框"），即实际显示/打印的区域
	//   TrimBox  - 页面的最终裁切区域（"成品尺寸"），用于印刷
	//   BleedBox - 出血区域，印刷时超出 TrimBox 的部分
	//   ArtBox   - 艺术品区域，页面中含意义内容的区域
	//
	// 当 MediaBox 与 CropBox 不一致时（常见于票据、发票等）：
	//   UseCropBox=nil 或 true  → 使用 CropBox 渲染，只生成裁剪区域的图片，内容充满整张图
	//   UseCropBox=false        → 使用 MediaBox 渲染，生成整张画布的图片，内容可能只占一小部分
	//
	// 推荐设置：
	//   票据/发票/回单 → nil 或 true（内容通常只占 MediaBox 的一小部分，需裁剪到 CropBox）
	//   普通文档/报告 → nil 或 true（多数文档 MediaBox == CropBox，设不设置效果一样）
	//   需要保留完整页面白边 → 显式设为 false
	//
	// 使用指针类型以便区分"未设置"（nil，默认 true）和"显式设为 false"
	UseCropBox *bool
}

// Pdf2ImageConverter PDF 转图片的转换器接口
// 方便以后支持非 pdftoppm 方案，如 go-pdfium 等
type Pdf2ImageConverter interface {
	// Convert 将 PDF 的指定页转为图片
	// pdfPath: PDF 文件路径
	// outDir:  输出目录
	// pages:   页码列表，1-indexed，nil 或空表示全部页
	// options: 可选参数，nil 使用默认值
	// 返回值：生成的图片文件路径列表，和错误（类型为 *mooonerror.CError）
	Convert(pdfPath string, outDir string, pages []int, options *Pdf2ImageOptions) ([]string, *mooonerror.CError)
}
