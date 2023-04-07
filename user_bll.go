package qf

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/UritMedical/qf/util"
	"github.com/UritMedical/qf/util/token"
	"strings"
)

//TODO 开发者密码要可以配置
var devUser = User{BaseModel: BaseModel{Id: 1, FullInfo: "{\"Name\":\"Developer\"}"},
	LoginId: "developer", Password: convertToMD5([]byte("lisurit"))}

const (
	ErrorCodeTokenInvalid = iota + 401
	ErrorCodeTokenExpires
	ErrorCodeLoginInvalid
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
	api.Reg(EApiKindSave, "logout", b.logout)
	api.Reg(EApiKindSave, "user/jwt/reset", b.resetJwtSecret) //刷新jwt密钥

	//用户增删改查
	api.Reg(EApiKindSave, "user", b.saveUser)
	api.Reg(EApiKindDelete, "user", b.deleteUser)
	api.Reg(EApiKindGetModel, "user", b.getUserModel)
	api.Reg(EApiKindGetList, "users", b.getAllUsers)
	api.Reg(EApiKindGetList, "user/orgs", b.getUserOrg)

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
	f.Reg(ErrorCodeTokenInvalid, "未登录或Token无效, 无法继续执行")
	f.Reg(ErrorCodeTokenExpires, "Token已过期, 请查询登陆")
	f.Reg(ErrorCodeLoginInvalid, "登陆失败, 用户名或密码不正确")
}

func (b *userBll) RegMsg(_ MessageMap) {

}

func (b *userBll) RegRef(_ RefMap) {
}

func (b *userBll) Init() error {
	b.initDefUser()
	token.InitJwtSecret()
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
	const adminId = 2
	if len(list) == 0 {
		_ = b.userDal.Save(&User{
			BaseModel: BaseModel{Id: adminId, FullInfo: "{\"Name\":\"Admin\"}"},
			LoginId:   "admin",
			Password:  convertToMD5([]byte("admin123"))})
	}
}

//
// resetJwtSecret
//  @Description: 重置密钥，然所有用户重新登录
//  @return interface{}
//  @return IError
//
func (b *userBll) resetJwtSecret(_ *Context) (interface{}, IError) {
	jwtStr := token.RandomString(32)
	token.JwtSecret = []byte(jwtStr)
	//将密钥进行AES加密后存入文件
	err := token.EncryptAndWriteToFile(jwtStr, token.JwtSecretFile, []byte(token.AESKey), []byte(token.IV))
	return jwtStr, Error(ErrorCodeTokenInvalid, err.Error())
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
		tkn, _ := token.GenerateToken(user.Id, role)

		//获取用户所在部门
		departs, _ := b.getDepartsByUserId(user.Id)

		//获取用户所拥有的角色
		roles, _ := b.getRolesByUserId(user.Id)

		return map[string]interface{}{
			"Token":    tkn,
			"Departs":  util.ToMaps(departs),
			"Roles":    util.ToMaps(roles),
			"UserInfo": util.ToMap(user),
		}, nil
	} else if params.LoginId == devUser.LoginId && params.Password == devUser.Password {
		//开发者账号
		tkn, _ := token.GenerateToken(devUser.Id, []uint64{})
		return map[string]interface{}{
			"Token":    tkn,
			"UserInfo": util.ToMap(devUser),
		}, nil
	} else {
		return nil, Error(ErrorCodeLoginInvalid, "loginId not exist or password error")
	}
}

func (b *userBll) logout(ctx *Context) (interface{}, IError) {
	return nil, nil
}

func (b *userBll) saveUser(ctx *Context) (interface{}, IError) {
	user := &User{}
	if err := ctx.Bind(user); err != nil {
		return nil, err
	}
	if !b.userDal.CheckExists(user.Id) {
		user.Password = convertToMD5([]byte(defPassword))
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
	if user.Id == 0 {
		return nil, Error(ErrorCodeRecordNotFound, "not found")
	}
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
	return nil, b.userDal.SetPassword(uId, convertToMD5([]byte(defPassword)))
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

//
// getUserOrg
//  @Description: 获取用户机构
//  @receiver b
//  @param ctx
//  @return interface{}
//  @return IError
//
func (b *userBll) getUserOrg(ctx *Context) (interface{}, IError) {
	userId := ctx.LoginUser().UserId
	return b.getOrg(userId)
}

// 计算a数组元素不在b数组之中的所有元素
func diffIntSet(a []uint64, b []uint64) []uint64 {
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

// 转换成MD5加密
func convertToMD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}
