package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
	"strings"
	"sync"
	"time"
)

var idx = 0

func main() {
	var wg sync.WaitGroup
	fmt.Println("Hello")
	wg.Add(1)
	// 8个协程调用generate()函数
	for i := 0; i < 8; i++ {
		go generate()
	}
	// 等待所有协程执行完毕
	wg.Wait()

}

func generate() {
	// 无限循环
	fmt.Println("generate")
	for {
		idx++
		if idx%10000 == 0 {
			fmt.Println("生成了", idx, "个地址")
		}
		// 利用go-ethereum 生成钱包信息,包括地址和私钥
		privateKey, err := crypto.GenerateKey()
		privateKeyBytes := crypto.FromECDSA(privateKey)
		pk := hexutil.Encode(privateKeyBytes)[2:]
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 2)
			continue
		}
		// 生成地址
		address := crypto.PubkeyToAddress(privateKey.PublicKey)
		// 打印地址, 私钥
		if checkRule(address) {
			fmt.Println("address:", address.Hex(), "privateKey:", pk)
			filePath := "./address2.txt"
			// 用系统API向文件追加地址和私钥
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(f)
			_, err = f.WriteString(address.Hex() + " " + pk + "\r\n")
			if err != nil {
				return
			}

		}
	}
}

func checkRule(address common.Address) bool {
	content := address.Hex()
	// 判断content是否都是数字8
	// 获取地址开头8个字符
	prefix := content[2:10]
	// 获取地址结尾8个字符
	suffix := content[len(content)-8:]
	// 判断prefix是否由一个字符构成
	if len(prefix) == 8 && strings.Count(prefix, prefix[0:1]) == 8 {
		return true
	}
	// 判断suffix是否由一个字符构成
	if len(suffix) == 8 && strings.Count(suffix, suffix[0:1]) == 8 {
		return true
	}
	return false
}
