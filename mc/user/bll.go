package user

import (
	"qf"
	uDal "qf/mc/user/dal"
	uModel "qf/mc/user/model"
	uUtils "qf/mc/user/utils"
)

//开发者内置账号
//TODO 开发者密码可以配置
const devId = 310202303 //此Id用于权限判断时识别
var devUser = uModel.User{Content: qf.Content{Id: devId}, LoginId: "developer", Password: uUtils.ConvertToMD5([]byte("lisurit"))}

type UserBll struct {
	qf.BaseBll
	userDal       *uDal.UserDal       //用户dal
	userRoleDal   *uDal.UserRoleDal   //用户-角色
	roleDal       *uDal.RoleDal       //角色dal
	roleRightsDal *uDal.RoleRightsDal //角色-权限
	rightsDal     *uDal.RoleRightsDal //权限dal
	rightsApiDal  *uDal.RightsApiDal  //权限-api
	dptDal        *uDal.DepartmentDal //部门dal
	dptUserDal    *uDal.DptUserDal    //部门-用户
	jwtSecret     []byte              //token密钥
}

func (u *UserBll) RegApi(api qf.ApiMap) {
	u.regUserApi(api)   //注册用户API
	u.regRoleApi(api)   //注册角色API
	u.regRightsApi(api) //注册权限组API
	u.regDptApi(api)    //注册部门组织API
}

func (u *UserBll) RegDal(dal qf.DalMap) {
	u.userDal = &uDal.UserDal{}
	dal.Reg(u.userDal, uModel.User{})

	u.userRoleDal = &uDal.UserRoleDal{}
	dal.Reg(u.userRoleDal, uModel.UserRole{})

	u.roleDal = &uDal.RoleDal{}
	dal.Reg(u.roleDal, uModel.Role{})

	u.roleRightsDal = &uDal.RoleRightsDal{}
	dal.Reg(u.roleRightsDal, uModel.RoleRights{})

	u.rightsDal = &uDal.RoleRightsDal{}
	dal.Reg(u.rightsDal, uModel.RightsGroup{})

	u.rightsApiDal = &uDal.RightsApiDal{}
	dal.Reg(u.rightsApiDal, uModel.RightsApi{})

	u.dptDal = &uDal.DepartmentDal{}
	dal.Reg(u.dptDal, uModel.Department{})

	u.dptUserDal = &uDal.DptUserDal{}
	dal.Reg(u.dptUserDal, uModel.DepartUser{})
}

func (u *UserBll) RegMsg(msg qf.MessageMap) {
}

func (u *UserBll) RefBll() []qf.IBll {
	return nil
}

func (u *UserBll) Init() error {
	u.initDefUser()
	//TODO 使用随机字符串初始化token密钥
	u.jwtSecret = []byte("asldkfvnkwejfioweklasjfowienalv234Sdf23")
	return nil
}

func (u *UserBll) Stop() {

}

//
// initDefUser
//  @Description: 当用户表数量为0时，初始化默认账号
//
func (u *UserBll) initDefUser() {
	//创建admin,developer账号
	list := make([]uModel.User, 0)
	err := u.userDal.GetList(0, 10, &list)
	if err != nil {
		panic("can't create default user")
	}
	const adminId = 1
	if len(list) == 0 {
		_ = u.userDal.Save(&uModel.User{
			Content:  qf.Content{Id: adminId},
			LoginId:  "admin",
			Password: uUtils.ConvertToMD5([]byte("admin123"))})

		//创建默认角色
		_ = u.roleDal.Save(&uModel.Role{Content: qf.Content{Id: adminId}, Name: "administrator"})

		//分配角色
		_ = u.userRoleDal.SetRoleUsers(adminId, []uint{adminId}) //admin 分配 administrator角色

	}
}
