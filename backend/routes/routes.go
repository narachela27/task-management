package routes

import (
	"database/sql"
	"net/http"

	"task-management-backend/internal/handler"
	"task-management-backend/internal/middleware"
	"task-management-backend/internal/repository"
	"task-management-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Task Management API V2 is running",
		})
	})

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	authService := service.NewAuthService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	taskService := service.NewTaskService(taskRepo, categoryRepo)
	adminService := service.NewAdminService(userRepo, taskRepo)

	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	taskHandler := handler.NewTaskHandler(taskService)
	adminHandler := handler.NewAdminHandler(adminService)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthRequired(userRepo), authHandler.Me)
		}

		categories := api.Group("/categories")
		categories.Use(middleware.AuthRequired(userRepo))
		{
			categories.GET("", categoryHandler.GetAll)
			categories.POST("", categoryHandler.Create)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		tasks := api.Group("/tasks")
		tasks.Use(middleware.AuthRequired(userRepo))
		{
			tasks.GET("", taskHandler.GetAll)
			tasks.GET("/summary", taskHandler.Summary)
			tasks.GET("/:id", taskHandler.GetByID)
			tasks.POST("", taskHandler.Create)
			tasks.PUT("/:id", taskHandler.Update)
			tasks.PATCH("/:id/status", taskHandler.UpdateStatus)
			tasks.DELETE("/:id", taskHandler.Delete)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired(userRepo), middleware.RequireRoles("ADMIN"))
		{
			admin.GET("/users", adminHandler.GetAllUsers)
			admin.GET("/users/:id", adminHandler.GetUserByID)
			admin.PATCH("/users/:id/role", adminHandler.UpdateUserRole)
			admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)
			admin.GET("/tasks", adminHandler.GetAllTasks)
		}
	}

	return r
}