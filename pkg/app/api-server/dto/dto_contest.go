package dto

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"zhku-oj-server/pkg/models"
)

// 创建竞赛请求
type CreateContestReq struct {
	ContestCode     string                  `json:"contest_code"`                    // 竞赛代码
	Name            string                  `json:"name" binding:"required"`         // 竞赛名称
	Description     string                  `json:"description"`                     // 竞赛描述
	StartTime       int64                   `json:"start_time" binding:"required"`   // 开始时间
	EndTime         int64                   `json:"end_time" binding:"required"`     // 结束时间
	ContestType     int                     `json:"contest_type" binding:"required"` // 竞赛类型
	AccessType      int                     `json:"access_type" binding:"required"`  // 参赛规则
	MaxParticipants int                     `json:"max_participants"`                // 最大参赛人数
	Problems        []models.ContestProblem `json:"problems"`                        // 关联题目列表
}

// 更新竞赛请求
type UpdateContestReq struct {
	Name            string                  `json:"name"`             // 竞赛名称
	Description     string                  `json:"description"`      // 竞赛描述
	StartTime       int64                   `json:"start_time"`       // 开始时间
	EndTime         int64                   `json:"end_time"`         // 结束时间
	ContestType     int                     `json:"contest_type"`     // 竞赛类型
	AccessType      int                     `json:"access_type"`      // 参赛规则
	MaxParticipants int                     `json:"max_participants"` // 最大参赛人数
	Problems        []models.ContestProblem `json:"problems"`         // 关联题目列表
	Status          int                     `json:"status"`           // 竞赛状态
}

// 获取竞赛列表请求
type GetContestListReq struct {
	Page        int    `form:"page"`         // 页码
	PageSize    int    `form:"page_size"`    // 每页数量
	Name        string `form:"name"`         // 竞赛名称
	ContestType int    `form:"contest_type"` // 竞赛类型
	Status      int    `form:"status"`       // 竞赛状态
	CreatorID   string `form:"creator_id"`   // 创建者ID
	AccessType  int    `form:"access_type"`  // 参赛规则
}

// 获取竞赛详情请求
type GetContestDetailReq struct {
	ID string `form:"id" binding:"required"` // 竞赛ID
}

// 删除竞赛请求
type DeleteContestReq struct {
	ID string `json:"id" binding:"required"` // 竞赛ID
}

// 添加参赛者请求
type AddParticipantReq struct {
	ContestID   string `json:"contest_id" binding:"required"`   // 竞赛ID
	StudentID   string `json:"student_id" binding:"required"`   // 学生ID
	StudentName string `json:"student_name" binding:"required"` // 学生姓名
}

// 批量添加参赛者请求
type BatchAddParticipantsReq struct {
	ContestID string        `json:"contest_id" binding:"required"` // 竞赛ID
	Students  []models.User `json:"students" binding:"required"`   // 学生列表
}

// 移除参赛者请求
type RemoveParticipantReq struct {
	ContestID string `json:"contest_id" binding:"required"` // 竞赛ID
	StudentID string `json:"student_id" binding:"required"` // 学生ID
}

// 获取参赛者列表请求
type GetParticipantListReq struct {
	ContestID string `form:"contest_id" binding:"required"` // 竞赛ID
	Page      int    `form:"page" binding:"required"`       // 页码
	PageSize  int    `form:"page_size" binding:"required"`  // 每页数量
	Status    int    `form:"status"`                        // 状态
}

// 审核参赛者请求
type AuditParticipantReq struct {
	ContestID string `json:"contest_id" binding:"required"` // 竞赛ID
	StudentID string `json:"student_id" binding:"required"` // 学生ID
	Status    int    `json:"status" binding:"required"`     // 状态：1-通过，2-拒绝
}

// 获取竞赛排名请求
type GetContestRankingReq struct {
	ContestID string `form:"contest_id" binding:"required"` // 竞赛ID
	Page      int    `form:"page"`                          // 页码
	PageSize  int    `form:"page_size"`                     // 每页数量
}

// 导出竞赛成绩请求
type ExportContestScoreReq struct {
	ContestID string `form:"contest_id" binding:"required"` // 竞赛ID
}

// 竞赛列表响应
type ContestListResp struct {
	Total int64            `json:"total"` // 总数
	List  []models.Contest `json:"list"`  // 列表
}

// 参赛者列表响应
type ParticipantListResp struct {
	Total int64                       `json:"total"` // 总数
	List  []models.ContestParticipant `json:"list"`  // 列表
}

// 竞赛排名响应
type ContestRankingResp struct {
	Total int64                   `json:"total"` // 总数
	List  []models.ContestRanking `json:"list"`  // 列表
}

// 绑定并校验请求参数
func BindAndValid(c *gin.Context, req interface{}) (int, string) {
	err := c.ShouldBind(req)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	return http.StatusOK, ""
}

// ApplyJoinContestReq 学生申请加入竞赛请求
type ApplyJoinContestReq struct {
	ContestID string `json:"contest_id" binding:"required"` // 竞赛ID
}

// 将CreateContestReq转换为Contest
func (req *CreateContestReq) ToContest(creatorID, creatorName string) *models.Contest {
	now := time.Now().Unix()

	// 创建一个 ObjectID 对象用于创建者ID
	creatorObjID, _ := primitive.ObjectIDFromHex(creatorID)

	return &models.Contest{
		ID:              primitive.NewObjectID(),
		ContestCode:     req.ContestCode,
		Name:            req.Name,
		Description:     req.Description,
		StartTime:       req.StartTime / 1000,
		EndTime:         req.EndTime / 1000,
		ContestType:     req.ContestType,
		AccessType:      req.AccessType,
		MaxParticipants: req.MaxParticipants,
		CreatorID:       creatorID,
		CreatorName:     creatorName,
		Problems:        req.Problems,
		Members: map[string][]models.User{
			"1": {{ID: creatorObjID, Username: creatorName}}, // 创建者默认为管理员
		},
		Stats: models.ContestStats{
			ParticipantCount: 0,
			SubmissionCount:  0,
			PassRate:         0,
			LastUpdated:      now,
		},
		Status: 1, // 默认未开始
		Ctime:  now,
		Mtime:  now,
	}
}

// 添加题目到竞赛请求
type AddProblemsToContestReq struct {
	Problems []models.ContestProblem `json:"problems" binding:"required,dive"` // 要添加的题目列表, dive 验证切片内元素
}

// 批量移除竞赛题目请求
type BatchRemoveProblemsReq struct {
	ProblemIDs []string `json:"problem_ids" binding:"required,min=1"` // 要移除的题目ID列表，不能为空
}
