package user

import (
	"github.com/gin-gonic/gin"
	"qf"
)

type GroupBll struct {
	qf.BaseBll
	groupDal    *GroupDal
	relationDal *RelationDal
	userBll     *UserBll
}

func (g *GroupBll) RegApi(api qf.ApiMap) {
	api.Reg(qf.EKindSave, "group", g.saveGroup)
	api.Reg(qf.EKindSave, "group/adduser", g.addUserToGroup)
	api.Reg(qf.EKindDelete, "group/remove/user", g.removeUserFromGroup)
	api.Reg(qf.EKindGetList, "groups/users/all", g.getGroupsUsers)
	api.Reg(qf.EKindGetList, "groups/all", g.getAllGroups)
	api.Reg(qf.EKindGetList, "users/by/group", g.getUsersByGroupId)
	api.Reg(qf.EKindGetList, "groups/by/user", g.getGroupsByUserId)
}

func (g *GroupBll) RegDal(dal qf.DalMap) {
	dal.Reg(g.groupDal, GroupDal{})
	dal.Reg(g.relationDal, RelationDal{})
}

func (g *GroupBll) RegMsg(msg qf.MessageMap) {
}

func (g *GroupBll) RefBll() []qf.IBll {
	return []qf.IBll{g.userBll}
}

func (g *GroupBll) Init() error {
	return nil
}

func (g *GroupBll) Stop() {
}

//创建或者更新组信息
func (g *GroupBll) saveGroup(ctx *qf.Context) (interface{}, error) {
	var params Group
	ctx.BindModel(&params)
	_, err := g.groupDal.Save(params)
	return params, err
}

//向指定的组添加用户
func (g *GroupBll) addUserToGroup(ctx *qf.Context) (interface{}, error) {
	params := struct {
		GroupId uint   `json:"groupId"`
		UserIds []uint `json:"userIds"`
	}{}
	ctx.BindModel(&params)
	err := g.relationDal.AddUsersToGroup(params.GroupId, params.UserIds)
	return nil, err
}

//把组中指定的用户移除
func (g *GroupBll) removeUserFromGroup(ctx *qf.Context) (interface{}, error) {
	groupId := strToInt(ctx.GetValue("groupId"))
	params := struct {
		UserIds []uint `json:"userIds"`
	}{}
	ctx.BindModel(&params)
	err := g.relationDal.RemoveUsersFromGroup(groupId, params.UserIds)
	return nil, err
}

//获取所有组以及每个组所有的用户
func (g *GroupBll) getGroupsUsers(ctx *qf.Context) (interface{}, error) {
	//获取所有组
	all, err := g.groupDal.GetAll()
	if err != nil {
		return nil, err
	}
	list := make([]gin.H, 0)
	for _, group := range all {
		uIds, e := g.relationDal.GetUsersByGroup(group.Id)
		if e != nil {
			return nil, err
		}
		users := g.userBll.GetUsersByIds(uIds)
		item := gin.H{
			"groupId":   group.Id,
			"groupName": group.Name,
			"userList":  users,
		}
		list = append(list, item)
	}
	return list, err
}

//获取所有组
func (g *GroupBll) getAllGroups(ctx *qf.Context) (interface{}, error) {
	return g.groupDal.GetAll()
}

//获取指定组的所有用户
func (g *GroupBll) getUsersByGroupId(ctx *qf.Context) (interface{}, error) {
	gId := strToInt(ctx.GetValue("groupId"))

	//获取用户id数据
	uIds, err := g.relationDal.GetUsersByGroup(gId)
	if err != nil {
		return nil, err
	}
	//根据组Id列表获取组信息
	users := g.userBll.GetUsersByIds(uIds)
	return users, nil
}

//获取指定用户所在的组列表
func (g *GroupBll) getGroupsByUserId(ctx *qf.Context) (interface{}, error) {
	uId := strToInt(ctx.GetValue("userId"))
	gIds, err := g.relationDal.GetGroupsByUser(uId)
	if err != nil {
		return nil, err
	}
	return g.groupDal.GetGroupsByIds(gIds)
}
