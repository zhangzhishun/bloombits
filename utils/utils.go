package utils

import (
	"bufio"
	"fmt"
	"github.com/portto/solana-go-sdk/types"
	"os"
)

func GetTestAddress() string {
	privateKey := types.NewAccount()
	publicKey := privateKey.PublicKey
	return publicKey.ToBase58()
}

func ReadAddressFile(path string) (result []string, err error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("文件打开失败：", err)
		return []string{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("文件关闭失败：", err)
		}
	}(file)

	// 创建扫描器
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	// 检查扫描器是否出错
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	return result, nil
}
