package patient

import (
	"qf"
	"qf/helper/content"
)

var (
	RouterPatient = qf.ApiRouter{Id: "patient", Explain: "患者信息业务"}
	RouterDict    = qf.ApiRouter{Id: "dict", Explain: "字典业务"}
)

type Bll struct {
	qf.BaseBll
	dal *dal
}

// RegApis 注册需要暴露的方法
func (b *Bll) RegApis(apis *qf.Apis) {
	apis.Reg(RouterPatient, b.submit, b.delete, b.getModel, fake())
	apis.CustomReg(RouterPatient, qf.EApiKindGet, "list", "按条件查询一组列表", b.getList, "")
}

// RegMessages 注册需要发送的自定义消息
func (b *Bll) RegMessages(messages *qf.Messages) {
	// 其他自定义消息
	//messages.Reg(RouterPatient, "CheckSample", "样本审核")
}

// RegReferences 注册需要引用的其他业务方法
func (b *Bll) RegReferences(references *qf.References) {
	// 通过字典业务获取字典数据
	references.Reg(RouterDict, qf.EApiKindGet, "")
}

// Init 业务初始化
func (b *Bll) Init() (err error) {
	b.dal = newDal(b.DB)
	return nil
}

// Stop 业务释放
func (b *Bll) Stop() {

}

func (b *Bll) submit(c content.Content) (interface{}, error) {
	// 保存内容
	nc, err := b.Content.Save(c)
	if err != nil {
		return nil, err
	}
	// 保存索引
	err = b.dal.updateIndexes(c, nc)
	if err != nil {
		return nil, err
	}
	return nc, nil
}

func (b *Bll) delete(c content.Content) (interface{}, error) {
	// 删除索引
	err := b.Content.Delete(c.ID)
	if err != nil {
		return nil, err
	}
	// 保存索引
	err = b.dal.updateIndexes(c, content.Content{})
	if err != nil {
		return nil, err
	}
	return c.ID, nil
}

func (b *Bll) getModel(c content.Content) (interface{}, error) {
	return b.Content.GetModel(c.ID)
}

func (b *Bll) getList(c content.Content) (interface{}, error) {
	query := map[string]interface{}{}
	err := c.BindQuery(query)
	if err != nil {
		return nil, err
	}
	return b.dal.selectList(query)
}
