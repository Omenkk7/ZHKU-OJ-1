package service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// CreateContest 创建竞赛
func (s *Service) CreateContest(ctx context.Context, req *dto.CreateContestReq, userID, userName string) (string, error) {
	lg := utils.GetDefaultLogger()
	lg.Infof("创建竞赛请求: %+v", req)

	// 验证访问类型
	if req.AccessType != utils.ContestAccessPublic && req.AccessType != utils.ContestAccessPrivate {
		lg.Errorf("无效的竞赛访问类型: %d", req.AccessType)
		return "", fmt.Errorf("无效的竞赛访问类型: %d，只能是1(公开)或2(私有)", req.AccessType)
	}

	// 验证时间
	if req.StartTime >= req.EndTime {
		lg.Errorf("时间验证失败: 开始时间 %d >= 结束时间 %d", req.StartTime, req.EndTime)
		return "", errors.New("开始时间必须早于结束时间")
	}

	// 如果没有提供竞赛代码，则自动生成
	if req.ContestCode == "" {
		// 生成竞赛代码
		req.ContestCode = utils.GenerateContestCode(req.ContestType)

		// 检查竞赛代码是否已存在
		existingContest, err := s.dao.GetContestByCode(ctx, req.ContestCode)
		if err != nil && err != mongo.ErrNoDocuments {
			lg.Errorf("检查竞赛代码失败: %v", err)
			return "", err
		}

		// 如果竞赛代码已存在，则重新生成
		for existingContest != nil {
			req.ContestCode = utils.GenerateContestCode(req.ContestType)
			existingContest, err = s.dao.GetContestByCode(ctx, req.ContestCode)
			if err != nil && err != mongo.ErrNoDocuments {
				lg.Errorf("检查竞赛代码失败: %v", err)
				return "", err
			}
		}
	} else {
		// 如果提供了竞赛代码，检查是否已存在
		existingContest, err := s.dao.GetContestByCode(ctx, req.ContestCode)
		if err != nil && err != mongo.ErrNoDocuments {
			lg.Errorf("检查竞赛代码失败: %v", err)
			return "", err
		}
		if existingContest != nil {
			lg.Errorf("竞赛代码已存在: %s", req.ContestCode)
			return "", errors.New("竞赛代码已存在")
		}
	}

	// 验证题目
	if len(req.Problems) > 0 {
		for i, p := range req.Problems {
			// 验证题目ID是否有效
			_, err := primitive.ObjectIDFromHex(p.ProblemID)
			if err != nil {
				lg.Errorf("无效的题目ID: %s, 错误: %v", p.ProblemID, err)
				return "", errors.New("无效的题目ID: " + p.ProblemID)
			}
			// 设置默认值
			if p.Order == 0 {
				req.Problems[i].Order = i + 1
			}
			if p.Score == 0 {
				req.Problems[i].Score = 10 // 默认分值
			}
			req.Problems[i].Status = 1 // 默认可用
		}
	}

	// 转换为Contest对象
	contest := req.ToContest(userID, userName)

	// 记录最终的竞赛对象
	lg.Infof("最终创建的竞赛对象: %+v", contest)

	// 创建竞赛
	contestID, err := s.dao.CreateContest(ctx, contest)
	if err != nil {
		lg.Errorf("创建竞赛失败: %v", err)
		return "", err
	}

	// 初始化竞赛排行榜（空的）
	err = s.CreateContestRanking(ctx, contestID)
	if err != nil {
		// 记录错误但不影响竞赛创建
		lg.Errorf("初始化竞赛排行榜失败: %v", err)
	}

	return contestID, nil
}

// CreateContestRanking 初始化竞赛排行榜
func (s *Service) CreateContestRanking(ctx context.Context, contestID string) error {
	// 获取竞赛信息以获取竞赛名称
	contest, err := s.GetContestByID(ctx, contestID)
	if err != nil {
		return err
	}

	ranking := &models.ContestRanking{
		ID:            primitive.NewObjectID(),
		ContestID:     contestID,
		ContestName:   contest.Name,
		Rankings:      make([]models.RankingItem, 0),
		GeneratedTime: time.Now().Unix(),
	}

	return s.dao.UpdateContestRanking(ctx, ranking)
}

// GetContestByID 根据ID获取竞赛
func (s *Service) GetContestByID(ctx context.Context, id string) (*models.Contest, error) {
	// 验证ID
	contestID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("无效的竞赛ID")
	}

	// 获取竞赛
	contest, err := s.dao.GetContestByID(ctx, contestID.Hex())
	if err != nil {
		return nil, err
	}

	return contest, nil
}

// GetContestList 获取竞赛列表
func (s *Service) GetContestList(ctx context.Context, req *dto.GetContestListReq, userID string, userRole int) (*dto.ContestListResp, error) {
	query := bson.M{}

	// 默认过滤已删除的竞赛
	if req.Status == 0 {
		query["status"] = bson.M{"$ne": utils.ContestStatusDeleted}
	}

	// 其他过滤条件
	if req.Name != "" {
		query["name"] = bson.M{"$regex": req.Name, "$options": "i"}
	}
	if req.ContestType > 0 {
		query["contest_type"] = req.ContestType
	}
	if req.Status > 0 {
		query["status"] = req.Status
	}
	if req.AccessType > 0 {
		query["access_type"] = req.AccessType
	}

	// 根据用户角色添加不同的查询条件
	if userRole == utils.RoleTeacher {
		// 教师只能查看自己创建的或参与的竞赛
		// 使用 $or 查询：创建者是自己 或 成员中包含自己
		query["$or"] = []bson.M{
			{"creator_id": userID},
			{"members." + userID: bson.M{"$exists": true}}, // 教师角色为2
		}
	} else if userRole == utils.RoleAssistant {
		// 助教只能查看自己参与的竞赛
		query["members."+userID] = bson.M{"$exists": true} // 助教角色为3
	} else if userRole == utils.RoleStudent {
		// 学生只能查看自己参与的竞赛和公开竞赛
		query["$or"] = []bson.M{
			{"access_type": 1}, // 公开竞赛
			{"members." + userID: bson.M{"$exists": true}}, // 学生角色为4
		}
	}
	// 管理员和超级管理员可以查看所有竞赛，不需要额外条件

	// 如果明确指定了创建者ID，则覆盖上面的条件
	if req.CreatorID != "" {
		query["creator_id"] = req.CreatorID
	}

	lg := utils.GetDefaultLogger()
	lg.Infof("竞赛列表查询条件: %+v", query)

	// 查询列表
	contests, total, err := s.dao.GetContestList(ctx, query, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	lg.Infof("竞赛列表查询结果: total=%d, contests=%+v", total, contests)

	// 构建响应
	resp := &dto.ContestListResp{
		Total: total,
		List:  make([]models.Contest, 0, len(contests)),
	}
	for _, contest := range contests {
		resp.List = append(resp.List, *contest)
	}

	return resp, nil
}

// UpdateContest 更新竞赛
func (s *Service) UpdateContest(ctx context.Context, id string, req *dto.UpdateContestReq, userID string) error {
	lg := utils.GetDefaultLogger()

	// 验证ID
	contestID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		lg.Errorf("无效的竞赛ID: %s, 错误: %v", id, err)
		return errors.New("无效的竞赛ID")
	}

	// 获取竞赛
	contest, err := s.dao.GetContestByID(ctx, contestID.Hex())
	if err != nil {
		lg.Errorf("获取竞赛失败: %v", err)
		return err
	}
	if contest == nil {
		lg.Error("竞赛不存在")
		return errors.New("竞赛不存在")
	}

	// 记录当前时间和竞赛开始时间
	now := time.Now().Unix()
	lg.Infof("当前时间: %v (%s), 竞赛开始时间: %v (%s)",
		now, time.Unix(now, 0).Format("2006-01-02 15:04:05"),
		contest.StartTime, time.Unix(contest.StartTime, 0).Format("2006-01-02 15:04:05"))

	// 记录请求中的开始时间和结束时间
	if req.StartTime > 0 {
		startTime := req.StartTime / 1000 // 转换为秒级时间戳
		lg.Infof("请求中的开始时间: %v (%s)",
			startTime, time.Unix(startTime, 0).Format("2006-01-02 15:04:05"))
	}
	if req.EndTime > 0 {
		endTime := req.EndTime / 1000 // 转换为秒级时间戳
		lg.Infof("请求中的结束时间: %v (%s)",
			endTime, time.Unix(endTime, 0).Format("2006-01-02 15:04:05"))
	}

	// 验证权限
	hasPermission, err := s.CheckContestPermission(ctx, id, userID, []int{utils.ContestRoleAdmin, utils.RoleTeacher})
	if err != nil {
		lg.Errorf("检查权限失败: %v", err)
		return err
	}
	if !hasPermission {
		lg.Error("无权限更新此竞赛")
		return errors.New("无权限更新此竞赛")
	}

	// 检查竞赛是否已开始
	// 修改：只有当竞赛状态为"进行中"或"已结束"时才认为已开始
	// 不再直接比较时间戳
	if contest.Status == utils.ContestStatusRunning || contest.Status == utils.ContestStatusEnded {
		lg.Errorf("不能修改已经开始的竞赛，当前状态: %d", contest.Status)
		return errors.New("不能修改已经开始的竞赛")
	}

	// 构建更新内容
	update := bson.M{"mtime": time.Now().Unix()}

	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.StartTime > 0 {
		update["start_time"] = req.StartTime / 1000 // 转换为秒级时间戳
	}
	if req.EndTime > 0 {
		update["end_time"] = req.EndTime / 1000 // 转换为秒级时间戳
	}
	if req.ContestType > 0 {
		update["contest_type"] = req.ContestType
	}
	if req.AccessType > 0 {
		update["access_type"] = req.AccessType
	}
	if req.MaxParticipants > 0 {
		update["max_participants"] = req.MaxParticipants
	}
	if req.Problems != nil && len(req.Problems) > 0 {
		update["problems"] = req.Problems
	}
	if req.Status > 0 {
		update["status"] = req.Status
	}

	lg.Infof("更新内容: %+v", update)

	// 更新竞赛
	return s.dao.UpdateContest(ctx, contestID.Hex(), update)
}

// DeleteContest 删除竞赛
func (s *Service) DeleteContest(ctx context.Context, id string, userID string) error {
	// 验证ID
	contestID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 获取竞赛
	contest, err := s.dao.GetContestByID(ctx, contestID.Hex())
	if err != nil {
		return err
	}

	// 验证权限
	hasPermission, err := s.CheckContestPermission(ctx, id, userID, []int{utils.ContestRoleAdmin, utils.RoleTeacher})
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("无权限删除此竞赛")
	}

	// 验证状态
	if contest.Status != utils.ContestStatusNotStart {
		return errors.New("只能删除未开始的竞赛")
	}

	// 删除竞赛
	return s.dao.DeleteContest(ctx, contestID.Hex())
}

// ArchiveContest 归档竞赛
func (s *Service) ArchiveContest(ctx context.Context, id string, userID string) error {
	// 验证ID
	contestID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 获取竞赛
	contest, err := s.dao.GetContestByID(ctx, contestID.Hex())
	if err != nil {
		return err
	}

	// 验证权限
	hasPermission, err := s.CheckContestPermission(ctx, id, userID, []int{utils.ContestRoleAdmin, utils.RoleTeacher})
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("无权限归档此竞赛")
	}

	// 如果竞赛未结束，先将其设置为已结束状态
	if contest.Status != utils.ContestStatusEnded && contest.Status != utils.ContestStatusArchived {
		// 更新竞赛状态为已结束
		err = s.dao.UpdateContestStatus(ctx, contestID.Hex(), utils.ContestStatusEnded)
		if err != nil {
			return errors.New("将竞赛设置为已结束状态失败: " + err.Error())
		}

		// 记录日志
		log.Printf("竞赛 %s 已自动设置为已结束状态", contestID.Hex())
	}

	// 如果竞赛已经是归档状态，直接返回成功
	if contest.Status == utils.ContestStatusArchived {
		return nil
	}

	// 更新状态为归档
	return s.dao.UpdateContestStatus(ctx, contestID.Hex(), utils.ContestStatusArchived)
}

// UpdateContestStatus 更新竞赛状态
func (s *Service) UpdateContestStatus(ctx context.Context, id string, status int, userID string) error {
	// 验证ID
	contestID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 验证权限
	hasPermission, err := s.CheckContestPermission(ctx, id, userID, []int{utils.ContestRoleAdmin})
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("无权限更新此竞赛状态")
	}

	// 更新状态
	return s.dao.UpdateContestStatus(ctx, contestID.Hex(), status)
}

// AddParticipant 添加参赛者
func (s *Service) AddParticipant(ctx context.Context, req *dto.AddParticipantReq, userID string) error {
	// 验证竞赛ID
	contestID, err := primitive.ObjectIDFromHex(req.ContestID)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 获取竞赛
	contest, err := s.dao.GetContestByID(ctx, contestID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("竞赛不存在或已被删除")
		}
		return err
	}

	if contest.Status != utils.ContestStatusNotStart {
		return errors.New("只能在未开始的竞赛添加参赛者")
	}

	// 验证参赛人数
	if contest.MaxParticipants > 0 && contest.Stats.ParticipantCount >= contest.MaxParticipants {
		return errors.New("参赛人数已达上限")
	}

	// 验证权限
	if contest.AccessType == utils.ContestAccessPrivate {
		// 私有竞赛需要管理员或教师权限
		hasPermission, err := s.CheckContestPermission(ctx, req.ContestID, userID, []int{utils.ContestRoleAdmin, utils.RoleTeacher})
		if err != nil {
			return err
		}
		if !hasPermission {
			return errors.New("无权限添加参赛者")
		}
	}

	// 创建参赛者
	participant := &models.ContestParticipant{
		ID:           primitive.NewObjectID(),
		ContestID:    contestID.Hex(),
		ContestName:  contest.Name,
		StudentID:    req.StudentID,
		StudentName:  req.StudentName,
		RegisterTime: time.Now().Unix(),
		Status:       1, // 默认通过
		Ctime:        time.Now().Unix(),
		Mtime:        time.Now().Unix(),
	}

	// 添加参赛者
	_, err = s.dao.AddParticipant(ctx, participant)
	if err != nil {
		return err
	}

	studentObjID, err := primitive.ObjectIDFromHex(req.StudentID)
	if err != nil {
		return errors.New("无效的学生ID")
	}
	user := models.User{
		ID:       studentObjID,
		Username: req.StudentName,
	}
	return s.dao.AddMember(ctx, contestID.Hex(), strconv.Itoa(utils.RoleStudent), user)
}

// BatchAddParticipants 批量添加参赛者
func (s *Service) BatchAddParticipants(ctx context.Context, req *dto.BatchAddParticipantsReq, userID string) error {
	// 验证竞赛ID
	contestID, err := primitive.ObjectIDFromHex(req.ContestID)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 获取竞赛
	contest, err := s.dao.GetContestByID(ctx, contestID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("竞赛不存在或已被删除") // 修改这里的错误提示
		}
		return err
	}

	if contest.Status != utils.ContestStatusNotStart {
		return errors.New("只能在未开始的竞赛添加参赛者")
	}

	// 验证参赛人数
	if contest.MaxParticipants > 0 && contest.Stats.ParticipantCount+len(req.Students) > contest.MaxParticipants {
		return errors.New("参赛人数将超过上限")
	}

	// 验证权限
	hasPermission, err := s.CheckContestPermission(ctx, req.ContestID, userID, []int{utils.ContestRoleAdmin, utils.RoleTeacher})
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("无权限批量添加参赛者")
	}

	// 批量添加参赛者
	successCount := 0
	failCount := 0
	for _, student := range req.Students {
		// 验证学生ID
		if student.ID.IsZero() {
			return errors.New("学生ID不能为空")
		}

		participant := &models.ContestParticipant{
			ContestID:    req.ContestID,
			ContestName:  contest.Name,
			StudentID:    student.ID.Hex(),
			StudentName:  student.Username,
			RegisterTime: time.Now().Unix(),
			Status:       1, // 设置为待审核状态
			Ctime:        time.Now().Unix(),
			Mtime:        time.Now().Unix(),
		}

		// 添加参赛者
		_, err = s.dao.AddParticipant(ctx, participant)
		if err != nil {
			failCount++
			log.Printf("添加参赛者失败: contestID=%s, studentID=%s, error=%v",
				req.ContestID, student.ID.Hex(), err)
			continue // 忽略错误，继续添加
		}

		// 添加到竞赛成员
		err = s.dao.AddMember(ctx, contestID.Hex(), strconv.Itoa(utils.RoleStudent), student)
		if err != nil {
			log.Printf("添加竞赛成员失败: contestID=%s, studentID=%s, error=%v",
				req.ContestID, student.ID.Hex(), err)
			continue
		}

		successCount++
	}

	log.Printf("批量添加参赛者完成: 成功=%d, 失败=%d", successCount, failCount)

	// 如果全部失败但没有返回错误，这里可以选择返回一个错误
	if failCount > 0 && successCount == 0 {
		return errors.New(fmt.Sprintf("所有参赛者(%d人)添加失败，请检查日志", failCount))
	}

	return nil
}

// RemoveParticipant 移除参赛者
func (s *Service) RemoveParticipant(ctx context.Context, req *dto.RemoveParticipantReq, userID string) error {
	// 验证竞赛ID
	contestID, err := primitive.ObjectIDFromHex(req.ContestID)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 获取竞赛
	contest, err := s.dao.GetContestByID(ctx, contestID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("竞赛不存在或已被删除")
		}
		return err
	}

	// 验证竞赛状态
	if contest.Status != utils.ContestStatusNotStart {
		return errors.New("只能在未开始的竞赛移除参赛者")
	}

	// 验证权限
	hasPermission, err := s.CheckContestPermission(ctx, req.ContestID, userID, []int{utils.ContestRoleAdmin, utils.RoleTeacher})
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("无权限移除参赛者")
	}

	// 验证参赛者是否存在
	_, err = s.dao.GetParticipantByID(ctx, contestID.Hex(), req.StudentID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("参赛者不存在或已被删除")
		}
		return err
	}

	// 移除参赛者
	return s.dao.RemoveParticipant(ctx, contestID.Hex(), req.StudentID)
}

// GetParticipantList 获取参赛者列表
func (s *Service) GetParticipantList(ctx context.Context, req *dto.GetParticipantListReq) (*dto.ParticipantListResp, error) {
	// 验证竞赛ID
	contestID, err := primitive.ObjectIDFromHex(req.ContestID)
	if err != nil {
		return nil, errors.New("无效的竞赛ID")
	}

	// 构建查询条件
	query := bson.M{"contest_id": contestID.Hex()}
	if req.Status > 0 {
		query["status"] = req.Status
	}

	// 查询列表
	participants, total, err := s.dao.GetParticipantList(ctx, query, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 构建响应
	resp := &dto.ParticipantListResp{
		Total: total,
		List:  make([]models.ContestParticipant, 0, len(participants)),
	}
	for _, participant := range participants {
		resp.List = append(resp.List, *participant)
	}

	return resp, nil
}

// AuditParticipant 审核参赛者
func (s *Service) AuditParticipant(ctx context.Context, req *dto.AuditParticipantReq, userID string) error {
	// 验证竞赛ID
	contestID, err := primitive.ObjectIDFromHex(req.ContestID)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 验证权限
	hasPermission, err := s.CheckContestPermission(ctx, req.ContestID, userID, []int{utils.ContestRoleAdmin, utils.RoleTeacher})
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("无权限审核参赛者")
	}

	// 更新参赛者状态
	return s.dao.UpdateParticipantStatus(ctx, contestID.Hex(), req.StudentID, req.Status)
}

// GetContestRanking 获取竞赛排名
func (s *Service) GetContestRanking(ctx context.Context, req *dto.GetContestRankingReq) (*dto.ContestRankingResp, error) {
	// 验证竞赛ID
	contestID, err := primitive.ObjectIDFromHex(req.ContestID)
	if err != nil {
		return nil, errors.New("无效的竞赛ID")
	}

	// 查询排名
	rankings, total, err := s.dao.GetContestRanking(ctx, contestID.Hex(), req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 构建响应
	resp := &dto.ContestRankingResp{
		Total: total,
		List:  make([]models.ContestRanking, 0, len(rankings)),
	}
	for _, ranking := range rankings {
		resp.List = append(resp.List, *ranking)
	}

	return resp, nil
}

// ExportContestScore 导出竞赛成绩
func (s *Service) ExportContestScore(ctx context.Context, req *dto.ExportContestScoreReq) ([]byte, error) {
	// 验证竞赛ID
	contestID, err := primitive.ObjectIDFromHex(req.ContestID)
	if err != nil {
		return nil, errors.New("无效的竞赛ID")
	}

	// 检查竞赛是否存在
	_, err = s.dao.GetContestByID(ctx, contestID.Hex())
	if err != nil {
		return nil, err
	}

	// 获取所有排名
	rankings, _, err := s.dao.GetContestRanking(ctx, contestID.Hex(), 1, 1000) // 假设最多1000人
	if err != nil {
		return nil, err
	}

	// 生成CSV数据
	csvData := []byte("学生ID,学生姓名,解决题目数,总分数,总用时(秒),排名\n")

	// 遍历每个竞赛排名记录
	for _, ranking := range rankings {
		// 遍历每个排名项
		for _, item := range ranking.Rankings {
			line := []byte(item.StudentID + "," + item.StudentName + "," +
				strconv.Itoa(item.SolvedCount) + "," +
				strconv.Itoa(item.TotalScore) + "," +
				strconv.Itoa(item.TotalTime) + "," + // 注意：TotalTime现在是int而不是int64
				strconv.Itoa(item.Rank) + "\n")
			csvData = append(csvData, line...)
		}
	}

	return csvData, nil
}

// ApplyJoinContest 学生申请加入竞赛
func (s *Service) ApplyJoinContest(ctx context.Context, req *dto.ApplyJoinContestReq, studentID, studentName string) error {
	lg := utils.GetDefaultLogger()

	// 验证竞赛ID
	contestID, err := primitive.ObjectIDFromHex(req.ContestID)
	if err != nil {
		lg.Errorf("无效的竞赛ID: %s, 错误: %v", req.ContestID, err)
		return errors.New("无效的竞赛ID")
	}

	// 获取竞赛信息
	contest, err := s.dao.GetContestByID(ctx, contestID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("竞赛不存在或已被删除")
		}
		lg.Errorf("获取竞赛失败: %v", err)
		return err
	}

	// 检查竞赛是否为私有竞赛
	if contest.AccessType != utils.ContestAccessPrivate {
		return errors.New("只能申请加入私有竞赛")
	}

	// 检查竞赛是否已开始
	if contest.Status == utils.ContestStatusRunning || contest.Status == utils.ContestStatusEnded {
		return errors.New("竞赛已开始或已结束，无法申请加入")
	}

	// 检查是否已经是参赛者
	isParticipant, err := s.dao.IsContestParticipant(ctx, req.ContestID, studentID)
	if err != nil {
		lg.Errorf("检查参赛者失败: %v", err)
		return err
	}
	if isParticipant {
		return errors.New("您已经是该竞赛的参赛者")
	}

	// 检查是否已经申请过
	hasApplied, err := s.dao.HasAppliedContest(ctx, req.ContestID, studentID)
	if err != nil {
		lg.Errorf("检查申请记录失败: %v", err)
		return err
	}
	if hasApplied {
		return errors.New("您已经申请过该竞赛，请等待审核")
	}

	// 创建申请记录
	participant := &models.ContestParticipant{
		ID:          primitive.NewObjectID(),
		ContestID:   req.ContestID,
		ContestName: contest.Name,
		StudentID:   studentID,
		StudentName: studentName,
		Status:      utils.ContestParticipantStatusPending, // 待审核状态
		Ctime:       time.Now().Unix(),
		Mtime:       time.Now().Unix(),
	}

	// 保存申请记录
	_, err = s.dao.CreateContestParticipant(ctx, participant)
	if err != nil {
		lg.Errorf("创建申请记录失败: %v", err)
		return err
	}

	return nil
}

// 检查用户是否有权限
func (s *Service) hasPermission(contest *models.Contest, userID string, allowedRoles []string) bool {
	// 创建者有权限
	if contest.CreatorID == userID {
		return true
	}

	// 检查用户角色
	for _, role := range allowedRoles {
		members, ok := contest.Members[role]
		if !ok {
			continue
		}
		for _, member := range members {
			// 使用正确的字段名
			if member.ID.Hex() == userID {
				return true
			}
		}
	}

	return false
}

// CheckContestPermission 检查用户是否有权限操作竞赛
func (s *Service) CheckContestPermission(ctx context.Context, contestID string, userID string, allowedRoles []int) (bool, error) {
	// 验证ID
	objID, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		return false, errors.New("无效的竞赛ID")
	}

	// 获取竞赛信息
	contest, err := s.dao.GetContestByID(ctx, objID.Hex())
	if err != nil {
		return false, err
	}
	if contest == nil {
		return false, errors.New("竞赛不存在")
	}

	// 检查用户是否有权限操作竞赛
	hasPermission := false

	// 检查是否是系统管理员
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err == nil {
		selector := bson.M{
			"_id": userObjID,
		}
		user, err := s.dao.GetOneUser(ctx, selector)
		if err == nil && user != nil && user.Role == 1 { // 系统管理员
			hasPermission = true
			return hasPermission, nil
		}
	}

	// 检查是否是创建者
	if contest.CreatorID == userID {
		hasPermission = true
		return hasPermission, nil
	}

	// 检查用户角色
	for _, role := range allowedRoles {
		// 将int角色转为string
		roleStr := strconv.Itoa(role)
		members, ok := contest.Members[roleStr]
		if !ok {
			continue
		}
		for _, member := range members {
			if member.ID.Hex() == userID {
				hasPermission = true
				return hasPermission, nil
			}
		}
	}

	return hasPermission, nil
}

// AddProblemsToContest 向竞赛添加题目
func (s *Service) AddProblemsToContest(ctx context.Context, contestID string, req *dto.AddProblemsToContestReq, userID string, userRole int) error {
	lg := utils.GetDefaultLogger()
	lg.Infof("用户 %s (角色 %d) 尝试向竞赛 %s 添加题目: %+v", userID, userRole, contestID, req.Problems)

	//  获取竞赛信息
	contest, err := s.dao.GetContestByID(ctx, contestID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			lg.Warnf("竞赛 %s 不存在", contestID)
			return errors.New("竞赛不存在")
		}
		lg.Errorf("获取竞赛 %s 失败: %v", contestID, err)
		return fmt.Errorf("获取竞赛信息失败: %w", err)
	}

	//  权限检查 (使用 Casbin)
	roleName := utils.GetRoleName(userRole)
	canAddProblem := utils.CheckPermission(roleName, utils.ObjContest, utils.ActAddProblem)

	if !canAddProblem {
		lg.Warnf("用户 %s (角色 %s) Casbin 权限检查失败 (Obj: %s, Act: %s)", userID, roleName, utils.ObjContest, utils.ActAddProblem)
		return utils.ErrNoPermission
	}

	// 查看竞赛状态，不许向已结束或已归档的竞赛添加题目
	if contest.Status == utils.ContestStatusEnded || contest.Status == utils.ContestStatusArchived {
		lg.Warnf("无法向已结束或已归档的竞赛 %s 添加题目", contestID)
		return errors.New("无法向已结束或已归档的竞赛添加题目")
	}

	//  数据验证和处理
	if len(req.Problems) == 0 {
		return errors.New("题目列表不能为空")
	}

	// 获取现有题目ID集合，用于检查重复
	existingProblemIDs := make(map[string]bool)
	maxOrder := 0
	for _, p := range contest.Problems {
		existingProblemIDs[p.ProblemID] = true
		if p.Order > maxOrder {
			maxOrder = p.Order
		}
	}

	problemsToAdd := make([]models.ContestProblem, 0, len(req.Problems))
	for i, p := range req.Problems {
		// 验证题目ID和题目名称是否为空
		if p.ProblemID == "" {
			lg.Errorf("题目ID不能为空")
			return errors.New("题目ID不能为空")
		}
		if p.Title == "" {
			lg.Errorf("题目名称不能为空")
			return errors.New("题目名称不能为空")
		}

		// 验证题目ID格式
		if _, err := primitive.ObjectIDFromHex(p.ProblemID); err != nil {
			lg.Errorf("无效的题目ID格式: %s", p.ProblemID)
			return fmt.Errorf("无效的题目ID格式: %s", p.ProblemID)
		}

		// 检查题目是否已存在于竞赛中
		if existingProblemIDs[p.ProblemID] {
			lg.Warnf("题目 %s 已存在于竞赛 %s 中，跳过添加", p.ProblemID, contestID)
			continue // 跳过已存在的题目
		}

		// 检查待添加列表中是否有重复
		for j := 0; j < i; j++ {
			if req.Problems[j].ProblemID == p.ProblemID {
				lg.Warnf("待添加列表中存在重复题目ID: %s", p.ProblemID)
				return fmt.Errorf("待添加列表中存在重复题目ID: %s", p.ProblemID)
			}
		}

		// TODO: 验证 ProblemID 是否在 problems 集合中真实存在
		//problemExists, err := s.dao.CheckProblemExists(ctx, p.ProblemID)
		//if err != nil { ... }
		//if !problemExists { return fmt.Errorf("题目 %s 不存在", p.ProblemID) }

		// 设置默认值和状态
		newProblem := p
		if newProblem.Order == 0 {
			maxOrder++
			newProblem.Order = maxOrder
		} else {
			// 如果指定了顺序，需要检查是否与现有或其他新题目的顺序冲突 (简化处理：暂不处理复杂排序逻辑，仅追加)
			maxOrder++
			newProblem.Order = maxOrder // 强制追加顺序
		}
		if newProblem.Score == 0 {
			newProblem.Score = 10 // 默认分值
		}
		newProblem.Status = 1 // 默认可用

		problemsToAdd = append(problemsToAdd, newProblem)
	}

	if len(problemsToAdd) == 0 {
		lg.Infof("没有新的题目需要添加到竞赛 %s", contestID)
		return nil // 没有实际需要添加的题目
	}

	//  调用 DAO 添加题目
	err = s.dao.AddProblemsToContest(ctx, contestID, problemsToAdd)
	if err != nil {
		lg.Errorf("向竞赛 %s 添加题目失败: %v", contestID, err)
		return fmt.Errorf("添加题目到竞赛失败: %w", err)
	}

	lg.Infof("成功向竞赛 %s 添加 %d 个题目", contestID, len(problemsToAdd))
	return nil
}

// BatchRemoveProblemsFromContest 批量从竞赛移除题目 (新增)
func (s *Service) BatchRemoveProblemsFromContest(ctx context.Context, contestID string, req *dto.BatchRemoveProblemsReq, userID string, userRole int) error {
	lg := utils.GetDefaultLogger()
	lg.Infof("用户 %s (角色 %d) 尝试从竞赛 %s 批量移除题目: %v", userID, userRole, contestID, req.ProblemIDs)

	//  验证竞赛ID
	_, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		lg.Errorf("无效的竞赛ID格式: %s, 错误: %v", contestID, err)
		return errors.New("无效的竞赛ID")
	}

	//  获取竞赛信息
	contest, err := s.dao.GetContestByID(ctx, contestID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			lg.Warnf("尝试移除题目的竞赛 %s 不存在", contestID)
			return errors.New("竞赛不存在")
		}
		lg.Errorf("获取竞赛 %s 信息失败: %v", contestID, err)
		return fmt.Errorf("获取竞赛信息失败: %w", err)
	}

	//  权限检查 (超级管理员、普通管理员 或 竞赛创建者/教师)
	isSuperAdmin := (userRole == utils.RoleSuperAdmin) // 超级管理员
	isAdmin := (userRole == utils.ClassRoleAdmin)      // 普通管理员 (新增检查)
	isTeacher := (userRole == utils.RoleTeacher)       // 教师
	isCreator := (contest.CreatorID == userID)         // 是否为竞赛创建者

	// 超级管理员和普通管理员可以移除任何竞赛的题目
	// 教师只能移除自己创建的竞赛的题目
	if !(isSuperAdmin || isAdmin || (isTeacher && isCreator)) { // 修改条件，加入 isAdmin
		lg.Warnf("用户 %s (角色 %d) 权限不足，无法从竞赛 %s 移除题目", userID, userRole, contestID)
		return utils.ErrNoPermission // 使用预定义的权限错误
	}

	//  验证要移除的题目ID列表
	if len(req.ProblemIDs) == 0 {
		lg.Warnf("尝试从竞赛 %s 移除题目，但列表为空", contestID)
		return errors.New("要移除的题目ID列表不能为空")
	}

	//  调用 DAO 移除题目
	err = s.dao.BatchRemoveProblemsFromContest(ctx, contestID, req.ProblemIDs)
	if err != nil {
		lg.Errorf("从竞赛 %s 移除题目失败: %v", contestID, err)
		// 不直接返回 DAO 错误，包装一下
		return fmt.Errorf("从竞赛移除题目时发生错误")
	}

	lg.Infof("用户 %s 成功从竞赛 %s 移除 %d 个题目", userID, contestID, len(req.ProblemIDs))
	return nil
}
