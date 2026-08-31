package main

import (
	"math/big"
	"os"
	"strings"
	"testing"

	"none.com/encrypt_pro/lib"
)

// 从key.txt读取位数最多的那一行作为密钥(617位),返回内容和行号
func longestKeyLine(t *testing.T) (string, int) {
	t.Helper()
	f, err := os.ReadFile(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	best, bestLn := "", 0
	for i, l := range strings.Split(string(f), "\n") {
		l = strings.TrimSpace(l)
		if len(l) > len(best) {
			best, bestLn = l, i+1
		}
	}
	return best, bestLn
}

// 双重加解密往返测试:明文 -> 加密2次 -> 解密2次 必须还原原文
func TestDoubleRoundTrip(t *testing.T) {
	plain := "hello, 双重加密测试 12345!\n第二行内容\nend"
	kStr, kLn := longestKeyLine(t)
	rStr := "1" + strings.Repeat("0", 500)

	var k big.Int
	if _, ok := k.SetString(kStr, 10); !ok {
		t.Fatal("密钥解析失败")
	}
	var r big.Int
	if _, ok := r.SetString(rStr, 10); !ok {
		t.Fatal("r解析失败")
	}
	kr := new(big.Int).Mul(&k, &r) //与JiaMi一致:加密用k*r

	//第一次加密:明文 -> B
	ms := lib.GetM1String(plain, maxLen)
	b := jiami(ms, kr, kLn)

	//第二次加密:B -> C
	ms2 := lib.GetM1String(b, maxLen)
	c := jiami(ms2, kr, kLn)

	if c == b {
		t.Fatal("两次加密结果不应相同")
	}

	//第一次解密:C -> B(解密用原密钥k,即GetK2的行为)
	ns, which := lib.GetNString(c)
	if which != kLn {
		t.Fatalf("密文行号不一致: got %d want %d", which, kLn)
	}
	b2 := lib.GetM2(jiemi(ns, &k))
	if b2 != b {
		t.Fatalf("第一次解密未还原出密文B:\n got %q\nwant %q", b2, b)
	}

	//第二次解密:B -> 明文
	ns2, which2 := lib.GetNString(b2)
	if which2 != kLn {
		t.Fatalf("第二层行号不一致: got %d want %d", which2, kLn)
	}
	out := lib.GetM2(jiemi(ns2, &k))
	if out != plain {
		t.Fatalf("双重加解密往返失败:\n got %q\nwant %q", out, plain)
	}
}

// 大文件(多块)也必须正确往返,依赖jiemi对多个块的拼接
func TestDoubleRoundTripLarge(t *testing.T) {
	plain := strings.Repeat("这是一个用来测试多重分组的消息段,确保超过maxLen后被拆成多个密文块。", 30)
	kStr, kLn := longestKeyLine(t)
	rStr := "1" + strings.Repeat("0", 500)

	var k big.Int
	if _, ok := k.SetString(kStr, 10); !ok {
		t.Fatal("密钥解析失败")
	}
	var r big.Int
	if _, ok := r.SetString(rStr, 10); !ok {
		t.Fatal("r解析失败")
	}
	kr := new(big.Int).Mul(&k, &r)

	ms := lib.GetM1String(plain, maxLen)
	b := jiami(ms, kr, kLn)
	c := jiami(lib.GetM1String(b, maxLen), kr, kLn)
	if len(lib.GetM1String(b, maxLen)) < 2 {
		t.Fatalf("大文件应被拆成多个块,got %d 块", len(lib.GetM1String(b, maxLen)))
	}

	ns, which := lib.GetNString(c)
	if which != kLn {
		t.Fatalf("行号不一致: got %d want %d", which, kLn)
	}
	b2 := lib.GetM2(jiemi(ns, &k))
	ns2, which2 := lib.GetNString(b2)
	if which2 != kLn {
		t.Fatalf("第二层行号不一致: got %d want %d", which2, kLn)
	}
	out := lib.GetM2(jiemi(ns2, &k))
	if out != plain {
		t.Fatalf("大文件双重加解密往返失败:\n got %q\nwant %q", out, plain)
	}
}