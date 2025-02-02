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
	"zhku-oj-server/pkg/utils"
	"zhku-oj-server/pkg/utils/middleware"
)

type Server struct {
	lg     logrus.FieldLogger
	app    *gin.Engine
	svc    *service.Service
	opts   *CmdOptions
	res    *utils.Result
	stopCh <-chan struct{}
}

// RegisterUser 路由器 ——————user_manager
func (s *Server) RegisterUser(g *gin.RouterGroup) {
	userGroup := g.Group("/user")
	{
		userGroup.POST("/", s.PostUser)   // 注册
		userGroup.POST("/login", s.Login) // 登录
	}

	// 对以下路由应用JWT拦截器，并且只有管理员可以执行以下操作
	securedGroup := userGroup.Group("/").Use(middleware.JWTInterceptor())
	{
		securedGroup.GET("/", s.GetSomeUser)      //查一堆
		securedGroup.GET("/:id", s.GetOneUser)    //查一个
		userGroup.PUT("/:id", s.PutUser)          //改一个
		securedGroup.DELETE("/:id", s.DeleteUser) //删一个
	}
}

// RegisterLabel 路由器 ——————label_manager
func (s *Server) RegisterLabel(g *gin.RouterGroup) {
	//查询操作不拦截
	userGroup := g.Group("/label")
	{
		userGroup.GET("/", s.GetSomeLabel)   //查一堆
		userGroup.GET("/:id", s.GetOneLabel) //查一个
	}
	// 对以下路由应用JWT拦截器，并且只有管理员可以执行以下操作
	securedGroup := userGroup.Group("/").Use(middleware.JWTInterceptor())
	{
		securedGroup.PUT("/:id", s.PutLabel)       //改一个
		securedGroup.DELETE("/:id", s.DeleteLabel) //删一个
		securedGroup.POST("/", s.PostLabel)        //增一个
	}
}

// RegisterProblem 路由器 ——————label_problem
func (s *Server) RegisterProblem(g *gin.RouterGroup) {
	//查询操作不拦截
	//TODO 只有管理员才可以查到私密题库  或扩展业务，充VIP得到付费资源
	userGroup := g.Group("/problem")
	{
		userGroup.GET("/", s.GetSomeProblem)   //查一堆
		userGroup.GET("/:id", s.GetOneProblem) //查一个
	}
	// 对以下路由应用JWT拦截器，并且只有管理员可以执行以下操作
	securedGroup := userGroup.Group("/").Use(middleware.JWTInterceptor())
	{
		securedGroup.PUT("/:id", s.PutProblem)       //改一个
		securedGroup.DELETE("/:id", s.DeleteProblem) //删一个
		securedGroup.POST("/", s.PostProblem)        //增一个
	}
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

// RegisterRoutes 注册路由
func (s *Server) RegisterRoutes() {
	v1 := s.app.Group("/api/v1")
	s.RegisterUser(v1) //调用middleware的路由组
	s.RegisterLabel(v1)
	s.RegisterProblem(v1)
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
