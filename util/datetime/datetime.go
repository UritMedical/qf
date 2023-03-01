package datetime

import (
	"strings"
	"time"
)

//
// Parse
//  @Description: 字符串转换为时间
//  @param valueStr 时间
//  @param formatStr 格式化串 可以是yyyy/MM/dd HH:mm:ss 或者 yyyy-MM-dd HH:mm:ss:fff等
//  @return time.Time
//  @return error
//
func Parse(valueStr, formatStr string) (time.Time, error) {
	layout := getLayout(formatStr)
	return time.ParseInLocation(layout, valueStr, time.Now().Location())
}

//
// ToString
//  @Description: 转化为字符串
//  @param value 时间
//  @param formatStr 格式化串 可以是yyyy/MM/dd HH:mm:ss 或者 yyyy-MM-dd HH:mm:ss:fff等
//  @return string
//
func ToString(value time.Time, formatStr string) string {
	layout := getLayout(formatStr)
	return value.Format(layout)
}

func getLayout(formatStr string) string {
	//"2006-01-02 15:04:05"
	if strings.Contains(formatStr, "yyyy") {
		formatStr = strings.Replace(formatStr, "yyyy", "2006", 1)
	}
	if strings.Contains(formatStr, "yy") {
		formatStr = strings.Replace(formatStr, "yy", "06", 1)
	}
	if strings.Contains(formatStr, "YYYY") {
		formatStr = strings.Replace(formatStr, "YYYY", "2006", 1)
	}
	if strings.Contains(formatStr, "YY") {
		formatStr = strings.Replace(formatStr, "YY", "06", 1)
	}
	if strings.Contains(formatStr, "MM") {
		formatStr = strings.Replace(formatStr, "MM", "01", 1)
	}
	if strings.Contains(formatStr, "M") {
		formatStr = strings.Replace(formatStr, "M", "1", 1)
	}
	if strings.Contains(formatStr, "DD") {
		formatStr = strings.Replace(formatStr, "DD", "02", 1)
	}
	if strings.Contains(formatStr, "D") {
		formatStr = strings.Replace(formatStr, "D", "2", 1)
	}
	if strings.Contains(formatStr, "dd") {
		formatStr = strings.Replace(formatStr, "dd", "02", 1)
	}
	if strings.Contains(formatStr, "d") {
		formatStr = strings.Replace(formatStr, "d", "2", 1)
	}
	if strings.Contains(formatStr, "HH") {
		formatStr = strings.Replace(formatStr, "HH", "15", 1)
	}
	if strings.Contains(formatStr, "H") {
		formatStr = strings.Replace(formatStr, "H", "15", 1)
	}
	if strings.Contains(formatStr, "hh") {
		formatStr = strings.Replace(formatStr, "hh", "15", 1)
	}
	if strings.Contains(formatStr, "h") {
		formatStr = strings.Replace(formatStr, "h", "15", 1)
	}
	if strings.Contains(formatStr, "mm") {
		formatStr = strings.Replace(formatStr, "mm", "04", 1)
	}
	if strings.Contains(formatStr, "m") {
		formatStr = strings.Replace(formatStr, "m", "4", 1)
	}
	if strings.Contains(formatStr, "ss") {
		formatStr = strings.Replace(formatStr, "ss", "05", 1)
	}
	if strings.Contains(formatStr, "s") {
		formatStr = strings.Replace(formatStr, "s", "5", 1)
	}
	if strings.Contains(formatStr, "fff") {
		formatStr = strings.Replace(formatStr, "fff", "000", 1)
	}
	return formatStr
}
