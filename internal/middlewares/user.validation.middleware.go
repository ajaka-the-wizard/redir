package middlewares

import (
	"net/http"
	"strings"

	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
)

func RegisterValidationMiddleware(c *gin.Context) {
	logger := utils.GetLogger(c)
	logger.Info("Validating user registration details")
	var RegisterRequestBody domain.
		CreateUserDetails
	var response domain.
		CreateUserResponse
	var errors []string
	response.Success = false
	err := c.ShouldBindJSON(&RegisterRequestBody)
	if err != nil {
		logger.Warn("User provided bad data on registration")
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, &response)
		c.Abort()
		return
	}
	trimmedName := strings.TrimSpace(RegisterRequestBody.FullName)
	names := strings.Fields(trimmedName)
	if len(names) < 2 {
		errors = append(errors, "Fullname should contain at least two names")
	}
	if len(RegisterRequestBody.Password) < 8 {
		errors = append(errors, "Password should be at least 8 characters")
	}
	if len(errors) > 0 {
		logger.Warn("Validating user registration details failed")
		response.Errors = errors
		c.JSON(http.StatusBadRequest, &response)
		c.Abort()
		return
	}
	c.Set("reg", &RegisterRequestBody)
	c.Next()
}

func LoginValidationMiddleware(c *gin.Context) {
	logger := utils.GetLogger(c)
	logger.Info("Validating user login details")
	var LoginRequestBody domain.
		LoginUserDetails
	var response domain.
		LoginResponse
	var errors []string
	response.Success = false

	err := c.ShouldBindJSON(&LoginRequestBody)
	if err != nil {
		logger.Warn("User provided bad data on login")
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, &response)
		c.Abort()
		return
	}
	if len(LoginRequestBody.Password) < 8 {
		errors = append(errors, "Password should be at least 8 characters")
	}
	if len(errors) > 0 {
		logger.Warn("Validating user login details failed")
		response.Errors = errors
		c.JSON(http.StatusBadRequest, &response)
		c.Abort()
		return
	}
	c.Set("login", &LoginRequestBody)
	c.Next()
}

func ProductValidationMiddleware(c *gin.Context) {
	var request domain.CreateProductDetails
	logger := utils.GetLogger(c)
	logger.Info("Validating user product creation details")
	err := c.ShouldBindJSON(&request)
	if err != nil {
		logger.Warn("User provided bad data on product creation ")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}
	c.Set("product", &request)
	c.Next()
}
