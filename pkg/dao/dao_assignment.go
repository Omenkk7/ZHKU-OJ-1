package dao

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"time"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// CreateAssignment 创建作业
func (d *Dao) CreateAssignment(ctx context.Context, assignment *models.Assignment) (id string, err error) {
	// 设置创建和修改时间
	now := time.Now().Unix()
	assignment.Ctime = now
	assignment.Mtime = now

	// 插入数据库
	return d.CreateOne(ctx, utils.AssignmentTable, assignment)
}

// GetAssignmentByID 根据ID获取作业信息
func (d *Dao) GetAssignmentByID(ctx context.Context, id string) (assignment *models.Assignment, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objID}
	log.Printf("查询条件： %v", query)

	// 初始化一个新的Assignment对象
	assignment = &models.Assignment{}

	// 使用FindOne方法查询
	err = d.mongo.FindOne(ctx, utils.AssignmentTable, query, assignment)

	// 检查错误类型
	if err != nil {
		// 如果是"没有文档"的错误，返回nil表示没找到
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("作业ID %s 不存在", id)
			return nil, nil
		}

		// 其他错误直接返回
		log.Printf("查询作业ID %s 时发生错误: %v", id, err)
		return nil, err
	}

	// 检查ID是否为空，如果为空说明没有找到记录
	if assignment.ID.IsZero() {
		log.Printf("作业ID %s 不存在 (ID为空)", id)
		return nil, nil
	}

	log.Printf("成功查询到作业: %s", assignment.Title)
	return assignment, nil
}

// GetAssignmentByCode 根据作业代码获取作业信息
func (d *Dao) GetAssignmentByCode(ctx context.Context, assignmentCode string) (assignment *models.Assignment, err error) {
	// 打印查询条件，便于调试
	log.Printf("查询条件： %v", bson.M{"assignment_code": assignmentCode})

	// 初始化一个新的Assignment对象
	assignment = &models.Assignment{}

	// 使用自定义的FindOne方法
	err = d.mongo.FindOne(ctx, utils.AssignmentTable, bson.M{"assignment_code": assignmentCode}, assignment)

	// 检查错误类型
	if err != nil {
		// 如果是"没有文档"的错误，返回nil表示没找到
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("作业代码 %s 不存在", assignmentCode)
			return nil, nil
		}

		// 其他错误直接返回
		log.Printf("查询作业代码 %s 时发生错误: %v", assignmentCode, err)
		return nil, err
	}

	// 检查ID是否为空，如果为空说明没有找到记录
	if assignment.ID.IsZero() {
		log.Printf("作业代码 %s 不存在 (ID为空)", assignmentCode)
		return nil, nil
	}

	log.Printf("成功查询到作业: %s", assignment.Title)
	return assignment, nil
}

// UpdateAssignment 更新作业信息
func (d *Dao) UpdateAssignment(ctx context.Context, id string, update bson.M) (err error) {
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
	_, err = d.Update(ctx, utils.AssignmentTable, selector, update)
	return err
}

// DeleteAssignment 删除作业
func (d *Dao) DeleteAssignment(ctx context.Context, id string) (err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	selector := bson.M{"_id": objID}
	_, err = d.DeleteOne(ctx, utils.AssignmentTable, selector)
	return err
}

// ArchiveAssignment 归档作业（将状态设为0）
func (d *Dao) ArchiveAssignment(ctx context.Context, id string) (err error) {
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
	_, err = d.Update(ctx, utils.AssignmentTable, selector, update)
	return err
}

// GetAssignmentList 获取作业列表
func (d *Dao) GetAssignmentList(ctx context.Context, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	return d.GetSome(ctx, utils.AssignmentTable, comQuery)
}

// GetAssignmentsByCourseID 获取课程下的作业列表
func (d *Dao) GetAssignmentsByCourseID(ctx context.Context, courseID string, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	if comQuery.Filters == nil {
		comQuery.Filters = make(map[string]any)
	}
	comQuery.Filters["course_id"] = courseID

	return d.GetSome(ctx, utils.AssignmentTable, comQuery)
}

// GetAssignmentsByClassID 获取班级关联的作业列表
func (d *Dao) GetAssignmentsByClassID(ctx context.Context, classID string, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	if comQuery.Filters == nil {
		comQuery.Filters = make(map[string]any)
	}
	comQuery.Filters["class_ids.class_id"] = classID

	return d.GetSome(ctx, utils.AssignmentTable, comQuery)
}

// SubmitAssignment 学生提交作业
func (d *Dao) SubmitAssignment(ctx context.Context, assignmentID string, studentID string, submission *models.Submission) error {
	objID, err := primitive.ObjectIDFromHex(assignmentID)
	if err != nil {
		return err
	}

	// 设置提交时间
	submission.SubmitTime = time.Now().Unix()
	submission.Status = 1 // 已提交未批改

	// 更新学生的提交信息
	update := bson.M{
		"$set": bson.M{
			"members.4.$[elem].submission": submission,
			"mtime":                        time.Now().Unix(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.user_id": studentID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	selector := bson.M{"_id": objID}
	_, err = d.mongo.UpdateOne(ctx, utils.AssignmentTable, selector, update, opts)
	if err != nil {
		return err
	}

	// 更新统计信息
	return d.UpdateAssignmentStats(ctx, assignmentID)
}

// GradeAssignment 教师批改作业
func (d *Dao) GradeAssignment(ctx context.Context, assignmentID string, studentID string, graderID string, graderName string, score int, comment string) error {
	objID, err := primitive.ObjectIDFromHex(assignmentID)
	if err != nil {
		return err
	}

	// 更新学生的提交信息
	update := bson.M{
		"$set": bson.M{
			"members.4.$[elem].submission.status":      2, // 已批改
			"members.4.$[elem].submission.score":       score,
			"members.4.$[elem].submission.comment":     comment,
			"members.4.$[elem].submission.grader_id":   graderID,
			"members.4.$[elem].submission.grader_name": graderName,
			"members.4.$[elem].submission.grade_time":  time.Now().Unix(),
			"mtime": time.Now().Unix(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.user_id": studentID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	selector := bson.M{"_id": objID}
	_, err = d.mongo.UpdateOne(ctx, utils.AssignmentTable, selector, update, opts)
	if err != nil {
		return err
	}

	// 更新统计信息
	return d.UpdateAssignmentStats(ctx, assignmentID)
}

// RejectAssignment 教师打回作业
func (d *Dao) RejectAssignment(ctx context.Context, assignmentID string, studentID string, graderID string, graderName string, comment string) error {
	objID, err := primitive.ObjectIDFromHex(assignmentID)
	if err != nil {
		return err
	}

	// 更新学生的提交信息
	update := bson.M{
		"$set": bson.M{
			"members.4.$[elem].submission.status":      3, // 被打回
			"members.4.$[elem].submission.comment":     comment,
			"members.4.$[elem].submission.grader_id":   graderID,
			"members.4.$[elem].submission.grader_name": graderName,
			"members.4.$[elem].submission.grade_time":  time.Now().Unix(),
			"mtime": time.Now().Unix(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.user_id": studentID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	selector := bson.M{"_id": objID}
	_, err = d.mongo.UpdateOne(ctx, utils.AssignmentTable, selector, update, opts)
	if err != nil {
		return err
	}

	// 更新统计信息
	return d.UpdateAssignmentStats(ctx, assignmentID)
}

// UpdateAssignmentStats 更新作业统计信息
func (d *Dao) UpdateAssignmentStats(ctx context.Context, assignmentID string) error {
	// 获取作业信息
	assignment, err := d.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("作业不存在")
	}

	// 计算统计信息
	totalStudents := 0
	submittedCount := 0
	gradedCount := 0
	totalScore := 0

	if students, ok := assignment.Members["4"]; ok {
		totalStudents = len(students)
		for _, student := range students {
			if student.Submission != nil {
				if student.Submission.Status > 0 { // 已提交
					submittedCount++
				}
				if student.Submission.Status == 2 { // 已批改
					gradedCount++
					totalScore += student.Submission.Score
				}
			}
		}
	}

	// 计算平均分
	averageScore := 0
	if gradedCount > 0 {
		averageScore = totalScore / gradedCount
	}

	// 更新统计信息
	stats := models.AssignmentStats{
		TotalStudents:  totalStudents,
		SubmittedCount: submittedCount,
		GradedCount:    gradedCount,
		AverageScore:   averageScore,
	}

	update := bson.M{
		"$set": bson.M{
			"stats": stats,
			"mtime": time.Now().Unix(),
		},
	}

	objID, _ := primitive.ObjectIDFromHex(assignmentID)
	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.AssignmentTable, selector, update)
	return err
}

// AddStudentsToAssignment 批量添加学生到作业
func (d *Dao) AddStudentsToAssignment(ctx context.Context, assignmentID string, students []models.Member) error {
	if len(students) == 0 {
		return nil
	}

	objID, err := primitive.ObjectIDFromHex(assignmentID)
	if err != nil {
		return err
	}

	// 添加学生到作业
	update := bson.M{
		"$push": bson.M{
			"members.4": bson.M{
				"$each": students,
			},
		},
		"$set": bson.M{
			"mtime": time.Now().Unix(),
		},
	}

	selector := bson.M{"_id": objID}
	_, err = d.Update(ctx, utils.AssignmentTable, selector, update)
	if err != nil {
		return err
	}

	// 更新统计信息
	return d.UpdateAssignmentStats(ctx, assignmentID)
}

// GetStudentSubmission 获取学生的作业提交信息
func (d *Dao) GetStudentSubmission(ctx context.Context, assignmentID string, studentID string) (*models.Submission, error) {
	objID, err := primitive.ObjectIDFromHex(assignmentID)
	if err != nil {
		return nil, err
	}

	// 构建聚合管道
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": objID},
		},
		{
			"$project": bson.M{
				"student": bson.M{
					"$filter": bson.M{
						"input": "$members.4",
						"as":    "student",
						"cond":  bson.M{"$eq": []interface{}{"$$student.user_id", studentID}},
					},
				},
			},
		},
	}

	// 执行聚合查询
	var result []bson.M
	err = d.mongo.Aggregate(ctx, utils.AssignmentTable, pipeline, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 || len(result[0]["student"].(primitive.A)) == 0 {
		return nil, nil
	}

	// 提取学生提交信息
	studentDoc := result[0]["student"].(primitive.A)[0].(bson.M)
	if submission, ok := studentDoc["submission"]; ok && submission != nil {
		// 将bson.M转换为models.Submission
		submissionDoc := submission.(bson.M)
		sub := &models.Submission{}

		if status, ok := submissionDoc["status"]; ok {
			sub.Status = int(status.(int32))
		}
		if textContent, ok := submissionDoc["text_content"]; ok {
			sub.TextContent = textContent.(string)
		}
		if fileURL, ok := submissionDoc["file_url"]; ok {
			sub.FileURL = fileURL.(string)
		}
		if score, ok := submissionDoc["score"]; ok {
			sub.Score = int(score.(int32))
		}
		if comment, ok := submissionDoc["comment"]; ok {
			sub.Comment = comment.(string)
		}
		if graderID, ok := submissionDoc["grader_id"]; ok {
			sub.GraderID = graderID.(string)
		}
		if graderName, ok := submissionDoc["grader_name"]; ok {
			sub.GraderName = graderName.(string)
		}
		if submitTime, ok := submissionDoc["submit_time"]; ok {
			sub.SubmitTime = submitTime.(int64)
		}
		if gradeTime, ok := submissionDoc["grade_time"]; ok {
			sub.GradeTime = gradeTime.(int64)
		}

		return sub, nil
	}

	return nil, nil
}

// ExportAssignmentGrades 导出作业成绩
func (d *Dao) ExportAssignmentGrades(ctx context.Context, assignmentID string) ([]map[string]interface{}, error) {
	// 获取作业信息
	assignment, err := d.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, errors.New("作业不存在")
	}

	// 准备导出数据
	var grades []map[string]interface{}

	if students, ok := assignment.Members["4"]; ok {
		for _, student := range students {
			grade := map[string]interface{}{
				"assignment_code":  assignment.AssignmentCode,
				"assignment_title": assignment.Title,
				"course_name":      assignment.CourseName,
				"student_id":       student.UserID,
				"student_name":     student.UserName,
				"status":           "未提交",
				"score":            0,
				"comment":          "",
				"submit_time":      "",
				"grade_time":       "",
			}

			if student.Submission != nil {
				switch student.Submission.Status {
				case 1:
					grade["status"] = "已提交未批改"
				case 2:
					grade["status"] = "已批改"
					grade["score"] = student.Submission.Score
				case 3:
					grade["status"] = "被打回"
				}

				grade["comment"] = student.Submission.Comment

				if student.Submission.SubmitTime > 0 {
					grade["submit_time"] = time.Unix(student.Submission.SubmitTime, 0).Format("2006-01-02 15:04:05")
				}

				if student.Submission.GradeTime > 0 {
					grade["grade_time"] = time.Unix(student.Submission.GradeTime, 0).Format("2006-01-02 15:04:05")
				}
			}

			grades = append(grades, grade)
		}
	}

	return grades, nil
}

// GetStudentAssignments 获取学生的作业列表
func (d *Dao) GetStudentAssignments(ctx context.Context, studentID string, comQuery *utils.CommonQuery, status int) (*utils.RespPageQuery, error) {
	log.Printf("开始执行GetStudentAssignments，studentID=%s, status=%d", studentID, status)
	log.Printf("comQuery.Filters=%v", comQuery.Filters)

	// 使用聚合管道查询学生的作业
	pipeline := []bson.M{
		{
			"$match": comQuery.Filters,
		},
		{
			"$addFields": bson.M{
				"student": bson.M{
					"$filter": bson.M{
						"input": "$members.4",
						"as":    "student",
						"cond":  bson.M{"$eq": []interface{}{"$$student.user_id", studentID}},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"student": bson.M{"$ne": []interface{}{}},
			},
		},
	}

	// 如果有提交状态筛选
	if status >= 0 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"student.0.submission.status": status,
			},
		})
	}

	// 计算总数
	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)
	countPipeline = append(countPipeline, bson.M{"$count": "total"})
	var countResult []bson.M
	err := d.mongo.Aggregate(ctx, utils.AssignmentTable, countPipeline, &countResult)
	if err != nil {
		return nil, err
	}

	total := int64(0)
	if len(countResult) > 0 {
		total = countResult[0]["total"].(int64)
	}

	// 添加分页
	pageNum := 1
	pageSize := 10

	// 从 comQuery 中获取分页参数
	if v, ok := comQuery.Filters["page"].(int); ok {
		pageNum = v
	}
	if v, ok := comQuery.Filters["page_size"].(int); ok {
		pageSize = v
	}

	pipeline = append(pipeline,
		bson.M{"$skip": (pageNum - 1) * pageSize},
		bson.M{"$limit": pageSize},
	)

	// 在执行聚合查询前添加日志
	log.Printf("准备执行聚合查询，pipeline长度=%d", len(pipeline))

	// 使用d.mongo.Aggregate方法执行查询
	var results []bson.M
	err = d.mongo.Aggregate(ctx, utils.AssignmentTable, pipeline, &results)
	if err != nil {
		log.Printf("聚合查询失败: %v", err)
		return nil, err
	}

	log.Printf("聚合查询成功，results长度=%d", len(results))

	// 构建响应
	items := make([]*map[string]interface{}, len(results))
	for i, result := range results {
		item := make(map[string]interface{})
		for k, v := range result {
			if k != "_id" {
				item[k] = v
			} else {
				item["id"] = v.(primitive.ObjectID).Hex()
			}
		}
		itemPtr := &item
		items[i] = itemPtr
	}

	return &utils.RespPageQuery{
		Items: items,
		Total: total,
	}, nil
}
