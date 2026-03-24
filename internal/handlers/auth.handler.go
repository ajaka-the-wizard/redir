package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
		_, err = repository.CreateUser(pool, RegisterRequestBody, cfg)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				response.Message = "Email is already registered"
				c.JSON(http.StatusConflict, &response)
				logger.Warn("Email unique constraint conflict", "error", err.Error())
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
		logger.Info("Starting user login attempt")
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

		cookie := utils.PerformLoginActivity(mmap, cfg, &lightUser)
		c.SetCookieData(cookie)
		response.Success = true
		response.Message = "User logged in successfully"
		c.JSON(http.StatusOK, &response)
		logger.Info("User logged in successfully")
	}
}

func HandleLogout(mmap *memory.AuthMemoryMap, cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get("sessionId")
		sessionId, _ := val.(string)
		mmap.RevokeUser(sessionId)
		exp := time.Now().Add(-1 * time.Second)
		cookie := utils.SetAndGetCookieDetails("sessionId", "", cfg.PRODUCTION, exp)
		c.SetCookieData(cookie)
		c.Status(http.StatusNoContent)
	}
}

type GoogleOauth struct {
	o *oauth2.Config
}

func InitGoogleOauth(cfg *configs.EnvData) *GoogleOauth {
	c := oauth2.Config{
		ClientID:     cfg.GOOGLE_CLIENT_ID,
		ClientSecret: cfg.GOOGLE_CLIENT_SECRET,
		RedirectURL:  cfg.GOOGLE_REDIRECT_URL,
		Scopes:       []string{"https://www.googleapis.com"},
		Endpoint:     google.Endpoint,
	}
	return &GoogleOauth{
		o: &c,
	}
}

func (g *GoogleOauth) HandleRedirectToGoogle(cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		// logger := utils.GetLogger(c)
		state := utils.GenUUID()
		exp := time.Now().Add(15 * time.Minute)
		cookie := utils.SetAndGetCookieDetails("state", state, cfg.PRODUCTION, exp)
		c.SetCookieData(cookie)
		c.Redirect(http.StatusFound, g.o.AuthCodeURL(state))
	}
}

func (g *GoogleOauth) HandleGoogleCallback(pool *pgxpool.Pool, cfg *configs.EnvData, mmap *memory.AuthMemoryMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		// logger := utils.GetLogger(c)
		state, err := c.Cookie("state")
		returnedState := c.Query("state")
		if err != nil || state != returnedState {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid or expired state token"})
			return
		}
		code := c.Query("code")
		token, err := g.o.Exchange(c.Request.Context(), code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		client := g.o.Client(c.Request.Context(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to get user info"})
			return
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		var user domain.GoogleUser
		if err := json.Unmarshal(data, &user); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var cUser domain.CreateUserDetails
		cUser.Email = user.Email
		cUser.Password = ""
		cUser.FullName = user.Name
		newUser, err := repository.CreateUser(pool, &cUser, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		cookie := utils.PerformLoginActivity(mmap, cfg, newUser)
		c.SetCookieData(cookie)
		c.Redirect(http.StatusFound, cfg.CLIENT_DASHBOARD)
	}
}
