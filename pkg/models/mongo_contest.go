package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Contest 竞赛表
type Contest struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ContestCode     string             `json:"contest_code,omitempty" bson:"contest_code,omitempty"`         // 竞赛代码：竞赛类型缩写+年份+随机4位数字，如WS(周赛)、MS(月赛)
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`                         // 竞赛名称
	Description     string             `json:"description,omitempty" bson:"description,omitempty"`           // 竞赛描述
	StartTime       int64              `json:"start_time,omitempty" bson:"start_time,omitempty"`             // 开始时间
	EndTime         int64              `json:"end_time,omitempty" bson:"end_time,omitempty"`                 // 结束时间
	ContestType     int                `json:"contest_type,omitempty" bson:"contest_type,omitempty"`         // 竞赛类型：1(周赛)、2(月赛)、3(挑战赛)等
	AccessType      int                `json:"access_type,omitempty" bson:"access_type,omitempty"`           // 参赛规则：1(公开)、0(私有)
	MaxParticipants int                `json:"max_participants,omitempty" bson:"max_participants,omitempty"` // 最大参赛人数，0表示不限制
	CreatorID       string             `json:"creator_id,omitempty" bson:"creator_id,omitempty"`             // 创建者ID
	CreatorName     string             `json:"creator_name,omitempty" bson:"creator_name,omitempty"`         // 创建者姓名(冗余字段，便于查询)
	Problems        []ContestProblem   `json:"problems,omitempty" bson:"problems,omitempty"`                 // 关联题目列表
	Members         map[string][]User  `json:"members,omitempty" bson:"members,omitempty"`                   // 权限管理
	Stats           ContestStats       `json:"stats,omitempty" bson:"stats,omitempty"`                       // 竞赛统计信息
	Status          int                `json:"status,omitempty" bson:"status,omitempty"`                     // 竞赛状态：1-未开始，2-进行中，3-已结束，4-已归档，0-已删除
	Ctime           int64              `json:"ctime,omitempty" bson:"ctime,omitempty"`                       // 创建时间
	Mtime           int64              `json:"mtime,omitempty" bson:"mtime,omitempty"`                       // 修改时间
}

// ContestProblem 竞赛题目
type ContestProblem struct {
	ProblemID string `json:"problem_id,omitempty" bson:"problem_id,omitempty"` // 题目ID
	Order     int    `json:"order,omitempty" bson:"order,omitempty"`           // 题目在竞赛中的顺序
	Score     int    `json:"score,omitempty" bson:"score,omitempty"`           // 题目分值
	Status    int    `json:"status,omitempty" bson:"status,omitempty"`         // 状态：1-可用，0-不可用(题目被隐藏时)
}

// ContestStats 竞赛统计信息
type ContestStats struct {
	ParticipantCount int     `json:"participant_count,omitempty" bson:"participant_count,omitempty"` // 参赛人数
	SubmissionCount  int     `json:"submission_count,omitempty" bson:"submission_count,omitempty"`   // 提交次数
	PassRate         float64 `json:"pass_rate,omitempty" bson:"pass_rate,omitempty"`                 // 通过率
	LastUpdated      int64   `json:"last_updated,omitempty" bson:"last_updated,omitempty"`           // 最后更新时间
}

// ContestParticipant 参赛者表
type ContestParticipant struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ContestID    string             `json:"contest_id,omitempty" bson:"contest_id,omitempty"`       // 竞赛ID
	ContestName  string             `json:"contest_name,omitempty" bson:"contest_name,omitempty"`   // 竞赛名称(冗余字段，便于查询)
	StudentID    string             `json:"student_id,omitempty" bson:"student_id,omitempty"`       // 学生ID
	StudentName  string             `json:"student_name,omitempty" bson:"student_name,omitempty"`   // 学生姓名(冗余字段，便于查询)
	RegisterTime int64              `json:"register_time,omitempty" bson:"register_time,omitempty"` // 报名时间
	Score        int                `json:"score,omitempty" bson:"score,omitempty"`                 // 总分数
	Rank         int                `json:"rank,omitempty" bson:"rank,omitempty"`                   // 排名
	Status       int                `json:"status,omitempty" bson:"status,omitempty"`               // 状态：0-待审核，1-已通过，2-已拒绝
	Ctime        int64              `json:"ctime,omitempty" bson:"ctime,omitempty"`                 // 创建时间
	Mtime        int64              `json:"mtime,omitempty" bson:"mtime,omitempty"`                 // 修改时间
	Result       ContestResult      `json:"result" bson:"result"`                                   // 竞赛结果
}

// ContestRanking 竞赛排名记录
type ContestRanking struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ContestID     string             `json:"contest_id,omitempty" bson:"contest_id,omitempty"`
	ContestName   string             `json:"contest_name,omitempty" bson:"contest_name,omitempty"`
	Rankings      []RankingItem      `json:"rankings,omitempty" bson:"rankings,omitempty"`
	GeneratedTime int64              `json:"generated_time,omitempty" bson:"generated_time,omitempty"`
}

// RankingItem 排名项
type RankingItem struct {
	Rank        int    `json:"rank,omitempty" bson:"rank,omitempty"`
	StudentID   string `json:"student_id,omitempty" bson:"student_id,omitempty"`
	StudentName string `json:"student_name,omitempty" bson:"student_name,omitempty"`
	SolvedCount int    `json:"solved_count,omitempty" bson:"solved_count,omitempty"`
	TotalScore  int    `json:"total_score,omitempty" bson:"total_score,omitempty"`
	TotalTime   int    `json:"total_time,omitempty" bson:"total_time,omitempty"`
}

// ContestResult 竞赛结果
type ContestResult struct {
	SolvedProblems int                      `json:"solved_problems" bson:"solved_problems"`   // 解决的题目数量
	TotalScore     int                      `json:"total_score" bson:"total_score"`           // 总分数
	TotalTime      int64                    `json:"total_time" bson:"total_time"`             // 总用时（毫秒）
	LastSubmitTime int64                    `json:"last_submit_time" bson:"last_submit_time"` // 最后提交时间
	ProblemResults map[string]ProblemResult `json:"problem_results" bson:"problem_results"`   // 每道题的结果，key为题目ID
}

// ProblemResult 题目结果
type ProblemResult struct {
	Status         int   `json:"status" bson:"status"`                     // 状态：0-未提交，1-已解决，2-尝试过但未解决
	Score          int   `json:"score" bson:"score"`                       // 得分
	SubmitCount    int   `json:"submit_count" bson:"submit_count"`         // 提交次数
	SolvedTime     int64 `json:"solved_time" bson:"solved_time"`           // 解决时间（相对于竞赛开始时间的毫秒数）
	FirstACTime    int64 `json:"first_ac_time" bson:"first_ac_time"`       // 首次通过时间
	LastSubmitTime int64 `json:"last_submit_time" bson:"last_submit_time"` // 最后提交时间
}
