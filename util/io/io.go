package io

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// PathExists 判断文件或文件夹是否存在
func PathExists(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// GetCurrentDirectory 获取程序的当前工作目录
func GetCurrentDirectory() (string, error) {
	if runtime.GOOS == "windows" {
		file, err := exec.LookPath(os.Args[0])
		if err != nil {
			return "", err
		}
		return filepath.Abs(file)
	}
	return filepath.Abs(os.Args[0])
}

// GetFullPath 获取完整路径
func GetFullPath(path string) string {
	// 将\\转为/
	path = strings.Replace(path, "\\", "/", -1)
	path = strings.Replace(path, "//", "/", -1)
	full, _ := filepath.Abs(path)
	return full
}

// CreateDirectory 创建文件夹
func CreateDirectory(path string) string {
	path = GetFullPath(path)
	if PathExists(path) == false {
		_ = os.MkdirAll(path, 0777)
	}
	return path
}

// IsDirectory 判断所给路径是否为文件夹
func IsDirectory(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDirectory(path)
}
