package patient

import "qf"

type Bll struct {
	qf.BaseBll
	dal *Dal
}

func (b *Bll) RegApi(api qf.ApiMap) {
	//api.Reg(qf.EKindSave, "", b.dal.Save)            // Post http://.../.../patient
	//api.Reg(qf.EKindDelete, "", b.dal.Delete)        // Delete http://.../.../patient
	//api.Reg(qf.EKindGetModel, "", b.dal.GetModel)    // Get http://.../.../patient?id=1234
	//api.Reg(qf.EKindGetList, "list", b.dal.GetList)  // Get http://.../.../patients/list?startid=9999&maxcount=-1000
	//api.Reg(qf.EKindGetList, "search", b.dal.Search) // Get http://.../.../patients/search?patdh=654321&patname=张三
}

func (b *Bll) RegDal(dal qf.DalMap) {

}

func (b *Bll) RefBll() []qf.IBll {
	return nil
}

func (b *Bll) Init() error {
	return nil
}

func (b *Bll) Stop() {

}
