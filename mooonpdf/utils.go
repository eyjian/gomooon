// Package mooonpdf
// Wrote by yijian on 2024/12/12
package mooonpdf

import (
	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
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

	err := pdfapi.ValidateFile(filepath, pdfmodel.NewDefaultConfiguration())
	return err == nil
}

// GetPdfPageCount 获取 pdf 文件页数
func GetPdfPageCount(filepath string) (int, error) {
	return pdfapi.PageCountFile(filepath) // 也可通过 ReadContextFile 间接获得页数
}

func ValidatePdf(filepath string) error {
	config := pdfmodel.NewDefaultConfiguration()
	return pdfapi.ValidateFile(filepath, config)
}

func OptimizePdf(filepath, newFilepath string) error {
	config := pdfmodel.NewDefaultConfiguration()
	return pdfapi.OptimizeFile(filepath, newFilepath, config)
}