package middlewares

import (
	"net/http"
	"strings"

	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/gin-gonic/gin"
)

func RegisterValidationMiddleware(c *gin.Context) {
	var RegisterRequestBody domain.
		CreateUserDetails
	var response domain.
		CreateUserResponse
	var errors []string
	response.Success = false
	err := c.ShouldBindJSON(&RegisterRequestBody)
	if err != nil {

		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, &response)
		c.Abort()
		return
	}
	names := strings.Split(RegisterRequestBody.FullName, " ")
	if len(names) < 2 {
		errors = append(errors, "Fullname should contain at least two names")
	}
	if len(RegisterRequestBody.Password) < 8 {
		errors = append(errors, "Password should be at least 8 characters")
	}
	if len(errors) > 0 {
		response.Errors = errors
		c.JSON(http.StatusBadRequest, &response)
		c.Abort()
		return
	}
	c.Set("reg", &RegisterRequestBody)
	c.Next()
}

func LoginValidationMiddleware(c *gin.Context) {
	var LoginRequestBody domain.
		LoginUserDetails
	var response domain.
		LoginResponse
	var errors []string
	response.Success = false

	err := c.ShouldBindJSON(&LoginRequestBody)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, &response)
		c.Abort()
		return
	}
	if len(LoginRequestBody.Password) < 8 {
		errors = append(errors, "Password should be at least 8 characters")
	}
	if len(errors) > 0 {
		response.Errors = errors
		c.JSON(http.StatusBadRequest, &response)
		c.Abort()
		return
	}
	c.Set("login", &LoginRequestBody)
	c.Next()
}
