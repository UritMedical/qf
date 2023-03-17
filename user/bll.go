package user

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/user/dal"
	"github.com/UritMedical/qf/user/model"
	"github.com/UritMedical/qf/util"
)

//TODO 开发者密码要可以配置
var devUser = model.User{BaseModel: qf.BaseModel{Id: 202303, FullInfo: "{\"Name\":\"Developer\"}"},
	LoginId: "developer", Password: util.ConvertToMD5([]byte("lisurit"))}

const (
	ErrorCodeToken = iota + 400
	ErrorCodeLogin
)

type Bll struct {
	qf.BaseBll
	userDal           *dal.UserDal           //用户dal
	userRoleDal       *dal.UserRoleDal       //用户-角色
	roleDal           *dal.RoleDal           //角色dal
	rolePermissionDal *dal.RolePermissionDal //角色-权限
	permissionDal     *dal.PermissionDal     //权限dal
	permissionApiDal  *dal.PermissionApiDal  //权限-api
	dptDal            *dal.DepartmentDal     //部门dal
	dptUserDal        *dal.DptUserDal        //部门-用户
}

func (b *Bll) RegApi(api qf.ApiMap) {
	b.regUserApi(api)       //注册用户API
	b.regRoleApi(api)       //注册角色API
	b.regPermissionApi(api) //注册权限组API
	b.regDptApi(api)        //注册部门组织API

	api.Reg(qf.EApiKindSave, "user/jwt/reset", b.resetJwtSecret)  //刷新jwt密钥
	api.Reg(qf.EApiKindSave, "user/parseToken", b.testParseToken) //测试token
}

func (b *Bll) RegDal(regDal qf.DalMap) {
	b.userDal = &dal.UserDal{}
	regDal.Reg(b.userDal, model.User{})

	b.userRoleDal = &dal.UserRoleDal{}
	regDal.Reg(b.userRoleDal, model.UserRole{})

	b.roleDal = &dal.RoleDal{}
	regDal.Reg(b.roleDal, model.Role{})

	b.rolePermissionDal = &dal.RolePermissionDal{}
	regDal.Reg(b.rolePermissionDal, model.RolePermission{})

	b.permissionDal = &dal.PermissionDal{}
	regDal.Reg(b.permissionDal, model.Permission{})

	b.permissionApiDal = &dal.PermissionApiDal{}
	regDal.Reg(b.permissionApiDal, model.PermissionApi{})

	b.dptDal = &dal.DepartmentDal{}
	regDal.Reg(b.dptDal, model.Department{})

	b.dptUserDal = &dal.DptUserDal{}
	regDal.Reg(b.dptUserDal, model.DepartUser{})
}

func (b *Bll) RegFault(f qf.FaultMap) {
	f.Reg(ErrorCodeToken, "用户Token故障")
	f.Reg(ErrorCodeLogin, "用户登陆故障")
}

func (b *Bll) RegMsg(msg qf.MessageMap) {

}

func (b *Bll) RegRef(ref qf.RefMap) {
}

func (b *Bll) Init() error {
	b.initDefUser()
	util.InitJwtSecret()
	return nil
}

func (b *Bll) Stop() {

}

//
// initDefUser
//  @Description: 当用户表数量为0时，初始化默认账号
//
func (b *Bll) initDefUser() {
	//创建admin,developer账号
	list := make([]model.User, 0)
	err := b.userDal.GetList(0, 10, &list)
	if err != nil {
		panic("can't create default user")
	}
	const adminId = 1
	if len(list) == 0 {
		_ = b.userDal.Save(&model.User{
			BaseModel: qf.BaseModel{Id: adminId, FullInfo: "{\"Name\":\"Admin\"}"},
			LoginId:   "admin",
			Password:  util.ConvertToMD5([]byte("admin123"))})

		//创建默认角色
		_ = b.roleDal.Save(&model.Role{BaseModel: qf.BaseModel{Id: adminId, FullInfo: "{\"Name\":\"administrator\"}"}, Name: "administrator"})

		//分配角色
		_ = b.userRoleDal.SetRoleUsers(adminId, []uint64{adminId}) //admin 分配 administrator角色

	}
}

//
//  resetJwtSecret
//  @Description: 重置密钥，然所有用户重新登录
//  @receiver b
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) resetJwtSecret(ctx *qf.Context) (interface{}, qf.IError) {
	jwtStr := util.RandomString(32)
	util.JwtSecret = []byte(jwtStr)
	//将密钥进行AES加密后存入文件
	err := util.EncryptAndWriteToFile(jwtStr, util.JwtSecretFile, []byte(util.AESKey), []byte(util.IV))
	return jwtStr, qf.Error(ErrorCodeToken, err.Error())
}

func (b *Bll) testParseToken(ctx *qf.Context) (interface{}, qf.IError) {
	token := ctx.GetStringValue("token")
	claims, err := util.ParseToken(token)
	if err != nil {
		return nil, qf.Error(ErrorCodeToken, err.Error())
	}
	return claims, nil
}
