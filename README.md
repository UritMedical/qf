# QuickFrame

`QuickFrame` 是一个轻量化、模块化的快速医疗信息开发平台

# Catalogue

[TOC]

# Start

## Installation

```go
go get github.com/UritMedical/qf
```



## Usage

```go
import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/mc/patient"
	"github.com/UritMedical/qf/user"
)

func main() {
	// 启动
	qf.Run(regBll, nil)
}

func regBll(s *qf.Service) {
	// 注册相关业务
	s.RegBll(&your.Bll{}, "")
	...
    
    // 如果需要自定义扩展路由组，则使用如下方法
    s.RegBll(&your.Bll{}, "custom")
    // 默认：.../api/...  ->  设置后：.../api/custom/...
}
```

## Business Define

```go
// 定义一个业务结构体，继承qf.BaseBll，然后实现qf.IBLL接口
type CustomBll struct {
	qf.BaseBll
}

// -------------------------------------------------------------
// 然后实现qf.IBLL接口剩余方法
//   小技巧：Idea，将光标移动到结构体名称上面，按下Ctrl+I 或者 Alt+Enter选择实现接口
//          在弹出的窗口中输入IBLL，就可以快速补齐代码
func (b *CustomBll) RegApi(a qf.ApiMap) {
	// 注册需要暴露的API
}

func (b *CustomBll) RegDal(d qf.DalMap) {
	// 注册数据访问层
}

func (b *CustomBll) RegFault(f qf.FaultMap) {
	// 注册自定义故障码
}

func (b *CustomBll) RegMsg(m qf.MessageMap) {
	// 注册需要发送的消息
}

func (b *CustomBll) RegRef(r qf.RefMap) {
	// 注册需要引用的其他业务
}

func (b *CustomBll) Init() error {
	// 自定义初始化
	return nil
}

func (b *CustomBll) Stop() {
	// 自定义资源释放
}
// -------------------------------------------------------------
```

### 1、RegApi

```go
func (b *CustomBll) RegApi(a qf.ApiMap) {
	// a.Reg(类型, 相对路由, 路由方法)
	// 类型：EApiKindSave：提交类
	//      EApiKindDelete：删除类
	//      EApiKindGetModel：查询并返回单条
	//      EApiKindGetList：查询并返回列表
    
	// 生成路由 POST：api/custom
	a.Reg(qf.EApiKindSave, "custom", b.yourSave)
	// 生成路由 DELETE：api/custom
	a.Reg(qf.EApiKindDelete, "custom", b.yourDelete)
	// 生成路由 GET：api/custom/one
	a.Reg(qf.EApiKindGetModel, "custom/one", b.yourGetModel)
	// 生成路由 GET：api/custom/list
	a.Reg(qf.EApiKindGetList, "custom/list", b.yourGetList)
}

func (b *CustomBll) yourSave(ctx *qf.Context) (interface{}, qf.IError) {
	// ...
}

...
```

### 2、RegDal

```go
// 定义一个结构体，继承至qf.BaseDal，或者直接定义个qf.BaseDal对象
type CustomDal struct {
    qf.BaseDal
}

// 定义对应的实体层（即需要操作的数据表）
type CustomModel struct {
    qf.BaseModel
}

// 注册数据库访问层
type CustomBll struct {
	...
	dal1 *CustomDal
	dal2 *CustomDal
	dal3 *qf.BaseDal
}
func (b *CustomBll) RegDal(d qf.DalMap) {
	// 注册dal1，然后绑定结构体CustomModel
	// 绑定后，qf会自动创建一个CustomModel表结构
	b.dal1 = &CustomDal{}
	d.Reg(b.dal1, CustomModel{})
    
	// 如果不需要创建数据表，仅需要创建一个查询用的数据层，则传nil即可
	b.dal2 = &qf.CustomDal{}
	d.Reg(b.dal2, nil)
    
	// 如果仅仅想用框架提供的方法，则直接定义一个qf.BaseDal对象即可，不用继承
	b.dal3 = &qf.BaseDal{}
	d.Reg(b.dal3, CustomModel{})
}
```

### 3、RegFault

```go
// 定义故障码请从1000之后开始，1000之前是qf框架的通用故障码
const (
	ErrorCodeXXX1 = iota + 1000
	ErrorCodeXXX2 
)

// 注册故障码及对应的描述
func (b *CustomBll) RegFault(f qf.FaultMap) {
	f.Reg(ErrorCodeXXX1, "异常描述1")
	f.Reg(ErrorCodeXXX2, "异常描述2")
}
```

### 4、RegMsg

```go
// 注：此方法一般用于外挂业务或者第三方对接业务使用


// 注册消息执行方法，执行后，当被注册的API执行成功后，就会立即执行绑定的方法
func (b *CustomBll) RegMsg(m qf.MessageMap) {
	// m.Reg(其他业务暴露的类型, 其他业务暴露相对路由, 需要执行的方法)
    
	// 即：当 POST：api/ohterbll/edit api执行成功后，会立即执行yourFunc方法
	m.Reg(qf.EApiKindSave, "ohterbll/edit", b.yourFunc)
}
```

### 5、RegRef

```go
// qf原则上是不允许业务之间相互引用的，每个业务应该是独立
// 需要通过qf间接调用其他业务包


type CustomBll struct {
	...
	getDict qf.ApiHandler
}

// 注册第三方包
func (b *CustomBll) RegRef(r qf.RefMap) {
	b.getDict = r.Load(qf.EApiKindGetModel, "dict/getmodel")
}

// 使用第三方包
func (b *CustomBll) yourFunc(ctx *qf.Context) (interface{}, qf.IError) {
	// 直接调用方法并返回结果
	result, err := b.getDict(ctx.NewContext(yourInputParams))
}

```

### 6、Init and Stop

```go
func (b *CustomBll) Init() error {
	// 自定义初始化
    ...
	return nil
}

func (b *CustomBll) Stop() {
	// 自定义资源释放
    ...
}
```

# Example

## bll

```go
type Bll struct {
	qf.BaseBll
	infoDal *InfoDal
}

func (b *Bll) RegApi(a qf.ApiMap) {
	a.Reg(qf.EApiKindSave, "patient", b.SavePatient)         // 保存患者基本信息
	a.Reg(qf.EApiKindDelete, "patient", b.DeletePatient)     // 删除患者基本信息
	a.Reg(qf.EApiKindGetModel, "patient", b.GetPatient)      // 获取患者基本信息
}

func (b *Bll) RegDal(d qf.DalMap) {
	b.infoDal = &InfoDal{}
	d.Reg(b.infoDal, Patient{})
}

//
// SavePatient
//  @Description: 保存患者基本信息
//  @param ctx 输入结构(body)
//		{
//			"Id": 0,
//			"HisId": "院内唯一号",
//			"Name": "患者姓名",
//			"Sex": "患者性别",
//			"Birth": "出生日期",
//			"Phone": "联系电话",
//			"IDCard": "身份证号"
//		}
//  @return interface{} 患者唯一号
//  @return error 异常
func (b *Bll) SavePatient(ctx *qf.Context) (interface{}, qf.IError) {
	model := &Patient{}
	// 将前端提交的内容绑定到结构体
	if err := ctx.Bind(model); err != nil {
		return nil, err
	}
	// 获取唯一Id
	if model.Id == 0 {
		model.Id = ctx.NewId(model)
	}
	// 提交，如果HisId重复，则返回失败
	err := b.infoDal.Save(model)
	if err != nil {
		return 0, qf.Error(qf.ErrorCodeSaveFailure, err.Error())
	}
	return model.Id, nil
}

//
// DeletePatient
//  @Description: 删除患者基本信息
//  @param ctx 输入结构(query)
//		{
//			"Id": 需要删除的患者唯一Id,
//		}
//  @return interface{}
//  @return error
//
func (b *Bll) DeletePatient(ctx *qf.Context) (interface{}, qf.IError) {
	// 删除患者信息
	err := b.infoDal.Delete(ctx.GetId())
	return nil, qf.Error(qf.ErrorCodeDeleteFailure, err.Error())
}

//
// GetFull
//  @Description: 获取患者基本信息
//  @param ctx 输入结构(query)
//		{
//			"Id": 患者唯一Id,
//		}
//  @return interface{}
//  @return error
//
func (b *Bll) GetPatient(ctx *qf.Context) (interface{}, qf.IError) {
	// 通过ID检索
	patInfo := Patient{}
	err := b.infoDal.GetModel(ctx.GetId(), &patInfo)
	if err != nil || patInfo.Id == 0 {
		return nil, err
	}
	// 返回
	return util.ToMap(patInfo), nil
}
```

## dal

```go
type InfoDal struct {
	qf.BaseDal
}

// 通过条件检索
func (dal *InfoDal) GetListByKey(key string, dest interface{}) qf.IError {
	err := dal.DB().Where("HisId = ? or Name LIKE ?", key, "%"+key+"%").Find(dest).Error
	if err != nil {
		return qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return nil
}
```

## model

```go
type Patient struct {
	qf.BaseModel
	// 患者姓名，索引
	Name string `gorm:"index"`
	// HIS唯一号，唯一索引
	HisId *string `gorm:"unique"`
}
```

