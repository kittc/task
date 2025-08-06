package main

import (
	"log"
	"os"

	"task-management-platform/internal/database"
	"task-management-platform/internal/handlers"
	"task-management-platform/internal/middleware"
	"task-management-platform/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func main() {
	// 初始化数据库
	dbPath := getEnv("DB_PATH", "./task_management.db")
	if err := database.InitDatabase(dbPath); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化服务
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-here")
	authService := services.NewAuthService(jwtSecret)
	taskService := services.NewTaskService()
	orgService := services.NewOrganizationService()

	// 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)
	orgHandler := handlers.NewOrganizationHandler(orgService)

	// 设置 Gin 模式
	if getEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := gin.Default()

	// 设置 CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	router.Use(func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		ctx.Next()
	})

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Task Management Platform is running",
		})
	})

	// API 路由组
	api := router.Group("/api/v1")

	// 认证路由（无需认证）
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(authMiddleware.RequireAuth())
	{
		// 用户相关
		protected.GET("/profile", authHandler.GetProfile)
		protected.POST("/change-password", authHandler.ChangePassword)
		protected.POST("/logout", authHandler.Logout)

		// 任务相关
		tasks := protected.Group("/tasks")
		{
			tasks.GET("", taskHandler.GetTasks)
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			
			// 任务成员管理
			tasks.POST("/:id/members", taskHandler.AddTaskMember)
			tasks.DELETE("/:id/members/:member_id", taskHandler.RemoveTaskMember)
			
			// 任务清单管理
			tasks.POST("/:id/checklists", taskHandler.CreateChecklist)
			tasks.PUT("/checklists/:checklist_id", taskHandler.UpdateChecklist)
			tasks.DELETE("/checklists/:checklist_id", taskHandler.DeleteChecklist)
		}

		// 组织架构相关
		orgs := protected.Group("/organizations")
		{
			orgs.GET("", orgHandler.GetOrganizations)
			orgs.POST("", orgHandler.CreateOrganization)
			orgs.GET("/:id", orgHandler.GetOrganization)
			orgs.PUT("/:id", orgHandler.UpdateOrganization)
			orgs.DELETE("/:id", orgHandler.DeleteOrganization)
			orgs.GET("/:id/stats", orgHandler.GetOrganizationStats)
		}

		// 部门相关
		depts := protected.Group("/departments")
		{
			depts.GET("", orgHandler.GetDepartments)
			depts.POST("", orgHandler.CreateDepartment)
			depts.GET("/:id", orgHandler.GetDepartment)
			depts.PUT("/:id", orgHandler.UpdateDepartment)
			depts.DELETE("/:id", orgHandler.DeleteDepartment)
			depts.GET("/:id/stats", orgHandler.GetDepartmentStats)
		}

		// 用户管理相关（需要管理权限）
		users := protected.Group("/users")
		users.Use(authMiddleware.RequireAdminOrManagerOrSupervisor())
		{
			users.GET("", orgHandler.GetUsers)
			users.POST("", orgHandler.CreateUser)
			users.GET("/:id", orgHandler.GetUser)
			users.PUT("/:id", orgHandler.UpdateUser)
			users.DELETE("/:id", orgHandler.DeleteUser)
			users.POST("/:id/transfer", orgHandler.TransferUser)
		}

		// 组织架构概览
		protected.GET("/organization-structure", orgHandler.GetOrganizationStructure)
	}

	// 启动定时任务
	go startScheduledTasks(taskService)

	// 启动服务器
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// startScheduledTasks 启动定时任务
func startScheduledTasks(taskService *services.TaskService) {
	// 这里可以添加定时任务，比如检查超期任务
	// 暂时用简单的轮询实现，生产环境建议使用 cron
	log.Println("Scheduled tasks started")
}