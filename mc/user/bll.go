package user

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/helper"
	"github.com/UritMedical/qf/mc/user/dal"
	"github.com/UritMedical/qf/mc/user/model"
	utils "github.com/UritMedical/qf/mc/user/utils"
)

//TODO 开发者密码可以配置
const DeveloperId = 202303 //开发者内存Id
var devUser = model.User{BaseModel: qf.BaseModel{Id: DeveloperId}, LoginId: "developer", Password: utils.ConvertToMD5([]byte("lisurit"))}

type Bll struct {
	qf.BaseBll
	userDal       *dal.UserDal        //用户dal
	userRoleDal   *dal.UserRoleDal    //用户-角色
	roleDal       *dal.RoleDal        //角色dal
	roleRightsDal *dal.RoleRightsDal  //角色-权限
	rightsDal     *dal.RightsGroupDal //权限dal
	rightsApiDal  *dal.RightsApiDal   //权限-api
	dptDal        *dal.DepartmentDal  //部门dal
	dptUserDal    *dal.DptUserDal     //部门-用户
}

func (b *Bll) RegApi(api qf.ApiMap) {
	b.regUserApi(api)   //注册用户API
	b.regRoleApi(api)   //注册角色API
	b.regRightsApi(api) //注册权限组API
	b.regDptApi(api)    //注册部门组织API

	api.Reg(qf.EApiKindSave, "jwt/reset", b.resetJwtSecret)  //刷新jwt密钥
	api.Reg(qf.EApiKindSave, "parseToken", b.testParseToken) //测试token
}

func (b *Bll) RegDal(regDal qf.DalMap) {
	b.userDal = &dal.UserDal{}
	regDal.Reg(b.userDal, model.User{})

	b.userRoleDal = &dal.UserRoleDal{}
	regDal.Reg(b.userRoleDal, model.UserRole{})

	b.roleDal = &dal.RoleDal{}
	regDal.Reg(b.roleDal, model.Role{})

	b.roleRightsDal = &dal.RoleRightsDal{}
	regDal.Reg(b.roleRightsDal, model.RoleRights{})

	b.rightsDal = &dal.RightsGroupDal{}
	regDal.Reg(b.rightsDal, model.RightsGroup{})

	b.rightsApiDal = &dal.RightsApiDal{}
	regDal.Reg(b.rightsApiDal, model.RightsApi{})

	b.dptDal = &dal.DepartmentDal{}
	regDal.Reg(b.dptDal, model.Department{})

	b.dptUserDal = &dal.DptUserDal{}
	regDal.Reg(b.dptUserDal, model.DepartUser{})
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
	list := make([]model.User, 0)
	err := b.userDal.GetList(0, 10, &list)
	if err != nil {
		panic("can't create default user")
	}
	const adminId = 1
	if len(list) == 0 {
		_ = b.userDal.Save(&model.User{
			BaseModel: qf.BaseModel{Id: adminId, FullInfo: "{\"LoginId\":\"admin\",\"Name\":\"Admin\"}"},
			LoginId:   "admin",
			Password:  utils.ConvertToMD5([]byte("admin123"))})

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
func (b *Bll) resetJwtSecret(ctx *qf.Context) (interface{}, error) {
	jwtStr := utils.RandomString(32)
	helper.JwtSecret = []byte(jwtStr)
	//将密钥进行AES加密后存入文件
	err := helper.EncryptAndWriteToFile(jwtStr, helper.JwtSecretFile, []byte(helper.AESKey), []byte(helper.IV))
	return jwtStr, err
}

func (b *Bll) testParseToken(ctx *qf.Context) (interface{}, error) {
	token := ctx.GetStringValue("token")
	return helper.ParseToken(token)
}
