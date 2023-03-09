package uUtils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

//ConvertToMD5 转换成MD5加密
func ConvertToMD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

//
// DiffIntSet
//  @Description: 计算a数组元素不在b数组之中的所有元素
//  @param a
//  @param b
//  @return []uint64
//
func DiffIntSet(a []uint64, b []uint64) []uint64 {
	c := make([]uint64, 0)
	temp := map[uint64]struct{}{}
	//把b所有的值作为key存入temp
	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}
	//如果a中的值作为key在temp中找不到，说明它不在b中
	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}
	return c
}

// RandomString
//  @Description: 生成随机字符串
//  @param length
//  @return string
//
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
