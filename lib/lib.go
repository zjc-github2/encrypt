package lib

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"strconv"
)

const subKey = 1856823 //在把文字变成数字是还能再加密一下

func printErr(thing ...any) {
	fmt.Println(thing...)
	os.Exit(-1)
}

func GetK1(path string) (big.Int, int) {
	file, err := os.Open(path)
	if err != nil {
		printErr("读取密钥时发生错误:", err.Error(), "请确保你已经把密钥输入到了", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	which := rand.Intn(99) + 1 //1~99
	ln := 0
	var k big.Int
	for scanner.Scan() {
		ln++
		if ln == which {
			_, success := k.SetString(scanner.Text(), 10)
			if !success {
				printErr("读取密钥时发生错误.请确保你输入了正确的密钥.")
			}
		}
	}

	fmt.Println("->读取密钥成功")
	return k, which
}

func GetK2(path string, which int) big.Int {
	//不用看了,抄GetK1的
	file, err := os.Open(path)
	if err != nil {
		printErr("读取密钥时发生错误:", err.Error(), "请确保你已经把密钥输入到了", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ln := 0
	var k big.Int //TODO:成0了
	for scanner.Scan() {
		ln++
		if ln == which {
			_, success := k.SetString(scanner.Text(), 10)
			if !success {
				printErr("读取密钥时发生错误.请确保你输入了正确的密钥.")
			}
		}
	}

	fmt.Println("->读取密钥成功")
	return k
}

//GetR在另外一个文件里

// 读取文件的n个字，utf-8专用的
func readWords(n int, reader *bufio.Reader) ([]rune, bool) {
	var (
		res []rune
		r   rune
		err error
	)
	for range n + 1 {
		r, _, err = reader.ReadRune()
		if err != nil {
			break // 如果遇到错误（如EOF），则跳出循环
		}
		res = append(res, r)
	}

	if err != nil && err != io.EOF {
		printErr("读取文件时遇到错误:", err.Error())
	} else if err == io.EOF {
		return res, true
	}

	return res, false
}

// 把文件变成一个很大的数字。
func changetoNum(file []rune) big.Int {

	var (
		chStr  string
		numStr string
	)
	//先把每一个unicode加上subKey，然后把整个文章变成一个很大的数字
	for _, ch := range file {
		ch += subKey //确保最前面不是0，否则会在加密时被去掉。unicode最大1114111（恶臭）
		chStr = strconv.Itoa(int(ch))
		numStr += chStr
	}

	var (
		bigInt  big.Int
		success bool
	)
	_, success = bigInt.SetString(numStr, 10)
	if !success {
		printErr("这个程序出了点小问题。")
	}

	//整体加密给我删了，反正他没法分解2048位的素数

	return bigInt
}

// 由于强制m<k*r,所以一个长的m可能要加密多次
func GetM1(path string, maxLen int) []big.Int {

	file, err := os.Open(path)
	if err != nil {
		printErr("读取位于", path, "的文件时发生错误：", err.Error())
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	/*
		if len(file) <= maxLen {
			return []big.Int{changetoNum(file
		}*/

	var (
		res     []big.Int
		lastEof bool
	)
	//下面这行大概就是循环截取文章中maxLen个字
	for subFile, eof := readWords(maxLen, reader); !lastEof; subFile, eof = readWords(maxLen, reader) {
		subFileBigInt := changetoNum(subFile)
		res = append(res, subFileBigInt)
		if eof {
			lastEof = true
		}
	}
	return res
}

// 解码，仍然是加密.go抄的
func GetM2(file string) string {
	var (
		chStr string
		chInt int
		res   string
	)
	for i := 0; i < len(file); i += 7 {
		if i+7 <= len(file) {
			chStr = string(file[i : i+7])
		} else { //凑不满7个一组
			printErr("密文损坏或者程序有BUG")
		}
		chInt, _ = strconv.Atoi(chStr)
		chInt -= subKey
		res += string(chInt)
	}
	return res
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

func GetN(path string) ([]big.Int, int) {
	file, err := os.Open(path)
	if err != nil {
		printErr("打开文件时发生错误：", err.Error())
	}
	defer file.Close()

	//先把文件拆开，分别解密，再合起来
	//这里只需要把文件拆开
	ln := bufio.NewScanner(file)
	ln.Scan()
	which, err := strconv.Atoi(ln.Text())
	if err != nil {
		printErr("密文损坏") //TODO:这里
	}

	var (
		res []big.Int
	)
	for ln.Scan() {
		lnBigInt, success := new(big.Int).SetString(ln.Text(), 10)
		if !success {
			printErr("密文损坏")
		}
		res = append(res, *lnBigInt)
	}

	return res, which
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
