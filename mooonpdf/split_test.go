// Package mooonpdf
// Wrote by yijian on 2025/06/14
package mooonpdf

import (
	"github.com/eyjian/gomooon/mooonutils"
	"log"
	"os"
	"testing"
)

// go test -v -run="TestSplitFile"
func TestSplitFile(t *testing.T) {
	outDir := "out_dir"

	// 如果目录已存在则删除
	exists, isDir, err := mooonutils.PathExists(outDir)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		if !isDir {
			log.Fatalf("`%s` is not a directory", outDir)
		}
		err = os.RemoveAll(outDir)
		if err != nil {
			t.Error(err)
		}
	}

	// 创建目录
	err = os.Mkdir(outDir, 0755) // 权限位：rwxr-xr-x
	if err != nil {
		t.Error(err)
		return
	}

	// 将 pdf 按页拆分为一个个 pdf 文件
	err = SplitFile("test.pdf", outDir, 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
