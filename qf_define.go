package qf

import (
	"fmt"
	"github.com/UritMedical/qf/util/qdate"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"time"
)

//
// IBll
//  @Description: 通用业务接口方法
//
type IBll interface {
	RegApi(a ApiMap)     // 注册需要暴露的API方法
	RegDal(d DalMap)     // 注册需要框架初始化的数据访问层对象
	RegFault(f FaultMap) // 注册故障码
	RegMsg(m MessageMap) // 注册需要接收处理的消息
	RegRef(r RefMap)     // 注册引用
	Init() error         // 业务自己的初始化方法
	Stop()               // 业务自己的资源释放方法
	// 框架内部实现的方法
	set(sub IBll, qfGroup, subGroup string)               // 将主服务的部分对象设置被基础业务
	key() string                                          // 获取业务唯一编号
	regApi(bind func(key string, handler ApiHandler))     // 框架注册方法
	regMsg(bind func(key string, handler MessageHandler)) // 框架注册方法
	regDal(db *gorm.DB)                                   // 框架注册方法
	regError(bind func(code int, err string))             // 框架注册方法
	regRef(getApi func(key string) ApiHandler)            // 框架注册方法
	Debug(content string)                                 // 调试日志
	//GetConfig() map[string]interface{}                              // 获取配置
	//SetConfig(value map[string]interface{}) (bool, error)           // 写入配置
}

//
// IDal
//  @Description: 数据访问层接口
//
type IDal interface {
	DB() *gorm.DB                                                   // 返回数据库对象
	Create(content interface{}) IError                              // 执行新增操作
	Save(content interface{}) IError                                // 执行新增或修改操作
	Delete(id uint64) IError                                        // 执行删除操作
	GetModel(id uint64, dest interface{}) IError                    // 根据Id获取单条信息
	GetList(startId uint64, maxCount uint, dest interface{}) IError // 根据起始Id和最大数量，获取一组信息
	GetCount(query interface{}, args ...interface{}) int64          // 根据条件获取数量
	CheckExists(id uint64) bool                                     // 检测Id是否存在
	// 框架内部实现的方法
	init(db *gorm.DB, model interface{})
}

// EApiKind 行为类别
type EApiKind string

var (
	EApiKindSave     EApiKind = "Save"     // 新增或修改
	EApiKindDelete   EApiKind = "Delete"   // 删除
	EApiKindGetModel EApiKind = "GetModel" // 获取单条
	EApiKindGetList  EApiKind = "GetList"  // 获取多条
)

//
// HttpMethod
//  @Description: 返回Http方式名称
//  @return string
//
func (kind EApiKind) HttpMethod() string {
	if kind == EApiKindSave {
		return "POST"
	}
	if kind == EApiKindDelete {
		return "DELETE"
	}
	return "GET"
}

// 注册相关函数定义
type (
	ApiHandler     func(ctx *Context) (interface{}, IError) // API函数指针
	MessageHandler func(ctx *Context) IError                // 消息函数指针
)

// 注册相关结构定义
type (
	ApiMap struct {
		bllName string
		dict    map[EApiKind]map[string]ApiHandler
	} // API字典
	DalMap struct {
		bllName string
		dict    map[IDal]interface{}
	} // 数据访问字典
	FaultMap struct {
		bllName string
		dict    map[int]string
	} // 异常字典
	MessageMap struct {
		bllName string
		dict    map[EApiKind]map[string]MessageHandler
	} // 消息字典
	RefMap struct {
		bllName string
		getApi  func(kind EApiKind, router string) ApiHandler
	} // 外部引用字典
)

//
// Reg
//  @Description: 注册对外暴露的方法
//  @param kind 类型
//  @param router 路由相对路径
//  @param handler 函数指针
//
func (api ApiMap) Reg(kind EApiKind, router string, handler ApiHandler) {
	if api.dict[kind] == nil {
		api.dict[kind] = make(map[string]ApiHandler)
	}
	if _, ok := api.dict[kind][router]; ok == false {
		api.dict[kind][router] = handler
	} else {
		panic(fmt.Sprintf("this router (%s,%s) already exists", kind, router))
	}
}

//
// Reg
//  @Description: 注册数据访问层，并初始化数据库
//  @param iDal 访问层对象
//  @param model 数据库表实体
//
func (d DalMap) Reg(iDal IDal, model interface{}) {
	t := reflect.TypeOf(iDal)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("this dal (%s) must be of type pointer", t.String()))
	}
	t = t.Elem()
	v := reflect.ValueOf(iDal)
	if v.IsNil() {
		panic(fmt.Sprintf("this dal (%s) has not been initialized", t.String()))
	}
	if _, ok := d.dict[iDal]; ok == false {
		d.dict[iDal] = model
	} else {
		panic(fmt.Sprintf("this dal (%s) already exists", t.String()))
	}
}

//
// Reg
//  @Description: 注册异常
//  @param code
//  @param desc
//
func (err FaultMap) Reg(code int, desc string) {
	if _, ok := err.dict[code]; ok == false {
		err.dict[code] = desc
	} else {
		panic(fmt.Sprintf("this fault code (%d,%s) already exists", code, desc))
	}
}

//
// Reg
//  @Description: 注册消息方法
//  @param kind 类型
//  @param router 路由相对路径
//  @param handler 函数指针
//
func (msg MessageMap) Reg(kind EApiKind, router string, handler MessageHandler) {
	if msg.dict[kind] == nil {
		msg.dict[kind] = make(map[string]MessageHandler)
	}
	if _, ok := msg.dict[kind][router]; ok == false {
		msg.dict[kind][router] = handler
	} else {
		panic(fmt.Sprintf("this msg (%s,%s )already exists", kind, router))
	}
}

//
// Load
//  @Description: 加载外部方法
//  @param kind 类型
//  @param router 路由相对路径
//  @return ApiHandler 方法函数指针
//
func (ref RefMap) Load(kind EApiKind, router string) ApiHandler {
	return ref.getApi(kind, router)
}

//
// BaseModel
//  @Description: 基础实体对象
//
type BaseModel struct {
	Id       uint64   `gorm:"primaryKey"` // 唯一号
	LastTime DateTime `gorm:"index"`      // 最后操作时间时间
	FullInfo string   // 内容
}

//
// LoginUser
//  @Description: 登陆用户信息
//
type LoginUser struct {
	UserId         uint64           // 登陆用户唯一号
	UserName       string           // 登陆用户名字
	LoginId        string           // 登陆用户账号
	Roles          []RoleInfo       // 所属的角色列表
	Departments    []DepartmentInfo // 所属部门列表
	DepartmentTree []DepartNode     // 所属部门树
	apis           map[string]byte  // 允许操作的api列表
	userBll        *userBll
}

func (u LoginUser) GetDps() {

	if u.userBll == nil {
		return
	}

	// 获取完整部门树
	tree := u.userBll.buildTree()

	// 获取当前用户所属部门
	dps, _ := u.userBll.getDepartsByUserId(u.UserId)

	fmt.Println(tree, dps)
}

func (u LoginUser) copyTo() LoginUser {
	user := LoginUser{
		UserId:         u.UserId,
		UserName:       u.UserName,
		LoginId:        u.LoginId,
		userBll:        u.userBll,
		Roles:          make([]RoleInfo, len(u.Roles)),
		Departments:    make([]DepartmentInfo, len(u.Departments)),
		DepartmentTree: make([]DepartNode, len(u.DepartmentTree)),
		apis:           map[string]byte{},
	}
	for i, r := range u.Roles {
		user.Roles[i] = r
	}
	for i, d := range u.Departments {
		user.Departments[i] = d
	}
	for i, d := range u.DepartmentTree {
		user.DepartmentTree[i] = d
	}
	for k, v := range u.apis {
		user.apis[k] = v
	}
	return user
}

//
// RoleInfo
//  @Description: 角色信息
//
type RoleInfo struct {
	Id   uint64
	Name string
}

//
// DepartmentInfo
//  @Description: 机构信息
//
type DepartmentInfo struct {
	Id       uint64
	Name     string
	ParentId uint64
}

//
// DepartNode
//  @Description: 部门树节点
//
type DepartNode struct {
	Id       uint64
	Name     string
	ParentId uint64
	Children []*DepartNode
}

var (
	dateFormat     string // 日期掩码
	dateTimeFormat string // 日期时间掩码
)

type Date uint32

//
// FromTime
//  @Description: 通过原生的time赋值
//  @param time
//
//goland:noinspection GoMixedReceiverTypes
func (d *Date) FromTime(time time.Time) {
	t := time.Local()
	s := fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
	v, _ := strconv.ParseUint(s, 10, 32)
	*d = Date(v)
}

//
// ToString
//  @Description: 根据全局format格式化输出
//  @return string
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) ToString() string {
	return qdate.ToString(d.ToTime(), dateFormat)
}

//
// ToTime
//  @Description: 转为原生时间对象
//  @return time.Time
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) ToTime() time.Time {
	str := fmt.Sprintf("%d", d)
	if len(str) != 8 {
		return time.Time{}
	}
	year, _ := strconv.Atoi(str[0:4])
	month, _ := strconv.Atoi(str[4:6])
	day, _ := strconv.Atoi(str[6:8])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

//
// MarshalJSON
//  @Description: 复写json转换
//  @return []byte
//  @return error
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", d.ToString())
	return []byte(str), nil
}

//
// UnmarshalJSON
//  @Description: 复写json转换
//  @param data
//  @return error
//
//goland:noinspection GoMixedReceiverTypes
func (d *Date) UnmarshalJSON(data []byte) error {
	v, err := qdate.ToNumber(string(data), dateFormat)
	if err == nil {
		*d = Date(v)
	}
	return err
}

type DateTime uint64

//
// FromTime
//  @Description: 通过原生的time赋值
//  @param time
//
//goland:noinspection GoMixedReceiverTypes
func (d *DateTime) FromTime(time time.Time) {
	t := time.Local()
	s := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	v, _ := strconv.ParseUint(s, 10, 64)
	*d = DateTime(v)
}

//
// ToString
//  @Description: 根据全局format格式化输出
//  @return string
//
//goland:noinspection GoMixedReceiverTypes
func (d DateTime) ToString() string {
	return qdate.ToString(d.ToTime(), dateTimeFormat)
}

//
// ToTime
//  @Description: 转为原生时间对象
//  @return time.Time
//
//goland:noinspection GoMixedReceiverTypes
func (d DateTime) ToTime() time.Time {
	str := fmt.Sprintf("%d", d)
	if len(str) != 14 {
		return time.Time{}
	}
	year, _ := strconv.Atoi(str[0:4])
	month, _ := strconv.Atoi(str[4:6])
	day, _ := strconv.Atoi(str[6:8])
	hour, _ := strconv.Atoi(str[8:10])
	minute, _ := strconv.Atoi(str[10:12])
	second, _ := strconv.Atoi(str[12:14])
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

//
// MarshalJSON
//  @Description: 复写json转换
//  @return []byte
//  @return error
//
//goland:noinspection GoMixedReceiverTypes
func (d DateTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", d.ToString())
	return []byte(str), nil
}

//
// UnmarshalJSON
//  @Description: 复写json转换
//  @param data
//  @return error
//
//goland:noinspection GoMixedReceiverTypes
func (d *DateTime) UnmarshalJSON(data []byte) error {
	v, err := qdate.ToNumber(string(data), dateTimeFormat)
	if err == nil {
		*d = DateTime(v)
	}
	return err
}

//
// IError
//  @Description: 异常
//
type IError interface {
	//
	// Code
	//  @Description: 获取故障码
	//  @return int
	//
	Code() int
	//
	// Error
	//  @Description: 获取异常描述
	//  @return string
	//
	Error() string
}

const (
	ErrorCodeParamInvalid        = iota + 100 // 传入参数无效
	ErrorCodePermissionDenied                 // 权限不足，拒绝访问
	ErrorCodeRecordNotFound                   // 未找到记录
	ErrorCodeRecordExist                      // 记录已经存在
	ErrorCodeSaveFailure                      // 保存失败
	ErrorCodeDeleteFailure                    // 删除失败
	ErrorCodeFileNotFound                     // 文件不存在
	ErrorCodeUploadedFileNull                 // 未上传任何文件
	ErrorCodeUploadedFileInvalid              // 上传文件解析失败
)

const (
	ErrorCodeOSError = 900 // 系统故障
	ErrorCodeUnknown = 999 // 未知异常
)

var errorCodeTextMap = map[int]string{
	ErrorCodeParamInvalid:        "无效的参数",
	ErrorCodePermissionDenied:    "权限不足，拒绝访问",
	ErrorCodeRecordNotFound:      "未找到记录",
	ErrorCodeRecordExist:         "记录已经存在",
	ErrorCodeSaveFailure:         "保存失败",
	ErrorCodeDeleteFailure:       "删除失败",
	ErrorCodeFileNotFound:        "指定文件不存在",
	ErrorCodeUploadedFileNull:    "未上传任何文件",
	ErrorCodeUploadedFileInvalid: "上传文件解析失败",
	ErrorCodeOSError:             "系统运行故障",
	ErrorCodeUnknown:             "其他未知故障",
}

type errorInfo struct {
	code  int
	error string
}

func (e errorInfo) Code() int {
	return e.code
}

func (e errorInfo) Error() string {
	return e.error
}

//
// File
//  @Description: 文件
//
type File struct {
	Name string // 文件名
	Size int64  // 文件大小
	Data []byte // 内容
}
