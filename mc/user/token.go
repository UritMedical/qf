package user

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

//AESKey 注意：此密钥长度必须为16、32、64，否则生成加密模块时会报错
const AESKey = "wwuritcomlisurit"
const IV = "wwwuritcom123456"
const JwtSecretFile = "jwtSecret" //密钥存储的文件名称

//token有效期
const tokenExpireDuration = time.Hour * 24 * 3

//Claims token 信息
type Claims struct {
	Id      uint64
	RoleIds []uint64
	jwt.StandardClaims
}

//generateToken 生成token
func generateToken(id uint64, roleIds []uint64, jwtSecret []byte) (string, error) {
	claims := Claims{
		id,
		roleIds,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtSecret)
}

//parseToken 验证token的函数
func parseToken(token string, jwtSecret []byte) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
