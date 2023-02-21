package content

import (
	"encoding/json"
	"time"
)

type Content struct {
	ID         uint      `json:"id,omitempty"` // 内容唯一好
	LatestTime time.Time `json:"time"`         // 操作时间
	IsDelete   uint
	Info       string `json:"info"` // 内容完整信息，一般为json
}

func (c Content) BindJson(model interface{}) error {
	return json.Unmarshal([]byte(c.Info), &model)
}

func (c Content) BindQuery(model interface{}) error {
	return json.Unmarshal([]byte(c.Info), &model)
}
