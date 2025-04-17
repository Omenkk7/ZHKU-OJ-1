package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"math"
	"os"
	"strconv"
	"strings"
	"zhku-oj-server/pkg/dao"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"

	"io/ioutil"
	"net/http"
)

type Cmd struct {
	Args          []string        `json:"args,omitempty"`
	Env           []string        `json:"env,omitempty"`
	Files         []File          `json:"files,omitempty"`
	CPULimit      int64           `json:"cpuLimit,omitempty"`
	MemoryLimit   int64           `json:"memoryLimit,omitempty"`
	ProcLimit     int             `json:"procLimit,omitempty"`
	CopyIn        map[string]File `json:"copyIn,omitempty"`
	CopyOut       []string        `json:"copyOut,omitempty"`
	CopyOutCached []string        `json:"copyOutCached,omitempty"`
}

type File struct {
	Content string `json:"content,omitempty"`
	Name    string `json:"name,omitempty"`
	Max     int    `json:"max,omitempty"`
	FileId  string `json:"fileId,omitempty"`
}

type RequestBody struct {
	Cmd []Cmd `json:"cmd,omitempty"`
}

// 定义响应体的结构
type ResponseBody struct {
	Status     string            `json:"status,omitempty"`
	ExitStatus int               `json:"exitStatus,omitempty"`
	Time       int64             `json:"time,omitempty"`
	Memory     int64             `json:"memory,omitempty"`
	RunTime    int64             `json:"runTime,omitempty"`
	Files      map[string]string `json:"files,omitempty"`
	FileIds    map[string]string `json:"fileIds,omitempty"`
}

// 用于编译代码，拿到编译后的filedId
func getFieldId(language string, code string) (fileId string) {
	lg := utils.GetDefaultLogger()
	//通过读取配置文件的形式，找到调用沙箱的请求方式和url，找到jdk/g++等环境的路径
	cfg, _ := utils.LoadConfig("conf/config.yaml")
	judgeCfg := cfg.GetJudgeConfig()
	sandboxCfg := cfg.GetSandboxConfig()

	//创建请求体
	var requestBody RequestBody
	//TODO 1、支持其他语言
	//TODO 2、减少硬编码，如 a.java，content：“1 1"
	//TODO 3、构建请求体，构建json抽取为工具/方法
	switch language {
	case "java":
		jdk11cfg, _ := judgeCfg.Java["JDK11"]
		requestBody = RequestBody{
			Cmd: []Cmd{
				{
					Args:        []string{jdk11cfg.BaseArgs},
					Env:         []string{jdk11cfg.Env},
					Files:       []File{{Content: "1 1"}, {Name: "stdout", Max: 10240000}, {Name: "stderr", Max: 10240000}},
					CPULimit:    968750000000,
					MemoryLimit: 1048576000000,
					ProcLimit:   50,
					CopyIn: map[string]File{
						"a.java": {
							Content: code,
						},
					},
					CopyOut:       []string{"stdout", "stderr"},
					CopyOutCached: []string{"a.java"},
				},
			},
		}
	case "go":
		goCfg, _ := judgeCfg.Go["Go1.23.5"]
		requestBody = RequestBody{
			Cmd: []Cmd{
				{
					Args:        []string{goCfg.BaseArgs, "run", "main.go"},
					Env:         []string{goCfg.Env},
					Files:       []File{{Content: "1 1"}, {Name: "stdout", Max: 10240000}, {Name: "stderr", Max: 10240000}},
					CPULimit:    968750000000,
					MemoryLimit: 1048576000000,
					ProcLimit:   50,
					CopyIn: map[string]File{
						"main.go": {
							Content: code,
						},
					},
					CopyOut:       []string{"stdout", "stderr"},
					CopyOutCached: []string{"main.go"},
				},
			},
		}
	case "python":
		pyCfg, _ := judgeCfg.Python["python310"]
		requestBody = RequestBody{
			Cmd: []Cmd{
				{
					Args: []string{pyCfg.BaseArgs, "main.py"},
					Env: []string{pyCfg.Env,
						"PYTHONHASHSEED=0",
						"PYTHONIOENCODING=UTF-8"}, //不设置这个的话，go-judge会报错failed to get random numbers to initialize Python
					Files:       []File{{Content: "1 1"}, {Name: "stdout", Max: 10240000}, {Name: "stderr", Max: 10240000}},
					CPULimit:    968750000000,
					MemoryLimit: 1048576000000,
					ProcLimit:   50,
					CopyIn: map[string]File{
						"main.py": {
							Content: code,
						},
					},
					CopyOut:       []string{"stdout", "stderr"},
					CopyOutCached: []string{"main.py"},
				},
			},
		}
	}

	// 将请求体编码为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		lg.Infoln("Error encoding JSON:", err)
		return ""
	}

	// 创建HTTP请求
	req, err := http.NewRequest(sandboxCfg.Method, sandboxCfg.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		lg.Infoln("Error creating request:", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		lg.Infoln("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		lg.Infoln("Error reading response body:", err)
		return ""
	}

	// 解析响应体
	var response []ResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		lg.Infoln("Error decoding JSON:", err)
		return ""
	}
	//返回文件id
	switch language {
	case "java":
		return response[0].FileIds["a.java"]
	case "go":
		return response[0].FileIds["main.go"]
	case "python":
		return response[0].FileIds["main.py"]
	}
	return ""
}

// 判题
func judge(fileId string, language string, example string) (result string) {
	lg := utils.GetDefaultLogger()
	cfg, _ := utils.LoadConfig("conf/config.yaml")
	sandboxCfg := cfg.GetSandboxConfig()
	judgeCfg := cfg.GetJudgeConfig()

	//创建请求体
	var requestBody RequestBody
	switch language {
	case "java":
		jdk11cfg, _ := judgeCfg.Java["JDK11"]
		requestBody = RequestBody{
			Cmd: []Cmd{
				{
					Args:        []string{jdk11cfg.BaseArgs, "a.java", "UTF-8"},
					Env:         []string{jdk11cfg.Env},
					Files:       []File{{Content: example}, {Name: "stdout", Max: 10240000}, {Name: "stderr", Max: 10240000}},
					CPULimit:    9687500000,
					MemoryLimit: 10485760000,
					ProcLimit:   50,
					CopyIn: map[string]File{
						"a.java": {
							FileId: fileId,
						},
					},
				},
			},
		}
	case "go":
		goCfg, _ := judgeCfg.Go["Go1.23.5"]
		requestBody = RequestBody{
			Cmd: []Cmd{
				{
					Args:        []string{goCfg.BaseArgs, "run", "main.go"},
					Env:         []string{goCfg.Env},
					Files:       []File{{Content: example}, {Name: "stdout", Max: 10240000}, {Name: "stderr", Max: 10240000}},
					CPULimit:    9687500000000,
					MemoryLimit: 10485760000000,
					ProcLimit:   50,
					CopyIn: map[string]File{
						"main.go": {
							FileId: fileId,
						},
					},
				},
			},
		}
	case "python":
		pyCfg, _ := judgeCfg.Python["Python310"]
		requestBody = RequestBody{
			Cmd: []Cmd{
				{
					Args: []string{pyCfg.BaseArgs, "main.py"},
					Env: []string{pyCfg.Env,
						"PYTHONHASHSEED=0",
						"PYTHONIOENCODING=UTF-8"}, //不设置这个的话，go-judge会报错failed to get random numbers to initialize Python
					Files:       []File{{Content: example}, {Name: "stdout", Max: 10240000}, {Name: "stderr", Max: 10240000}},
					CPULimit:    9687500000000,
					MemoryLimit: 10485760000000,
					ProcLimit:   50,
					CopyIn: map[string]File{
						"main.py": {
							FileId: fileId,
						},
					},
				},
			},
		}
	}

	// 将请求体编码为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		lg.Infoln("Error encoding JSON:", err)
		return ""
	}

	// 创建HTTP请求
	req, err := http.NewRequest(sandboxCfg.Method, sandboxCfg.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		lg.Infoln("Error creating request:", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		lg.Infoln("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		lg.Infoln("Error reading response body:", err)
		return ""
	}

	// 解析响应体
	var response []ResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		lg.Infoln("Error decoding JSON:", err)
		return ""
	}

	// 检查状态和输出
	if len(response) > 0 {
		if stdout, ok := response[0].Files["stdout"]; ok {
			return stdout
		} else {
			lg.Infoln("stdout is empty or not found in response")
		}
	}
	lg.Infoln("Response is empty")
	return response[0].Files["stderr"]
}

func getTestExampleUrl(problemId string) (url string) {
	lg := utils.GetDefaultLogger()
	//1.template集合内查模板，problemID——>模板
	objectID, _ := primitive.ObjectIDFromHex(problemId)
	query := bson.M{
		"_id": objectID,
	}
	daoProblem, err := dao.NewDao().GetOneProblem(context.Background(), query)
	if err != nil {
		lg.Info("find template error:", err)
		return ""
	}
	return daoProblem.URL
}

// 通过url读取测试用例
// TODO 目前是读本地txt文件的形式，后续要更改
func fetchTestCases(filePath string) ([]string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 按行分割内容
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	// 验证格式正确性
	if len(lines)%2 != 0 {
		return nil, errors.New("测试用例文件格式错误：行数必须为偶数")
	}

	return lines, nil
}

// 判断是否相等
func isEqual(expected, actual string) bool {
	expectedFloat, err1 := strconv.ParseFloat(expected, 64)
	actualFloat, err2 := strconv.ParseFloat(actual, 64)

	if err1 != nil || err2 != nil {
		// 如果无法解析为浮点数，直接进行字符串比较
		return expected == actual
	}

	// 设置精度阈值
	const epsilon = 1e-6
	return math.Abs(expectedFloat-actualFloat) < epsilon
}

// 调用Oj
func invokeSandbox(task *models.LocalTask) {
	lg := utils.GetDefaultLogger()

	// 获取编译后的文件ID
	filedId := getFieldId(task.Language, task.Code)
	lg.Infoln("文件id：" + filedId)

	// 获取测试用例URL
	url := getTestExampleUrl(task.ProblemId)
	lg.Infof("测试用例URL: %s", url)

	// 获取测试用例
	testCases, err := fetchTestCases(url)
	if err != nil {
		lg.Errorf("获取测试用例失败: %v", err)
		return
	}

	// 执行判题
	for i := 0; i < len(testCases); i += 2 {
		input := strings.TrimSpace(testCases[i])
		expected := strings.TrimSpace(testCases[i+1])

		// 执行判题
		actual := judge(filedId, task.Language, input)
		actual = strings.TrimSpace(actual)

		lg.Infof("测试用例 %d: 输入=%s 预期=%s 实际=%s",
			i/2+1, input, expected, actual)

		// 结果比对
		if !isEqual(expected, actual) {
			lg.Errorf("判题失败！失败用例：输入=%s（预期：%s，实际：%s）",
				input, expected, actual)
			//把submits的evaluation，改为{input： ，expected:   ,actual:  }表示有测试用例不通过判题
			evaluation := "input:" + input + ",expected:" + expected + ",actual:" + actual
			//selector是筛选条件，update是要更新的内容
			selector := bson.M{
				"_id": task.ID,
			}
			update := bson.M{"$set": bson.M{"status": utils.BadJudge,
				"evaluation": evaluation,
			}}
			_, err = dao.NewDao().UpdateSubmit(context.Background(), selector, update)
			if err != nil {
				lg.Info(utils.UpdateErr, err)
				return
			}
			return
		}
	}

	lg.Info("所有测试用例通过，判题成功！")
	//把mongo的submits，status改为utils.SuccessJudge
	//selector是筛选条件，update是要更新的内容
	selector := bson.M{
		"_id": task.ID,
	}
	update := bson.M{"$set": bson.M{"status": utils.SuccessJudge}}
	_, err = dao.NewDao().UpdateSubmit(context.Background(), selector, update)
	if err != nil {
		lg.Info(utils.UpdateErr, err)
		return
	}

}

func (s *Service) InvokeSandbox(task *models.LocalTask) {
	invokeSandbox(task)
}
