package user

import (
	"qf"
	uDal "qf/mc/user/dal"
	uModel "qf/mc/user/model"
)

type UserBll struct {
	qf.BaseBll
	userDal       *uDal.UserDal       //用户dal
	userRole      *uDal.UserRoleDal   //用户-角色
	roleDal       *uDal.RoleDal       //角色dal
	roleRightsDal *uDal.RoleRightsDal //角色-权限
	rightsDal     *uDal.RoleRightsDal //权限dal
	rightsApiDal  *uDal.RightsApiDal  //权限-api
	dptDal        *uDal.DepartmentDal //部门dal
	dptUserDal    *uDal.DptUserDal    //部门-用户
}

func (u *UserBll) RegApi(api qf.ApiMap) {
	//注册用户API
	u.regUserApi(api)

	////注册角色API
	//u.regRoleApi(api)
	//
	////注册权限组API
	//u.regRightsApi(api)

}

func (u *UserBll) RegDal(dal qf.DalMap) {
	u.userDal = &uDal.UserDal{}
	dal.Reg(u.userDal, uModel.User{})
}

func (u *UserBll) RegMsg(msg qf.MessageMap) {
}

func (u *UserBll) RefBll() []qf.IBll {
	return nil
}

func (u *UserBll) Init() error {
	//u.initDefUser()
	return nil
}

func (u *UserBll) Stop() {

}
