package main

import (
	"fmt"
	"math/big"

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
	maxLen    = 617 / 3 //密钥长度617，除以3是因为每个字母会加密成3位的数字，这里转换成字母的长度
)

func JiaMi(path1, path2 string) {
	ms := lib.GetM1(path1, maxLen)
	k := lib.GetK(keyPath)
	r := lib.GetR(keyLength)

	k.Mul(&k, &r) //n=k*r+m,所以现在是n=k+m

	var (
		n   big.Int
		res string
	)
	for _, m := range ms {
		n.Add(&k, &m)
		res += (n.String() + "\n")
	}

	lib.Write(path2, res)
	fmt.Println("加密成功,")
}

func JieMi(path1, path2 string, notWrite bool) {
	k := lib.GetK(keyPath)
	ns := lib.GetN(path1)

	var (
		m    big.Int
		mStr string
	)
	//n mod k = m
	for _, n := range ns {
		m.Mod(&n, &k)
		mStr += m.String()
	}

	res := lib.GetM2(mStr)
	if !notWrite {
		lib.Write(path2, res)
		fmt.Println("解密成功,")
	} else {
		fmt.Print("----------------------------------------------------------------\n\n")
		fmt.Println(res)
	}
}

/*
func main() {
	JiaMi("testFile/testi.txt")
	JieMi("testFile/testo1.txt")
}
*/

func main() {
	fmt.Print("****************************加密小程序****************************\n\n\n")
	fmt.Print("加密还是解密？[1/2]:")
	mode := 0
	fmt.Scanln(&mode)
	fmt.Print("请输入需要操作的文件的路径:")
	path1 := ""
	fmt.Scanln(&path1)

	if mode == 1 { //加密
	coverStart1:

		fmt.Print("是否覆盖？[Y/n]:")
		cover := ""
		fmt.Scanln(&cover)

		if cover == "y" || cover == "Y" || cover == "" { //覆盖
			JiaMi(path1, path1)
		} else if cover == "n" || cover == "N" { //不覆盖
			fmt.Print("请输入新文件的路径:")
			path2 := ""
			fmt.Scanln(&path2)
			JiaMi(path1, path2)
		} else {
			fmt.Println("输入错误,请重新输入命令.")
			goto coverStart1
		}

	} else if mode == 2 {
		fmt.Print("仅在终端里阅读还是写入到文件？[1/2]")
		notWrite := 0
		fmt.Scanln(&notWrite)

		if notWrite == 1 { //不写
			JieMi(path1, path1, true) //写入路径随便写，反正用不上
		} else if notWrite == 2 {
			//----------------------------------------------------抄上面的
		coverStart2: //为什么不缩进啊难受

			fmt.Print("是否覆盖？[Y/n]:")
			cover := ""
			fmt.Scanln(&cover)

			if cover == "y" || cover == "Y" || cover == "" { //覆盖
				JieMi(path1, path1, false)
			} else if cover == "n" || cover == "N" { //不覆盖
				fmt.Print("请输入新文件的路径:")
				path2 := ""
				fmt.Scanln(&path2)
				JieMi(path1, path2, false)
			} else {
				fmt.Println("输入错误,请重新输入命令.")
				goto coverStart2
			}
			//----------------------------------------------------
		}
	}
}
