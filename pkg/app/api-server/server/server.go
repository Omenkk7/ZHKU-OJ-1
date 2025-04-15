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
	/*classGroup := g.Group("/class").Use(middleware.JWTMiddleware())
	{
		// 班级管理
		classGroup.POST("", s.createClass)
		classGroup.PUT("/:id", s.updateClass)
		classGroup.DELETE("/:id", s.deleteClass)
		classGroup.PUT("/:id/archive", s.archiveClass)
		classGroup.GET("/:id", s.getClassByID)
		classGroup.GET("/code/:code", s.getClassByCode)
		classGroup.GET("", s.getClassList)

		// 班级学生管理
		classGroup.POST("/:id/student", s.addStudentToClass)
		classGroup.POST("/:id/students", s.batchAddStudentsToClass)
		classGroup.DELETE("/:id/student/:studentId", s.removeStudentFromClass)
		classGroup.GET("/:id/students", s.getClassStudents)

		// 班级课程管理
		classGroup.POST("/:id/course", s.addCourseToClass)
		classGroup.DELETE("/:id/course/:courseId", s.removeCourseFromClass)
		classGroup.PUT("/:id/course/:courseId/status", s.updateCourseStatusInClass)
		classGroup.GET("/course/:courseId", s.getClassesByCourseID)

		// 加入申请管理
		classGroup.POST("/join", s.createJoinRequest)
		classGroup.PUT("/join/:id/review", s.reviewJoinRequest)
		classGroup.GET("/join", s.getJoinRequestList)

		// 学生班级查询
		classGroup.GET("/student/:studentId", s.getStudentClasses)

		// 班级成员管理
		classGroup.GET("/:id/members", s.getClassMembers)
		classGroup.POST("/:id/members", s.addClassMember)
		classGroup.DELETE("/:id/members", s.removeClassMember)
	}*/
	classGroup := g.Group("/class").Use(middleware.JWTMiddleware())
	{
		// 班级管理
		classGroup.POST("", utils.CasbinMiddleware(utils.ObjClass, utils.ActCreate), s.createClass)
		classGroup.PUT("/:id", s.updateClass)          // 使用自定义权限检查
		classGroup.DELETE("/:id", s.deleteClass)       // 使用自定义权限检查
		classGroup.PUT("/:id/archive", s.archiveClass) // 使用自定义权限检查
		classGroup.GET("/:id", utils.CasbinMiddleware(utils.ObjClass, utils.ActRead), s.getClassByID)
		classGroup.GET("/code/:code", utils.CasbinMiddleware(utils.ObjClass, utils.ActRead), s.getClassByCode)
		classGroup.GET("", utils.CasbinMiddleware(utils.ObjClass, utils.ActRead), s.getClassList)

		// 班级学生管理
		classGroup.POST("/:id/student", s.addStudentToClass)                   // 使用自定义权限检查
		classGroup.POST("/:id/students", s.batchAddStudentsToClass)            // 使用自定义权限检查
		classGroup.DELETE("/:id/student/:studentId", s.removeStudentFromClass) // 使用自定义权限检查
		classGroup.GET("/:id/students", utils.CasbinMiddleware(utils.ObjClassStudent, utils.ActRead), s.getClassStudents)

		// 班级课程管理
		classGroup.POST("/:id/course", utils.CasbinMiddleware(utils.ObjClassCourse, utils.ActUpdate), s.addCourseToClass)
		classGroup.DELETE("/:id/course/:courseId", s.removeCourseFromClass)         // 使用自定义权限检查
		classGroup.PUT("/:id/course/:courseId/status", s.updateCourseStatusInClass) // 使用自定义权限检查
		classGroup.GET("/course/:courseId", utils.CasbinMiddleware(utils.ObjClassCourse, utils.ActRead), s.getClassesByCourseID)

		// 加入申请管理
		classGroup.POST("/join", s.createJoinRequest)
		classGroup.PUT("/join/:id/review", s.reviewJoinRequest) // 使用自定义权限检查
		classGroup.GET("/join", utils.CasbinMiddleware(utils.ObjJoinRequest, utils.ActRead), s.getJoinRequestList)

		// 学生班级查询
		classGroup.GET("/student/:studentId", utils.CasbinMiddleware(utils.ObjClass, utils.ActRead), s.getStudentClasses)
	}
}

func (s *Server) RegisterCourse(g *gin.RouterGroup) {
	courseGroup := g.Group("/courses").Use(middleware.JWTMiddleware())
	{
		// 课程基础管理
		courseGroup.POST("", s.createCourse)
		courseGroup.PUT("/:id", s.updateCourse)
		courseGroup.DELETE("/:id", s.deleteCourse)
		courseGroup.PUT("/:id/archive", s.archiveCourse)
		courseGroup.GET("/:id", s.getCourse)
		courseGroup.GET("/code", s.getCourseByCode)
		courseGroup.GET("", s.getCourseList)

		// 课程成员管理
		courseGroup.POST("/:id/members", s.addCourseMember)
		courseGroup.DELETE("/:id/members", s.removeCourseMember)
		courseGroup.GET("/:id/members", s.getCourseMembers)

		// 加入申请管理
		courseGroup.POST("/join", s.createJoinCourseRequest)
		courseGroup.PUT("/join/:id", s.reviewJoinCourseRequest)
		courseGroup.GET("/:id/join-requests", s.getJoinCourseRequestList)
	}
}

func (s *Server) RegisterAssignment(r *gin.RouterGroup) {
	assignment := r.Group("/assignment").Use(middleware.JWTMiddleware())
	{
		assignment.POST("", s.CreateAssignment)             // 创建作业
		assignment.PUT("/:id", s.UpdateAssignment)          // 更新作业
		assignment.DELETE("/:id", s.DeleteAssignment)       // 删除作业
		assignment.PUT("/:id/archive", s.ArchiveAssignment) // 归档作业
		assignment.GET("/:id", s.GetAssignmentDetail)       // 获取作业详情
		assignment.GET("", s.GetAssignmentList)             // 获取作业列表

		/*学生提交部分*/

		//assignment.GET("/student", s.GetStudentAssignments) // 获取学生作业列表（未完成）
		assignment.POST("/:id/submit", s.SubmitAssignment) // 提交作业
		//assignment.GET("/submission", s.GetStudentSubmission) // 获取学生提交信息（未完成）

		// 教师、管理员操作
		assignment.POST("/grade", s.GradeAssignment)            // 批改作业
		assignment.POST("/reject", s.RejectAssignment)          // 打回作业
		assignment.GET("/export/:id", s.ExportAssignmentGrades) // 导出成绩
		assignment.POST("/students", s.AddStudentsToAssignment) // 添加学生到作业
	}
}

// 竞赛模块路由
func (s *Server) RegisterContest(r *gin.RouterGroup) {
	contestGroup := r.Group("/contest").Use(middleware.JWTMiddleware())
	{
		// 创建竞赛
		contestGroup.POST("/create", s.CreateContest)
		// 更新竞赛
		contestGroup.PUT("/:id", s.UpdateContest)
		// 删除竞赛
		contestGroup.DELETE("/:id", s.DeleteContest)
		// 归档竞赛
		contestGroup.PUT("/:id/archive", s.ArchiveContest)
		// 更新竞赛状态
		contestGroup.PUT("/:id/status", s.UpdateContestStatus)

		// 添加参赛者
		contestGroup.POST("/participant/add", s.AddParticipant)
		// 批量添加参赛者
		contestGroup.POST("/participant/batch-add", s.BatchAddParticipants)
		// 移除参赛者
		contestGroup.DELETE("/participant/remove", s.RemoveParticipant)
		// 审核参赛者
		contestGroup.PUT("/participant/audit", s.AuditParticipant)
		// 导出竞赛成绩
		contestGroup.GET("/:id/export", s.ExportContestScore) //TODO 未完成

		//申请加入竞赛
		contestGroup.POST("/apply", s.ApplyJoinContest)

		// 获取竞赛列表
		contestGroup.GET("/list", s.GetContestList)
		// 获取竞赛详情
		contestGroup.GET("/:id", s.GetContestDetail)
		// 获取参赛者列表
		contestGroup.GET("/participant/list", s.GetParticipantList)
		// 获取竞赛排名
		contestGroup.GET("/:id/ranking", s.GetContestRanking)
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
	s.RegisterAssignment(v1)
	s.RegisterContest(v1)
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
