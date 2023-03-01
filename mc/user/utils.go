package user

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

//UE是UserError的缩写
const (
	ErrUserNotExist  = "UE101" //用户不存在
	ErrPasswordError = "UE102" //密码不正确
	ErrTokenInvalid  = "UE103" //授权失败
)

//StrToInt 数字字符串转 uint
func strToInt(str string) uint {
	i, _ := strconv.Atoi(str)
	return uint(i)
}

//转换成MD5加密
func convertToMD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

//
// diff
//  @Description: 计算a数组元素不在b数组之中的所有元素
//  @param a
//  @param b
//  @return []uint
//
func diffIntSet(a []uint, b []uint) []uint {
	c := make([]uint, 0)
	temp := map[uint]struct{}{}
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
