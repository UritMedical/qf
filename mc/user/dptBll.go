package user

import (
	"fmt"
	"qf"
	uModel "qf/mc/user/model"
)

// DepartNode
// @Description: 部门树节点
//
type DepartNode struct {
	Id       uint64
	Name     string
	ParentId uint64
	Children []*DepartNode
}

//注册部门相关API
func (u *UserBll) regDptApi(api qf.ApiMap) {
	//部门
	api.Reg(qf.EKindSave, "dpt", u.saveDpt)        //添加部门
	api.Reg(qf.EKindDelete, "dpt", u.deleteDpt)    //删除部门
	api.Reg(qf.EKindGetModel, "dpt", u.getDptTree) //获取部门组织树

	//部门-用户
	api.Reg(qf.EKindSave, "dpt/users", u.addDptUsers)    //批量添加用户
	api.Reg(qf.EKindDelete, "dpt/user", u.deleteDptUser) //从部门中删除单个用户

}

func (u *UserBll) saveDpt(ctx *qf.Context) (interface{}, error) {
	dpt := uModel.Department{}
	if err := ctx.Bind(&dpt); err != nil {
		return nil, err
	}
	return nil, u.dptDal.Save(&dpt)
}

func (u *UserBll) deleteDpt(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	err := u.dptDal.Delete(uId)
	return nil, err
}

//
// getDptTree
//  @Description: 获取部门树
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) getDptTree(ctx *qf.Context) (interface{}, error) {
	//获取素有部门
	dptList := make([]uModel.Department, 0)
	err := u.dptDal.GetList(0, 100, &dptList)
	if err != nil {
		return nil, err
	}
	nodes := make([]*DepartNode, 0)
	for _, department := range dptList {
		nodes = append(nodes, &DepartNode{
			Id:       uint64(department.Id),
			Name:     department.Name,
			ParentId: department.ParentId,
			Children: nil,
		})
	}
	dptTree := u.buildTree(nodes)
	return dptTree, nil
}

//
// buildTree
//  @Description: 创建部门树
//  @param departments
//  @return []*DepartNode
//
func (u *UserBll) buildTree(departments []*DepartNode) []*DepartNode {
	lookup := make(map[uint64]*DepartNode)
	for _, department := range departments {
		lookup[department.Id] = department
		department.Children = []*DepartNode{}
	}

	var rootNodes []*DepartNode
	for _, department := range departments {
		if department.ParentId == 0 {
			rootNodes = append(rootNodes, department)
		} else {
			parent, ok := lookup[department.ParentId]
			if !ok {
				fmt.Printf("Invalid department: %v\n", department)
			} else {
				parent.Children = append(parent.Children, department)
			}
		}
	}
	return rootNodes
}

//
// addDptUsers
//  @Description: 向指定部门批量添加用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) addDptUsers(ctx *qf.Context) (interface{}, error) {
	params := struct {
		DepartId uint64
		UserIds  []uint64
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, u.dptUserDal.AddUsers(params.DepartId, params.UserIds)
}

func (u *UserBll) deleteDptUser(ctx *qf.Context) (interface{}, error) {
	params := struct {
		DepartId uint64
		UserId   uint64
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, u.dptUserDal.RemoveUser(params.DepartId, params.UserId)
}
