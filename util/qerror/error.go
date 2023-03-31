package qerror

import (
	"fmt"
	"github.com/UritMedical/qf/util/qdate"
	"github.com/UritMedical/qf/util/qio"
	"github.com/fatih/color"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//
// Recover
//  @Description: Panic的异常收集
//
func Recover(after func(err string)) {
	if r := recover(); r != nil {
		// 获取异常
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		stackInfo := string(buf[:n])

		// 输出异常
		log := ""
		fmt.Println("")
		color.New(color.FgWhite).PrintfFunc()(qdate.ToString(time.Now(), "yyyy-MM-dd HH:mm:ss"))
		log += qdate.ToString(time.Now(), "yyyy-MM-dd HH:mm:ss")
		color.New(color.FgRed, color.Bold).PrintfFunc()(" [ERROR] %s", r)
		log += fmt.Sprintf(" [ERROR] %s\n", r)
		fmt.Println("")
		lines := strings.Split(stackInfo, "\n")
		for i := 0; i < len(lines); i++ {
			line := strings.Replace(lines[i], "\t", "", -1)
			if strings.HasPrefix(line, "panic") {
				errStr := ""
				if i+3 < len(lines) {
					sp := strings.Split(strings.Replace(lines[i+3], "\t", "", -1), "+")
					errStr += fmt.Sprintf("   %-5s -> %s in %s\n", "curr", filepath.Base(lines[i+2]), sp[0])
				}
				if i+5 < len(lines) {
					sp := strings.Split(strings.Replace(lines[i+5], "\t", "", -1), "+")
					errStr += fmt.Sprintf("   %-5s -> %s in %s\n", "upper", filepath.Base(lines[i+4]), sp[0])
				}
				color.New(color.FgMagenta).PrintfFunc()("%s\n", errStr)
			}
			log += fmt.Sprintf("%s\n", lines[i])
		}
		// 写入日志
		logFile := fmt.Sprintf("%s/%s_Error.log", "./log", qdate.ToString(time.Now(), "yyyy-MM-dd"))
		logFile = qio.GetFullPath(logFile)
		log += "----------------------------------------------------------------------------------------------\n\n"
		_ = qio.WriteString(logFile, log, true)

		// 执行外部方法
		if after != nil {
			after(fmt.Sprintf("%s", r))
		}
	}
}
