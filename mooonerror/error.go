// Package mooonerror
// Wrote by yijian on 2026/05/29
package mooonerror

// CError gomooon 通用错误类型，包含错误码和错误消息
// 错误码方便调用者程序化处理，错误消息方便人阅读
type CError struct {
ErrCode int    // 错误码，方便调用者程序化处理
ErrMsg  string // 错误消息，方便人阅读
}

// Error 实现 error 接口
func (e *CError) Error() string {
return e.ErrMsg
}

// NewError 创建一个 CError
func NewError(errCode int, errMsg string) *CError {
return &CError{ErrCode: errCode, ErrMsg: errMsg}
}
