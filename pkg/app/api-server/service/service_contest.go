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
	// 验证时间
	if req.StartTime >= req.EndTime {
		return "", errors.New("开始时间必须早于结束时间")
	}

	// 验证题目
	if len(req.Problems) > 0 {
		for i, p := range req.Problems {
			// 验证题目ID是否有效
			_, err := primitive.ObjectIDFromHex(p.ProblemID)
			if err != nil {
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

	// 创建竞赛
	contestID, err := s.dao.CreateContest(ctx, contest)
	if err != nil {
		return "", err
	}

	// 初始化竞赛排行榜（空的）
	err = s.CreateContestRanking(ctx, contestID)
	if err != nil {
		// 记录错误但不影响竞赛创建
		lg := utils.GetDefaultLogger()
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
func (s *Service) GetContestList(ctx context.Context, req *dto.GetContestListReq) (*dto.ContestListResp, error) {
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
	if req.CreatorID != "" {
		query["creator_id"] = req.CreatorID
	}
	if req.AccessType > 0 { // 只有当明确设置为公开(1)时才过滤
		query["access_type"] = req.AccessType
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
	// 验证ID
	contestID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 获取竞赛
	contest, err := s.dao.GetContestByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("竞赛不存在或已被删除")
		}
		return err
	}

	hasPermission, err := s.CheckContestPermission(ctx, id, userID, []string{utils.ContestRoleAdmin, utils.ContestRoleTeacher})
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("无权限更新此竞赛")
	}

	// 验证状态
	if contest.Status != utils.ContestStatusNotStart {
		return errors.New("只能修改未开始的竞赛")
	}

	// 构建更新对象
	update := bson.M{"mtime": time.Now().Unix()}

	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.StartTime > 0 && req.EndTime > 0 {
		if req.StartTime >= req.EndTime {
			return errors.New("开始时间必须早于结束时间")
		}
		update["start_time"] = req.StartTime
		update["end_time"] = req.EndTime
	} else if req.StartTime > 0 {
		if req.StartTime >= contest.EndTime {
			return errors.New("开始时间必须早于结束时间")
		}
		update["start_time"] = req.StartTime
	} else if req.EndTime > 0 {
		if contest.StartTime >= req.EndTime {
			return errors.New("开始时间必须早于结束时间")
		}
		update["end_time"] = req.EndTime
	}
	if req.ContestType > 0 {
		update["contest_type"] = req.ContestType
	}
	if req.AccessType >= 0 {
		update["access_type"] = req.AccessType
	}
	if req.MaxParticipants >= 0 {
		update["max_participants"] = req.MaxParticipants
	}
	if req.Status > 0 {
		update["status"] = req.Status
	}

	// 更新题目
	if len(req.Problems) > 0 {
		problems := make([]models.ContestProblem, 0, len(req.Problems))
		for _, p := range req.Problems {
			problemID, err := primitive.ObjectIDFromHex(p.ProblemID)
			if err != nil {
				return errors.New("无效的题目ID: " + p.ProblemID)
			}
			problems = append(problems, models.ContestProblem{
				ProblemID: problemID.Hex(),
				Order:     p.Order,
				Score:     p.Score,
				Status:    1, // 默认可用
			})
		}
		update["problems"] = problems
	}

	// 更新到数据库
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
	hasPermission, err := s.CheckContestPermission(ctx, id, userID, []string{utils.ContestRoleAdmin, utils.ContestRoleTeacher})
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
	hasPermission, err := s.CheckContestPermission(ctx, id, userID, []string{utils.ContestRoleAdmin, utils.ContestRoleTeacher})
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("无权限归档此竞赛")
	}

	// 验证状态
	if contest.Status != utils.ContestStatusEnded {
		return errors.New("只能归档已结束的竞赛")
	}

	// 更新状态
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
	hasPermission, err := s.CheckContestPermission(ctx, id, userID, []string{utils.ContestRoleAdmin})
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
		hasPermission, err := s.CheckContestPermission(ctx, req.ContestID, userID, []string{utils.ContestRoleAdmin, utils.ContestRoleTeacher})
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
	return s.dao.AddMember(ctx, contestID.Hex(), utils.ContestRoleStudent, user)
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
	hasPermission, err := s.CheckContestPermission(ctx, req.ContestID, userID, []string{utils.ContestRoleAdmin, utils.ContestRoleTeacher})
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
		err = s.dao.AddMember(ctx, contestID.Hex(), utils.ContestRoleStudent, student)
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
	hasPermission, err := s.CheckContestPermission(ctx, req.ContestID, userID, []string{utils.ContestRoleAdmin, utils.ContestRoleTeacher})
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
	hasPermission, err := s.CheckContestPermission(ctx, req.ContestID, userID, []string{utils.ContestRoleAdmin, utils.ContestRoleTeacher})
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
func (s *Service) CheckContestPermission(ctx context.Context, contestID string, userID string, allowedRoles []string) (bool, error) {
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
		members, ok := contest.Members[role]
		if !ok {
			continue
		}
		for _, member := range members {
			if member.ID.Hex() == userID { // 修正这里，使用 ID.Hex() 而不是 UserID
				hasPermission = true
				return hasPermission, nil
			}
		}
	}

	return hasPermission, nil
}
