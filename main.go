package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"slices"
)

// 512 MiB
const chunkSize = 536870912

var (
	action  string
	infile  string
	outfile string
)

func split(infile string) {
	if infile == "" {
		panic("请输入正确的文件名")
	}

	fileInfo, err := os.Stat(infile)
	if err != nil {
		if os.IsNotExist(err) {
			panic("文件不存在")
		}
		panic(err)
	}

	num := int64(math.Ceil(float64(fileInfo.Size()) / chunkSize))

	fi, err := os.OpenFile(infile, os.O_RDONLY, os.ModePerm)
	defer func(fi *os.File) {
		err := fi.Close()
		if err != nil {
			fmt.Println("文件关闭失败")
		}
	}(fi)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("要拆分成%d份\n", num)
	b := make([]byte, chunkSize)
	fileLen := fileInfo.Size()
	var i int64 = 1
	for ; i <= num; i++ {
		nowOffset := (i - 1) * chunkSize
		_, err = fi.Seek(nowOffset, 0)
		if err != nil {
			fmt.Println("设置文件读取偏移量失败")
			return
		}
		if len(b) > int(fileLen-nowOffset) {
			b = make([]byte, fileLen-nowOffset)
		}
		_, err := fi.Read(b)
		if err != nil {
			fmt.Println("文件读取失败")
			return
		}
		outFilePath := fmt.Sprintf("./out/%010d.part", i)
		fmt.Printf("生成%s\n", outFilePath)
		f, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		_, err = f.Write(b)
		if err != nil {
			fmt.Println(outFilePath, "文件写入失败")
			return
		}
		err = f.Close()
		if err != nil {
			fmt.Println(outFilePath, "文件关闭失败")
		}
	}
	fmt.Println("拆分完成")
}

func merge(outfile string) {
	fii, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	partList, err := filepath.Glob("./*.part")
	slices.Sort[[]string](partList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("要把%v份合并成一个文件%s\n", partList, outfile)
	i := 0
	for _, v := range partList {
		println(v)
		f, err := os.OpenFile(v, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = fii.Write(b)
		if err != nil {
			fmt.Println("写入失败")
			return
		}
		err = f.Close()
		if err != nil {
			fmt.Println("文件关闭失败")
			return
		}
		i++
		fmt.Printf("合并%d个\n", i)
	}
	err = fii.Close()
	if err != nil {
		fmt.Printf("总文件关闭失败")
		return
	}
	fmt.Println("合并成功")
}

func main() {
	flag.StringVar(&action, "a", "split", "请输入用途：split/merge 默认是split")
	flag.StringVar(&infile, "f", "", "请输入文件名")
	flag.StringVar(&outfile, "o", "azhang.mp4", "请输入要合并的文件名")
	flag.Parse()
	if action == "split" {
		split(infile)
	} else if action == "merge" {
		merge(outfile)
	} else {
		panic("-a只能输入split/merge")
	}
}
