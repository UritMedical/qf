package user

import (
	"errors"
	"qf"
	"sort"
	"strings"
)

//defPassword 默认密码
const defPassword = "123456"

//默认账号
var defUsers = [...]User{
	{LoginId: "visitor", Name: "visitor", Role: RoleVisitor, Password: defPassword},     //访客
	{LoginId: "admin", Name: "admin", Role: RoleAdmin, Password: "admin123"},            //admin
	{LoginId: "developer", Name: "developer", Role: RoleDeveloper, Password: "lisurit"}, //开发者
}

type UserBll struct {
	qf.BaseBll
	userDal   *UserDal
	groupBll  *GroupBll
	jwtSecret []byte        //token加密密钥
	cache     map[uint]User //用户缓存
}

func (u *UserBll) RegApi(api qf.ApiMap) {
	api.Reg(qf.EKindSave, "login", u.login)            //登录
	api.Reg(qf.EKindSave, "user/create", u.insertUser) //创建账号
	api.Reg(qf.EKindSave, "user/set/role", u.setRole)  //设置账号角色
	api.Reg(qf.EKindDelete, "user/delete", u.deleteUser)
	api.Reg(qf.EKindSave, "user/password/reset", u.resetPassword)
	api.Reg(qf.EKindSave, "user/modify/account", u.changeAccountInfo)
	api.Reg(qf.EKindGetModel, "user", u.getUserInfo)
	api.Reg(qf.EKindSave, "user/misc", u.updateUserMisc)
	api.Reg(qf.EKindSave, "user/password", u.changePassword)
	api.Reg(qf.EKindGetModel, "user/all", u.getAllUsers)
}

func (u *UserBll) RegDal(dal qf.DalMap) {
	dal.Reg(u.userDal, User{})
}

func (u *UserBll) RegMsg(msg qf.MessageMap) {
	//TODO implement me
	panic("implement me")
}

func (u *UserBll) RefBll() []qf.IBll {
	return []qf.IBll{u.groupBll}
}

func (u *UserBll) Init() error {
	list, _ := u.userDal.GetAllUser()
	u.cache = make(map[uint]User)
	for _, user := range *list {
		u.cache[user.Id] = user
	}

	//创建默认账号
	for _, user := range defUsers {
		_ = u.insertWithPwd(user.LoginId, user.Name, user.Password, user.Role)
	}
	return nil
}

func (u *UserBll) Stop() {}

//
// insertWithPwd
//  @Description: 创建用户账号
//  @param loginId 登录账号
//  @param name 用户姓名
//  @param role 用户权限
//  @return error
//
func (u *UserBll) insertWithPwd(loginId, name, password string, role byte) error {
	//判断用户是否已经存在
	if _, ok := u.isExist(loginId, name); ok {
		return errors.New("loginId or name exist")
	}

	user := User{
		LoginId:  loginId,
		Name:     name,
		Role:     role,
		Password: convertToMD5([]byte(password)),
	}

	//插入到数据库
	err := u.userDal.Insert(&user)
	if err == nil {
		u.cache[user.Id] = user
	}
	return err
}

//isExist 判断用户是否已经存在
func (u *UserBll) isExist(loginId, name string) (User, bool) {
	for _, user := range u.cache {
		if user.LoginId == loginId || user.Name == name {
			return user, true
		}
	}
	return User{}, false
}

//
// login
//  @Description: 登录
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) login(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		LoginId  string `json:"loginId"`
		Password string `json:"password"`
	}{}
	ctx.BindModel(&params)
	params.LoginId = strings.Replace(params.LoginId, " ", "", -1)

	//验证账号是否存在
	user, ok := u.isExist(params.LoginId, "")
	if !ok {
		return "", errors.New("account not exist")
	}
	//验证密码是否正确
	if user.Password != convertToMD5([]byte(params.Password)) {
		return "", errors.New("password error")
	}
	//验证通过，返回新的token
	return GenerateToken(user.Id, user.Role, u.jwtSecret)
}

//
// insertUser
//  @Description: 创建用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) insertUser(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		LoginId string `json:"loginId"`
		Name    string `json:"name"`
	}{}
	ctx.BindModel(&params)
	err := u.insertWithPwd(params.LoginId, params.Name, defPassword, RoleDef)
	return nil, err
}

//
// setRole
//  @Description: 设置账号权限
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) setRole(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		UserId uint `json:"userId"` //目标用户Id
		Role   byte `json:"role"`   //角色
	}{}
	ctx.BindModel(&params)
	//禁止将权限设置为开发者
	if params.Role == byte(RoleDeveloper) {
		return nil, errors.New("can't set role as developer")
	}
	err := u.userDal.SetRole(params.UserId, params.Role)
	if err == nil {
		user, ok := u.cache[params.UserId]
		if ok {
			user.Role = params.Role
			u.cache[user.Id] = user
		}
	}
	return nil, err
}

//
// deleteUser
//  @Description: 删除用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) deleteUser(ctx *qf.Context) (interface{}, error) {
	uId := strToInt(ctx.GetValue("userId"))
	err := u.userDal.Remove(uId)
	//同步删除缓存数据
	if err == nil {
		delete(u.cache, uId)
	}
	return nil, err
}

//
// resetPassword
//  @Description: 重置密码
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) resetPassword(ctx *qf.Context) (interface{}, error) {
	uId := strToInt(ctx.GetValue("userId"))
	user, ok := u.cache[uId]
	if !ok {
		return uId, errors.New("user no found")
	}
	//将数据库中密码重置
	defPwd := convertToMD5([]byte(defPassword))
	err := u.userDal.ChangePassword(uId, defPwd)
	if err == nil {
		//更新缓存
		user.Password = convertToMD5([]byte(defPwd))
		u.cache[user.Id] = user
	}
	return uId, err
}

//
// changeAccountInfo
//  @Description: 修改登录账号、用户名称
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) changeAccountInfo(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		UserId     uint   `json:"userId"`     //目标用户Id
		NewLoginId string `json:"newLoginId"` //新的登录Id
		NewName    string `json:"newName"`    //新的名称
	}{}
	ctx.BindModel(&params)
	if _, ok := u.isExist(params.NewLoginId, params.NewName); ok {
		return nil, errors.New("loginId or name has exist")
	}
	err := u.userDal.ChangeAccountInfo(params.UserId, params.NewLoginId, params.NewName)
	if err == nil {
		user, ok := u.cache[params.UserId]
		if ok {
			user.LoginId = params.NewLoginId
			user.Name = params.NewName
			u.cache[user.Id] = user
		}
	}
	return nil, err
}

//
// getUserInfo
//  @Description: 获取用户信息
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) getUserInfo(ctx *qf.Context) (interface{}, error) {
	user, ok := u.cache[ctx.UserId]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

//
// updateUserMisc
//  @Description: 更新Misc各种信息
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) updateUserMisc(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		Misc string `json:"misc"`
	}{}
	ctx.BindModel(&params)
	err := u.userDal.UpdateMisc(ctx.UserId, params.Misc)
	if err == nil {
		user, ok := u.cache[ctx.UserId]
		if ok {
			user.Misc = params.Misc
			u.cache[user.Id] = user
		}
	}
	return nil, err
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
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}{}
	ctx.BindModel(&params)
	user, ok := u.cache[ctx.UserId]
	if !ok {
		return nil, errors.New("user not found")
	}
	if user.Password != convertToMD5([]byte(params.OldPassword)) {
		return nil, errors.New("old password not correct")
	} else {
		newPwd := convertToMD5([]byte(params.NewPassword))
		err := u.userDal.ChangePassword(user.Id, newPwd)
		if err == nil {
			user.Password = newPwd
			u.cache[user.Id] = user
		}
		return nil, err
	}
}

//
// getAllUsers
//  @Description: 获取全部有效账号
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) getAllUsers(ctx *qf.Context) (interface{}, error) {
	list := make([]User, 0)
	for _, user := range u.cache {
		//默认的这3个账号不返回
		if user.LoginId != "visitor" &&
			user.LoginId != "admin" &&
			user.LoginId != "developer" {
			list = append(list, user)
		}
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Id < list[j].Id
	})
	return list, nil
}

//
// GetUsersByIds
//  @Description: 获取所有用户列表
//  @return *[]User 所有用户列表
//  @return error
//
func (u *UserBll) GetUsersByIds(uIds []uint) []User {
	list := make([]User, 0)
	for _, user := range u.cache {
		if user.LoginId != "visitor" &&
			user.LoginId != "admin" &&
			user.LoginId != "developer" {
			for _, id := range uIds {
				if id == user.Id {
					list = append(list, user)
				}
			}
		}
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Id < list[j].Id
	})
	return list
}
