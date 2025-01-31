package myUtils

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"zhku-oj-server/pkg/models"

	"github.com/dgrijalva/jwt-go"
)

// JWTCustomClaims 自定义JWT负载结构
type JWTCustomClaims struct {
	ID       primitive.ObjectID `json:"_id"`
	Username string             `json:"username"`
	Role     int32              `json:"role"` // 0-> admin; 1-> member
	jwt.StandardClaims
}

// GenerateToken 生成JWT
func GenerateToken(user *models.User, secret string, expirationTime time.Duration) (string, error) {
	// 设置过期时间
	expiration := time.Now().Add(expirationTime)

	// 创建自定义负载
	claims := &JWTCustomClaims{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	// 使用HS256算法创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥对token进行签名
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT
func ParseToken(tokenString string, secret string) (*JWTCustomClaims, error) {
	claims := &JWTCustomClaims{}

	// 解析token
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("token格式错误")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// 令牌已过期或未生效
				return nil, errors.New("token已过期或未生效")
			} else {
				return nil, errors.New("无法解析token")
			}
		}
	}

	return claims, nil
}
