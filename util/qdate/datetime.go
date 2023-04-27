package qdate

import (
	"fmt"
	"regexp"
	"strconv"
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
	valueStr = strings.Trim(valueStr, "\"")
	layout := getLayout(formatStr)
	return time.ParseInLocation(layout, valueStr, time.Now().Location())
}

//
// ToNumber
//  @Description: 将时间字符串转为数值形式，如20230101
//  @return uint64
//  @return error
//
func ToNumber(valueStr, formatStr string) (uint64, error) {
	valueStr = strings.Trim(valueStr, "\"")
	// 才分格式化串
	layouts := strings.FieldsFunc(formatStr, func(r rune) bool {
		return r == '-' || r == '/' || r == ' ' || r == ':' ||
			r == '年' || r == '月' || r == '日' ||
			r == '时' || r == '分' || r == '秒'
	})
	layLen := len(layouts)
	// 然后从时间字符串中提取所有数值
	numMap := map[string]string{}
	exp := regexp.MustCompile(`\d+`).FindAllStringSubmatch(valueStr, -1)
	for i := 0; i < layLen; i++ {
		f := "%0" + fmt.Sprintf("%d", len(layouts[i])) + "d"
		if i < len(exp) && len(exp[i]) > 0 {
			n, err := strconv.Atoi(exp[i][0])
			if err == nil {
				numMap[layouts[i]] = fmt.Sprintf(f, n)
			}
		} else {
			numMap[layouts[i]] = fmt.Sprintf(f, 0)
		}
	}
	// 按照年月日时分秒重新排序
	final := ""
	for _, sort := range []string{"yyyy", "yy", "MM", "dd", "HH", "hh", "mm", "ss"} {
		if m, ok := numMap[sort]; ok {
			final += m
		}
	}
	// 最后返回数值
	if final == "00010101" {
		return 0, nil
	}
	return strconv.ParseUint(final, 10, 64)
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
