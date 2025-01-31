/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 13:09
@Name: server.go
@Description:
*/

package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"zhku-oj-server/pkg/app/api-server/service"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
	"zhku-oj-server/pkg/utils/middleware"
)

type Server struct {
	lg     logrus.FieldLogger
	app    *gin.Engine
	svc    *service.Service
	opts   *CmdOptions
	res    *models.Result
	stopCh <-chan struct{}
}

func NewServer(lg logrus.FieldLogger, svc *service.Service, opts *CmdOptions, stopCh <-chan struct{}) *Server {
	app := gin.Default()
	app.Use(middleware.CorsHandler()) // set cors
	app.Use(middleware.LoggerHandler(utils.GetLogger(context.Background()), true))
	app.Use(gin.Recovery()) // panic recovery
	gin.SetMode(gin.DebugMode)
	return &Server{
		lg:     lg,
		app:    app,
		svc:    svc,
		opts:   opts,
		stopCh: stopCh,
	}
}

func (s *Server) Init() {
	s.RegisterRoutes()
}

func (s *Server) RegisterRoutes() {
	v1 := s.app.Group("/api/v1")
	s.RegisterUser(v1)
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.opts.Host, s.opts.Port),
		Handler: s.app,
	}

	go func() {
		select {
		case <-s.stopCh:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			// Shutdown server
			s.lg.Infoln("Shutdown server error: ", srv.Shutdown(ctx))
			return
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			// 主动关闭server,如接收到控制台退出信号等...
			s.lg.Println("Http Server closed, bye!")
			return err
		}
		s.lg.Errorf("listen: %s", err.Error())
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	s.svc.Close(ctx)
}

type CmdOptions struct {
	Host string
	Port string
}

func NewCmdOptions(host, port string) *CmdOptions {
	return &CmdOptions{
		Host: host,
		Port: port,
	}
}
