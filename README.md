# QuickFrame

`QuickFrame` 是一个轻量化、模块化的快速医疗信息开发平台

# Start

获取包

```go
go get github.com/UritMedical/qf
```

快速开始

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
}
```

自定义路由组

```go
// 默认路由组为：.../api/...
s.RegBll(&your.Bll{}, "")

// 自定义路由组为：.../custom/...
s.RegBll(&your.Bll{}, "custom")
```

