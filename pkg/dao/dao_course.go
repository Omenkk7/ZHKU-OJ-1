package dao

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// CreateCourse 创建课程
func (d *Dao) CreateCourse(ctx context.Context, course *models.Course) (id string, err error) {
	// 设置创建和修改时间
	now := time.Now().Unix()
	course.Ctime = now
	course.Mtime = now

	// 插入数据库
	return d.CreateOne(ctx, utils.CourseTable, course)
}

// GetCourseByID 根据ID获取课程信息
func (d *Dao) GetCourseByID(ctx context.Context, id string) (course *models.Course, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objID}
	log.Printf("查询条件： %v", query)

	// 初始化一个新的Course对象
	course = &models.Course{}

	// 使用FindOne方法查询
	err = d.mongo.FindOne(ctx, utils.CourseTable, query, course)

	// 检查错误类型
	if err != nil {
		// 如果是"没有文档"的错误，返回nil表示没找到
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("课程ID %s 不存在", id)
			return nil, nil
		}

		// 其他错误直接返回
		log.Printf("查询课程ID %s 时发生错误: %v", id, err)
		return nil, err
	}

	// 检查ID是否为空，如果为空说明没有找到记录
	if course.ID.IsZero() {
		log.Printf("课程ID %s 不存在 (ID为空)", id)
		return nil, nil
	}

	log.Printf("成功查询到课程: %s", course.Name)
	return course, nil
}

// GetCourseByCode 根据课程代码获取课程信息
func (d *Dao) GetCourseByCode(ctx context.Context, courseCode string) (course *models.Course, err error) {
	// 打印查询条件，便于调试
	log.Printf("查询条件： %v", bson.M{"course_code": courseCode})

	// 初始化一个新的Course对象
	course = &models.Course{}

	// 使用自定义的FindOne方法
	err = d.mongo.FindOne(ctx, utils.CourseTable, bson.M{"course_code": courseCode}, course)

	// 检查错误类型
	if err != nil {
		// 如果是"没有文档"的错误，返回nil表示没找到
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("课程代码 %s 不存在", courseCode)
			return nil, nil
		}

		// 其他错误直接返回
		log.Printf("查询课程代码 %s 时发生错误: %v", courseCode, err)
		return nil, err
	}

	// 检查ID是否为空，如果为空说明没有找到记录
	if course.ID.IsZero() {
		log.Printf("课程代码 %s 不存在 (ID为空)", courseCode)
		return nil, nil
	}

	log.Printf("成功查询到课程: %s", course.Name)
	return course, nil
}

// UpdateCourse 更新课程信息
func (d *Dao) UpdateCourse(ctx context.Context, id string, update bson.M) (err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// 添加更新时间
	if update["$set"] == nil {
		update["$set"] = bson.M{}
	}
	update["$set"].(bson.M)["mtime"] = time.Now().Unix()

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.CourseTable, selector, update)
	return err
}

// DeleteCourse 删除课程
func (d *Dao) DeleteCourse(ctx context.Context, id string) (err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	selector := bson.M{"_id": objID}
	_, err = d.DeleteOne(ctx, utils.CourseTable, selector)
	return err
}

// ArchiveCourse 归档课程（将状态设为0）
func (d *Dao) ArchiveCourse(ctx context.Context, id string) (err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status": 0,
			"mtime":  time.Now().Unix(),
		},
	}

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.CourseTable, selector, update)
	return err
}

// GetCourseList 获取课程列表
func (d *Dao) GetCourseList(ctx context.Context, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	return d.GetSome(ctx, utils.CourseTable, comQuery)
}

// AddCourseMember 添加课程成员
func (d *Dao) AddCourseMember(ctx context.Context, courseID string, member models.CourseMember, role int) error {
	objID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return err
	}

	roleStr := strconv.Itoa(role)
	update := bson.M{
		"$push": bson.M{
			"members." + roleStr: member,
		},
		"$set": bson.M{
			"mtime": time.Now().Unix(),
		},
	}

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.CourseTable, selector, update)
	return err
}

// RemoveCourseMember 移除课程成员
func (d *Dao) RemoveCourseMember(ctx context.Context, courseID string, userID string, role int) error {
	objID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return err
	}

	roleStr := strconv.Itoa(role)
	update := bson.M{
		"$pull": bson.M{
			"members." + roleStr: bson.M{"user_id": userID},
		},
		"$set": bson.M{
			"mtime": time.Now().Unix(),
		},
	}

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.CourseTable, selector, update)
	return err
}

// UpdateCourseMemberStatus 更新课程成员状态
func (d *Dao) UpdateCourseMemberStatus(ctx context.Context, courseID string, userID string, role int, status int) error {
	objID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return err
	}

	roleStr := strconv.Itoa(role)
	update := bson.M{
		"$set": bson.M{
			"members." + roleStr + ".$[elem].status": status,
			"mtime":                                  time.Now().Unix(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.user_id": userID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	selector := bson.M{"_id": objID}
	_, err = d.mongo.UpdateOne(ctx, utils.CourseTable, selector, update, opts)
	return err
}

// CheckUserCourseRole 检查用户在课程中的角色
func (d *Dao) CheckUserCourseRole(ctx context.Context, courseID string, userID string) ([]int, error) {
	_, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, err
	}

	// 获取课程信息
	course, err := d.GetCourseByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("课程不存在")
	}

	// 检查用户是否是创建者
	if course.Creator == userID {
		return []int{1}, nil // 创建者默认为管理员
	}

	// 检查用户在各角色中的存在情况
	var roles []int
	for roleStr, members := range course.Members {
		role, _ := strconv.Atoi(roleStr)
		for _, member := range members {
			if member.UserID == userID && member.Status == 1 { // 只检查状态为正常的成员
				roles = append(roles, role)
				break
			}
		}
	}

	return roles, nil
}

// CreateJoinCourseRequest 创建加入课程申请
func (d *Dao) CreateJoinCourseRequest(ctx context.Context, request *models.CourseJoinRequest) (id string, err error) {
	// 设置创建和修改时间
	now := time.Now().Unix()
	request.Ctime = now
	request.Mtime = now

	// 插入数据库
	return d.CreateOne(ctx, utils.CourseJoinRequestTable, request)
}

// GetJoinCourseRequestByID 根据ID获取加入课程申请
func (d *Dao) GetJoinCourseRequestByID(ctx context.Context, id string) (request *models.CourseJoinRequest, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	request = &models.CourseJoinRequest{}
	err = d.mongo.FindOne(ctx, utils.CourseJoinRequestTable, bson.M{"_id": objID}, request)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return request, nil
}

// UpdateJoinCourseRequest 更新加入课程申请
func (d *Dao) UpdateJoinCourseRequest(ctx context.Context, id string, update bson.M) (err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// 添加更新时间
	if update["$set"] == nil {
		update["$set"] = bson.M{}
	}
	update["$set"].(bson.M)["mtime"] = time.Now().Unix()

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.CourseJoinRequestTable, selector, update)
	return err
}

// GetJoinCourseRequestList 获取加入课程申请列表
func (d *Dao) GetJoinCourseRequestList(ctx context.Context, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	return d.GetSome(ctx, utils.CourseJoinRequestTable, comQuery)
}

// UpdateCourseStudentCount 更新课程学生数量
func (d *Dao) UpdateCourseStudentCount(ctx context.Context, courseID string) error {
	// 获取课程信息
	course, err := d.GetCourseByID(ctx, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("课程不存在")
	}

	// 计算学生数量
	studentCount := 0
	if students, ok := course.Members["4"]; ok {
		for _, student := range students {
			if student.Status == 1 { // 只计算状态为正常的学生
				studentCount++
			}
		}
	}

	// 更新课程的学生数量
	objID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"student_count": studentCount,
			"mtime":         time.Now().Unix(),
		},
	}

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.CourseTable, selector, update)
	return err
}

// GetCoursesByUserID 获取用户参与的课程列表
func (d *Dao) GetCoursesByUserID(ctx context.Context, userID string, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	// 构建查询条件，查找用户是创建者或成员的课程
	if comQuery.Filters == nil {
		comQuery.Filters = make(map[string]any)
	}

	// 使用$or操作符查询用户是创建者或任意角色成员的课程
	comQuery.Filters["$or"] = []bson.M{
		{"creator": userID},
		{"members.1.user_id": userID, "members.1.status": 1}, // 管理员
		{"members.2.user_id": userID, "members.2.status": 1}, // 教师
		{"members.3.user_id": userID, "members.3.status": 1}, // 助教
		{"members.4.user_id": userID, "members.4.status": 1}, // 学生
	}

	return d.GetSome(ctx, utils.CourseTable, comQuery)
}
