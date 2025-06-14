// Package mooonpdf
// Wrote by yijian on 2025/06/14
package mooonpdf

import (
	"github.com/eyjian/gomooon/mooonutils"
	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// SplitFile 将 pdf 文件的每页都拆分为一个 pdf 文件
// outDir 存放拆分后新生成的 pdf 文件目录，应确保为一个空的目录
// span 拆分跨度，1 表示每页拆分为一个 pdf 文件，2 表示每 2 页拆分为一个 pdf 文件，依次类推。但 0 表示书签进行拆分。
func SplitFile(inFile, outDir string, span int) ([]string, error) {
	config := pdfmodel.NewDefaultConfiguration()
	err := pdfapi.SplitFile(inFile, outDir, span, config)
	if err != nil {
		return nil, err
	}
	return mooonutils.GetFilesBySuffix(outDir, []string{".pdf"})
}
