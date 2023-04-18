package qf

import (
	"fmt"
	"github.com/UritMedical/qf/util/token"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

//
// verifyToken
//  @Description: 验证token有效性
//  @param ctx
//  @param url
//  @return LoginUser
//  @return IError
//
func (b *userBll) verifyToken(ctx *gin.Context, url string) (LoginUser, IError) {
	// 验证登陆
	login, err := b.doVerify(ctx, url)
	// 如果在白名单，则跳过验证
	if b.tokenWhiteList[url] == 1 {
		return login, nil
	}
	return login, err
}

//
// doVerify
//  @Description: 验证token有效性
//  @param ctx
//  @param url
//  @return LoginUser
//  @return IError
//
func (b *userBll) doVerify(ctx *gin.Context, url string) (LoginUser, IError) {
	login := LoginUser{}
	// 如果是登陆，则跳过
	if strings.Contains(url, "/qf/login") {
		return login, nil
	}
	// 获取Token值
	tokenStr := ctx.GetHeader("Token")
	// 解析token
	claims, err := token.ParseToken(tokenStr)
	// 判断token是否有效
	if err != nil {
		return login, Error(ErrorCodeTokenInvalid, err.Error())
	}
	// 判断是否过期
	if time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		return login, Error(ErrorCodeTokenExpires, err.Error())
	}
	// 生成用户
	u, exist := b.getMap(tokenStr)
	if exist == false {
		login = b.saveToken(claims.Id, tokenStr)
	} else {
		login = u
	}
	// 特殊放行
	if tokenStr == b.tokenSkipVerify || ctx.Query("Bi") == b.tokenSkipVerify {
		return login, nil
	}
	// 权限验证
	if login.UserId > 2 {
		// 获取用户权限
		if _, exist := login.apis[url]; exist == false {
			return login, Error(ErrorCodePermissionDenied, fmt.Sprintf("the user does not have %s permission", url))
		}
	}
	return login, nil
}

//
// saveToken
//  @Description: 保存token到内存
//  @param id
//  @param tokenStr
//
func (b *userBll) saveToken(id uint64, tokenStr string) LoginUser {
	login := LoginUser{}
	if user, err := b.getFullUser(id); err == nil {
		login = LoginUser{
			UserId:   user.UserId,
			UserName: user.UserName,
			LoginId:  user.LoginId,
			roles:    user.roles,
			deptTree: user.deptTree,
			apis:     b.getUserAllApis(user.roles),
		}
		b.setMap(tokenStr, login)
	}
	return login
}

//
// getMap
//  @Description: 从内存获取
//  @param tokenStr
//  @return LoginUser
//  @return bool
//
func (b *userBll) getMap(tokenStr string) (LoginUser, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	user, exist := b.tokenLoginUser[tokenStr]
	return user, exist
}

//
// setMap
//  @Description: 写入内存
//  @param tokenStr
//  @param login
//
func (b *userBll) setMap(tokenStr string, login LoginUser) {
	mutex.Lock()
	defer mutex.Unlock()
	b.tokenLoginUser[tokenStr] = login
}

//
// removeToken
//  @Description: 登出移除内存
//  @param tokenStr
//
func (b *userBll) removeToken(tokenStr string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(b.tokenLoginUser, tokenStr)
}

//
// removeTokenById
//  @Description: 登出移除内存
//  @param id
//
func (b *userBll) removeTokenById(id uint64) {
	mutex.Lock()
	defer mutex.Unlock()
	tokenStr := ""
	for k, v := range b.tokenLoginUser {
		if v.UserId == id {
			tokenStr = k
		}
	}
	delete(b.tokenLoginUser, tokenStr)
}
