package controllers

import (
	"em-test/config"
	"em-test/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Получение списка всех задач
// @Summary Get all tasks
// @Description Get a list of all tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks [get]
func GetTasksHandler(c *gin.Context) {
	var tasks []models.Task
	result := config.DB.Find(&tasks)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// Получение задачи по id
// @Summary Get task by ID
// @Description Get a single task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id} [get]
func GetTaskHandler(c *gin.Context) {
	id := c.Param("id")
	var task models.Task 

	result := config.DB.First(&task, id)

	if result.Error != nil {
    if result.Error == gorm.ErrRecordNotFound {
        c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
    } else {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
    }
    return
  }

	c.JSON(http.StatusOK, task)
}

// Создание новой задачи
// @Summary Create a new Task
// @Description Create a new user with the input payload
// @Tags tasks
// @Accept json
// @Produce json
// @Param user body models.Task true "Task JSON"
// @Success 201 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks [post]
func CreateTaskHandler(c *gin.Context, validate *validator.Validate) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return 
	}

	if err := validate.Struct(&task); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make([]string, len(validationErrors))

		for i, fieldError := range validationErrors {
			errors[i] = fieldError.Error()
		}

		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}

	result := config.DB.Create(&task)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}