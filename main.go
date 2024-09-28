package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"none.com/encrypt_pro/lib"
)

/*
定义一个密钥k,一个很大的素数
定义一个随机数r,也是一个很大的素数
把需要加密的信息记作m,加密后的信息记作n,则：
	k*r + m = n    加密  (m < a*r)
	n mod k = m      解密
*/

const (
	keyPath   = "key.txt"
	keyLength = 2048
	//密钥长度617，除以7是因为utf-8最大值是1114111(恶臭)，这里转换成十进制的长度。
	maxLen = 617/7 - 1
)

func JiaMi(path1, path2 string) {
	ms := lib.GetM1(path1, maxLen)
	k, rn := lib.GetK1(keyPath)
	lib.Write(path2, strconv.Itoa(rn)+"\n", true)
	r := lib.GetR(keyLength)

	k.Mul(&k, &r) //n=k*r+m,所以现在是n=k+m

	var (
		n big.Int
	)
	for _, m := range ms {
		n.Add(&k, &m)
		lib.Write(path2, n.String()+"\n", false)
	}
}

func JieMi(path1, path2 string, notWrite bool) {
	ns, which := lib.GetN(path1)
	k := lib.GetK2(keyPath, which)

	var (
		m    big.Int
		mStr string
	)
	//n mod k = m
	for _, n := range ns {
		m.Mod(&n, &k)
		mStr = m.String()
	}

	res := lib.GetM2(mStr)
	if !notWrite {
		lib.Write(path2, res, true)
	} else {
		fmt.Print("----------------------------------------------------------------\n\n")
		fmt.Println(res)
		fmt.Println("\n----------------------------------------------------------------")
	}
}

func wrongArgs() {
	fmt.Println("参数错误.输入-h获取帮助.")
	os.Exit(-1)
}

// 去掉文件路径的双引号
func trim(inp string) string {
	if inp[0] == '"' { //有开双引号
		if inp[len(inp)-1] != '"' {
			return inp[1:] //没有闭双引号，让他去吧
		} else { //最后一个字符=="
			return inp[1 : len(inp)-1]
		}
	}
	return inp
}

const help = `参数：
1.-e或-d,加密或解密
2.文件路径
3.-c覆盖原文件或者新文件路径,解密时-nw仅在终端中阅读
[在第一个参数中使用-h显示帮助,你显然已经这样做了]`

func main() {
	defer func() {
		fmt.Println("按下回车退出...")
		fmt.Scanln()
	}()

	var args []string
	if len(os.Args) <= 1 {
		fmt.Print("输入参数:")
		inp := ""
		fmt.Scanln(&inp)
		args = append([]string{""}, strings.Fields(inp)...)

	} else {
		args = os.Args
	}

	if args[1] == "-h" && len(args) == 2 {
		fmt.Println(help)
		return
	} else if len(args) == 4 {
		if args[1] == "-e" && args[2] != "" {

			if args[3] == "-c" {
				fmt.Printf("加密%s并覆盖\n", args[2])
				JiaMi(args[2], args[2])
				fmt.Println("加密成功,文件已覆盖")

			} else if args[3] != "" {
				fmt.Printf("加密%s并写入到%s\n", args[2], args[3])
				JiaMi(args[2], args[3])
				fmt.Println("加密成功")

			} else {
				wrongArgs()
			}

		} else if args[1] == "-d" && args[2] != "" {

			if args[3] == "-c" {
				fmt.Printf("解密%s并覆盖\n", args[2])
				JieMi(trim(args[2]), trim(args[2]), false)
				fmt.Println("解密成功,文件已覆盖")

			} else if args[3] == "-nw" {
				fmt.Printf("解密%s并在终端中阅读\n", args[2])
				JieMi(trim(args[2]), "none", true) //反正也不写入,path2瞎填就行了

			} else if args[3] != "" {
				fmt.Printf("解密%s并写入到%s\n", args[2], args[3])
				JieMi(trim(args[2]), trim(args[3]), false)
				fmt.Println("解密成功")

			} else {
				wrongArgs()
			}
		} else {
			wrongArgs()
		}
	} else {
		wrongArgs() //太抽象了
	}
}
