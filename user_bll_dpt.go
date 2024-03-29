package qf

import (
	"github.com/UritMedical/qf/util"
	"sort"
)

const maxCount = 100

//注册部门相关API
func (b *userBll) regDptApi(api ApiMap) {
	//部门
	api.Reg(EApiKindSave, "dpt", b.saveDpt)             //添加部门
	api.Reg(EApiKindDelete, "dpt", b.deleteDpt)         //删除部门
	api.Reg(EApiKindGetList, "dpts", b.getDpts)         //获取所有部门
	api.Reg(EApiKindGetModel, "dpt/tree", b.getDptTree) //获取部门组织树

	//部门-用户
	api.Reg(EApiKindSave, "dpt/users", b.setDptUsers)               //批量添加用户
	api.Reg(EApiKindDelete, "dpt/user", b.deleteDptUser)            //从部门中删除单个用户
	api.Reg(EApiKindGetList, "dpt/users", b.getDptAndChildDptUsers) //获取指定部门以及子部门的所有用户
	api.Reg(EApiKindGetList, "dpt/cur/users", b.getCurDptUsersOnly) //仅获取当前部门的用户

}

func (b *userBll) saveDpt(ctx *Context) (interface{}, IError) {
	dpt := &Dept{}
	if err := ctx.Bind(dpt); err != nil {
		return nil, err
	}
	if dpt.Id == 0 {
		dpt.Id = ctx.NewId(dpt)
	}
	return nil, b.dptDal.Save(&dpt)
}

func (b *userBll) deleteDpt(ctx *Context) (interface{}, IError) {
	uId := ctx.GetId()
	return nil, b.dptDal.Delete(uId)
}

//
// getDptTree
//  @Description: 获取部门树
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) getDptTree(ctx *Context) (interface{}, IError) {
	dptList := make([]Dept, 0)
	err := b.dptDal.GetList(0, maxCount, &dptList)
	if err != nil {
		return nil, nil
	}
	return b.buildDpTree(dptList), nil
}

////
//// buildTree
////  @Description: 创建部门树
////  @param departments
////  @return []*DeptNode
////
//func (b *userBll) buildTree() []DeptNode {
//	//获取所有部门
//	dptList := make([]Dept, 0)
//	err := b.dptDal.GetList(0, maxCount, &dptList)
//	if err != nil {
//		return nil
//	}
//	//转换成DepartNode数据格式
//	nodes := make([]DeptNode, 0)
//	for _, department := range dptList {
//		nodes = append(nodes, DeptNode{
//			Id:       department.Id,
//			Name:     department.Name,
//			ParentId: department.ParentId,
//			Children: nil,
//		})
//	}
//
//	//生成部门树
//	lookup := make(map[uint64]DeptNode)
//	for _, department := range nodes {
//		lookup[department.Id] = department
//		department.Children = []DeptNode{}
//	}
//
//	rootNodes := make([]DeptNode, 0)
//	for _, department := range nodes {
//		if department.ParentId == 0 {
//			rootNodes = append(rootNodes, department)
//		} else {
//			parent, ok := lookup[department.ParentId]
//			if !ok {
//				fmt.Printf("Invalid department: %v\n", department)
//			} else {
//				parent.Children = append(parent.Children, department)
//			}
//		}
//	}
//	return rootNodes
//}

//
// setDptUsers
//  @Description: 向指定部门批量添加用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) setDptUsers(ctx *Context) (interface{}, IError) {
	params := struct {
		DepartId uint64
		UserIds  []uint64
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, b.userDpDal.SetDptUsers(params.DepartId, params.UserIds)
}

//
// deleteDptUser
//  @Description: 删除部门中的用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) deleteDptUser(ctx *Context) (interface{}, IError) {
	DepartId := ctx.GetUIntValue("DeptId")
	UserId := ctx.GetUIntValue("UserId")
	return nil, b.userDpDal.RemoveUser(DepartId, UserId)
}

//
// getDpts
//  @Description: 获取所有部门
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) getDpts(ctx *Context) (interface{}, IError) {
	list := make([]Dept, 0)
	err := b.dptDal.GetList(0, maxCount, &list)
	return util.ToMaps(list), err
}

//func (b *userBll) getDptInfo(id uint64) (Dept, IError) {
//	dpt := Dept{}
//	err := b.dptDal.GetModel(id, &dpt)
//	if err != nil {
//		return Dept{}, err
//	}
//	return dpt, nil
//}

//func (b *userBll) getDptList(pId uint64) ([]Dept, IError) {
//	list, err := b.dptDal.GetAll()
//	if err != nil {
//		return make([]Dept, 0), err
//	}
//	if pId == 0 {
//		return list, nil
//	}
//	finals := make([]Dept, 0)
//	for _, l := range list {
//		if l.ParentId == pId {
//			finals = append(finals, l)
//		}
//	}
//	return finals, nil
//}

func (b *userBll) getCurDptUsersOnly(ctx *Context) (interface{}, IError) {
	departId := ctx.GetUIntValue("DepartId")
	uIds, _ := b.userDpDal.GetUsersByDptId(departId)
	list, _ := b.userDal.GetUsersByIds(uIds)

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
	return result, nil
}

//
// getDptAndChildDptUsers
//  @Description: 获取部门的用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) getDptAndChildDptUsers(ctx *Context) (interface{}, IError) {
	departId := ctx.GetUIntValue("DeptId")
	list, err := b.getDptAndSubDptUsers(departId)
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
// getDptAndSubDptUsers
//  @Description: 获取部门节点以及子部门的所有用户
//  @param dptId
//  @return []uint64
//
func (b *userBll) getDptAndSubDptUsers(departId uint64) ([]User, IError) {
	dptList := make([]Dept, 0)
	err := b.dptDal.GetList(0, maxCount, &dptList)
	if err != nil {
		return nil, nil
	}

	tree := b.buildDpTree(dptList)
	//通过递归找到对应的部门节点
	node := tree.FindDept(departId)
	if node == nil {
		return nil, Error(ErrorCodeRecordNotFound, "can't find department")
	}

	//通过递归找到此部门节点下所有用户
	uIdMap := make(map[uint64]string, 0) //利用map去重
	b.findChildrenUserIds(uIdMap, node)

	//map转换成切片
	userIds := make([]uint64, 0)
	for k := range uIdMap {
		userIds = append(userIds, k)
	}

	//排序
	sort.Slice(userIds, func(i, j int) bool {
		return userIds[i] < userIds[j]
	})

	//userId
	return b.userDal.GetUsersByIds(userIds)
}

//递归查找用户
func (b *userBll) findChildrenUserIds(uIdMap map[uint64]string, dptNode *DeptNode) {
	ids, _ := b.userDpDal.GetUsersByDptId(dptNode.Id)
	for _, id := range ids {
		uIdMap[id] = ""
	}
	if len(dptNode.Children) > 0 {
		for _, child := range dptNode.Children {
			b.findChildrenUserIds(uIdMap, child)
		}
	}
}

//
// getDepartsByUserId
//  @Description: 获取用户的所在部门
//  @receiver b
//  @param userId
//  @return []Dept
//  @return error
//
func (b *userBll) getDepartsByUserId(userId uint64) ([]Dept, IError) {
	dptIds, err := b.userDpDal.GetDptsByUserId(userId)
	if err != nil {
		return nil, err
	}
	list := make([]Dept, 0)
	err = b.dptDal.GetListByIN(dptIds, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

//
// getOrg
//  @Description: 获取用户所在的组织机构
//  @receiver b
//  @param userId
//  @return interface{}
//  @return IError
//
func (b *userBll) getOrg(userId uint64) ([]Dept, IError) {
	//获取用户的所在部门
	dptIds, err := b.userDpDal.GetDptsByUserId(userId)
	if err != nil {
		return nil, err
	}
	//获取所有部门
	dptList := make([]Dept, 0)
	err = b.dptDal.GetList(0, maxCount, &dptList)
	if err != nil {
		return nil, err
	}
	//转成map，便于做递归判断
	allDptMap := make(map[uint64]Dept, 0)
	for _, dpt := range dptList {
		allDptMap[dpt.Id] = dpt
	}

	//获取此用户所在机构列表
	orgMap := make(map[uint64]Dept, 0)
	for _, id := range dptIds {
		b.findParentDpt(id, orgMap, allDptMap)
	}

	//转换成数组
	ret := make([]Dept, 0)
	for _, department := range orgMap {
		ret = append(ret, department)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Id < ret[j].Id
	})
	return ret, nil
}

func (b *userBll) findParentDpt(dptId uint64, orgMap map[uint64]Dept, allDptMap map[uint64]Dept) {
	dpt, ok := allDptMap[dptId]
	if !ok {
		return
	}

	if dpt.ParentId == 0 {
		orgMap[dpt.Id] = dpt
	} else {
		orgMap[dpt.Id] = dpt
		b.findParentDpt(dpt.ParentId, orgMap, allDptMap)
	}
}
