package user

import (
    "errors"
    "qf"
    uModel "qf/mc/user/model"
    uUtils "qf/mc/user/utils"
    "strings"
)

//defPassword 默认密码
const defPassword = "123456"

func (u *UserBll) regUserApi(api qf.ApiMap) {
    //登录
    api.Reg(qf.EApiKindSave, "login", u.login)

    //用户增删改查
    api.Reg(qf.EApiKindSave, "", u.saveUser)
    api.Reg(qf.EApiKindDelete, "", u.deleteUser)
    api.Reg(qf.EApiKindGetModel, "", u.getUserModel)
    api.Reg(qf.EApiKindGetList, "", u.getAllUsers)

    //密码重置、修改
    api.Reg(qf.EApiKindSave, "pwd/reset", u.resetPassword)
    api.Reg(qf.EApiKindSave, "pwd", u.changePassword)
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
    user.BaseModel = u.BuildBaseModel(user)
    //创建用户
    return nil, u.userDal.Save(user)
}

func (u *UserBll) deleteUser(ctx *qf.Context) (interface{}, error) {
    uId := ctx.GetUIntValue("Id")
    ret, err := u.userDal.Delete(uId)
    return ret, err
}

func (u *UserBll) getUserModel(ctx *qf.Context) (interface{}, error) {
    var user uModel.User
    //获取用户角色
    roleIds, err := u.userRoleDal.GetRolesByUserId(uint64(ctx.UserId))
    if err != nil {
        return nil, err
    }
    roles, err := u.roleDal.GetRolesByIds(roleIds)
    if err != nil {
        return nil, err
    }
    err = u.userDal.GetModel(uint64(ctx.UserId), &user)
    ret := map[string]interface{}{
        "info":  u.Map(user),
        "roles": u.Maps(roles),
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
    list, err := u.userDal.GetAllUsers()
    return u.Maps(list), err
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
