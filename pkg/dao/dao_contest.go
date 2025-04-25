package dao

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// CreateContest 创建竞赛
func (d *Dao) CreateContest(ctx context.Context, contest *models.Contest) (string, error) {
	return d.CreateOne(ctx, contestTable, contest)
}

// UpdateContest 更新竞赛
func (d *Dao) UpdateContest(ctx context.Context, id string, update bson.M) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateDoc := bson.M{"$set": update}

	_, err = d.Update(ctx, contestTable, bson.M{"_id": objectID}, updateDoc)
	return err
}

// GetContestByID 根据ID获取竞赛
func (d *Dao) GetContestByID(ctx context.Context, id string) (*models.Contest, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var contest models.Contest
	err = d.mongo.FindOne(ctx, contestTable, bson.M{"_id": objectID, "status": bson.M{"$ne": 0}}, &contest)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(utils.FindContestErr)
		}
		return nil, err
	}

	return &contest, nil
}

// GetContestList 获取竞赛列表
func (d *Dao) GetContestList(ctx context.Context, query bson.M, page, pageSize int) ([]*models.Contest, int64, error) {
	// 设置分页
	skip := (page - 1) * pageSize
	limit := int64(pageSize)

	// 设置排序
	opts := options.Find().SetSkip(int64(skip)).SetLimit(limit).SetSort(bson.M{"ctime": -1})

	// 查询总数
	total, err := d.mongo.Count(ctx, contestTable, query)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	var contests []*models.Contest
	err = d.mongo.Find(ctx, contestTable, query, &contests, opts)
	if err != nil {
		return nil, 0, err
	}

	return contests, total, nil
}

// DeleteContest 删除竞赛（软删除）
func (d *Dao) DeleteContest(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// 使用 $set 操作符包装更新内容
	updateDoc := bson.M{"$set": bson.M{"status": 0, "mtime": time.Now().Unix()}}

	_, err = d.Update(ctx, contestTable, bson.M{"_id": objectID}, updateDoc)
	return err
}

// AddParticipant 添加参赛者
func (d *Dao) AddParticipant(ctx context.Context, participant *models.ContestParticipant) (string, error) {
	// 检查是否已经存在
	var existingParticipant models.ContestParticipant
	result, err := d.GetOne(ctx, contestParticipantTable, &existingParticipant, bson.M{
		"contest_id": participant.ContestID,
		"student_id": participant.StudentID,
	})

	// 如果已经存在，则返回错误
	if result != nil {
		return "", errors.New("participant already exists")
	}

	// 如果是其他错误（不是未找到文档的错误），则返回错误
	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}

	// 插入新参赛者
	id, err := d.CreateOne(ctx, contestParticipantTable, participant)
	if err != nil {
		return "", err
	}

	// 更新竞赛的参赛人数
	objectID, err := primitive.ObjectIDFromHex(participant.ContestID)
	if err != nil {
		return "", err
	}
	_, err = d.Update(ctx, contestTable, bson.M{"_id": objectID}, bson.M{"$inc": bson.M{"stats.participant_count": 1}})
	if err != nil {
		return "", err
	}

	return id, nil
}

// RemoveParticipant 移除参赛者
func (d *Dao) RemoveParticipant(ctx context.Context, contestID, studentID string) error {
	// 删除参赛者
	_, err := d.DeleteOne(ctx, contestParticipantTable, bson.M{
		"contest_id": contestID,
		"student_id": studentID,
	})
	if err != nil {
		return err
	}

	// 更新竞赛的参赛人数
	objectID, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		return err
	}
	_, err = d.Update(ctx, contestTable, bson.M{"_id": objectID}, bson.M{"$inc": bson.M{"stats.participant_count": -1}})
	if err != nil {
		return err
	}

	// 从竞赛的members字段中删除学生
	err = d.RemoveMember(ctx, contestID, strconv.Itoa(utils.RoleStudent), studentID)
	if err != nil {
		// 记录错误但不中断流程
		lg := utils.GetDefaultLogger()
		lg.Errorf("从竞赛members中移除学生失败: %v", err)
	}

	return nil
}

// GetParticipantList 获取参赛者列表
func (d *Dao) GetParticipantList(ctx context.Context, query bson.M, page, pageSize int) ([]*models.ContestParticipant, int64, error) {
	// 设置分页
	skip := (page - 1) * pageSize
	limit := int64(pageSize)

	// 设置排序
	opts := options.Find().SetSkip(int64(skip)).SetLimit(limit).SetSort(bson.M{"register_time": -1})

	// 查询总数
	total, err := d.mongo.Count(ctx, contestParticipantTable, query)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	var participants []*models.ContestParticipant
	err = d.mongo.Find(ctx, contestParticipantTable, query, &participants, opts)
	if err != nil {
		return nil, 0, err
	}

	return participants, total, nil
}

// UpdateParticipantStatus 更新参赛者状态
func (d *Dao) UpdateParticipantStatus(ctx context.Context, contestID, studentID string, status int) error {
	_, err := d.Update(ctx, contestParticipantTable, bson.M{
		"contest_id": contestID,
		"student_id": studentID,
	}, bson.M{
		"$set": bson.M{
			"status": status,
			"mtime":  time.Now().Unix(),
		},
	})
	return err
}

// GetContestRanking 获取竞赛排名
func (d *Dao) GetContestRanking(ctx context.Context, contestID string, page, pageSize int) ([]*models.ContestRanking, int64, error) {
	// 设置分页
	skip := (page - 1) * pageSize
	limit := int64(pageSize)

	// 设置排序（先按解决题目数降序，再按总用时升序）
	opts := options.Find().SetSkip(int64(skip)).SetLimit(limit).SetSort(bson.D{
		{Key: "solved_problems", Value: -1},
		{Key: "total_time", Value: 1},
	})

	// 查询条件：排除系统生成的占位记录
	query := bson.M{
		"contest_id": contestID,
		"student_id": bson.M{"$ne": ""}, // 排除空学生ID的记录
	}

	// 查询总数
	total, err := d.mongo.Count(ctx, contestRankingTable, query)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	var rankings []*models.ContestRanking
	err = d.mongo.Find(ctx, contestRankingTable, query, &rankings, opts)
	if err != nil {
		return nil, 0, err
	}

	return rankings, total, nil
}

// UpdateContestRanking 更新竞赛排名
func (d *Dao) UpdateContestRanking(ctx context.Context, ranking *models.ContestRanking) error {
	// 检查是否已经存在
	var existingRanking models.ContestRanking
	result, err := d.GetOne(ctx, contestRankingTable, &existingRanking, bson.M{
		"contest_id": ranking.ContestID,
	})

	if result != nil {
		// 更新整个排名记录
		_, err = d.Update(ctx, contestRankingTable, bson.M{
			"contest_id": ranking.ContestID,
		}, bson.M{
			"$set": bson.M{
				"contest_name":   ranking.ContestName,
				"rankings":       ranking.Rankings,
				"generated_time": ranking.GeneratedTime,
				"mtime":          time.Now().Unix(),
			},
		})
		return err
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// 插入新排名
	_, err = d.CreateOne(ctx, contestRankingTable, ranking)
	return err
}

// UpdateContestStats 更新竞赛统计信息
func (d *Dao) UpdateContestStats(ctx context.Context, contestID string, stats models.ContestStats) error {
	objectID, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		return err
	}
	_, err = d.Update(ctx, contestTable, bson.M{"_id": objectID}, bson.M{"stats": stats, "mtime": time.Now().Unix()})
	return err
}

// UpdateContestStatus 更新竞赛状态
func (d *Dao) UpdateContestStatus(ctx context.Context, contestID string, status int) error {
	objectID, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		return err
	}
	// 使用 $set 操作符包装更新内容
	updateDoc := bson.M{"$set": bson.M{"status": status, "mtime": time.Now().Unix()}}
	_, err = d.Update(ctx, contestTable, bson.M{"_id": objectID}, updateDoc)
	return err
}

// AddMember 添加竞赛成员
func (d *Dao) AddMember(ctx context.Context, contestID string, role string, user models.User) error {
	objectID, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		return err
	}
	_, err = d.Update(ctx, contestTable, bson.M{"_id": objectID}, bson.M{"$push": bson.M{"members." + role: user}})
	return err
}

// RemoveMember 移除竞赛成员
func (d *Dao) RemoveMember(ctx context.Context, contestID string, role string, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		return err
	}

	// 使用utils包中的日志函数记录日志
	lg := utils.GetDefaultLogger()

	// 先查询一下当前的members结构
	var contest models.Contest
	err = d.mongo.FindOne(ctx, contestTable, bson.M{"_id": objectID}, &contest)
	if err != nil {
		lg.Errorf("查询竞赛失败: %v", err)
		return err
	}

	lg.Infof("当前竞赛members结构: %+v", contest.Members)

	// 检查members字段是否为空
	if contest.Members == nil || len(contest.Members) == 0 {
		lg.Warnf("竞赛 %s 的members字段为空", contestID)
		return nil
	}

	// 将userID转换为ObjectID，因为members中存储的是ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		lg.Errorf("用户ID转换为ObjectID失败: %v", err)
		return err
	}

	// 执行移除操作，使用ID字段而不是user_id字段
	result, err := d.Update(ctx, contestTable, bson.M{"_id": objectID}, bson.M{
		"$pull": bson.M{"members." + role: bson.M{"_id": userObjID}}})
	if err != nil {
		lg.Errorf("移除成员失败: %v", err)
		return err
	}

	lg.Infof("移除成员结果: %v", result)

	return nil
}

// GetParticipantByID 根据ID获取参赛者
func (d *Dao) GetParticipantByID(ctx context.Context, contestID, studentID string) (*models.ContestParticipant, error) {
	// 使用日志记录查询条件
	lg := utils.GetDefaultLogger()
	lg.Infof("查询条件： map[contest_id:%s student_id:%s]", contestID, studentID)

	// 创建一个空的参赛者对象
	var participant models.ContestParticipant

	// 直接使用 FindOne 方法查询并解析到结构体
	err := d.mongo.FindOne(
		ctx,
		contestParticipantTable,
		bson.M{
			"contest_id": contestID,
			"student_id": studentID,
		},
		&participant,
	)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		lg.Errorf("查询参赛者失败: %v", err)
		return nil, err
	}

	return &participant, nil
}

// UpdateParticipantResult 更新参赛者成绩
func (d *Dao) UpdateParticipantResult(ctx context.Context, contestID, studentID string, result models.ContestResult) error {
	_, err := d.Update(ctx, contestParticipantTable, bson.M{
		"contest_id": contestID,
		"student_id": studentID,
	}, bson.M{
		"result": result,
		"mtime":  time.Now().Unix(),
	})
	return err
}

// GetContestsByStatus 根据状态获取竞赛列表
func (d *Dao) GetContestsByStatus(ctx context.Context, status int) ([]*models.Contest, error) {
	var contests []*models.Contest
	err := d.mongo.Find(ctx, contestTable, bson.M{"status": status}, &contests, nil)
	if err != nil {
		return nil, err
	}
	return contests, nil
}

// GetContestsByCreator 根据创建者获取竞赛列表
func (d *Dao) GetContestsByCreator(ctx context.Context, creatorID string, page, pageSize int) ([]*models.Contest, int64, error) {
	return d.GetContestList(ctx, bson.M{"creator_id": creatorID, "status": bson.M{"$ne": 0}}, page, pageSize)
}

// GetPublicContests 获取公开竞赛列表
func (d *Dao) GetPublicContests(ctx context.Context, page, pageSize int) ([]*models.Contest, int64, error) {
	return d.GetContestList(ctx, bson.M{"access_type": utils.ContestAccessPublic, "status": bson.M{"$ne": utils.ContestStatusDeleted}}, page, pageSize)
}

// GetContestByCode 根据竞赛代码获取竞赛
func (d *Dao) GetContestByCode(ctx context.Context, code string) (*models.Contest, error) {
	var contest models.Contest
	result, err := d.GetOne(ctx, contestTable, &contest, bson.M{"contest_code": code, "status": bson.M{"$ne": 0}})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, mongo.ErrNoDocuments
	}
	return result.(*models.Contest), nil
}

// CreateContestParticipant 创建竞赛参与者记录
func (d *Dao) CreateContestParticipant(ctx context.Context, participant *models.ContestParticipant) (string, error) {
	return d.CreateOne(ctx, contestParticipantTable, participant)
}

// IsContestParticipant 检查用户是否已是竞赛参与者
func (d *Dao) IsContestParticipant(ctx context.Context, contestID, studentID string) (bool, error) {
	count, err := d.mongo.Count(ctx, contestParticipantTable, bson.M{
		"contest_id": contestID,
		"student_id": studentID,
		"status":     utils.ContestParticipantStatusApproved, // 已通过状态
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasAppliedContest 检查用户是否已申请参加竞赛
func (d *Dao) HasAppliedContest(ctx context.Context, contestID, studentID string) (bool, error) {
	count, err := d.mongo.Count(ctx, contestParticipantTable, bson.M{
		"contest_id": contestID,
		"student_id": studentID,
		"status":     utils.ContestParticipantStatusPending, // 待审核状态
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// AddProblemsToContest 向竞赛添加题目
func (d *Dao) AddProblemsToContest(ctx context.Context, contestID string, problemsToAdd []models.ContestProblem) error {
	objectID, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 检查题目列表是否为空
	if len(problemsToAdd) == 0 {
		return errors.New("要添加的题目列表不能为空")
	}

	// 构造更新操作
	update := bson.M{
		"$push": bson.M{
			"problems": bson.M{
				"$each": problemsToAdd,
			},
		},
		"$set": bson.M{ // 同时更新修改时间
			"mtime": time.Now().Unix(),
		},
	}

	// 执行更新
	_, err = d.mongo.UpdateOne(ctx, contestTable, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("向竞赛添加题目失败: %w", err)
	}

	return nil
}

// RemoveProblemsFromContest 从竞赛中移除题目
func (d *Dao) RemoveProblemsFromContest(ctx context.Context, contestID string, problemIDs []string) error {
	objectID, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	// 检查题目ID列表是否为空
	if len(problemIDs) == 0 {
		return errors.New("要移除的题目ID列表不能为空")
	}

	// 构造更新操作 - 使用$pull操作符移除匹配的题目
	update := bson.M{
		"$pull": bson.M{
			"problems": bson.M{
				"problem_id": bson.M{
					"$in": problemIDs,
				},
			},
		},
	}

	// 执行更新操作
	_, err = d.Update(ctx, contestTable, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("移除题目失败: %w", err)
	}

	return nil
}

// BatchRemoveProblemsFromContest 批量从竞赛移除题目
func (d *Dao) BatchRemoveProblemsFromContest(ctx context.Context, contestID string, problemIDs []string) error {
	objectID, err := primitive.ObjectIDFromHex(contestID)
	if err != nil {
		return errors.New("无效的竞赛ID")
	}

	if len(problemIDs) == 0 {
		return errors.New("要移除的题目ID列表不能为空")
	}

	// 构造更新操作，使用 $pull 从数组中移除匹配的元素
	update := bson.M{
		"$pull": bson.M{
			"problems": bson.M{
				"problem_id": bson.M{
					"$in": problemIDs, // 移除 problem_id 在给定列表中的所有题目
				},
			},
		},
	}

	// 执行更新
	// 注意：这里使用 UpdateOne，因为我们是针对单个竞赛文档进行操作
	_, err = d.mongo.UpdateOne(ctx, contestTable, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("从竞赛移除题目失败: %w", err)
	}
	return nil
}
