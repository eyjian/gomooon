// Package mooonpdf
// Wrote by yijian on 2024/12/12
package mooonpdf

import (
	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"strings"
)

// IsPdfFile 判断文件是否为 pdf 文件
func IsPdfFile(filepath string) bool {
	if filepath == "" {
		return false
	}
	if len(filepath) <= 4 {
		return false
	}
	if strings.ToLower(filepath[len(filepath)-4:]) != ".pdf" {
		return false
	}

	// 加载 PDF 上下文
	_, err := pdfapi.ReadContextFile(filepath)
	return err == nil
}

// GetPdfPageCount 获取 pdf 文件页数
func GetPdfPageCount(filepath string) (int, error) {
	modelCtx, err := pdfapi.ReadContextFile(filepath)
	if err != nil {
		return 0, err
	}
	return modelCtx.PageCount, nil
}
