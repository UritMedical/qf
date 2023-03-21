package qf

import (
	"github.com/UritMedical/qf/util"
	"strings"
)

//TODO 开发者密码要可以配置
var devUser = User{BaseModel: BaseModel{Id: 202303, FullInfo: "{\"Name\":\"Developer\"}"},
	LoginId: "developer", Password: util.ConvertToMD5([]byte("lisurit"))}

const (
	ErrorCodeToken = iota + 400
	ErrorCodeLogin
)

const (
	defPassword = "123456" //默认密码
)

type userBll struct {
	BaseBll
	userDal           *userDal           //用户dal
	userRoleDal       *userRoleDal       //用户-角色
	roleDal           *roleDal           //角色dal
	rolePermissionDal *rolePermissionDal //角色-权限
	permissionDal     *permissionDal     //权限dal
	permissionApiDal  *permissionApiDal  //权限-api
	dptDal            *departmentDal     //部门dal
	dptUserDal        *dptUserDal        //部门-用户
}

func (b *userBll) RegApi(api ApiMap) {
	//登录
	api.Reg(EApiKindSave, "login", b.login)
	api.Reg(EApiKindSave, "user/jwt/reset", b.resetJwtSecret)  //刷新jwt密钥
	api.Reg(EApiKindSave, "user/parseToken", b.testParseToken) //测试token

	//用户增删改查
	api.Reg(EApiKindSave, "user", b.saveUser)
	api.Reg(EApiKindDelete, "user", b.deleteUser)
	api.Reg(EApiKindGetModel, "user", b.getUserModel)
	api.Reg(EApiKindGetList, "users", b.getAllUsers)

	//密码重置、修改
	api.Reg(EApiKindSave, "user/pwd/reset", b.resetPassword)
	api.Reg(EApiKindSave, "user/pwd", b.changePassword)

	b.regRoleApi(api) //注册角色API
	b.regDptApi(api)  //注册部门组织API
}

func (b *userBll) RegDal(regDal DalMap) {
	b.userDal = &userDal{}
	regDal.Reg(b.userDal, User{})

	b.userRoleDal = &userRoleDal{}
	regDal.Reg(b.userRoleDal, UserRole{})

	b.roleDal = &roleDal{}
	regDal.Reg(b.roleDal, Role{})

	b.rolePermissionDal = &rolePermissionDal{}
	regDal.Reg(b.rolePermissionDal, RolePermission{})

	b.permissionDal = &permissionDal{}
	regDal.Reg(b.permissionDal, Permission{})

	b.permissionApiDal = &permissionApiDal{}
	regDal.Reg(b.permissionApiDal, PermissionApi{})

	b.dptDal = &departmentDal{}
	regDal.Reg(b.dptDal, Department{})

	b.dptUserDal = &dptUserDal{}
	regDal.Reg(b.dptUserDal, DepartUser{})
}

func (b *userBll) RegFault(f FaultMap) {
	f.Reg(ErrorCodeToken, "未登录或Token无效, 无法继续执行")
	f.Reg(ErrorCodeLogin, "登陆失败, 用户名或密码不正确")
}

func (b *userBll) RegMsg(_ MessageMap) {

}

func (b *userBll) RegRef(_ RefMap) {
}

func (b *userBll) Init() error {
	b.initDefUser()
	util.InitJwtSecret()
	return nil
}

func (b *userBll) Stop() {

}

//
// initDefUser
//  @Description: 当用户表数量为0时，初始化默认账号
//
func (b *userBll) initDefUser() {
	//创建admin,developer账号
	list := make([]User, 0)
	err := b.userDal.GetList(0, 10, &list)
	if err != nil {
		panic("can't create default user")
	}
	const adminId = 1
	if len(list) == 0 {
		_ = b.userDal.Save(&User{
			BaseModel: BaseModel{Id: adminId, FullInfo: "{\"Name\":\"Admin\"}"},
			LoginId:   "admin",
			Password:  util.ConvertToMD5([]byte("admin123"))})

		//创建默认角色
		_ = b.roleDal.Save(&Role{BaseModel: BaseModel{Id: adminId, FullInfo: "{\"Name\":\"administrator\"}"}, Name: "administrator"})

		//分配角色
		_ = b.userRoleDal.SetRoleUsers(adminId, []uint64{adminId}) //admin 分配 administrator角色

	}
}

//
// resetJwtSecret
//  @Description: 重置密钥，然所有用户重新登录
//  @return interface{}
//  @return IError
//
func (b *userBll) resetJwtSecret(_ *Context) (interface{}, IError) {
	jwtStr := util.RandomString(32)
	util.JwtSecret = []byte(jwtStr)
	//将密钥进行AES加密后存入文件
	err := util.EncryptAndWriteToFile(jwtStr, util.JwtSecretFile, []byte(util.AESKey), []byte(util.IV))
	return jwtStr, Error(ErrorCodeToken, err.Error())
}

//
// testParseToken
//  @Description:
//  @param ctx
//  @return interface{}
//  @return IError
//
func (b *userBll) testParseToken(ctx *Context) (interface{}, IError) {
	token := ctx.GetStringValue("token")
	claims, err := util.ParseToken(token)
	if err != nil {
		return nil, Error(ErrorCodeToken, err.Error())
	}
	return claims, nil
}

//
// login
//  @Description: 用户登录
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) login(ctx *Context) (interface{}, IError) {
	var params = struct {
		LoginId  string
		Password string //md5
	}{}

	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	params.LoginId = strings.Replace(params.LoginId, " ", "", -1)
	if user, ok := b.userDal.CheckLogin(params.LoginId, params.Password); ok {
		role, _ := b.userRoleDal.GetUsersByRoleId(user.Id)
		token, _ := util.GenerateToken(user.Id, role)

		//获取用户所在部门
		departs, _ := b.getDepartsByUserId(user.Id)

		//获取用户所拥有的角色
		roles, _ := b.getRolesByUserId(user.Id)

		return map[string]interface{}{
			"Token":    token,
			"Departs":  util.ToMaps(departs),
			"Roles":    util.ToMaps(roles),
			"UserInfo": util.ToMap(user),
		}, nil
	} else if params.LoginId == devUser.LoginId && params.Password == devUser.Password {
		//开发者账号
		token, _ := util.GenerateToken(devUser.Id, []uint64{})
		return map[string]interface{}{
			"Token":    token,
			"UserInfo": util.ToMap(devUser),
		}, nil
	} else {
		return nil, Error(ErrorCodeLogin, "loginId not exist or password error")
	}
}

func (b *userBll) saveUser(ctx *Context) (interface{}, IError) {
	user := &User{}
	if err := ctx.Bind(user); err != nil {
		return nil, err
	}
	if !b.userDal.CheckExists(user.Id) {
		user.Password = util.ConvertToMD5([]byte(defPassword))
	}
	if user.Id == 0 {
		user.Id = ctx.NewId(user)
	}
	//创建用户
	return nil, b.userDal.Save(user)
}

func (b *userBll) deleteUser(ctx *Context) (interface{}, IError) {
	uId := ctx.GetId()
	return nil, b.userDal.Delete(uId)
}

func (b *userBll) getUserModel(ctx *Context) (interface{}, IError) {
	var user User
	userId := ctx.LoginUser().UserId

	//获取用户所在部门
	departs, _ := b.getDepartsByUserId(userId)

	//获取用户所拥有的角色
	roles, _ := b.getRolesByUserId(userId)

	err := b.userDal.GetModel(userId, &user)
	ret := map[string]interface{}{
		"Info":        util.ToMap(user),
		"Roles":       util.ToMaps(roles),
		"Departments": util.ToMaps(departs),
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
func (b *userBll) getAllUsers(_ *Context) (interface{}, IError) {
	list, err := b.userDal.GetAllUsers()
	result := make([]map[string]interface{}, 0)
	for _, user := range list {
		//获取用户所在部门
		departs, _ := b.getDepartsByUserId(user.Id)

		//获取用户所拥有的角色
		roles, _ := b.getRolesByUserId(user.Id)

		ret := map[string]interface{}{
			"UserInfo":    util.ToMap(user),
			"Roles":       util.ToMaps(roles),
			"Departments": util.ToMaps(departs),
		}
		result = append(result, ret)
	}
	return result, err
}

//
// resetPassword
//  @Description: 重置密码
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) resetPassword(ctx *Context) (interface{}, IError) {
	uId := ctx.GetId()
	return nil, b.userDal.SetPassword(uId, util.ConvertToMD5([]byte(defPassword)))
}

//
// changePassword
//  @Description: 修改密码
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) changePassword(ctx *Context) (interface{}, IError) {
	var params = struct {
		OldPassword string
		NewPassword string
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	if !b.userDal.CheckOldPassword(ctx.LoginUser().UserId, params.OldPassword) {
		return nil, Error(ErrorCodeSaveFailure, "old password is incorrect")
	}
	return nil, b.userDal.SetPassword(ctx.LoginUser().UserId, params.NewPassword)
}
