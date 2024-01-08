package main

import (
	"fmt"
	"os"
	"time"
)

// 512 MiB
// const chunkSize = 536870912

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  > split file:")
	fmt.Println("    litsplitor s|split -f file_to_splitor -o where_split_files_storaged [-c chunkSizeInByte]")
	fmt.Println("  > merge file:")
	fmt.Println("    litsplitor m|merge -f where_split_files_storaged -o where_merged_file_storaged")
}

func main() {
	appStartTime := time.Now().UnixMilli()
	if len(os.Args) <= 1 {
		printHelp()
		return
	}
	a := os.Args[1]
	fmt.Println("action", a)
	if a == "split" || a == "s" {
		f, o, c, err := getSplitParams(os.Args[2:])
		if err != nil {
			handleShouldToExitErr(err)
			return
		}
		split(f, o, c)
	} else if a == "merge" || a == "m" {
		f, o, err := getMergeParams(os.Args[2:])
		if err != nil {
			handleShouldToExitErr(err)
			return
		}
		merge(f, o)
	} else {
		fmt.Println("命令只能是 split(s) 或者 merge(m)")
	}
	appEndTime := time.Now().UnixMilli()
	appProcessTime := appEndTime - appStartTime
	appProcessTimeSeconds := appProcessTime / 1000
	appProcessTimeMilliSeconds := appProcessTime - appProcessTimeSeconds*1000
	fmt.Println("运行时间: ", appProcessTimeSeconds, "秒", appProcessTimeMilliSeconds, "毫秒")
}
