package controllers

import (
	"em-test/config"
	"em-test/models"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Получение списка всех пользователей
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /users [get]
func GetUsersHandler(c *gin.Context) {
    var users []models.User

		name := c.Query("name")
    surname := c.Query("surname")
    address := c.Query("address")
		pageStr := c.Query("page")
		pageSizeStr := c.Query("page_size")

		page := 1
		pageSize := 10

		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		if pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
				pageSize = ps
			}
		}

		query := config.DB

		if name != "" {
        query = query.Where("name ILIKE ?", "%"+name+"%")
    }
    if surname != "" {
        query = query.Where("surname ILIKE ?", "%"+surname+"%")
    }
    if address != "" {
        query = query.Where("address ILIKE ?", "%"+address+"%")
    }

		offset := (page - 1) * pageSize
		result := query.Limit(pageSize).Offset(offset).Find(&users)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

// Получение пользователя по id
// @Summary Get user by ID
// @Description Get a single user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [get]
func GetUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user models.User 

	result := config.DB.First(&user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		}
	}

	c.JSON(http.StatusOK, user)
}

// Создание нового пользователя
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User JSON"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [post]
func CreateUserHandler(c *gin.Context, validate *validator.Validate) {
    var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
        return
    }

		if err := validate.Struct(&user); err != nil {
        validationErrors := err.(validator.ValidationErrors)
        errors := make([]string, len(validationErrors))
        for i, fieldError := range validationErrors {
            errors[i] = fieldError.Error()
        }
        c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
        return
    }

    result := config.DB.Create(&user)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}

// Удаление пользователя
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	result := config.DB.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		}
		return
	}

	result = config.DB.Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}


// Изменение данных пользователя
// @Summary Update a user
// @Description Update user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User data"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [put]
func UpdateUserHandler(c *gin.Context, validate *validator.Validate) {
	id := c.Param("id")
	var user models.User

	result := config.DB.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		}
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if err := validate.Struct(&user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make([]string, len(validationErrors))
		for i, fieldError := range validationErrors {
			errors[i] = fieldError.Error()
		}
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}

	result = config.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
