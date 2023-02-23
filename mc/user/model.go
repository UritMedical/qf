package user

//User 用户信息
type User struct {
	Id       uint   `json:"id" gorm:"index"`
	LoginId  string `json:"loginId" gorm:"unique"` //登录账号
	Name     string `json:"name" gorm:"unique"`    //真实姓名
	Misc     string `json:"misc"`                  //各种其他用户信息
	Password string `json:"-"`                     //密码
	Role     byte   `json:"role"`                  //账号权限
	Remove   byte   `json:"-"`                     // 0-正常，1-删除
}

//Group 用户组
type Group struct {
	Id     uint   `json:"id" gorm:"index"`
	Name   string `json:"name" gorm:"unique"` //组名称
	Remove byte   `json:"-"`                  // 0-正常，1-删除
}

//Relation 用户与组的关系
type Relation struct {
	Id      uint `json:"id" gorm:"index"`
	GroupId uint `json:"groupId"` //组Id
	UserId  uint `json:"userId"`  //用户Id
}

//设置几个默认的角色，其他角色根据项目需要设置
//例如：医生、护士等
//不允许前端将权限设置为102
const (
	RoleDef       = 0   //默认
	RoleVisitor   = 100 //游客
	RoleAdmin     = 101 //管理员
	RoleDeveloper = 102 //开发者
)
