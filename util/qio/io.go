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

//
// PathExists
//  @Description: 判断文件或文件夹是否存在
//  @param path 路径
//  @return bool
//
func PathExists(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//
// GetCurrentDirectory
//  @Description: 获取程序的当前工作目录
//  @return string
//  @return error
//
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

//
// GetFullPath
//  @Description: 获取完整路径
//  @param path 相对路径
//  @return string
//
func GetFullPath(path string) string {
	// 将\\转为/
	full, _ := filepath.Abs(formatPath(path))
	return full
}

//
// GetFileName
//  @Description: 获取文件名
//  @param path
//  @return string
//
func GetFileName(path string) string {
	return filepath.Base(formatPath(path))
}

//
// GetFileExt
//  @Description: 获取文件的后缀名
//  @param path
//  @return string
//
func GetFileExt(path string) string {
	return filepath.Ext(formatPath(path))
}

//
// GetFileNameWithoutExt
//  @Description: 获取没有后缀名的文件名
//  @param path
//  @return string
//
func GetFileNameWithoutExt(path string) string {
	fileName := filepath.Base(formatPath(path))
	fileExt := filepath.Ext(fileName)
	return fileName[0 : len(fileName)-len(fileExt)]
}

//
// CreateDirectory
//  @Description: 创建文件夹
//  @param path
//  @return string
//
func CreateDirectory(path string) string {
	path = GetFullPath(path)
	if PathExists(path) == false {
		_ = os.MkdirAll(path, 0777)
	}
	return path
}

//
// DeleteFile
//  @Description: 删除文件
//  @param filename
//  @return error
//
func DeleteFile(filename string) error {
	return os.Remove(filename)
}

//
// IsDirectory
//  @Description: 判断所给路径是否为文件夹
//  @param path
//  @return bool
//
func IsDirectory(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

//
// IsFile
//  @Description: 判断所给路径是否为文件
//  @param path
//  @return bool
//
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

func formatPath(path string) string {
	path = filepath.Clean(path)
	path = strings.Replace(path, "\\", "/", -1)
	path = strings.Replace(path, "//", "/", -1)
	return path
}
