package router

import (
	"fmt"
	"os"

	"em-test/controllers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(validate *validator.Validate) *gin.Engine {
	router := gin.Default()

	serviceAddress := os.Getenv("SERVICE_ADDRESS")

	router.GET("/users", controllers.GetUsersHandler)
	router.GET("/users/:id", controllers.GetUserHandler)
	router.POST("/users", func(c *gin.Context) {
		controllers.CreateUserHandler(c, validate)
	})
	router.PUT("/users/:id", func(c *gin.Context) {
		controllers.UpdateUserHandler(c, validate)
	})
	router.DELETE("/users/:id", controllers.DeleteUserHandler)

	router.GET("/tasks", controllers.GetTasksHandler)
	router.GET("/tasks/:id", controllers.GetTaskHandler)
	router.POST("/tasks", func(c *gin.Context) {
		controllers.CreateTaskHandler(c, validate)
	})

	router.GET("/tasklogs", controllers.GetTaskLogsHandler)
	router.GET("/tasklogs/:id", controllers.GetTaskLogHandler)
	router.POST("/tasklogs", func(c *gin.Context) {
		controllers.CreateAndStartTaskLog(c, validate)
	})
	router.PUT("/tasklogs/:id/complete", controllers.CompleteTaskLogHandler)

	router.GET("/tasktimes", controllers.GetUserTaskTimes)

	swaggerAddress := fmt.Sprintf("%s/swagger/doc.json", serviceAddress)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(swaggerAddress)))

	return router
}
