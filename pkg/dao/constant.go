/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 21:04
@Name: constant.go
@Description:
*/

package dao

const (
	// mongo指令常量
	mongoID        = "_id"
	mongoSet       = "$set"
	mongoAddToSet  = "$addToSet"
	mongoPull      = "$pull"
	mongoGTE       = "$gte"
	mongoOr        = "$or"
	mongoIn        = "$in"
	mongoElemMatch = "$elemMatch"
	mongoRegex     = "$regex"
	// mongo field常量

	// mongo collection常量
	userTable = "users"
)
