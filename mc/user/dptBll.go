package user

import (
	"errors"
	"fmt"
	"qf"
	uModel "qf/mc/user/model"
	"sort"
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

const maxCount = 100

//注册部门相关API
func (b *Bll) regDptApi(api qf.ApiMap) {
	//部门
	api.Reg(qf.EApiKindSave, "dpt", b.saveDpt)        //添加部门
	api.Reg(qf.EApiKindDelete, "dpt", b.deleteDpt)    //删除部门
	api.Reg(qf.EApiKindGetList, "dpt", b.getDpts)     //获取所有部门
	api.Reg(qf.EApiKindGetModel, "dpt", b.getDptTree) //获取部门组织树

	//部门-用户
	api.Reg(qf.EApiKindSave, "dpt/users", b.addDptUsers)    //批量添加用户
	api.Reg(qf.EApiKindDelete, "dpt/user", b.deleteDptUser) //从部门中删除单个用户
	api.Reg(qf.EApiKindGetList, "dpt/users", b.getDptUsers) //获取指定部门的所有用户

}

func (b *Bll) saveDpt(ctx *qf.Context) (interface{}, error) {
	dpt := uModel.Department{}
	if err := ctx.Bind(&dpt); err != nil {
		return nil, err
	}
	return nil, b.dptDal.Save(&dpt)
}

func (b *Bll) deleteDpt(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	ret, err := b.dptDal.Delete(uId)
	return ret, err
}

//
// getDptTree
//  @Description: 获取部门树
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) getDptTree(ctx *qf.Context) (interface{}, error) {
	return b.buildTree(), nil
}

//
// buildTree
//  @Description: 创建部门树
//  @param departments
//  @return []*DepartNode
//
func (b *Bll) buildTree() []*DepartNode {
	//获取所有部门
	dptList := make([]uModel.Department, 0)
	err := b.dptDal.GetList(0, maxCount, &dptList)
	if err != nil {
		return nil
	}
	//转换成DepartNode数据格式
	nodes := make([]*DepartNode, 0)
	for _, department := range dptList {
		nodes = append(nodes, &DepartNode{
			Id:       department.Id,
			Name:     department.Name,
			ParentId: department.ParentId,
			Children: nil,
		})
	}

	//生成部门树
	lookup := make(map[uint64]*DepartNode)
	for _, department := range nodes {
		lookup[department.Id] = department
		department.Children = []*DepartNode{}
	}

	rootNodes := make([]*DepartNode, 0)
	for _, department := range nodes {
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
func (b *Bll) addDptUsers(ctx *qf.Context) (interface{}, error) {
	params := struct {
		DepartId uint64
		UserIds  []uint64
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, b.dptUserDal.AddUsers(params.DepartId, params.UserIds)
}

//
// deleteDptUser
//  @Description: 删除部门中的用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) deleteDptUser(ctx *qf.Context) (interface{}, error) {
	DepartId := ctx.GetUIntValue("DepartId")
	UserId := ctx.GetUIntValue("UserId")
	return nil, b.dptUserDal.RemoveUser(DepartId, UserId)
}

//
// getDpts
//  @Description: 获取所有部门
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) getDpts(ctx *qf.Context) (interface{}, error) {
	list := make([]uModel.Department, 0)
	err := b.dptDal.GetList(0, maxCount, &list)
	return b.Maps(list), err
}

//
// getDptUsers
//  @Description: 获取部门的用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) getDptUsers(ctx *qf.Context) (interface{}, error) {
	departId := ctx.GetUIntValue("DepartId")
	users, err := b.getDptAndSubDptUsers(departId)
	return b.Maps(users), err
}

//
// getDptAndSubDptUsers
//  @Description: 获取部门节点以及子部门的所有用户
//  @param dptId
//  @return []uint64
//
func (b *Bll) getDptAndSubDptUsers(departId uint64) ([]uModel.User, error) {
	dptNodes := b.buildTree()
	//通过递归找到对应的部门节点
	node := b.findChildrenDpt(departId, dptNodes)

	if node == nil {
		return nil, errors.New("can't find node'")
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
func (b *Bll) findChildrenUserIds(uIdMap map[uint64]string, dptNode *DepartNode) {
	ids, _ := b.dptUserDal.GetUsersByDptId(dptNode.Id)
	for _, id := range ids {
		uIdMap[id] = ""
	}
	if len(dptNode.Children) > 0 {
		for _, child := range dptNode.Children {
			b.findChildrenUserIds(uIdMap, child)
		}
	}
}

//递归查找部门
func (b *Bll) findChildrenDpt(departId uint64, dptNodes []*DepartNode) *DepartNode {
	var targetNode *DepartNode
	for _, node := range dptNodes {
		if node.Id == departId {
			//	找到了部门节点
			targetNode = node
			break
		} else {
			targetNode = b.findChildrenDpt(departId, node.Children)
			if targetNode != nil {
				break
			}
		}
	}
	return targetNode
}
