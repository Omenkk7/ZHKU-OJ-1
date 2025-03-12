package service

import (
	"bytes"
	"encoding/json"
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
					CPULimit:    9687500000,
					MemoryLimit: 10485760000,
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
	return response[0].FileIds["a.java"]
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

func invokeSandbox(language string, code string) {
	lg := utils.GetDefaultLogger()

	//获取编译后的文件id
	filedId := getFieldId(language, code)
	lg.Infoln("文件id：" + filedId)

	//TODO 把测试用例丢进去example判题
	//"1 1"用于测试两数之和，模拟一个测试用例; 可以修改“1 1”进行各种测试
	result := judge(filedId, language, "111 123")
	lg.Infof("结果为：%s", result)

}

func (s *Service) InvokeSandbox(language string, code string) {
	invokeSandbox(language, code)
}
