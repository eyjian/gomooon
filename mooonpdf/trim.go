// Package mooonpdf
// Wrote by yijian on 2025/06/14
package mooonpdf

import (
	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// TrimFile 生成包含所有选定页面的 inFile 的修剪版本
// 如果 outFile 为空则直接修改 inFile
func TrimFile(inFile, outFile string, selectedPages []string) error {
	return pdfapi.TrimFile(inFile, outFile, selectedPages, pdfmodel.NewDefaultConfiguration())
}
