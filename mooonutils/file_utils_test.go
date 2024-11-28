// Package mooonutils
// Wrote by yijian on 2024/10/22
package mooonutils

import (
	"os"
	"testing"
)

// go test -v -run="TestMd5File"
func TestMd5File(t *testing.T) {
	md5Str, err := Md5File("str_utils.go")
	if err != nil {
		t.Errorf("md5 error: %s]\n", err.Error())
	} else {
		t.Logf("md5 ok: %s\n", md5Str)
	}
}

// go test -v -run="TestDeleteFile"
func TestDeleteFile(t *testing.T) {
	err := DeleteFile("str_utils.copy")
	if err != nil {
		t.Errorf("delete %s error: %s\n", "str_utils_copy.go", err.Error())
	} else {
		t.Logf("delete %s ok\n", "str_utils_copy.go")
	}

	err = DeleteFile("str_utils.copy")
	if err != nil {
		t.Errorf("delete %s error: %s\n", "str_utils_copy.go", err.Error())
	} else {
		t.Logf("delete %s ok\n", "str_utils_copy.go")
	}
}

// go test -v -run="TestCopyFile"
func TestCopyFile(t *testing.T) {
	srcFile := "str_utils.go"
	dstFile := "str_utils.copy"
	err := CopyFile(srcFile, dstFile, true)
	if err != nil {
		t.Errorf("copy %s to %s error: %s\n", srcFile, dstFile, err.Error())
	} else {
		t.Logf("copy %s to %s ok\n", srcFile, dstFile)
	}

	err = CopyFile(srcFile, dstFile, false)
	if err != nil {
		t.Errorf("copy %s to %s error: %s\n", srcFile, dstFile, err.Error())
	} else {
		t.Logf("copy %s to %s ok\n", srcFile, dstFile)
	}
}

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