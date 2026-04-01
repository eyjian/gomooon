// Package mooonpdf
// Wrote by yijian on 2024/10/24
package mooonpdf

import (
	"fmt"
	"github.com/eyjian/gomooon/mooonutils"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

// Jpg2Png 将 jpg 文件转为 png 文件
func Jpg2Png(jpgFilepath string) (string, error) {
	dirPath := filepath.Dir(jpgFilepath)
	baseName := mooonutils.ExtractFilenameWithoutExtension(jpgFilepath)
	pngFilepath := fmt.Sprintf("%s/%s.png", dirPath, baseName)

	// 打开 JPEG 文件
	file, err := os.Open(jpgFilepath)
	if err != nil {
		return "", fmt.Errorf("open JPEG %s error: %s", jpgFilepath, err.Error())
	}
	defer file.Close()

	// 解码 JPEG 文件
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("decode JPEG %s error: %s", jpgFilepath, err.Error())
	}

	// 创建 PNG 文件
	outFile, err := os.Create(pngFilepath)
	if err != nil {
		return "", fmt.Errorf("create PNG %s error: %s", pngFilepath, err.Error())
	}
	defer outFile.Close()

	// 编码 PNG 文件
	err = png.Encode(outFile, img)
	if err != nil {
		return "", fmt.Errorf("encode PNG %s error: %s", pngFilepath, err.Error())
	}

	return pngFilepath, nil
}
