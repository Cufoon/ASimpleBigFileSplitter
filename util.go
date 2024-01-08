package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

func checkDestDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				fmt.Println("无法创建目标目录", path)
				return false
			}
			return true
		}
		fmt.Println(path, "目录存在异常，请删除这个文件夹，重新创建一下。")
		return false
	}
	if s.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(path, "目录存在异常，请删除这个文件夹，重新创建一下。")
			return false
		}
		_, err = f.Readdir(1)
		if err == io.EOF {
			return true
		}
		fmt.Println(path, "目录是一个非空目录，请重新指定")
		return false
	}
	fmt.Println("存在一个文件名是", path, "的文件，请删除它或保存到其他地方")
	return false
}

var preNumber = []int64{9, 99, 999, 9999, 99999, 999999, 9999999, 99999999, 999999999, 9999999999, 99999999999, 999999999999, 9999999999999, 99999999999999, 999999999999999, 9999999999999999, 99999999999999999, 999999999999999999}

func getNumWidth(num int64) int {
	if num == math.MinInt64 {
		return 19
	}
	if num < 0 {
		num = -num
	}
	for i := 0; i < 18; i++ {
		if num <= preNumber[i] {
			return i + 1
		}
	}
	return 19
}
