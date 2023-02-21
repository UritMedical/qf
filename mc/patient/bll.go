package patient

import "qf"

type Bll struct {
	qf.BaseBll
	dal *Dal
	//dictBll *dict.Bll
}

func (b *Bll) SetApis(api qf.Api) {
	api.Set(qf.EKindSave, "", b.dal.Save)            // Post http://.../.../patient
	api.Set(qf.EKindDelete, "", b.dal.Delete)        // Delete http://.../.../patient
	api.Set(qf.EKindGetModel, "", b.dal.GetModel)    // Get http://.../.../patient?id=1234
	api.Set(qf.EKindGetList, "list", b.dal.GetList)  // Get http://.../.../patients/list?startid=9999&maxcount=-1000
	api.Set(qf.EKindGetList, "search", b.dal.Search) // Get http://.../.../patients/search?patdh=654321&patname=张三
}

func (b *Bll) SetDal(dal qf.Dal) {
	// 注册dal和需要的实体，由框架自动创建数据表，并且框架将前端提交的内容转为该结构
	dal.Set(b.dal, Patient{})
}

func (b *Bll) SetReference(ref qf.Reference) {
	// 申明需要引用的其他包业务
	//ref.Set(dictBll)
}

func (b *Bll) Init() error {
	return nil
}

func (b *Bll) Stop() {

}
