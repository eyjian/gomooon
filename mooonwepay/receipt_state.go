// Package mooonwepay
// Wrote by yijian on 2025/03/27
package mooonwepay

const (
	ReceiptStateGenerating = "GENERATING" // 电子回单已受理成功并在处理中
	ReceiptStateFinished   = "FINISHED"   // 电子回单已处理完成
	ReceiptStateFailed     = "FAILED"     // 电子回单生成失败，失败原因字段会返回具体的失败原因
)

const (
	HashTypeSM3    = "SM3"    // 国密SM3摘要算法
	HashTypeSHA256 = "SHA256" // SHA256摘要算法
)