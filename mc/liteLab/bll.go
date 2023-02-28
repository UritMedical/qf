/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/26 9:28
 */

package liteLab

import "qf"

type Bll struct {
	qf.BaseBll
	resDal *LabResultDal
}

func (bll *Bll) RegApi(api qf.ApiMap) {
	//TODO implement me
	panic("implement me")
}

func (bll *Bll) RegDal(dal qf.DalMap) {
	dal.Reg(bll.resDal, LabResult{})

}

func (bll *Bll) RegMsg(msg qf.MessageMap) {

}

func (bll *Bll) RefBll() []qf.IBll {
	return nil
}

func (bll *Bll) Init() error {
	return nil
}

func (bll *Bll) Stop() {
}
