package lib

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func printErr(thing ...any) {
	fmt.Println(thing...)
	os.Exit(-1)
	//没错，就2行我也要封装
}

func GetK(path string) big.Int {
	file, err := os.ReadFile(path)
	if err != nil {
		printErr("读取密钥时发生错误:", err.Error(), "请确保你已经把密钥输入到了", path)
	}

	var a big.Int
	a.SetString(string(file), 10)

	fmt.Println("->读取密钥成功")
	return a
}

//GetR在另外一个文件里

// 把文件变成一个很大的数字,加密.go抄的
func changetoNum(file []byte) big.Int {

	var (
		chInt   int
		chStr   string
		bigInt  big.Int
		numStr  string
		success bool
	)

	//先把每一个ASCII加上234，然后把整个文章变成一个很大的数字
	for _, ch := range file {
		chInt = int(ch)
		chInt += 234
		chStr = strconv.Itoa(chInt)
		numStr += chStr
	}
	_, success = bigInt.SetString(numStr, 10)
	if !success {
		printErr("这个程序出了点小问题。")
	}

	//整体加密给我删了，反正他没法分解2048位的素数

	return bigInt
}

// 由于强制m<k*r,所以一个长的m可能要加密多次
func GetM1(path string, maxLen int) []big.Int {

	file, err := os.ReadFile(path)
	if err != nil {
		printErr("读取文件时发生错误：", err.Error())
	}

	if len(file) <= maxLen {
		return []big.Int{changetoNum(file)}
	}

	var (
		res     []big.Int
		subFile []byte
	)
	for i := 0; i < len(file); i += maxLen {
		if i+maxLen >= len(file) { //len(file)不是maxLen的整倍数
			subFile = file[i:]
		} else {
			subFile = file[i : i+maxLen]
		}
		subFileBigInt := changetoNum(subFile)
		res = append(res, subFileBigInt)
	}
	return res
}

// 解码，仍然是加密.go抄的
func GetM2(file string) string {
	var (
		chStr string
		chInt int
		res   []byte
	)
	for i := 0; i < len(file); i += 3 {
		if i+3 <= len(file) {
			chStr = string(file[i : i+3])
		} else { //凑不满3个一组
			printErr("密文损坏或者程序有BUG")
		}
		chInt, _ = strconv.Atoi(chStr)
		chInt -= 234
		res = append(res, byte(chInt))
	}
	return string(res)
}

/*   发现了bufio.Scanner,不用了
func findLn(inp []byte, start int) int {
	for i := start + 1; i < len(inp); i++ {
		if inp[i] == '\n' {
			return i
		}
	}
	return -1
}
*/

func GetN(path string) []big.Int {
	file, err := os.Open(path)
	if err != nil {
		printErr("打开文件时发生错误：", err.Error())
	}

	//先把文件拆开，分别解密，再合起来
	//这里只需要把文件拆开
	ln := bufio.NewScanner(file)
	var (
		res []big.Int
	)
	for ln.Scan() {
		lnBigInt, success := new(big.Int).SetString(ln.Text(), 10)
		if !success {
			printErr("密文损坏或者程序有BUG")
		}
		res = append(res, *lnBigInt)
	}

	return res
}

func Write(path string, thing string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		printErr("读取文件时发生错误：", err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(thing)
	if err != nil {
		printErr("写入文件时发生错误：", err.Error())
	}
}
