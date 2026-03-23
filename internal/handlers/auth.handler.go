package handlers

import (
	"net/http"
	"strings"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(pool *pgxpool.Pool, cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := utils.GetLogger(c)
		logger.Info("Beginning registration for user")
		var response domain.
			CreateUserResponse
		response.Message = "Something went wrong"
		response.Success = false

		val, _ := c.Get("reg")
		RegisterRequestBody, ok := val.(*domain.
			CreateUserDetails)
		if !ok {
			c.JSON(http.StatusInternalServerError, response)
			logger.Error("Couldn't get user reg details from context")
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(RegisterRequestBody.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &response)
			logger.Error("Couldn't hash user password", "error", err.Error())
			return
		}
		RegisterRequestBody.Password = string(hash)
		err = repository.CreateUser(pool, RegisterRequestBody, cfg)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				response.Message = "Email is already registered"
				c.JSON(http.StatusConflict, &response)
				logger.Warn("Email unique contraint conflict", "error", err.Error())
				return
			}
			c.JSON(http.StatusInternalServerError, &response)
			logger.Error("Something went wrong", "error", err.Error())
			return
		}
		response.Success = true
		response.Message = "User registered successfully, please proceed to login"
		logger.Info("User registered successfully")
		c.JSON(http.StatusCreated, &response)
	}
}

func HandleLogin(pool *pgxpool.Pool, cfg *configs.EnvData, mmap *memory.AuthMemoryMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := utils.GetLogger(c)
		logger.Info("Starting user login in attempt")
		var response domain.
			LoginResponse
		response.Message = "Something went wrong"
		response.Success = false

		val, _ := c.Get("login")
		LoginRequestBody, ok := val.(*domain.
			LoginUserDetails)
		if !ok {
			c.JSON(http.StatusInternalServerError, response)
			logger.Error("Couldn't get user login details from context")
			return
		}
		user, err := repository.GetUserByEmail(pool, cfg, LoginRequestBody.Email)
		if err != nil {
			response.Message = "Invalid credentials"
			c.JSON(http.StatusUnauthorized, &response)
			logger.Warn("User provided invalid credentials", "error", err.Error())
			return
		}
		if !user.Verified {
			response.Message = "Please verify email before logging in"
			c.JSON(http.StatusForbidden, &response)
			logger.Warn("User login attempt without email verification rejected")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginRequestBody.Password))
		if err != nil {
			response.Message = "Invalid credentials"
			c.JSON(http.StatusUnauthorized, &response)
			logger.Warn("User provided invalid credentials", "error", err.Error())
			return
		}
		lightUser := domain.
			LightUser{
			Id:    user.Id,
			Email: user.Email,
			Admin: user.Admin,
			Paid:  user.Paid,
		}
		id := utils.GenCleanedUpUUid()
		mmap.SetUserOnline(id, &lightUser)
		response.Success = true
		response.Message = "User logged in successfully"
		cookie := utils.SetAndGetCookieDetails("sessionId", id, cfg.ENVIRONMENT == "production")
		c.SetCookieData(cookie)
		c.JSON(http.StatusOK, &response)
		logger.Info("User logged in successfully")
	}
}

func HandleLogout(mmap *memory.AuthMemoryMap, cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get("sessionId")
		sessionId, _ := val.(string)
		mmap.RevokeUser(sessionId)
		cookie := utils.SetAndGetCookieDetails("sessionId", "", cfg.ENVIRONMENT == "production")
		c.SetCookieData(cookie)
		c.Status(http.StatusNoContent)
	}
}
