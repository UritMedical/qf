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
