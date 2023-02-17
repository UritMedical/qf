package patient

import (
	"encoding/json"
	"time"
)

func fake() string {
	info, _ := json.Marshal(struct {
		PatId   string
		PatNo   string
		PatName string
		PatSex  string
		PatAge  string
	}{
		PatId:   "123456789",
		PatNo:   "65123",
		PatName: "张三",
		PatSex:  "男",
		PatAge:  "20岁",
	})
	content, _ := json.Marshal(struct {
		ID   uint
		Time time.Time
		User string
		Info string
	}{
		ID:   0,
		Time: time.Now().Local(),
		User: "admin",
		Info: string(info),
	})
	return string(content)
}
