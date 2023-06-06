package main

import (
	"bloombits/utils"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Program name:", os.Args)
	if len(os.Args) != 2 {
		panic(errors.New("缺少参数"))
	}
	number, err := strconv.ParseUint(os.Args[1], 10, 32)
	if err != nil {
		panic(errors.New("地址数量参数错误"))
	}

	f, err := os.OpenFile("address.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	for i := uint64(0); i < number; i++ {
		if i%10000 == 0 {
			fmt.Printf("生成地址进度: %.3f\n", float64(i+1)/float64(number))
		}
		str := utils.GetTestAddress() + "\n"
		if err != nil {
			panic(err)
		}
		if _, err := f.Write([]byte(str)); err != nil {
			panic(err)
		}
	}
	fmt.Println("生成地址完成. 文件名: address.txt")
}
