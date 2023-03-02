package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"qf"
	uModel "qf/mc/user/model"
	uUtils "qf/mc/user/utils"
	"strings"
)

//defPassword 默认密码
const defPassword = "123456"

func (u *UserBll) regUserApi(api qf.ApiMap) {
	api.Reg(qf.EKindSave, "login", u.login)
	api.Reg(qf.EKindSave, "", u.saveUser)
	api.Reg(qf.EKindDelete, "", u.deleteUser)
	api.Reg(qf.EKindGetModel, "", u.getUserModel)
	api.Reg(qf.EKindGetList, "", u.getAllUsers)
	api.Reg(qf.EKindSave, "pwd/reset", u.resetPassword)
	api.Reg(qf.EKindSave, "pwd", u.changePassword)
}

//
// login
//  @Description: 用户登录
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) login(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		LoginId  string
		Password string
	}{}

	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	params.LoginId = strings.Replace(params.LoginId, " ", "", -1)
	params.Password = uUtils.ConvertToMD5([]byte(params.Password))
	if user, ok := u.userDal.CheckLogin(params.LoginId, params.Password); ok {
		role, _ := u.userRoleDal.GetUsersByRoleId(user.Id)
		return GenerateToken(user.Id, role, u.jwtSecret)
	} else if params.LoginId == devUser.LoginId && params.Password == devUser.Password {
		return GenerateToken(devUser.Id, []uint64{}, u.jwtSecret)
	} else {
		return nil, errors.New("loginId not exist or password error")
	}
}

func (u *UserBll) saveUser(ctx *qf.Context) (interface{}, error) {
	user := &uModel.User{}
	if err := ctx.Bind(user); err != nil {
		return nil, err
	}

	if !u.userDal.CheckExists(user.Id) {
		user.Password = uUtils.ConvertToMD5([]byte(defPassword))
	}
	user.Content = u.BuildContent(user)
	//创建用户
	return nil, u.userDal.Save(user)
}

func (u *UserBll) deleteUser(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	err := u.userDal.Delete(uId)
	return nil, err
}

func (u *UserBll) getUserModel(ctx *qf.Context) (interface{}, error) {
	var user uModel.User
	//获取用户角色
	roles, err := u.userRoleDal.GetRolesByUserId(uint64(ctx.UserId))
	if err != nil {
		return nil, err
	}
	err = u.userDal.GetModel(uint64(ctx.UserId), &user)
	ret := gin.H{
		"info":  user,
		"roles": roles,
	}
	return ret, err
}

//
// getAllUsers
//  @Description: 获取所有用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) getAllUsers(ctx *qf.Context) (interface{}, error) {
	return u.userDal.GetAllUsers()
}

//
// resetPassword
//  @Description: 重置密码
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) resetPassword(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	return nil, u.userDal.SetPassword(uId, uUtils.ConvertToMD5([]byte(defPassword)))
}

//
// changePassword
//  @Description: 修改密码
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) changePassword(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		OldPassword string
		NewPassword string
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	if !u.userDal.CheckOldPassword(uint64(ctx.UserId), params.OldPassword) {
		return nil, errors.New("old password is incorrect")
	}
	return nil, u.userDal.SetPassword(uint64(ctx.UserId), uUtils.ConvertToMD5([]byte(params.NewPassword)))
}
