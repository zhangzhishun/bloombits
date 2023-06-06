package eth

import (
	"bloombits/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"testing"
)

func TestETHBloomBit(t *testing.T) {
	// 创建一个新的布隆过滤器
	filter := types.CreateBloom(nil)

	// 添加字符串到布隆过滤器中
	for i := 0; i < 1*1e4; i++ {
		filter.Add([]byte(utils.GetTestAddress()))
	}

	// 检查字符串是否存在于布隆过滤器中
	fmt.Println(filter.Test([]byte(utils.GetTestAddress()))) // true

	fmt.Println(len(filter.Bytes()))
}
