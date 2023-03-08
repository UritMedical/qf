package uUtils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

//ConvertToMD5 转换成MD5加密
func ConvertToMD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

//
// DiffIntSet
//  @Description: 计算a数组元素不在b数组之中的所有元素
//  @param a
//  @param b
//  @return []uint64
//
func DiffIntSet(a []uint64, b []uint64) []uint64 {
	c := make([]uint64, 0)
	temp := map[uint64]struct{}{}
	//把b所有的值作为key存入temp
	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}
	//如果a中的值作为key在temp中找不到，说明它不在b中
	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}
	return c
}

// RandomString
//  @Description: 生成随机字符串
//  @param length
//  @return string
//
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

//
// EncryptAndWriteToFile
//  @Description: 将字符串加密后写入文件
//  @param data
//  @param filename
//  @param key
//  @return error
//
func EncryptAndWriteToFile(data string, filename string, key, iv []byte) error {
	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	plaintext := []byte(data)

	// 创建AES块
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// 对明文进行填充
	plaintext = pad(plaintext)

	// 创建CBC模式加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密明文
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// 将加密后的数据写入文件
	if _, err := file.Write(ciphertext); err != nil {
		return err
	}
	return nil
}

//
// DecodeJwtFromFile
//  @Description: 从文件读取AES加密的内容进行解密
//  @param fileName
//  @param key
//  @param iv
//  @return string
//  @return error
//
func DecodeJwtFromFile(fileName string, key, iv []byte) (string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	// 创建AES块
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// 创建CBC模式解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密密文
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	// 去除填充
	decrypted = unpad(decrypted)

	return string(decrypted), nil
}

// 进行PKCS#7填充
func pad(plaintext []byte) []byte {
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// 去除PKCS#7填充
func unpad(plaintext []byte) []byte {
	length := len(plaintext)
	unpadding := int(plaintext[length-1])
	return plaintext[:(length - unpadding)]
}
