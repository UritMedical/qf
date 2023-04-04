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

// Role
//  @Description: 角色
//
type Role struct {
	BaseModel
	Name string `gorm:"unique"` // 角色名称
}

// RolePermission
//  @Description: 角色权限组关系
//
type RolePermission struct {
	BaseModel
	RoleId       uint64 `gorm:"index"`
	PermissionId uint64 `gorm:"index"`
}

// Permission
//  @Description: 权限组
//
type Permission struct {
	BaseModel
	Name string `gorm:"unique"` //权限组名称
}

// PermissionApi
//  @Description: 权限组与API的关系
//
type PermissionApi struct {
	BaseModel
	PermissionId uint64 `gorm:"index"`
	Group        string // api 所属模块
	ApiId        string //API key
}

type ApiInfo struct {
	ApiId string
	Group string
}

//
// Department
//  @Description: 部门
//
type Department struct {
	BaseModel
	Name     string `gorm:"unique"` // 部门名称
	ParentId uint64 `gorm:"index"`  //父级部门Id
}

// DepartUser
//  @Description: 用户组织关系表
//
type DepartUser struct {
	BaseModel
	DepartId uint64 `gorm:"index"`
	UserId   uint64 `gorm:"index"`
}
