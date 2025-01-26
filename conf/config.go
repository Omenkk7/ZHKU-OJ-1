/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 17:38
@Name: config.go
@Description:
*/

package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

var (
	AppMode      string        // 服务器启动模式，默认debug模式
	AppEnv       string        // 服务器部署环境，默认为开发环境模式
	Port         string        // 服务启动端口
	GracefulTime time.Duration // 服务优雅重启的时间，单位为秒，默认5s
	Config       *configYaml
	v            *viper.Viper
)

type configYaml struct {
	App struct {
		Mode         string
		Host         string
		Port         string
		Env          string
		GracefulTime time.Duration
	}
	Mongo struct {
		Uri string
	}
}

func Init() {
	// viper mapstructure 配置读取
	v = viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("../../../conf")
	v.AddConfigPath("../../conf")
	v.AddConfigPath("../conf")
	v.AddConfigPath("./conf")
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	v.AutomaticEnv()
	err := v.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(any(fmt.Errorf("fatal error config file: %s", err)))
	}
	v.AutomaticEnv()
	c := &configYaml{}
	// viper 配置解析成struct
	err = v.Unmarshal(&c)
	if err != nil {
		panic(any(fmt.Errorf("fatal error unmarshal config: %s", err)))
	}
	Config = c
	if Config.App.Env == "test" {
		fmt.Println("config: ", c)
	}
}

func ConfigUtils() *viper.Viper {
	return v
}
