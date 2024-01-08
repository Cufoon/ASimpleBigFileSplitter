package main

import (
	"flag"
	"fmt"
)

func getSplitParams(params []string) (f string, o string, c int64, err error) {
	sfs := flag.NewFlagSet("split", flag.ExitOnError)
	sfs.StringVar(&f, "f", "", "请输入想要拆分的文件的文件路径")
	sfs.StringVar(&o, "o", "", "请输入拆分后的文件的存放目录")
	sfs.Int64Var(&c, "c", 20180512, "请输入拆分成单个文件的大小，单位字节")
	err = sfs.Parse(params)
	if err != nil {
		return
	}
	if f == "" {
		fmt.Println("请提供想要拆分的文件的文件路径")
		err = ErrShouldToExit
	}
	if o == "" {
		fmt.Println("请提供拆分后的文件的存放目录")
		err = ErrShouldToExit
	}
	return
}

func getMergeParams(params []string) (f string, o string, err error) {
	mfs := flag.NewFlagSet("merge", flag.ExitOnError)
	mfs.StringVar(&f, "f", "", "请输入想要合并的已经拆分的文件的存放目录")
	mfs.StringVar(&o, "o", "", "请输入合并后文件的存放路径")
	err = mfs.Parse(params)
	if err != nil {
		return
	}
	if f == "" {
		fmt.Println("请提供想要合并的已经拆分的文件的存放目录")
		err = ErrShouldToExit
	}
	//if o == "" {
	//	fmt.Println("请提供合并后文件的存放路径")
	//	err = ErrShouldToExit
	//}
	return
}
