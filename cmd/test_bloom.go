package main

import (
	"bloombits/bloom"
	"bloombits/utils"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
	"unsafe"
)

func main() {
	fmt.Println("Program name:", os.Args)
	if len(os.Args) != 5 {
		panic(errors.New("缺少参数"))
	}
	size, err := strconv.ParseUint(os.Args[1], 10, 32)
	if err != nil {
		panic(errors.New("布隆过滤器长度参数错误"))
	}
	hashFunc, err := strconv.ParseUint(os.Args[2], 10, 32)
	if err != nil {
		panic(errors.New("哈希次数参数错误"))
	}
	if hashFunc == 0 {
		hashFunc = uint64(math.Ceil((float64(size) / 50) * math.Log(2)))
	}
	path := os.Args[3]
	testTimes, err := strconv.ParseUint(os.Args[4], 10, 32)
	if err != nil {
		panic(errors.New("测试次数参数错误"))
	}

	// 初始化布隆过滤器
	bf := bloom.NewBloomFilter(size, hashFunc)
	fmt.Printf("启动布隆过滤器, 参数如下: size=%d, hashFunc=%d\n", size, hashFunc)
	fmt.Printf("布隆过滤器大小: %d 字节\n", unsafe.Sizeof(*bf))

	// 加载地址
	address, err := utils.ReadAddressFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("加载地址完成，写入布隆过滤器的地址数量:%d\n", len(address))

	// 插入数据
	insertStartTime := time.Now()
	fmt.Printf("开始插入数据, 当前时间为:%s\n", insertStartTime.Format("2006-01-02 15:04:05"))
	addressMap := map[string]interface{}{}
	addressMemorySize := uint64(unsafe.Sizeof(address))
	for i := 0; i < len(address); i++ {
		if i%10000 == 0 {
			fmt.Printf("插入数据进度: %.3f\n", float64(i+1)/float64(len(address)))
		}
		data := []byte(address[i])
		bf.AddConcurrently(data)
		addressMap[address[i]] = 1
		addressMemorySize += uint64(utf8.RuneCountInString(address[i]))
	}
	insertEndTime := time.Now()
	fmt.Printf("插入数据完成, 当前时间为:%s\n", insertEndTime.Format("2006-01-02 15:04:05"))

	// 测试数据是否存在
	falsePositiveNum := 0
	for i := uint64(0); i < testTimes; i++ {
		if i%10000 == 0 {
			fmt.Printf("测试误判数据进度: %.3f\n", float64(i+1)/float64(testTimes))
		}
		data := []byte(utils.GetTestAddress())
		if bf.TestConcurrently(data) && addressMap[string(data)] != 1 {
			fmt.Printf("false positive1: %s\n", string(data))
			falsePositiveNum++
		}
	}
	fmt.Println("第一组测试完成")

	for i := 0; i < len(address); i++ {
		if i%10000 == 0 {
			fmt.Printf("测试已插入数据进度: %.3f\n", float64(i+1)/float64(len(address)))
		}
		if !bf.TestConcurrently([]byte(address[i])) {
			fmt.Printf("false positive2: %s\n", address[i])
			panic(errors.New("布隆过滤器异常，测试数据存在于布隆过滤器中，但是返回了false"))
		}
	}
	fmt.Println("第二组测试完成")

	fmt.Printf("-------------- 测试结果 --------------\n")
	fmt.Printf("> 布隆过滤器参数: \n"+
		"   - 数组长度=%d\n"+
		"   - 哈希次数=%d\n"+
		"   - 占用内存大小=%d字节\n", size, hashFunc, bf.Size())
	fmt.Printf("> 数据: \n"+
		"   - 插入数据占用内存(字符串内存)=%d字节\n"+
		"   - 插入数据量=%d\n", addressMemorySize, len(address))
	fmt.Printf("> 耗时: \n"+""+
		"   - 插入数据耗时=%s\n", insertEndTime.Sub(insertStartTime).String())
	fmt.Printf("> 误判: \n"+
		"   - 测试次数=%d\n"+
		"   - 误判次数=%d\n"+
		"   - 误判率=%f\n", testTimes, falsePositiveNum, float64(falsePositiveNum)/float64(testTimes))
	fmt.Printf("-------------- 测试结果 --------------\n")

}
