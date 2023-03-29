package qio

import (
	"errors"
	"io"
	"io/ioutil"
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

// DeleteFile 删除文件
func DeleteFile(filename string) error {
	return os.Remove(filename)
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

//
// ReadAllString
//  @Description: 读取文件内容到字符串
//  @param filename
//  @return string
//  @return error
//
func ReadAllString(filename string) (string, error) {
	if !PathExists(filename) {
		return "", errors.New("file '" + filename + "' is not exist")
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

//
// ReadAllBytes
//  @Description: 读取文件内容到数组
//  @param filename
//  @return []byte
//  @return error
//
func ReadAllBytes(filename string) ([]byte, error) {
	if !PathExists(filename) {
		return nil, errors.New("file '" + filename + "' is not exist")
	}
	return ioutil.ReadFile(filename)
}

//
// WriteAllBytes
//  @Description: 写入字节数组，如果文件不存在则创建
//  @param filename
//  @param bytes
//  @param isAppend
//  @return error
//
func WriteAllBytes(filename string, bytes []byte, isAppend bool) error {
	f, err := readyToWrite(filename, isAppend)
	if err != nil || f == nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bytes)
	return err
}

//
// WriteString
//  @Description: 写入字符串，如果文件不存在则创建
//  @param filename
//  @param str
//  @param isAppend
//  @return error
//
func WriteString(filename string, str string, isAppend bool) error {
	f, err := readyToWrite(filename, isAppend)
	if err != nil || f == nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, str) //写入文件(字符串)
	return err
}

func readyToWrite(filename string, isAppend bool) (f *os.File, e error) {
	filename = GetFullPath(filename)
	// 创建文件夹
	CreateDirectory(filepath.Dir(filename))
	// 如果文件存在
	if PathExists(filename) {
		if isAppend {
			return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		}
		err := DeleteFile(filename)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
}
