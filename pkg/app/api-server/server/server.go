/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 13:09
@Name: server.go
@Description:
*/

package server

import (
	"context"
	"errors"
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
	securedGroup := userGroup.Group("/").Use(middleware.JWTMiddleware())
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
	securedGroup := userGroup.Group("/").Use(middleware.JWTMiddleware())
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
	securedGroup := userGroup.Group("/").Use(middleware.JWTMiddleware())
	{
		securedGroup.PUT("/:id", s.PutProblem)       //改一个
		securedGroup.DELETE("/:id", s.DeleteProblem) //删一个
		securedGroup.POST("/", s.PostProblem)        //增一个
	}
}

// RegisterSubmit 路由器
func (s *Server) RegisterSubmit(g *gin.RouterGroup) {
	//操作不拦截
	userGroup := g.Group("/submit")
	{
		userGroup.POST("/", s.Submit)
	}
}

func (s *Server) RegisterClass(g *gin.RouterGroup) {
	classGroup := g.Group("/class")
	{
		// 班级管理
		classGroup.POST("", middleware.JWTMiddleware(), s.createClass)
		classGroup.PUT("/:id", middleware.JWTMiddleware(), s.updateClass)
		classGroup.DELETE("/:id", middleware.JWTMiddleware(), s.deleteClass)
		classGroup.PUT("/:id/archive", middleware.JWTMiddleware(), s.archiveClass)
		classGroup.GET("/:id", middleware.JWTMiddleware(), s.getClassByID)
		classGroup.GET("/code/:code", middleware.JWTMiddleware(), s.getClassByCode)
		classGroup.GET("", middleware.JWTMiddleware(), s.getClassList)

		// 班级学生管理
		classGroup.POST("/:id/student", middleware.JWTMiddleware(), s.addStudentToClass)
		classGroup.POST("/:id/students", middleware.JWTMiddleware(), s.batchAddStudentsToClass)
		classGroup.DELETE("/:id/student/:studentId", middleware.JWTMiddleware(), s.removeStudentFromClass)
		classGroup.GET("/:id/students", middleware.JWTMiddleware(), s.getClassStudents)

		// 班级课程管理
		classGroup.POST("/:id/course", middleware.JWTMiddleware(), s.addCourseToClass)
		classGroup.DELETE("/:id/course/:courseId", middleware.JWTMiddleware(), s.removeCourseFromClass)
		classGroup.PUT("/:id/course/:courseId/status", middleware.JWTMiddleware(), s.updateCourseStatusInClass)
		classGroup.GET("/course/:courseId", middleware.JWTMiddleware(), s.getClassesByCourseID)

		// 加入申请管理
		classGroup.POST("/join", middleware.JWTMiddleware(), s.createJoinRequest)
		classGroup.PUT("/join/:id/review", middleware.JWTMiddleware(), s.reviewJoinRequest)
		classGroup.GET("/join", middleware.JWTMiddleware(), s.getJoinRequestList)

		// 学生班级查询
		classGroup.GET("/student/:studentId", middleware.JWTMiddleware(), s.getStudentClasses)

		// 班级成员管理
		classGroup.GET("/:id/members", middleware.JWTMiddleware(), s.getClassMembers)
		classGroup.POST("/:id/members", middleware.JWTMiddleware(), s.addClassMember)
		classGroup.DELETE("/:id/members", middleware.JWTMiddleware(), s.removeClassMember)
	}
}

func (s *Server) RegisterCourse(g *gin.RouterGroup) {
	courseGroup := g.Group("/course")
	{
		// 课程基础管理
		courseGroup.POST("", middleware.JWTMiddleware(), s.createCourse)
		courseGroup.PUT("/:id", middleware.JWTMiddleware(), s.updateCourse)
		courseGroup.DELETE("/:id", middleware.JWTMiddleware(), s.deleteCourse)
		courseGroup.PUT("/:id/archive", middleware.JWTMiddleware(), s.archiveCourse)
		courseGroup.GET("/:id", middleware.JWTMiddleware(), s.getCourse)
		courseGroup.GET("/code", middleware.JWTMiddleware(), s.getCourseByCode)
		courseGroup.GET("", middleware.JWTMiddleware(), s.getCourseList)

		// 课程成员管理
		courseGroup.POST("/:id/members", middleware.JWTMiddleware(), s.addCourseMember)
		courseGroup.DELETE("/:id/members", middleware.JWTMiddleware(), s.removeCourseMember)
		courseGroup.GET("/:id/members", middleware.JWTMiddleware(), s.getCourseMembers)

		// 加入申请管理
		courseGroup.POST("/join", s.createJoinCourseRequest)
		courseGroup.PUT("/join/:id", middleware.JWTMiddleware(), s.reviewJoinCourseRequest)
		courseGroup.GET("/:id/join-requests", middleware.JWTMiddleware(), s.getJoinCourseRequestList)
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
	s.RegisterSubmit(v1)
	s.RegisterClass(v1)
	s.RegisterCourse(v1)
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
		if errors.Is(err, http.ErrServerClosed) {
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
