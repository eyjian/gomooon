// Package mooonpdf
// Wrote by yijian on 2026/05/29
package mooonpdf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/eyjian/gomooon/mooonerror"
	"github.com/eyjian/gomooon/mooonutils"
)

const pdftoppmCmdName = "pdftoppm"

// PdftoppmConverter 基于 pdftoppm 命令行的 PDF 转图片实现
type PdftoppmConverter struct{}

// NewPdftoppmConverter 创建一个 PdftoppmConverter
func NewPdftoppmConverter() *PdftoppmConverter {
	return &PdftoppmConverter{}
}

// Convert 将 PDF 的指定页转为图片
// pdfPath: PDF 文件路径
// outDir:  输出目录
// pages:   页码列表，1-indexed，nil 或空表示全部页
// options: 可选参数，nil 使用默认值
func (c *PdftoppmConverter) Convert(
	pdfPath string, outDir string, pages []int, options *Pdf2ImageOptions,
) ([]string, *mooonerror.CError) {
	// 1. 检查 pdftoppm 是否可用
	if _, err := exec.LookPath(pdftoppmCmdName); err != nil {
		var installHint string
		if mooonutils.IsLinux() {
			installHint = getLinuxInstallHint()
			installHint = "\nInstall command: " + installHint
		}
		return nil, mooonerror.NewError(mooonerror.ErrCodeToolNotFound,
			fmt.Sprintf("pdftoppm not available, please install poppler-utils: %s%s", err.Error(), installHint))
	}

	// 2. 检查 PDF 文件存在
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return nil, mooonerror.NewError(mooonerror.ErrCodeFileNotFound,
			fmt.Sprintf("PDF file not found: %s", pdfPath))
	}

	// 3. 检查输出目录
	exists, isDir, _ := mooonutils.PathExists(outDir)
	if !exists {
		return nil, mooonerror.NewError(mooonerror.ErrCodeFileNotFound,
			fmt.Sprintf("output directory not found: %s", outDir))
	}
	if !isDir {
		return nil, mooonerror.NewError(mooonerror.ErrCodeFileOperate,
			fmt.Sprintf("output path is not a directory: %s", outDir))
	}

	// 4. 填充默认参数
	if options == nil {
		options = &Pdf2ImageOptions{}
	}
	if options.DPI <= 0 {
		options.DPI = 150
	}
	if options.Format == "" {
		options.Format = ImageFormatPNG
	}

	// 5. 构造输出文件名前缀
	// pdftoppm 输出格式：前缀-1.png, 前缀-2.png, ...
	// 文件名为原pdf文件名（含.pdf后缀）作为前缀，减少冲突
	pdfBaseWithExt := filepath.Base(pdfPath) // 如 "test.pdf"
	outputPrefix := filepath.Join(outDir, pdfBaseWithExt)

	// 6. 构造命令参数
	args := []string{}
	args = append(args, fmt.Sprintf("-%s", string(options.Format))) // -png 或 -jpg
	args = append(args, "-r", strconv.Itoa(options.DPI))            // -r 150

	// 指定页码范围
	if len(pages) > 0 {
		minPage, maxPage := pages[0], pages[0]
		for _, p := range pages {
			if p < minPage {
				minPage = p
			}
			if p > maxPage {
				maxPage = p
			}
		}
		args = append(args, "-f", strconv.Itoa(minPage))
		args = append(args, "-l", strconv.Itoa(maxPage))
	}

	args = append(args, pdfPath)
	args = append(args, outputPrefix)

	// 7. 执行命令
	cmd := exec.Command(pdftoppmCmdName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, mooonerror.NewError(mooonerror.ErrCodeToolExecute,
			fmt.Sprintf("pdftoppm execution failed: %s, output: %s", err.Error(), string(output)))
	}

	// 8. 收集生成的图片文件
	ext := "." + string(options.Format)
	resultFiles, err2 := mooonutils.GetFilesBySuffix(outDir, []string{ext})
	if err2 != nil {
		return nil, mooonerror.NewError(mooonerror.ErrCodeFileOperate,
			fmt.Sprintf("failed to list output files: %s", err2.Error()))
	}

	// 过滤：只返回本次生成的文件（以输出前缀开头）
	var filteredFiles []string
	prefix := filepath.Join(outDir, pdfBaseWithExt)
	for _, f := range resultFiles {
		if strings.HasPrefix(f, prefix) {
			filteredFiles = append(filteredFiles, f)
		}
	}

	// 如果指定了特定页码，还需过滤只保留请求的页
	if len(pages) > 0 {
		pageSet := make(map[int]bool)
		for _, p := range pages {
			pageSet[p] = true
		}
		var selectedFiles []string
		for _, f := range filteredFiles {
			page := extractPageNumber(f, pdfBaseWithExt, ext)
			if pageSet[page] {
				selectedFiles = append(selectedFiles, f)
			}
		}
		filteredFiles = selectedFiles
	}

	// 排序
	sort.Strings(filteredFiles)

	return filteredFiles, nil
}

// extractPageNumber 从 pdftoppm 生成的文件名中提取页码
// 文件名格式：test.pdf-1.png，其中 pdfBaseWithExt 为 "test.pdf"，ext 为 ".png"
func extractPageNumber(filePath string, pdfBaseWithExt string, ext string) int {
	// filePath 形如 /tmp/output/test.pdf-1.png
	// 去掉目录和后缀
	base := filepath.Base(filePath)      // test.pdf-1.png
	base = strings.TrimSuffix(base, ext) // test.pdf-1
	prefix := pdfBaseWithExt + "-"       // test.pdf-
	if !strings.HasPrefix(base, prefix) {
		return -1
	}
	pageStr := strings.TrimPrefix(base, prefix) // 1
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return -1
	}
	return page
}

// getLinuxInstallHint 返回 Linux 各发行版安装 poppler-utils 的命令提示
// 通过检测已知包管理器来判断发行版，返回对应的安装命令
func getLinuxInstallHint() string {
	// 常见 Linux 发行版的安装命令
	type distroHint struct {
		detectFile string // 包管理器或发行版标识文件
		installCmd string // 安装命令
	}
	hints := []distroHint{
		{"/usr/bin/dnf", "dnf install poppler-utils"},         // Fedora / RHEL 8+ / CentOS 8+
		{"/usr/bin/yum", "yum install poppler-utils"},         // RHEL 7 / CentOS 7
		{"/usr/bin/apt-get", "apt-get install poppler-utils"}, // Debian / Ubuntu
		{"/usr/bin/zypper", "zypper install poppler-tools"},   // openSUSE（包名不同）
		{"/usr/bin/pacman", "pacman -S poppler"},              // Arch Linux（包名不同）
		{"/usr/bin/apk", "apk add poppler-utils"},             // Alpine Linux
		{"/usr/bin/emerge", "emerge app-text/poppler"},        // Gentoo
	}

	for _, h := range hints {
		if _, err := os.Stat(h.detectFile); err == nil {
			return h.installCmd
		}
	}

	// 无法识别发行版时，给出常见发行版的安装命令
	return "apt-get install poppler-utils (Debian/Ubuntu) | dnf install poppler-utils (Fedora/RHEL) | yum install poppler-utils (CentOS) | apk add poppler-utils (Alpine)"
}
