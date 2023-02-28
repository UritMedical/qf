package user

import (
	"errors"
	"qf"
	"strings"
)

//defPassword 默认密码
const defPassword = "123456"

//默认账号
var defUsers = [...]User{
	//admin
	{
		Content: qf.Content{
			Info: "{Name:\"Admin\"}",
		},
		LoginId:  "admin",
		Password: convertToMD5([]byte("admin123")),
	},
	//developer
	{
		Content: qf.Content{
			Info: "{Name:\"Developer\"}",
		},
		LoginId:  "developer",
		Password: convertToMD5([]byte("lisurit")),
	},
}

func (u *UserBll) regUserApi(api qf.ApiMap) {
	api.Reg(qf.EKindSave, "login", u.login)
	api.Reg(qf.EKindSave, "user", u.saveUser)
	api.Reg(qf.EKindDelete, "user", u.deleteUser)
	api.Reg(qf.EKindGetModel, "user", u.getUserModel)
	api.Reg(qf.EKindGetModel, "users", u.getUserList)
	api.Reg(qf.EKindSave, "user/pwd/reset", u.resetPassword)
	api.Reg(qf.EKindSave, "user/pwd", u.changePassword)
}

//
// initDefUser
//  @Description: 当用户表数量为0时，初始化默认账号
//
func (u *UserBll) initDefUser() {
	list := make([]User, 0)
	//u.userDal.GetList(0, 10, list)
	if len(list) == 0 {
		for _, defUser := range defUsers {
			err := u.userDal.Save(defUser)
			if err == nil {
				//TODO 分配角色
			}
		}
	}
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

	if err := ctx.BindModel(&params); err != nil {
		return nil, err
	}
	params.LoginId = strings.Replace(params.LoginId, " ", "", -1)
	params.Password = convertToMD5([]byte(params.Password))
	if u.userDal.CheckLogin(params.LoginId, params.Password) {
		//TODO 生成token
		return "tokentodo", nil
	} else {
		return nil, errors.New("loginId not exist or password error")
	}
}

func (u *UserBll) saveUser(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		Id      uint
		Info    string
		LoginId string
	}{}

	if err := ctx.BindModel(&params); err != nil {
		return nil, err
	}

	user := &User{
		Content: qf.Content{
			Info: params.Info,
		},
		LoginId:  params.LoginId,
		Password: "",
	}
	if params.Id == 0 {
		user.Password = convertToMD5([]byte(defPassword))
	}
	//创建用户
	return nil, u.userDal.Save(user)
}

func (u *UserBll) deleteUser(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	err := u.userDal.Delete(uId)
	return nil, err
}

func (u *UserBll) getUserModel(ctx *qf.Context) (interface{}, error) {
	var user User
	return u.userDal.GetModel(&user)
}

func (u *UserBll) getUserList(ctx *qf.Context) (interface{}, error) {
	return nil, nil
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
	return nil, u.userDal.SetPassword(uId, convertToMD5([]byte(defPassword)))
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
	if err := ctx.BindModel(&params); err != nil {
		return nil, err
	}
	if !u.userDal.CheckOldPassword(ctx.UserId, params.OldPassword) {
		return nil, errors.New("old password is incorrect")
	}
	return nil, u.userDal.SetPassword(ctx.UserId, convertToMD5([]byte(params.NewPassword)))
}
