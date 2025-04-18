// Package mooonpdf
// Wrote by yijian on 2024/10/24
package mooonpdf

import (
	"fmt"
	fitz "github.com/gen2brain/go-fitz"
	"image/jpeg"
	"image/png"
	"os"
)

type imageType int

const (
	pngImage imageType = iota
	jpgImage
)

func Pdf2Png(pdfFilepath string) ([]string, error) {
	return pdf2Image(pdfFilepath, pngImage)
}

func Pdf2Jpg(pdfFilepath string) ([]string, error) {
	return pdf2Image(pdfFilepath, jpgImage)
}

// pdf2Image 将 pdf 转为图片文件，每一页分别转成一个图片文件
// 返回值：图片文件数组，返回的图片文件名后缀格式为：.pdf.n.ext，其中 n 为从 0 开始的页号，ext 为对应的图片文件格式，如：png 或 jpg
func pdf2Image(pdfFilepath string, it imageType) ([]string, error) {
	doc, err := fitz.New(pdfFilepath)
	if err != nil {
		return nil, err
	}
	defer doc.Close()

	var ext string
	switch it {
	case pngImage:
		ext = "png"
	case jpgImage:
		ext = "jpg"
	default:
		return nil, fmt.Errorf("unsupported image type: %d", it)
	}

	//dirPath := filepath.Dir(pdfFilepath)
	//baseName := mooonutils.ExtractFilenameWithoutExtension(pdfFilepath)

	// Extract pages as images
	var imagePaths []string
	for i := 0; i < doc.NumPage(); i++ {
		img, err := doc.Image(i)
		if err != nil {
			return nil, err
		}

		imageFilepath := fmt.Sprintf("%s.%d.%s", pdfFilepath, i, ext)
		//imageFilepath := fmt.Sprintf("%s/%s_%d.%s", dirPath, baseName, i, ext)
		f, err := os.Create(imageFilepath)
		if err != nil {
			return nil, err
		}

		switch it {
		case pngImage:
			err = png.Encode(f, img)
		case jpgImage:
			err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		}
		if err != nil {
			return nil, err
		}
		f.Close()
		imagePaths = append(imagePaths, imageFilepath)
	}
	return imagePaths, nil
}