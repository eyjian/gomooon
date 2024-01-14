// Package main
// Wrote by yijian on 2024/01/14
package main

import (
	"flag"
	"fmt"
	"github.com/eyjian/gomooon/utils"
	"os"
)

var (
	help = flag.Bool("h", false, "Display a help message and exit.")
	key  = flag.String("key", "", "Key to sign.")
	str  = flag.String("str", "", "String to been signed.")
)

func main() {
	flag.Parse()
	if *help {
		usage()
		os.Exit(1)
	}
	if len(*key) == 0 {
		fmt.Println("Parameter[--key] is not set.")
		usage()
		os.Exit(1)
	}
	if len(*str) == 0 {
		fmt.Println("Parameter[--str] is not set.")
		usage()
		os.Exit(1)
	}

	signedStr, err := utils.UpperHmacSHA256Sign(*str, *key)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Signed: %s\n", signedStr)
	}
}

// 显示使用帮助函数
func usage() {
	flag.Usage()
}
