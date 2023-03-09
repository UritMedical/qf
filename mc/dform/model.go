package dform

import "github.com/Urit-Mediacal/qf"

//DynamicForm 动态表单
type DynamicForm struct {
	qf.BaseModel
	Name string `gorm:"unique"` //动态表单名称
}
