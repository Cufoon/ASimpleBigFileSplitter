package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
)

func merge(splitFileFolder, outfile string) {
	if outfile == "" {

		infoFile, err := os.OpenFile(path.Join(splitFileFolder, "info.json"), os.O_RDONLY, os.ModePerm)
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println("信息文件关闭失败！")
				fmt.Println(err)
				return
			}
		}(infoFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		infoFileStat, err := infoFile.Stat()
		if err != nil {
			fmt.Println("信息文件读取失败！")
			return
		}
		info := make([]byte, infoFileStat.Size())
		_, err = infoFile.Read(info)
		if err != nil {
			fmt.Println("信息文件读取失败！")
			return
		}
		var infoObject FileSplitInfo
		err = json.Unmarshal(info, &infoObject)
		if err != nil {
			fmt.Println("信息文件读取失败！")
			return
		}

		outfile = path.Join(splitFileFolder, infoObject.Name)
	}

	fii, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	partList, err := filepath.Glob(path.Join(splitFileFolder, "parts/*.part"))
	slices.Sort[[]string](partList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("要把%d份合并成一个文件%s\n", len(partList), outfile)
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
