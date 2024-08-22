import gmpy2  
import random  

length=2048
  
def generate_large_prime(bits):   
    max_int = 2 ** (bits+100) + 1  
    min_int=2**(bits-100)-1
    random_number = random.randint(min_int, max_int)  
      
    # 使用gmpy2的next_prime函数找到下一个素数  
    large_prime = gmpy2.next_prime(random_number)  
    return large_prime  
large_prime = generate_large_prime(length)  
#file.write(str(large_prime) + '\n')  # 写入素数并换行  
print(large_prime)  # 在控制台打印信息  