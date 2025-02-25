/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 09:54
@Name: runServer.go
@Description:
*/

package main

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"zhku-oj-server/conf"
	"zhku-oj-server/pkg/app/api-server/server"
	"zhku-oj-server/pkg/app/api-server/service"
	"zhku-oj-server/pkg/utils"
)

var runServerCfg = struct {
	serverHost string
	serverPort string
	logDir     string
	configDir  string
}{}

func init() {
	rootCmd.AddCommand(runServerCmd)
}

var runServerCmd = &cobra.Command{
	Use:   "run",
	Short: "run server",
	Long:  "Run biaDragon server, listening host:port and providing service",
	RunE:  runServe,
	// Hook before and after Run initialize and do something, respectively
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	//BashCompletionFunction: ``
}

func runServe(cmd *cobra.Command, args []string) error {
	//manager := routers.NewRouter()
	//userRouter := routers.NewUserRouter("user", &manager.RouterGroup)
	//manager.AddRouter("user", userRouter)
	//manager.Register()
	//manager.RunManager()
	if err := cmd.Flags().Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	// 初始化logger

	lg := utils.GetLogger(ctx)
	if lg == nil {
		// TODO: 实现自定义error
		return nil
	}
	// 加载配置文件
	conf.Init()
	conf.ConfigUtils().OnConfigChange(func(in fsnotify.Event) {
		// TODO: 配置文件热加载,服务优雅重启
		lg.Println("configFile changed ", in.Name, in.Op)
	})
	conf.ConfigUtils().WatchConfig()
	stopCh := make(chan struct{})
	signalCh := make(chan os.Signal, 1)
	//signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(signalCh, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	g := errgroup.Group{}
	svc := service.NewService()
	opts := server.NewCmdOptions(conf.Config.App.Host, conf.Config.App.Port)
	s := server.NewServer(lg, svc, opts, stopCh)

	//启动server
	g.Go(func() error {
		s.Init()
		return s.Run()
	})

	//启动server_judge
	j := server.NewJudge(utils.LocalJudge)
	g.Go(func() error {
		return j.RunJudge()
	})

	g.Go(func() error {
		// 监听信号,控制server的退出
		return waitForSignal(lg, signalCh, stopCh)

	})
	defer s.Shutdown(ctx)
	return g.Wait()
}

func waitForSignal(lg logrus.FieldLogger, signalCh <-chan os.Signal, stopCh chan struct{}) error {
	for {
		select {
		case s := <-signalCh:
			lg.Infof("Captured %s signal. Exiting ...\n", s)
			close(stopCh)
			lg.Infof("Server stopping ...")
			return nil
		}
	}
}
