package service

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// 写入内容到文件
func writeToFile(filePath string, content string) error { //TODO 写进工具类
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, []byte(content), 0644)
}

// PostProblem 新建题目
func (s *Service) PostProblem(reqProblem *dto.ReqProblem) (err error) {
	lg := utils.GetDefaultLogger()
	//1.题目名称和出题人不能为空
	if reqProblem.Title == "" || reqProblem.Creator == "" {
		lg.Info(utils.ProblemOrAuthorCannotBeNull)
		return errors.New(utils.ProblemOrAuthorCannotBeNull)
	}
	//2.检查题目是否存在
	query := bson.M{
		"title": reqProblem.Title,
	}
	dtoProblem, _ := s.dao.GetOneProblem(context.Background(), query)
	if dtoProblem != nil {
		lg.Infof("题目%s已存在", reqProblem.Title)
		return errors.New(utils.ProblemIsExist)
	}

	//TODO 5.文件存储  临时采用JSON传输模板和测试用例，加速前后端开发  2025/4/21
	//templateUrl/ 语言 / title  模板路径
	//测试用例路径
	cfg, _ := utils.LoadConfig("conf/config.yaml")
	testExampleCfg := cfg.GetTestExampleUrl()
	templateCfg := cfg.GetTemplateUrl()

	testExampleUrl := testExampleCfg + "\\" + reqProblem.Title + ".txt"
	templateGoUrl := templateCfg + "\\go\\" + reqProblem.Title + ".txt"
	templateJavaUrl := templateCfg + "\\java\\" + reqProblem.Title + ".txt"
	templatePythonUrl := templateCfg + "\\python\\" + reqProblem.Title + ".txt"

	//如果环境为linux
	osType := runtime.GOOS
	if osType == "linux" {
		testExampleUrl = testExampleCfg + "/" + reqProblem.Title + ".txt"
		templateGoUrl = templateCfg + "/go/" + reqProblem.Title + ".txt"
		templateJavaUrl = templateCfg + "/java/" + reqProblem.Title + ".txt"
		templatePythonUrl = templateCfg + "/python/" + reqProblem.Title + ".txt"
	}

	//生成文本文件
	if err = writeToFile(testExampleUrl, reqProblem.TestExample); err != nil {
		lg.Infoln("save file failed", err)
		return err
	}
	if reqProblem.GoTemplate != "" {
		if err = writeToFile(templateGoUrl, reqProblem.GoTemplate); err != nil {
			lg.Infoln("save file failed", err)
			return err
		}
	}
	if reqProblem.JavaTemplate != "" {
		if err = writeToFile(templateJavaUrl, reqProblem.JavaTemplate); err != nil {
			lg.Infoln("save file failed", err)
			return err
		}
	}
	if reqProblem.PythonTemplate != "" {
		if err = writeToFile(templatePythonUrl, reqProblem.PythonTemplate); err != nil {
			lg.Infoln("save file failed", err)
			return err
		}
	}

	fmt.Println("osType:", osType)
	fmt.Println("TemplateUrl:", templateGoUrl)
	fmt.Println("TemplateJavaUrl:", templateJavaUrl)
	fmt.Println("TemplatePythonUrl:", templatePythonUrl)
	fmt.Println("TestExampleUrl:", testExampleUrl)

	reqProblem.URL = testExampleUrl

	//TODO ————————————————————文件存储和入库，要保证事务性————————————————————————
	//3.题目不存在，一切正常，构建入库模型
	//TODO 写一个工具类，动态构建入库模型

	dtoProblem = &models.Problem{
		ID:            primitive.ObjectID{},
		Creator:       reqProblem.Creator,
		Description:   reqProblem.Description,
		InputExample:  reqProblem.InputExample,
		OutputExample: reqProblem.OutputExample,
		Status:        reqProblem.Status,
		SubmitNum:     reqProblem.SubmitNum,
		Difficulty:    reqProblem.Difficulty,
		Labels:        reqProblem.Labels,
		PassNum:       reqProblem.PassNum,
		Title:         reqProblem.Title,
		URL:           reqProblem.URL,
		Ctime:         time.Now().Unix(),
		Mtime:         time.Now().Unix(),
	}
	//4.入库
	_, err = s.dao.CreateProblem(context.Background(), dtoProblem)
	if err != nil {
		lg.Info(utils.ServerErr, err)
		return errors.New(utils.ServerErr)
	}
	//TODO —————————————————————————————保证事务一致性———————————————————————————————————————-
	return
}

// GetProblemList  查一堆数据
func (s *Service) GetProblemList(comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	items, err = s.dao.GetProblemList(context.Background(), comQuery)
	return
}

// GetOneProblem 查一个数据
func (s *Service) GetOneProblem(reqProblem *models.Problem) (daoProblem *models.Problem, err error) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query, err := bson.Marshal(reqProblem)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return nil, errors.New(utils.ConstructingBsonErr)
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询daoProblem
	daoProblem, err = s.dao.GetOneProblem(context.Background(), query)
	if daoProblem == nil {
		lg.Info(utils.ProblemNotExist, err)
		return nil, errors.New(utils.ProblemNotExist)
	}
	//响应
	return daoProblem, nil
}

// UpdateProblem 更新题目
func (s *Service) UpdateProblem(reqProblem *dto.ReqProblem) (id string, err error) {
	lg := utils.GetDefaultLogger()
	//selector是筛选条件，update是要更新的内容
	objectId, _ := primitive.ObjectIDFromHex(reqProblem.ID)
	selector := bson.M{
		"_id": objectId,
	}
	reqProblem.Mtime = time.Now().Unix()
	//动态构造bson
	update, err := s.dao.GenerateUpdateBson(reqProblem)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return "", errors.New(utils.ConstructingBsonErr)
	}
	lg.Infof("更新_id:%s\nselector:%s\n", reqProblem.ID, selector)
	_, err = s.dao.UpdateProblem(context.Background(), selector, update)
	if err != nil {
		lg.Info(utils.UpdateErr, err)
		return "", errors.New(utils.UpdateErr)
	}
	return reqProblem.ID, nil
}

// DeleteProblem 通过id删除题目
func (s *Service) DeleteProblem(id string) (_ string, err error) {
	lg := utils.GetDefaultLogger()
	//1.删之前先查询是否有该条数据
	objectId, _ := primitive.ObjectIDFromHex(id)
	selector := bson.M{
		"_id": objectId,
	}
	lg.Infof("查询_id: %s", objectId)

	//2.调用Dao查询daoProblem
	daoProblem, err := s.dao.GetOneProblem(context.Background(), selector)
	if daoProblem == nil {
		lg.Info(utils.ProblemNotExist, err)
		return "", errors.New(utils.ProblemNotExist)
	}

	//3.删除
	lg.Infof("删除id%s", id)
	_, err = s.dao.DeleteProblem(context.Background(), selector)
	if err != nil {
		lg.Info(utils.DeleteErr, err)
		return "", errors.New(utils.DeleteErr)
	}
	lg.Info(utils.DeleteSuccess)
	return id, nil
}
