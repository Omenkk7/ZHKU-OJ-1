package service

import (
	"bytes"
	"encoding/json"
	"zhku-oj-server/pkg/utils"

	"fmt"
	"io/ioutil"
	"net/http"
)

type Cmd struct {
	Args        []string        `json:"args"`
	Env         []string        `json:"env"`
	Files       []File          `json:"files"`
	CPULimit    int64           `json:"cpuLimit"`
	MemoryLimit int64           `json:"memoryLimit"`
	ProcLimit   int             `json:"procLimit"`
	CopyIn      map[string]File `json:"copyIn"`
	CopyOut     []string        `json:"copyOut"`
}

type File struct {
	Content string `json:"content,omitempty"`
	Name    string `json:"name,omitempty"`
	Max     int    `json:"max,omitempty"`
}

type RequestBody struct {
	Cmd []Cmd `json:"cmd"`
}

// 定义响应体的结构
type ResponseBody struct {
	Status     string            `json:"status"`
	ExitStatus int               `json:"exitStatus"`
	Time       int64             `json:"time"`
	Memory     int64             `json:"memory"`
	RunTime    int64             `json:"runTime"`
	Files      map[string]string `json:"files"`
}

type CmdResponse struct {
	Status     string            `json:"status"`
	ExitStatus int               `json:"exitStatus"`
	Time       int64             `json:"time"`
	Memory     int64             `json:"memory"`
	RunTime    int64             `json:"runTime"`
	Files      map[string]string `json:"files"`
}

func (s *Service) InvokeSandbox(language string, code string) {
	//通过读取配置文件的形式，找到调用沙箱的请求方式和url，找到jdk/g++等环境的路径
	cfg, _ := utils.LoadConfig("conf/config.yaml")
	judgeCfg := cfg.GetJudgeConfig()
	sandboxCfg := cfg.GetSandboxConfig()

	//创建请求体
	var requestBody RequestBody
	switch language {
	case "java":
		jdk11cfg, _ := judgeCfg.Java["JDK11"]
		requestBody = RequestBody{
			Cmd: []Cmd{
				{
					Args: []string{jdk11cfg.BaseArgs[0]},
					Env:  []string{jdk11cfg.Env},
					//TODO Content怎么用，如何进行判题？
					Files:       []File{{Content: "1 3"}, {Name: "stdout", Max: 10240000}, {Name: "stderr", Max: 10240000}},
					CPULimit:    9687500000,
					MemoryLimit: 10485760000,
					ProcLimit:   50,
					CopyIn: map[string]File{
						"a.java": {
							Content: code,
						},
					},
					CopyOut: []string{"stdout", "stderr"},
				},
			},
		}

	}

	// 将请求体编码为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 创建HTTP请求
	req, err := http.NewRequest(sandboxCfg.Method, sandboxCfg.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 解析响应体
	var response []ResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// 打印响应
	// 提取并打印 stdout 的值
	if len(response) > 0 {
		stdoutValue := response[0].Files["stdout"]
		fmt.Printf("stdout: %s\n", stdoutValue)
	} else {
		fmt.Println("No response data found")
	}
}
