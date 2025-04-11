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

// CreateClass 创建班级
func (d *Dao) CreateClass(ctx context.Context, class *models.Class) (id string, err error) {
	// 设置创建和修改时间
	now := time.Now().Unix()
	class.Ctime = now
	class.Mtime = now

	// 插入数据库
	return d.CreateOne(ctx, utils.ClassTable, class)
}

// GetClassByID 根据ID获取班级信息
func (d *Dao) GetClassByID(ctx context.Context, id string) (class *models.Class, err error) {
	/*class = &models.Class{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objID}
	result, err := d.GetOne(ctx, utils.ClassTable, class, query)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result.(*models.Class), nil*/
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objID}
	log.Printf("查询条件： %v", query)

	// 初始化一个新的Class对象
	class = &models.Class{}

	// 使用FindOne方法查询
	err = d.mongo.FindOne(ctx, utils.ClassTable, query, class)

	// 检查错误类型
	if err != nil {
		// 如果是"没有文档"的错误，返回nil表示没找到
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("班级ID %s 不存在", id)
			return nil, nil
		}

		// 其他错误直接返回
		log.Printf("查询班级ID %s 时发生错误: %v", id, err)
		return nil, err
	}

	// 检查ID是否为空，如果为空说明没有找到记录
	if class.ID.IsZero() {
		log.Printf("班级ID %s 不存在 (ID为空)", id)
		return nil, nil
	}

	log.Printf("成功查询到班级: %s", class.Name)
	return class, nil
}

// GetClassByCode 根据班级代码获取班级信息
func (d *Dao) GetClassByCode(ctx context.Context, classCode string) (class *models.Class, err error) {
	// 打印查询条件，便于调试
	log.Printf("查询条件： %v", bson.M{"class_code": classCode})

	// 初始化一个新的Class对象
	class = &models.Class{}

	// 使用自定义的FindOne方法
	err = d.mongo.FindOne(ctx, utils.ClassTable, bson.M{"class_code": classCode}, class)

	// 检查错误类型
	if err != nil {
		// 如果是"没有文档"的错误，返回nil表示没找到
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("班级代码 %s 不存在", classCode)
			return nil, nil
		}

		// 其他错误直接返回
		log.Printf("查询班级代码 %s 时发生错误: %v", classCode, err)
		return nil, err
	}

	// 检查ID是否为空，如果为空说明没有找到记录
	if class.ID.IsZero() {
		log.Printf("班级代码 %s 不存在 (ID为空)", classCode)
		return nil, nil
	}

	log.Printf("成功查询到班级: %s", class.Name)
	return class, nil
}

// UpdateClass 更新班级信息
func (d *Dao) UpdateClass(ctx context.Context, id string, update bson.M) (err error) {
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
	_, err = d.Update(ctx, utils.ClassTable, selector, update)
	return err
}

// DeleteClass 删除班级
func (d *Dao) DeleteClass(ctx context.Context, id string) (err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	selector := bson.M{"_id": objID}
	_, err = d.DeleteOne(ctx, utils.ClassTable, selector)
	return err
}

// ArchiveClass 归档班级（将状态设为0）
func (d *Dao) ArchiveClass(ctx context.Context, id string) (err error) {
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
	_, err = d.Update(ctx, utils.ClassTable, selector, update)
	return err
}

// GetClassList 获取班级列表
func (d *Dao) GetClassList(ctx context.Context, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	return d.GetSome(ctx, utils.ClassTable, comQuery)
}

// AddStudentToClass 添加学生到班级
func (d *Dao) AddStudentToClass(ctx context.Context, classStudent *models.ClassStudent) (id string, err error) {
	// 设置创建和修改时间
	now := time.Now().Unix()
	classStudent.Ctime = now
	classStudent.Mtime = now

	// 插入数据库
	return d.CreateOne(ctx, utils.ClassStudentTable, classStudent)
}

// RemoveStudentFromClass 从班级移除学生
func (d *Dao) RemoveStudentFromClass(ctx context.Context, classID, studentID string) (err error) {
	selector := bson.M{
		"class_id":   classID,
		"student_id": studentID,
		"status":     1, // 只移除状态为正常的
	}

	update := bson.M{
		"$set": bson.M{
			"status":     0,
			"leave_time": time.Now().Unix(),
			"mtime":      time.Now().Unix(),
		},
	}

	_, err = d.Update(ctx, utils.ClassStudentTable, selector, update)
	return err
}

// GetClassStudents 获取班级学生列表
func (d *Dao) GetClassStudents(ctx context.Context, classID string, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	// 添加班级ID筛选条件
	if comQuery.Filters == nil {
		comQuery.Filters = make(map[string]any)
	}
	comQuery.Filters["class_id"] = classID

	return d.GetSome(ctx, utils.ClassStudentTable, comQuery)
}

// CheckStudentInClass 检查学生是否在班级中
func (d *Dao) CheckStudentInClass(ctx context.Context, classID, studentID string) (bool, error) {
	query := bson.M{
		"class_id":   classID,
		"student_id": studentID,
		"status":     1, // 状态为正常
	}

	count, err := d.mongo.Count(ctx, utils.ClassStudentTable, query)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateJoinRequest 创建加入班级申请
func (d *Dao) CreateJoinRequest(ctx context.Context, request *models.ClassJoinRequest) (id string, err error) {
	// 设置创建和修改时间
	now := time.Now().Unix()
	request.Ctime = now
	request.Mtime = now

	// 插入数据库
	return d.CreateOne(ctx, utils.ClassJoinRequestTable, request)
}

// GetJoinRequestByID 根据ID获取加入申请
func (d *Dao) GetJoinRequestByID(ctx context.Context, id string) (request *models.ClassJoinRequest, err error) {
	request = &models.ClassJoinRequest{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objID}
	result, err := d.GetOne(ctx, utils.ClassJoinRequestTable, request, query)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result.(*models.ClassJoinRequest), nil
}

// UpdateJoinRequest 更新加入申请
func (d *Dao) UpdateJoinRequest(ctx context.Context, id string, update bson.M) (err error) {
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
	_, err = d.Update(ctx, utils.ClassJoinRequestTable, selector, update)
	return err
}

// GetJoinRequestList 获取加入申请列表
func (d *Dao) GetJoinRequestList(ctx context.Context, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	return d.GetSome(ctx, utils.ClassJoinRequestTable, comQuery)
}

// GetStudentClassCount 获取学生所在班级数量
func (d *Dao) GetStudentClassCount(ctx context.Context, studentID string) (int64, error) {
	query := bson.M{
		"student_id": studentID,
		"status":     1, // 状态为正常
	}

	return d.mongo.Count(ctx, utils.ClassStudentTable, query)
}

// UpdateClassStudentCount 更新班级学生数量
func (d *Dao) UpdateClassStudentCount(ctx context.Context, classID string) error {
	// 获取班级中状态为正常的学生数量
	query := bson.M{
		"class_id": classID,
		"status":   1, // 状态为正常
	}

	count, err := d.mongo.Count(ctx, utils.ClassStudentTable, query)
	if err != nil {
		return err
	}

	// 更新班级的学生数量
	objID, err := primitive.ObjectIDFromHex(classID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"student_count": count,
			"mtime":         time.Now().Unix(),
		},
	}

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.ClassTable, selector, update)
	return err
}

// BatchAddStudentsToClass 批量添加学生到班级
func (d *Dao) BatchAddStudentsToClass(ctx context.Context, classStudents []models.ClassStudent) error {
	if len(classStudents) == 0 {
		return nil
	}

	// 转换为interface{}切片
	docs := make([]interface{}, len(classStudents))
	for i, cs := range classStudents {
		// 设置创建和修改时间
		now := time.Now().Unix()
		cs.Ctime = now
		cs.Mtime = now
		docs[i] = cs
	}

	// 批量插入，使用BulkInsert替代InsertMany
	_, err := d.mongo.BulkInsert(ctx, utils.ClassStudentTable, true, docs...)
	return err
}

// GetClassesByStudentID 获取学生所在的班级列表
func (d *Dao) GetClassesByStudentID(ctx context.Context, studentID string, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	// 先查询学生所在的班级ID
	studentQuery := bson.M{
		"student_id": studentID,
		"status":     1, // 状态为正常
	}

	var classStudents []models.ClassStudent
	err = d.mongo.Find(ctx, utils.ClassStudentTable, studentQuery, &classStudents)
	if err != nil {
		return nil, err
	}

	if len(classStudents) == 0 {
		return &utils.RespPageQuery{
			Items: make([]*map[string]interface{}, 0),
			Total: 0,
		}, nil
	}

	// 提取班级ID
	classIDs := make([]primitive.ObjectID, 0, len(classStudents))
	for _, cs := range classStudents {
		objID, err := primitive.ObjectIDFromHex(cs.ClassID)
		if err != nil {
			continue
		}
		classIDs = append(classIDs, objID)
	}

	// 查询班级信息
	if comQuery.Filters == nil {
		comQuery.Filters = make(map[string]any)
	}
	comQuery.Filters["_id"] = bson.M{"$in": classIDs}

	return d.GetSome(ctx, utils.ClassTable, comQuery)
}

// GetClassesByCourseID 获取绑定了指定课程的班级列表
func (d *Dao) GetClassesByCourseID(ctx context.Context, courseID string, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	if comQuery.Filters == nil {
		comQuery.Filters = make(map[string]any)
	}
	comQuery.Filters["courses.course_id"] = courseID

	return d.GetSome(ctx, utils.ClassTable, comQuery)
}

// AddCoursesToClass 添加课程到班级
func (d *Dao) AddCoursesToClass(ctx context.Context, classID string, courses []models.CourseInfo) error {
	objID, err := primitive.ObjectIDFromHex(classID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"courses": bson.M{
				"$each": courses,
			},
		},
		"$set": bson.M{
			"mtime": time.Now().Unix(),
		},
	}

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.ClassTable, selector, update)
	return err
}

// RemoveCourseFromClass 从班级移除课程
func (d *Dao) RemoveCourseFromClass(ctx context.Context, classID string, courseID string) error {
	objID, err := primitive.ObjectIDFromHex(classID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$pull": bson.M{
			"courses": bson.M{
				"course_id": courseID,
			},
		},
		"$set": bson.M{
			"mtime": time.Now().Unix(),
		},
	}

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.ClassTable, selector, update)
	return err
}

// UpdateCourseStatusInClass 更新班级中课程的状态
func (d *Dao) UpdateCourseStatusInClass(ctx context.Context, classID string, courseID string, status int) error {
	objID, err := primitive.ObjectIDFromHex(classID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"courses.$[elem].status": status,
			"mtime":                  time.Now().Unix(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.course_id": courseID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	selector := bson.M{"_id": objID}
	_, err = d.mongo.UpdateOne(ctx, utils.ClassTable, selector, update, opts)
	return err
}

// AddClassMember 添加班级成员
func (d *Dao) AddClassMember(ctx context.Context, classID string, member models.ClassMember, role int) error {
	objID, err := primitive.ObjectIDFromHex(classID)
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
	_, err = d.Update(ctx, utils.ClassTable, selector, update)
	return err
}

// RemoveClassMember 移除班级成员
func (d *Dao) RemoveClassMember(ctx context.Context, classID string, userID string, role int) error {
	objID, err := primitive.ObjectIDFromHex(classID)
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
	_, err = d.Update(ctx, utils.ClassTable, selector, update)
	return err
}

// CheckUserClassRole 检查用户在班级中的角色
func (d *Dao) CheckUserClassRole(ctx context.Context, classID string, userID string) ([]int, error) {
	_, err := primitive.ObjectIDFromHex(classID)
	if err != nil {
		return nil, err
	}

	// 获取班级信息
	class, err := d.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("班级不存在")
	}

	// 检查用户是否是创建者
	if class.Creator == userID {
		return []int{1}, nil // 创建者默认为管理员
	}

	// 检查用户在各角色中的存在情况
	var roles []int
	for roleStr, members := range class.Members {
		role, _ := strconv.Atoi(roleStr)
		for _, member := range members {
			if member.UserID == userID {
				roles = append(roles, role)
				break
			}
		}
	}

	// 检查用户是否是学生
	isStudent, err := d.CheckStudentInClass(ctx, classID, userID)
	if err != nil {
		return nil, err
	}
	if isStudent {
		roles = append(roles, 4) // 4-学生
	}

	return roles, nil
}
