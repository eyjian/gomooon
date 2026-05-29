// Package mooonerror
// Wrote by yijian on 2026/05/29
package mooonerror

// 通用错误码（1-99）
const (
	ErrCodeSuccess      = 0 // 成功
	ErrCodeUnknown      = 1 // 未知错误
	ErrCodeInvalidParam = 2 // 参数错误
	ErrCodeFileNotFound = 3 // 文件不存在
	ErrCodeFileOperate  = 4 // 文件操作失败
	ErrCodeToolNotFound = 5 // 依赖工具不可用
	ErrCodeToolExecute  = 6 // 工具执行失败
)

// PDF 相关错误码（100-199）
const (
	ErrCodePdfInvalid   = 100 // 无效的 PDF 文件
	ErrCodePdfPageRange = 101 // 页码范围无效
	ErrCodePdfConvert   = 102 // PDF 转换失败
)
