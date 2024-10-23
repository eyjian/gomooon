// Package mooonutils
// Wrote by yijian on 2024/10/22
package mooonutils

import (
	"os"
	"testing"
)

// go test -v -run="TestUnzip"
func TestUnzip(t *testing.T) {
	destDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Getwd error: %s\n", err.Error())
		return
	}

	// 文件名不含中文
	zipFile := "Test_20231023.zip"
	paths, err := Unzip(zipFile, destDir)
	if err != nil {
		t.Errorf("unzip %s error: %s\n", zipFile, err.Error())
	} else {
		t.Logf("unzip %s ok:\n", zipFile)
		for _, path := range paths {
			t.Logf("=>%s\n", path)
		}
	}

	// 文件名含中文
	zipFile = "测试_文本文档_20241023.zip"
	paths, err = Unzip(zipFile, destDir)
	if err != nil {
		t.Errorf("unzip %s error: %s\n", zipFile, err.Error())
	} else {
		t.Logf("unzip %s ok:\n", zipFile)
		for _, path := range paths {
			t.Logf("=>%s\n", path)
		}
	}

	// 压缩包中含目录
	zipFile = "测试_20241023.zip"
	paths, err = Unzip(zipFile, destDir, true, true)
	if err != nil {
		t.Errorf("unzip %s error: %s\n", zipFile, err.Error())
	} else {
		t.Logf("unzip %s ok:\n", zipFile)
		for _, path := range paths {
			t.Logf("=>%s\n", path)
		}
	}
}
