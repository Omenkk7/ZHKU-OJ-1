package utils

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Config 结构体表示整个配置文件的结构
type Config struct {
	App     AppConfig     `yaml:"App"`
	Mongo   MongoConfig   `yaml:"Mongo"`
	Judge   JudgeConfig   `yaml:"Judge"`
	Sandbox SandboxConfig `yaml:"Sandbox"`
	Dify    DifyConfig    `yaml:"Dify"` // 新增Dify配置
}

// AppConfig 结构体表示App部分的配置
type AppConfig struct {
	Host string `yaml:"Host"`
	Port string `yaml:"Port"`
}

// MongoConfig 结构体表示Mongo部分的配置
type MongoConfig struct {
	Uri    string `yaml:"Uri"`
	DbName string `yaml:"DbName"`
}

// JudgeConfig 结构体表示Judge部分的配置
type JudgeConfig struct {
	Type string               `yaml:"Type"`
	Java map[string]JDKConfig `yaml:"Java"`
}

// JDKConfig 结构体表示JDK的配置
type JDKConfig struct {
	BaseArgs string `yaml:"BaseArgs"`
	Env      string `yaml:"Env"`
}

// SandboxConfig 结构体表示Sandbox部分的配置
type SandboxConfig struct {
	Url    string `yaml:"Url"`
	Method string `yaml:"Method"`
}

// DifyConfig 结构体表示Dify API的配置
type DifyConfig struct {
	APIKey string `yaml:"APIKey"`
	APIURL string `yaml:"APIURL"`
}

// LoadConfig 从指定的YAML文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	// 读取YAML文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解析YAML文件
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// GetAppConfig 返回App部分的配置
func (c *Config) GetAppConfig() AppConfig {
	return c.App
}

// GetMongoConfig 返回Mongo部分的配置
func (c *Config) GetMongoConfig() MongoConfig {
	return c.Mongo
}

// GetJudgeConfig 返回Judge部分的配置
func (c *Config) GetJudgeConfig() JudgeConfig {
	return c.Judge
}

// GetSandboxConfig 返回Sandbox部分的配置
func (c *Config) GetSandboxConfig() SandboxConfig {
	return c.Sandbox
}

// GetDifyConfig 返回Dify API的配置
func (c *Config) GetDifyConfig() DifyConfig {
	return c.Dify
}
