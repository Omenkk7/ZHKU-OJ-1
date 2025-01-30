/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 01:28
@Name: common.go
@Description:
*/

package utils

import (
	"database/sql/driver"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
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

// Bcrypt 用户密码取hash值保存
type Bcrypt struct {
	cost int
}

func (b *Bcrypt) Encode(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, b.cost)
}

func (b *Bcrypt) Match(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

var Encoder = Bcrypt{
	cost: bcrypt.DefaultCost,
}

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"` // 0-> admin; 1-> member
	jwt.StandardClaims
}

func GenerateStringToken(id string, username string, role time.Duration) (string, error) {
	claims := Claims{
		ID:       id,
		Username: username,
		Role:     int(role),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "urmsone",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(JwtTokenSecretKey))
	return token, err
}

func GetJwtTokenFromStringToken(token string) (c *jwt.Token, err error) {
	c, err = jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtTokenSecretKey), nil
	})
	return
}

type ContextUser struct {
	ID       string
	Username string
	Role     int
}
