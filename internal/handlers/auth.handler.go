package handlers

import (
	"log"
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
		var response domain.
			CreateUserResponse
		response.Message = "Something went wrong"
		response.Success = false

		val, _ := c.Get("reg")
		RegisterRequestBody, ok := val.(*domain.
			CreateUserDetails)
		if !ok {
			log.Println("from ok")
			log.Println("Val: ", val)
			log.Println("cast: ", RegisterRequestBody)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(RegisterRequestBody.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("from hash")
			c.JSON(http.StatusInternalServerError, &response)
			return
		}
		RegisterRequestBody.Password = string(hash)
		err = repository.CreateUser(pool, RegisterRequestBody, cfg)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				response.Message = "Email is already registered"
				c.JSON(http.StatusConflict, &response)
				return
			}
			log.Println("from create")
			c.JSON(http.StatusInternalServerError, &response)
			return
		}
		response.Success = true
		response.Message = "User registered successfully, please proceed to login"
		c.JSON(http.StatusCreated, &response)
	}
}

func HandleLogin(pool *pgxpool.Pool, cfg *configs.EnvData, mmap *memory.AuthMemoryMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		var response domain.
			LoginResponse
		response.Message = "Something went wrong"
		response.Success = false

		val, _ := c.Get("login")
		LoginRequestBody, ok := val.(*domain.
			LoginUserDetails)
		if !ok {
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		user, err := repository.GetUserByEmail(pool, cfg, LoginRequestBody.Email)
		if err != nil {
			response.Message = "Invalid credentials"
			c.JSON(http.StatusUnauthorized, &response)
			return
		}
		if !user.Verified {
			response.Message = "Please verify email before logging in"
			c.JSON(http.StatusForbidden, &response)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginRequestBody.Password))
		if err != nil {
			response.Message = "Invalid credentials"
			c.JSON(http.StatusUnauthorized, &response)
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
		cookie := utils.SetAndGetCookieDetails(id, cfg.ENVIROMENT == "production")
		c.SetCookieData(cookie)
		c.JSON(http.StatusOK, &response)
	}
}
