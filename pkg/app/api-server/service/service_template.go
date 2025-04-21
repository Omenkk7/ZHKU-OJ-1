package service

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

func (s *Service) MergeTemplate(task interface{}) (code interface{}, err error) {
	lg := utils.GetDefaultLogger()
	t, ok := task.(models.LocalTask)
	if !ok {
		return "", fmt.Errorf("invalid task type")
	}
	lg.Infof("题目id：%v 语言为：%v", t.ProblemId, t.Language)

	// 1. 获取题目基本信息
	problemId, err := primitive.ObjectIDFromHex(t.ProblemId)
	if err != nil {
		return "", fmt.Errorf("invalid problem id: %v", err)
	}

	query := bson.M{"_id": problemId}
	daoProblem, err := s.dao.GetOneProblem(context.Background(), query)
	if err != nil {
		lg.Errorf("find problem error: %v, problemId: %s", err, t.ProblemId)
		return "", fmt.Errorf("failed to find problem: %v", err)
	}

	// 2. 构建模板文件路径
	cfg, _ := utils.LoadConfig("conf/config.yaml")
	templateUrlCfg := cfg.GetTemplateUrl()
	templateDir := filepath.Clean(templateUrlCfg)
	languageDir := strings.ToLower(t.Language) // 统一使用小写避免大小写问题
	templateFile := filepath.Join(templateDir, languageDir, daoProblem.Title+".txt")

	// 检查模板文件是否存在
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		lg.Errorf("template file not found: %s", templateFile)
		return "", fmt.Errorf("template file not found for language: %s", t.Language)
	}

	// 3. 读取模板文件内容
	templateContent, err := ioutil.ReadFile(templateFile)
	if err != nil {
		lg.Errorf("read template file error: %v, path: %s", err, templateFile)
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	// 4. 处理换行符和合并代码
	templateStr := string(templateContent)
	templateStr = strings.ReplaceAll(templateStr, "\\n", "\n")
	userCode := strings.ReplaceAll(t.Code, "\\n", "\n")

	// 替换模板中的占位符
	mergedCode := strings.TrimSpace(strings.Replace(templateStr, "#function", userCode, 1))

	// 5. 验证合并后的代码
	if len(mergedCode) == 0 {
		lg.Error("merged code is empty")
		return "", fmt.Errorf("merged code is empty")
	}

	// 调试输出
	fmt.Println()
	fmt.Println()
	fmt.Println()

	lg.Infoln("\n=== 调试信息 ===")
	lg.Infoln("模板文件路径: %s", templateFile)
	lg.Infoln("模板代码:\n%s", templateStr)
	lg.Infoln("用户代码:\n%s", userCode)
	lg.Infoln("合并后代码:\n%s", mergedCode)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	return mergedCode, nil
}
