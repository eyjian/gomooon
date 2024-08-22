// Package mooonstr
// Wrote by yijian on 2024/08/22
package mooonstr

import (
	"fmt"
	"io"
	"os"
)

func WriteString2File(filepath string, str string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("create %s error: %s", filepath, err.Error())
	}
	defer file.Close()
	return WriteString2Writer(file, str)
}

func WriteString2Writer(writer io.Writer, str string) error {
	_, err := writer.Write([]byte(str))
	return err
}