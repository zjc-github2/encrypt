一种基于大素数分解的私钥加密

这听上去很奇怪,对8?它的原理是这样的: (main.go里面也有

定义一个密钥k,一个很大的素数

定义一个随机数r,也是一个很大的素数

把需要加密的信息记作m,加密后的信息记作n,则：

	k*r + m = n    加密  (m < a*r) 
 
	n mod k = m      解密


这东西做的很简陋,没办法,我只是个普普通通的初中生 (-_-)
