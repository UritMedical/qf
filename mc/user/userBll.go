package user

import (
	"errors"
	"github.com/Urit-Mediacal/qf"
	"github.com/Urit-Mediacal/qf/helper"
	"github.com/Urit-Mediacal/qf/mc/user/uModel"
	uUtils "github.com/Urit-Mediacal/qf/mc/user/utils"
	"strings"
)

//defPassword 默认密码
const defPassword = "123456"

func (b *Bll) regUserApi(api qf.ApiMap) {
	//登录
	api.Reg(qf.EApiKindSave, "login", b.login)

	//用户增删改查
	api.Reg(qf.EApiKindSave, "", b.saveUser)
	api.Reg(qf.EApiKindDelete, "", b.deleteUser)
	api.Reg(qf.EApiKindGetModel, "", b.getUserModel)
	api.Reg(qf.EApiKindGetList, "", b.getAllUsers)

	//密码重置、修改
	api.Reg(qf.EApiKindSave, "pwd/reset", b.resetPassword)
	api.Reg(qf.EApiKindSave, "pwd", b.changePassword)
}

//
// login
//  @Description: 用户登录
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) login(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		LoginId  string
		Password string
	}{}

	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	params.LoginId = strings.Replace(params.LoginId, " ", "", -1)
	params.Password = uUtils.ConvertToMD5([]byte(params.Password))
	if user, ok := b.userDal.CheckLogin(params.LoginId, params.Password); ok {
		role, _ := b.userRoleDal.GetUsersByRoleId(user.Id)
		return helper.GenerateToken(user.Id, role)
	} else if params.LoginId == devUser.LoginId && params.Password == devUser.Password {
		//开发者账号
		return helper.GenerateToken(devUser.Id, []uint64{})
	} else {
		return nil, errors.New("loginId not exist or password error")
	}
}

func (b *Bll) saveUser(ctx *qf.Context) (interface{}, error) {
	user := &uModel.User{}
	if err := ctx.Bind(user); err != nil {
		return nil, err
	}
	if !b.userDal.CheckExists(user.Id) {
		user.Password = uUtils.ConvertToMD5([]byte(defPassword))
	}
	//创建用户
	return nil, b.userDal.Save(user)
}

func (b *Bll) deleteUser(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetId()
	ret, err := b.userDal.Delete(uId)
	return ret, err
}

func (b *Bll) getUserModel(ctx *qf.Context) (interface{}, error) {
	var user uModel.User
	//获取用户角色
	roleIds, err := b.userRoleDal.GetRolesByUserId(uint64(ctx.UserId))
	if err != nil {
		return nil, err
	}
	roles, err := b.roleDal.GetRolesByIds(roleIds)
	if err != nil {
		return nil, err
	}
	err = b.userDal.GetModel(uint64(ctx.UserId), &user)
	ret := map[string]interface{}{
		"info":  b.Map(user),
		"roles": b.Maps(roles),
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
func (b *Bll) getAllUsers(ctx *qf.Context) (interface{}, error) {
	list, err := b.userDal.GetAllUsers()
	return b.Maps(list), err
}

//
// resetPassword
//  @Description: 重置密码
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) resetPassword(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetId()
	return nil, b.userDal.SetPassword(uId, uUtils.ConvertToMD5([]byte(defPassword)))
}

//
// changePassword
//  @Description: 修改密码
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) changePassword(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		OldPassword string
		NewPassword string
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	if !b.userDal.CheckOldPassword(uint64(ctx.UserId), params.OldPassword) {
		return nil, errors.New("old password is incorrect")
	}
	return nil, b.userDal.SetPassword(uint64(ctx.UserId), uUtils.ConvertToMD5([]byte(params.NewPassword)))
}
