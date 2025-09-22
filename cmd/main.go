/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 09:49
@Name: main.go
@Description:
*/

package main

import (
	"fmt"
	"log"
	"os"
	"zhku-oj-server/pkg/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bio-dragon",
	Short: "bio-dragon",
	Long:  "bio-dragon",
}



func main() {
	// 获取配置目录路径，支持从cmd目录或项目根目录运行
	configDir := "conf"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// 如果当前目录没有conf，尝试上级目录
		configDir = "../conf"
	}

	enforcer, err := utils.InitCasbin(configDir)
	if err != nil {
		log.Fatalf("初始化Casbin失败: %v", err)
	}
	// 确保 enforcer 不为 nil
	if enforcer != nil {
		// 重新加载策略
		err = utils.ReloadCasbinPolicy()
		if err != nil {
			log.Fatalf("重新加载Casbin策略失败: %v", err)
		}
		log.Println("Casbin策略加载成功")
	} else {
		log.Fatalf("Casbin enforcer 初始化失败")
	}

	// 初始化AI服务
	_, err = utils.InitAI(configDir)
	if err != nil {
		log.Printf("AI服务初始化失败: %v", err)
		log.Println("系统将在没有AI功能的情况下继续运行")
	} else {
		log.Println("AI服务初始化成功")
		utils.LogAIStatus() // 打印AI服务状态
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}


