/*
@Author: urmsone urmsone@163.com
@Date: 2025/2/6 00:21
@Name: proxy.go
@Description:
*/

package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type GoJudgeFile struct {
	Stderr string `json:"stderr"`
	Stdout string `json:"stdout"`
}

type GoJudgeRunResponse struct {
	Status     string `json:"status"`
	ExitStatus int    `json:"exit_status"`
	Time       int64  `json:"time"`
	Memory     int64  `json:"memory"`
	RunTime    int64  `json:"runTime"`

	Files []GoJudgeFile `json:"files"`
}

func NewGoJudgeRunResponse() *GoJudgeRunResponse {
	return &GoJudgeRunResponse{
		Status:     "",
		ExitStatus: -1,
		Time:       -1,
		Memory:     -1,
		RunTime:    -1,
		Files:      []GoJudgeFile{},
	}

}

type GoJudgeProxy struct {
	Lg      logrus.FieldLogger
	Headers map[string]string
	Client  *http.Client
}

func (gp *GoJudgeProxy) Do(method, url string, headers map[string]string, body io.Reader) (*GoJudgeRunResponse, error) {
	gp.Lg.Println("GoJudgeProxy request url: ", url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	if gp.Headers != nil {
		for k, v := range gp.Headers {
			req.Header.Set(k, v)
		}
	}
	gp.Lg.Println("GoJudgeProxy request headers: ", gp.Headers)
	resp, err := gp.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := NewGoJudgeRunResponse()
	gp.Lg.Infoln()
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	gp.Lg.Println("GoJudgeProxy response result: ", result)
	return result, nil
}

func NewGoJudgeProxy(lg logrus.FieldLogger, project, group string, headers map[string]string) (*GoJudgeProxy, error) {
	client := &http.Client{}
	if headers == nil {
		headers = map[string]string{}
	}

	return &GoJudgeProxy{
		Lg:      lg,
		Headers: headers,
		Client:  client,
	}, nil
}

func init() {
}
