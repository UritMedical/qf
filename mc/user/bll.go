package user

import (
	"github.com/Urit-Mediacal/qf"
	"github.com/Urit-Mediacal/qf/helper"
	uDal "github.com/Urit-Mediacal/qf/mc/user/dal"
	"github.com/Urit-Mediacal/qf/mc/user/uModel"
	uUtils "github.com/Urit-Mediacal/qf/mc/user/utils"
)

//TODO 开发者密码可以配置
const DeveloperId = 202303 //开发者内存Id
var devUser = uModel.User{BaseModel: qf.BaseModel{Id: DeveloperId}, LoginId: "developer", Password: uUtils.ConvertToMD5([]byte("lisurit"))}

type Bll struct {
	qf.BaseBll
	userDal       *uDal.UserDal        //用户dal
	userRoleDal   *uDal.UserRoleDal    //用户-角色
	roleDal       *uDal.RoleDal        //角色dal
	roleRightsDal *uDal.RoleRightsDal  //角色-权限
	rightsDal     *uDal.RightsGroupDal //权限dal
	rightsApiDal  *uDal.RightsApiDal   //权限-api
	dptDal        *uDal.DepartmentDal  //部门dal
	dptUserDal    *uDal.DptUserDal     //部门-用户
}

func (b *Bll) RegApi(api qf.ApiMap) {
	b.regUserApi(api)   //注册用户API
	b.regRoleApi(api)   //注册角色API
	b.regRightsApi(api) //注册权限组API
	b.regDptApi(api)    //注册部门组织API

	api.Reg(qf.EApiKindSave, "jwt/reset", b.resetJwtSecret)  //刷新jwt密钥
	api.Reg(qf.EApiKindSave, "parseToken", b.testParseToken) //测试token
}

func (b *Bll) RegDal(dal qf.DalMap) {
	b.userDal = &uDal.UserDal{}
	dal.Reg(b.userDal, uModel.User{})

	b.userRoleDal = &uDal.UserRoleDal{}
	dal.Reg(b.userRoleDal, uModel.UserRole{})

	b.roleDal = &uDal.RoleDal{}
	dal.Reg(b.roleDal, uModel.Role{})

	b.roleRightsDal = &uDal.RoleRightsDal{}
	dal.Reg(b.roleRightsDal, uModel.RoleRights{})

	b.rightsDal = &uDal.RightsGroupDal{}
	dal.Reg(b.rightsDal, uModel.RightsGroup{})

	b.rightsApiDal = &uDal.RightsApiDal{}
	dal.Reg(b.rightsApiDal, uModel.RightsApi{})

	b.dptDal = &uDal.DepartmentDal{}
	dal.Reg(b.dptDal, uModel.Department{})

	b.dptUserDal = &uDal.DptUserDal{}
	dal.Reg(b.dptUserDal, uModel.DepartUser{})
}

func (b *Bll) RegMsg(msg qf.MessageMap) {

}

func (b *Bll) RegRef(ref qf.RefMap) {
}

func (b *Bll) Init() error {
	b.initDefUser()
	helper.InitJwtSecret()
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
	list := make([]uModel.User, 0)
	err := b.userDal.GetList(0, 10, &list)
	if err != nil {
		panic("can't create default user")
	}
	const adminId = 1
	if len(list) == 0 {
		_ = b.userDal.Save(&uModel.User{
			BaseModel: qf.BaseModel{Id: adminId, FullInfo: "{\"LoginId\":\"admin\",\"Name\":\"Admin\"}"},
			LoginId:   "admin",
			Password:  uUtils.ConvertToMD5([]byte("admin123"))})

		//创建默认角色
		_ = b.roleDal.Save(&uModel.Role{BaseModel: qf.BaseModel{Id: adminId, FullInfo: "{\"Name\":\"administrator\"}"}, Name: "administrator"})

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
func (b *Bll) resetJwtSecret(ctx *qf.Context) (interface{}, error) {
	jwtStr := uUtils.RandomString(32)
	helper.JwtSecret = []byte(jwtStr)
	//将密钥进行AES加密后存入文件
	err := helper.EncryptAndWriteToFile(jwtStr, helper.JwtSecretFile, []byte(helper.AESKey), []byte(helper.IV))
	return jwtStr, err
}

func (b *Bll) testParseToken(ctx *qf.Context) (interface{}, error) {
	token := ctx.GetStringValue("token")
	return helper.ParseToken(token)
}
