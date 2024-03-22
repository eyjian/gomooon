// Package mooonhttp
// Wrote by yijian on 2024/03/09
package mooonhttp

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

// DownloadFile 通过 http 下载文件，本实现基于 AI 生成
// url 文件的链接
// localFilepath 本地文件路径
// 第一个返回值为 http 的响应代码，如果其值为 0 表示还没取得 http 的响应代码，是否出错应看第二个返回值是否为 nil
func DownloadFile(url, localFilepath string) (int, error) {
    // 创建一个新的 HTTP 客户端
    client := &http.Client{}

    // 发送 GET 请求
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return 0, fmt.Errorf("create %s error: %s", url, err.Error())
    }

    resp, err := client.Do(req)
    if err != nil {
        return 0, fmt.Errorf("request %s error: %s", url, err.Error())
    }
    defer resp.Body.Close()

    // 检查 HTTP 响应状态码
    if resp.StatusCode != http.StatusOK {
        return resp.StatusCode, fmt.Errorf("HTTP request error: %s", resp.Status)
    }

    // 创建一个新的文件
    file, err := os.Create(localFilepath)
    if err != nil {
        return 0, fmt.Errorf("create file://%s error: %s", localFilepath, err.Error())
    }
    defer file.Close()

    // 将响应体中的数据写入文件
    _, err = io.Copy(file, resp.Body) // 这里易遇到网络错误：unexpected EOF
    if err != nil {
        return resp.StatusCode, fmt.Errorf("write file://%s error: %s", localFilepath, err.Error())
    }

    return resp.StatusCode, nil
}
