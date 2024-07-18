package main

import (
	"fmt"
	"os"

	"em-test/config"
	"em-test/controllers"
	"em-test/validators"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	_ "em-test/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title User Management API
// @version 1.0
// @description This is a sample server for managing users.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	config.InitDB()

	validate := validator.New()

	validate.RegisterValidation("passport_number_format", validators.ValidatePassportNumberFormat)

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

	router.Run(":8080")
}