package main

import (
	"fmt"
	"strconv"
	"strings"
)
func departmentStatistics (employees []string, friendships []string) []string {
	// write your code here.
	type detail struct{
		FriendCount int
		Sum int
	}
	bgSum:=make(map[string]detail)
	IDBg:=make(map[int]string)
	for _,v:=range employees{
		tmp:=strings.Split(v,", ")
		tmpID,_:=strconv.Atoi(tmp[0])
		tmpBg:=tmp[2]
		IDBg[tmpID]=tmpBg
		if _,ok:=bgSum[tmpBg];ok{
			s:=bgSum[tmpBg]
			s.Sum+=1
			bgSum[tmpBg]=s
		}else{
			bgSum[tmpBg]=detail{0,1}
		}
	}

	for _,v:=range friendships{
		ss:=strings.Split(v,", ")
		for i:=0;i<len(ss);i++{
			a,_:=strconv.Atoi(ss[i])
			if bg,ok:=IDBg[a];ok{
				//bgSum[bg].FriendCount++
				s:=bgSum[bg]
				s.FriendCount+=1
				bgSum[bg]=s
				delete(IDBg,a)
			}
		}
	}
	res:=make([]string,0)
	for k,v:=range bgSum{
		tmp:=k+": "+strconv.Itoa(v.FriendCount)+" "+"of "+strconv.Itoa(v.Sum)
		res=append(res,tmp)
	}
	return res
}

func longestPalindromeBp(s string) string{
	n:= len(s)
	ans:=""
	dp:=make([][]bool,n)
	for i:=0;i<n;i++{
		dp[i]=make([]bool,n)
	}
	//回文子串的长度从短 到 长去遍历 k 为0时，会问子串就是遍历元素本身，k为当前遍历元素后面的第k个字符
	//k代表从i开始往后数k个作为j,也就是回文子串的长度为k+1
	for k:=0;k<n;k++{
		for i:=0;i+k<n;i++{
			j:=i+k
			//当i，j重合时，也就是某个元素本身，肯定是回文子串，因为就一个元素
			if k==0{
				dp[i][j]=true
				//也就是当相邻两个元素相同时
			}else if k==1{
				if s[i]==s[j]{
					dp[i][j]=true
				}
			}else{
				if s[i]==s[j]{
					dp[i][j]=dp[i+1][j-1]
				}
			}
			if dp[i][j]&&k+1> len(ans){
				ans=s[i:i+k+1]
			}
		}

	}
	return ans

}

func longestPalindrome(s string) string {
	if s==""{
		return ""
	}
	start,end:=0,0
	for i:=0;i< len(s);i++{
		left1,right1:=expand(s,i,i)
		left2,right2:=0,0
		//相邻的两个字符相等的时候
		if i<len(s)-1&&s[i]==s[i+1]{
			left2,right2=expand(s,i,i+1)
		}

		if right1-left1>end-start{
			start,end=left1,right1
		}
		if right2-left2>end-start{
			start,end=left2,right2
		}
	}
	return s[start:end+1]
}

func expand(s string,left,right int ) (int ,int ){
	for left>=0&&right<len(s)&&s[left]==s[right]{
		left,right=left-1,right+1
	}
	return left+1,right-1
}


func main(){
	/*cache,err:=bigcache.NewBigCache(bigcache.DefaultConfig(10*time.Minute))
	if err!=nil{
		log.Println(err)
		return
	}

	entry,err:=cache.Get("my-unique-key")
	if err!=nil{
		log.Println(err)
		return
	}

	if entry==nil{
		//如果从缓存中没有获取，则从数据源中获取（一般是数据库），然后设置到缓存
		entry = []byte("value")
		cache.Set("my-unique-key",entry)
	}
	log.Println(string(entry))*/


	fmt.Println(longestPalindromeBp("cbbd"))
}


