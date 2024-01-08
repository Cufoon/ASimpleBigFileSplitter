package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func split(infile, destFolder string, chunkSize int64) {
	if !checkDestDir(destFolder) {
		return
	}

	if infile == "" {
		fmt.Println("请输入正确的文件名")
		return
	}

	fileInfo, err := os.Stat(infile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("文件不存在")
		}
		fmt.Println("无法打开想要拆分的文件")
		fmt.Println(err)
		return
	}

	// https://blog.csdn.net/qq_41437512/article/details/128243628
	//splitNum := int64(math.Ceil(float64(fileInfo.Size()) / chunkSize))
	splitNum := (fileInfo.Size()-1)/chunkSize + 1
	splitFileNameFormat := fmt.Sprintf("%%0%dd.part", getNumWidth(splitNum))

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

	splitInfo, _ := json.Marshal(FileSplitInfo{
		Name:     filepath.Base(infile),
		Folder:   "parts",
		PartsNum: splitNum,
	})
	splitInfoFilePath := path.Join(destFolder, "info.json")
	sif, err := os.OpenFile(splitInfoFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer func(sif *os.File) {
		err := sif.Close()
		if err != nil {
			fmt.Println("信息文件关闭失败")
		}
	}(sif)
	if err != nil {
		fmt.Println("无法创建信息文件")
		return
	}
	_, err = sif.Write(splitInfo)
	if err != nil {
		fmt.Println("信息文件写入失败，请检查并重试")
		fmt.Println(err)
		return
	}

	realDestFolder := path.Join(destFolder, "parts")
	err = os.Mkdir(realDestFolder, os.ModePerm)
	if err != nil {
		fmt.Println("存储目录", realDestFolder, "创建失败，请检查并重试")
		fmt.Println(err)
		return
	}

	//sha256RoutineChan := make(chan bool, 17)

	fmt.Printf("要拆分成%d份\n", splitNum)
	b := make([]byte, chunkSize)
	fileLen := fileInfo.Size()
	var i int64 = 1
	for ; i <= splitNum; i++ {
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
		//sha256.Sum256(b)
		outFilePath := path.Join(realDestFolder, fmt.Sprintf(splitFileNameFormat, i))
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
