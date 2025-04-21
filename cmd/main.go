/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 09:49
@Name: main.go
@Description:
*/

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"zhku-oj-server/pkg/utils"
)

var rootCmd = &cobra.Command{
	Use:   "bio-dragon",
	Short: "bio-dragon",
	Long:  "bio-dragon",
}

func main() {

	enforcer, err := utils.InitCasbin("conf")
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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
