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
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "bio-dragon",
	Short: "bio-dragon",
	Long:  "bio-dragon",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
