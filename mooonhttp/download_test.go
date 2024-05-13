// Package mooonhttp
// Wrote by yijian on 2024/03/09
package mooonhttp

import (
    "testing"
)

// go test -v -run="TestDownloadFile" #-args url
func TestDownloadFile(t *testing.T) {
    url := "https://raw.githubusercontent.com/eyjian/mooon-district/main/example.csv"
    t.Logf("%s\n", url)

    localFilepath := "mooon.download"
    httpStatusCode, err := DownloadFile(url, localFilepath)
    if err != nil {
        t.Errorf("downlaod %s error: (%d) %s\n", url, httpStatusCode, err.Error())
    } else {
        t.Logf("download %s to file://%s ok\n", url, localFilepath)
    }
}
