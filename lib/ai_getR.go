package lib

//傻逼go，不让我同时导入math/rand和crypto/rand

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// 获取一个二进制n位的随机素数。文心一言写的。
func GetR(n int) big.Int {
	// 初始化一个足够大的数，保证至少有n位
	bitSize := n
	max := new(big.Int).Lsh(big.NewInt(1), uint(bitSize))
	max.Sub(max, big.NewInt(2)) // 2^bitSize - 2，确保是偶数

	var num *big.Int // 将num的声明移动到for循环外部

	// 不断尝试直到找到一个素数
	for {
		// 生成一个随机的奇数
		var err error
		num, err = rand.Int(rand.Reader, max)
		if err != nil {
			printErr("生成随机素数时发生错误：", err.Error())

		}
		// 确保是奇数（但rand.Int已经返回[0, max)内的随机数，这里可能需要检查是否已经是奇数）
		if num.BitLen()%2 == 0 { // 如果num是偶数
			num.Add(num, big.NewInt(1)) // 使其成为奇数
			if num.Cmp(max) >= 0 {      // 确保不超过max
				num.Sub(num, big.NewInt(2)) // 如果超过，则减去2
			}
		}

		// 使用Miller-Rabin测试
		if num.ProbablyPrime(100) {
			break
		}

	}

	fmt.Println("->生成密钥成功")
	return *num
}
