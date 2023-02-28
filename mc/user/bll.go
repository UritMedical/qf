package user

import (
	"qf"
	uDal "qf/mc/user/dal"
)

type UserBll struct {
	qf.BaseBll
	userDal       *uDal.UserDal
	userRole      *uDal.UserRoleDal
	roleDal       *uDal.RoleDal
	roleRightsDal *uDal.RoleRightsDal
	rightsDal     *uDal.RoleRightsDal
	rightsApiDal  *uDal.RightsApiDal
	dptDal        *uDal.DepartmentDal
	dptUserDal    *uDal.DptUserDal
}

func (u UserBll) RegApi(api qf.ApiMap) {
	//注册用户API
	u.regUserApi(api)

	//注册角色API
	u.regRoleApi(api)

}

func (u UserBll) RegDal(dal qf.DalMap) {
	u.userDal = &uDal.UserDal{}
	dal.Reg(u.userDal, uDal.UserDal{})
}

func (u UserBll) RegMsg(msg qf.MessageMap) {
}

func (u UserBll) RefBll() []qf.IBll {
	return nil
}

func (u UserBll) Init() error {
	u.initDefUser()
	return nil
}

func (u UserBll) Stop() {

}
