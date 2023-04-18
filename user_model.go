package qf

//
// User
//  @Description: 用户信息
//
type User struct {
	BaseModel
	LoginId  string `gorm:"unique"` //登录Id
	Password string `json:"-"`      //密码
}

// UserRole
//  @Description: 用户角色关系
//
type UserRole struct {
	BaseModel
	UserId uint64 `gorm:"index"`
	RoleId uint64 `gorm:"index"`
}

// UserDP
//  @Description: 用户组织关系表
//
type UserDP struct {
	BaseModel
	DpId   uint64 `gorm:"index"`
	UserId uint64 `gorm:"index"`
}

// Role
//  @Description: 角色
//
type Role struct {
	BaseModel
	Name string `gorm:"unique"` // 角色名称
}

// RoleApi
//  @Description: 角色Api关系
//
type RoleApi struct {
	BaseModel
	RoleId     uint64 `gorm:"index"`
	Url        string // API路由
	Permission string // API权限，r-只读，rw-读写
}

//
// Department
//  @Description: 部门
//
type Department struct {
	BaseModel
	Name     string `gorm:"unique"` // 部门名称
	ParentId uint64 `gorm:"index"`  // 父级部门Id
}
