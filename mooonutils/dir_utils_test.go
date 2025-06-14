// Package mooonutils
// Wrote by yijian on 2025/06/14
package mooonutils

import "testing"

// go test -v -run="TestPathExists"
func TestPathExists(t *testing.T) {
	// 不存在的
	dir := "acd"
	exists, isDir, err := PathExists(dir)
	if err != nil {
		t.Error(err)
	} else {
		if !exists {
			t.Log("not exists")
		} else if isDir {
			t.Log("dir exists")
		} else {
			t.Log("file exists")
		}
	}

	// 存在的目录
	dir = "."
	exists, isDir, err = PathExists(dir)
	if err != nil {
		t.Error(err)
	} else {
		if !exists {
			t.Log("not exists")
		} else if isDir {
			t.Log("dir exists")
		} else {
			t.Log("file exists")
		}
	}

	// 存在的文件
	dir = "dir_utils.go"
	exists, isDir, err = PathExists(dir)
	if err != nil {
		t.Error(err)
	} else {
		if !exists {
			t.Log("not exists")
		} else if isDir {
			t.Log("dir exists")
		} else {
			t.Log("file exists")
		}
	}
}
