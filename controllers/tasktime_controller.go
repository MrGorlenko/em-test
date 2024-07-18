package controllers

import (
	"em-test/config"
	"em-test/models"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Получение трудозатрат по пользователю за период
// @Summary Get user task times for a period
// @Description Get task times spent by a user for a given period, sorted by time spent in descending order
// @Tags tasktimes
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {array} models.TaskTime
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasktimes [get]
func GetUserTaskTimes(c *gin.Context) {
	userIDStr := c.Query("user_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if userIDStr == "" || startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "user_id, start_date, and end_date are required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user_id format"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid end_date format"})
		return
	}

	var taskLogs []models.TaskLog
	result := config.DB.Where("user_id = ? AND start_time >= ? AND end_time <= ?", userID, startDate, endDate).Find(&taskLogs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	taskTimeMap := make(map[uint]*models.TaskTime)

	for _, log := range taskLogs {
		if log.EndTime.IsZero() {
			continue
		}
		duration := log.EndTime.Sub(log.StartTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60

		if taskTime, exists := taskTimeMap[log.TaskID]; exists {
			taskTime.Hours += hours
			taskTime.Minutes += minutes
		} else {
			var task models.Task
			config.DB.First(&task, log.TaskID)
			taskTimeMap[log.TaskID] = &models.TaskTime{
				TaskID:  log.TaskID,
				Title:   task.Title,
				Hours:   hours,
				Minutes: minutes,
			}
		}
	}

	taskTimes := make([]models.TaskTime, 0, len(taskTimeMap))
	for _, taskTime := range taskTimeMap {
		// Convert minutes to hours if minutes >= 60
		taskTime.Hours += taskTime.Minutes / 60
		taskTime.Minutes = taskTime.Minutes % 60
		taskTimes = append(taskTimes, *taskTime)
	}

	// Sort taskTimes by hours and minutes in descending order
	sort.Slice(taskTimes, func(i, j int) bool {
		if taskTimes[i].Hours == taskTimes[j].Hours {
			return taskTimes[i].Minutes > taskTimes[j].Minutes
		}
		return taskTimes[i].Hours > taskTimes[j].Hours
	})

	c.JSON(http.StatusOK, taskTimes)
}
