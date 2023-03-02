package user

import (
	"fmt"
	"qf"
	uModel "qf/mc/user/model"
	uUtils "qf/mc/user/utils"
)

// DepartNode
// @Description: 部门树节点
//
type DepartNode struct {
	Id       uint
	Name     string
	ParentId uint
	Children []*DepartNode
}

//注册部门相关API
func (u *UserBll) regDptApi(api qf.ApiMap) {
	//部门
	api.Reg(qf.EKindSave, "dpt", u.saveDpt)        //添加部门
	api.Reg(qf.EKindDelete, "dpt", u.deleteDpt)    //删除部门
	api.Reg(qf.EKindGetModel, "dpt", u.getDptTree) //获取部门组织树

	//部门-用户
	api.Reg(qf.EKindSave, "dpt/user", u.setDptUser) //添加、删除用户
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
			Id:       department.Id,
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
	lookup := make(map[uint]*DepartNode)
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
// setDptUser
//  @Description: 前端传入调整后的人员，后端自己判断添加、删除的人员
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) setDptUser(ctx *qf.Context) (interface{}, error) {
	params := struct {
		DepartId uint
		UserIds  []uint
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	//获取此部门所有用户
	oldUsers, err := u.dptUserDal.GetUsersByDptId(params.DepartId)
	if err != nil {
		return nil, err
	}

	//分析需要新增、删除的用户Id
	newUsers := uUtils.DiffIntSet(params.UserIds, oldUsers)
	removeUsers := uUtils.DiffIntSet(oldUsers, params.UserIds)

	err = u.dptUserDal.AddRelation(params.DepartId, newUsers)
	if err != nil {
		return nil, err
	}
	err = u.dptUserDal.RemoveRelation(params.DepartId, removeUsers)
	return nil, err
}
