/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 01:28
@Name: common.go
@Description:
*/

package utils

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
	"zhku-oj-server/pkg/models"
)

// Mongo查询相关

// RespPageQuery getSome(分页查询)返回结构体
type RespPageQuery struct {
	Items []*map[string]interface{} `json:"items" bson:"items"`
	Total int64                     `json:"total" bson:"total"`
}

// CommonQuery 查询条件模版
type CommonQuery struct {
	PageNum   int64  `query:"_pageNum"`
	PageSize  int64  `query:"_pageSize"`
	Sort      string `query:"_sort"`
	Direction int64  `query:"_direction"`
	Filters   map[string]any
}

// Query gin的query Params的别名
type Query map[string][]string

// BuildCommonQuery 查询参数转CommonQuery
func BuildCommonQuery(query Query) *CommonQuery {
	// TODO://添加按关键字模糊查询的功能
	// 目前已自动支持分页，按关键字排序，按关键字查询的基础功能
	// 查询串demo：_pageNum=2&_pageSize=10&_sort=_id&_direction=a&name=test
	// _pageNum为当前页数，_pageSize为每页大小，_sort为排序的字段，_direction为排序方向
	// name=test，按关键字查询，name等于test的所有记录

	cq := &CommonQuery{
		Sort:      DefaultSort,
		Direction: DefaultDirection,
		PageSize:  DefaultPageSize,
		PageNum:   DefaultPageNum,
		Filters:   make(map[string]any)}
	for k, v := range query {
		if k == Sort && len(v) > 0 {
			cq.Sort = v[0]
		} else if k == PageSize && len(v) > 0 {
			pageSize, _ := strconv.Atoi(v[0])
			cq.PageSize = int64(pageSize)
		} else if k == PageNum && len(v) > 0 {
			pageNum, _ := strconv.Atoi(v[0])
			cq.PageNum = int64(pageNum)
		} else if k == Direction {
			if v[0] == Desc {
				cq.Direction = Descending
			} else {
				cq.Direction = Ascending
			}
		} else if len(v) > 0 {
			cq.Filters[k] = v[0]
		}
	}
	return cq
}

// JsonTime 用于model结构体的Ctime和Mtime
type JsonTime struct {
	time.Time
}

func NewJsonTimeWithTimestamp(t int64) JsonTime {
	return JsonTime{Time: time.Unix(t, 0)}
}

func NewJsonTimeWithTime(t time.Time) JsonTime {
	return JsonTime{
		Time: t,
	}
}

func NowJsonTime() JsonTime {
	return JsonTime{
		Time: time.Now(),
	}
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// HashPassword 对密码进行哈希加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash 验证密码是否与哈希值匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type JWTClaims struct {
	ID       primitive.ObjectID `json:"_id"`
	Username string             `json:"username"`
	Role     int32              `json:"role"` // 0-> user; 1-> admin
	jwt.StandardClaims
}

func GenerateStringToken(user *models.User) (string, error) {
	claims := JWTClaims{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "urmsone",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(JwtTokenSecretKey))
	return token, err
}

// ParseToken 解析JWT
func ParseToken(tokenString string, secret string) (*JWTClaims, error) {
	claims := &JWTClaims{}

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

type ContextUser struct {
	ID       string
	Username string
	Role     int
}
