package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
	uUtils "qf/mc/user/utils"
	"time"
)

//AESKey 注意：此密钥长度必须为16、32、64，否则生成加密模块时会报错
const AESKey = "wwuritcomlisurit"
const IV = "wwwuritcom123456"
const JwtSecretFile = "jwtSecret" //密钥存储的文件名称

//token有效期
const tokenExpireDuration = time.Hour * 24 * 3

var JwtSecret []byte //token密钥

//Claims token 信息
type Claims struct {
	Id      uint64
	RoleIds []uint64
	jwt.StandardClaims
}

//GenerateToken 生成token
func GenerateToken(id uint64, roleIds []uint64) (string, error) {
	claims := Claims{
		id,
		roleIds,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(JwtSecret)
}

//ParseToken 验证token的函数
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func InitJwtSecret() {
	//初始化token密钥
	jwt, err := DecodeJwtFromFile(JwtSecretFile, []byte(AESKey), []byte(IV))
	if err != nil || jwt == "" {
		jwtStr := uUtils.RandomString(32)
		JwtSecret = []byte(jwtStr)
		//将密钥进行AES加密后存入文件
		_ = EncryptAndWriteToFile(jwtStr, JwtSecretFile, []byte(AESKey), []byte(IV))
	} else {
		JwtSecret = []byte(jwt)
	}
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
