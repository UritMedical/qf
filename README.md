# QuickFrame

`QuickFrame` 是一个轻量化、模块化的快速医疗信息开发平台

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

