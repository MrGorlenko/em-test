package controllers

import (
	"em-test/config"
	"em-test/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Получение всех TaskLogs
// @Summary Get all task logs
// @Description Get a list of all task logs
// @Tags tasklogs
// @Accept json
// @Produce json
// @Success 200 {array} models.TaskLog
// @Failure 500 {object} models.ErrorResponse
// @Router /tasklogs [get]
func GetTaskLogsHandler(c *gin.Context) {
	var tasksLogs []models.TaskLog
	result := config.DB.Find(&tasksLogs)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, tasksLogs)
}

// Получение TaskLog по id
// @Summary Get task log by ID
// @Description Get a single task log by its ID
// @Tags tasklogs
// @Accept json
// @Produce json
// @Param id path int true "Task Log ID"
// @Success 200 {object} models.TaskLog
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasklogs/{id} [get]
func GetTaskLogHandler(c *gin.Context) {
	id := c.Param("id")
	var taskLog models.TaskLog

	result := config.DB.First(&taskLog, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task log not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, taskLog)
}

// Создание нового TaskLog и установка StartTime
// @Summary Create a new task log
// @Description Create a new task log with the input payload and set the start time
// @Tags tasklogs
// @Accept json
// @Produce json
// @Param tasklog body models.TaskLog true "Task Log JSON"
// @Success 201 {object} models.TaskLog
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasklogs [post]
func CreateAndStartTaskLog(c *gin.Context, validate *validator.Validate) {
	var taskLog models.TaskLog

	if err := c.ShouldBindJSON(&taskLog); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if err := validate.Struct(&taskLog); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make([]string, len(validationErrors))

		for i, fieldError := range validationErrors {
			errors[i] = fieldError.Error()
		}

		c.JSON(http.StatusBadRequest, gin.H{"validation_error": errors})
		return
	}

	taskLog.StartTime = time.Now()
	result := config.DB.Create(&taskLog)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, taskLog)
}

// Завершение TaskLog
// @Summary Complete a task log
// @Description Set the end time for a task log
// @Tags tasklogs
// @Accept json
// @Produce json
// @Param id path int true "Task Log ID"
// @Success 200 {object} models.TaskLog
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasklogs/{id}/complete [put]
func CompleteTaskLogHandler(c *gin.Context) {
	id := c.Param("id")
	var taskLog models.TaskLog
	result := config.DB.First(&taskLog, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task log not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		}
		return
	}

	taskLog.EndTime = time.Now()
	if err := config.DB.Save(&taskLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, taskLog)
}