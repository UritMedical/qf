# QuickFrame

`QuickFrame` 是一个轻量化、模块化的快速医疗信息开发平台

# 快速开始

获取包

```go
go get github.com/UritMedical/qf
```

快速开始

```go
import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/mc/patient"
	"github.com/UritMedical/qf/mc/user"
)

func main() {
	// 启动
	qf.Run(regBll, nil)
}

func regBll(s *qf.Service) {
	// 注册框架提供的通用业务，位于mc文件夹内
	s.RegBll(&user.Bll{}, "api")    // 用户业务
	s.RegBll(&patient.Bll{}, "api") // 患者信息业务

	// 注册自定义业务
	// ...
}
```

