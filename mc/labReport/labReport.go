/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/18 9:47
 */

package labReport

import "qf"

type Bll struct {
	qf.BaseBll
}

func (b *Bll) RegApis(apis *qf.Apis) {
	//TODO implement me
	panic("implement me")
	apis.Reg("lab_report", sample.sumbit, sample.delete)
}

func (b *Bll) RegMessages(messages *qf.Messages) {
	//TODO implement me
	panic("implement me")
}

func (b *Bll) RegReferences(references *qf.References) {
	//TODO implement me
	panic("implement me")
}

func (b *Bll) Init() (err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Bll) Stop() {
	//TODO implement me
	panic("implement me")
}
