#!/usr/bin/env python3
"""生成key.txt:99行,每行为500~700位十进制大素数。

main.go中GetK1会随机选取1~99行作为密钥,因此必须恰好生成99行。
解密时n mod k = m要求m < k,故取最短500位作为安全下限;
main.go的maxLen按最短密钥位设计,保证单块密文永远小于任何一把密钥。
"""
import gmpy2
import random
import sys

LINES = 99
MIN_LEN = 500  # 密钥最短十进制位数(决定maxLen安全上限)
MAX_LEN = 700  # 密钥最长十进制位数


def random_prime(dec_len: int) -> int:
    """生成恰好dec_len位十进制的大素数。"""
    while True:
        n = random.randint(10 ** (dec_len - 1), 10 ** dec_len - 1)
        n |= 1  # 确保为奇数
        p = int(gmpy2.next_prime(n))
        if MIN_LEN <= len(str(p)) <= MAX_LEN:
            return p


def main() -> None:
    out = sys.argv[1] if len(sys.argv) > 1 else "key.txt"
    with open(out, "w", encoding="ascii") as f:
        for _ in range(LINES):
            dec_len = random.randint(MIN_LEN, MAX_LEN)
            f.write(str(random_prime(dec_len)) + "\n")

    print(f"已生成{LINES}行密钥到 {out},每行{MIN_LEN}~{MAX_LEN}位十进制大素数")


if __name__ == "__main__":
    main()