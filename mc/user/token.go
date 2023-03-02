package user

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

//token有效期
const tokenExpireDuration = time.Hour * 24 * 3

//Claims token 信息
type Claims struct {
	Id      uint
	RoleIds []uint
	jwt.StandardClaims
}

//GenerateToken 生成token
func GenerateToken(id uint, roleIds []uint, jwtSecret []byte) (string, error) {
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

//ParseToken 验证token的函数
func (u *UserBll) ParseToken(token string, jwtSecret []byte) (*Claims, error) {
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
