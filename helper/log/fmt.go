package log

import (
	"fmt"
)

func NewLogByFmt() *ByFmt {
	return &ByFmt{}
}

type ByFmt struct {
}

func (l *ByFmt) Debug(title, content string) {
	fmt.Println(fmt.Sprintf("【Debug】%s:%s", title, content))
}

func (l *ByFmt) Info(title, content string) {
	fmt.Println(fmt.Sprintf("【Info】%s:%s", title, content))
}

func (l *ByFmt) Warn(title, content string) {
	fmt.Println(fmt.Sprintf("【Warn】%s:%s", title, content))
}

func (l *ByFmt) Error(title, content string) {
	fmt.Println(fmt.Sprintf("【Error】%s:%s", title, content))
}

func (l *ByFmt) Fatal(title, content string) {
	fmt.Println(fmt.Sprintf("【Fatal】%s:%s", title, content))
}

func (l *ByFmt) Any(logType, title, content string) {
	fmt.Println(fmt.Sprintf("【%s】%s:%s", logType, title, content))
}
