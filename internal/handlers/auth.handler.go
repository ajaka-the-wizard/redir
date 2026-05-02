package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func HandleRegister(cfg *configs.EnvData, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := utils.GetLogger(c)
		logger.Info("Beginning registration for user")
		val, _ := c.Get("reg")
		RegisterRequestBody, ok := val.(*domain.
			CreateUserDetails)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			logger.Error("Couldn't get user reg details from context")
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(RegisterRequestBody.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			logger.Error("Couldn't hash user password", "error", err.Error())
			return
		}
		RegisterRequestBody.Password = string(hash)
		err = store.CreateUser(c.Request.Context(), RegisterRequestBody, cfg)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email is already registered"})
				logger.Warn("Email unique constraint conflict", "error", err.Error())
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			logger.Error("Something went wrong", "error", err.Error())
			return
		}
		logger.Info("User registered successfully")
		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "User registered successfully, please proceed to login"})
	}
}

func HandleLogin(cfg *configs.EnvData, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := utils.GetLogger(c)
		logger.Info("Starting user login attempt")

		val, _ := c.Get("login")
		LoginRequestBody, ok := val.(*domain.
			LoginUserDetails)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			logger.Error("Couldn't get user login details from context")
			return
		}
		user, err := store.GetUserByEmail(c.Request.Context(), cfg, LoginRequestBody.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
			logger.Warn("User provided invalid credentials", "error", err.Error())
			return
		}
		if !user.Verified {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Please verify email before logging in"})
			logger.Warn("User login attempt without email verification rejected")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginRequestBody.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
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

		sessionCookie, err := utils.PerformLoginActivity(c.Request.Context(), store, cfg, &lightUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong, couldn't sign you in"})
			return
		}
		c.SetCookieData(sessionCookie)
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "User logged in successfully"})
		logger.Info("User logged in successfully")
	}
}

func HandleLogout(store *store.Store, cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get("sessionId")
		sessionId, _ := val.(string)
		err := store.RevokeUser(c.Request.Context(), sessionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		exp := time.Now().Add(-1 * time.Second)
		sessioncookie := utils.SetAndGetCookieDetails("sessionId", "", cfg.PRODUCTION, exp)
		c.SetCookieData(sessioncookie)
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
		Scopes:       []string{"openid", "email", "profile"},
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

func (g *GoogleOauth) HandleGoogleCallback(cfg *configs.EnvData, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := "google"
		// logger := utils.GetLogger(c)
		state, err := c.Cookie("state")
		returnedState := c.Query("state")
		if err != nil || state != returnedState {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid or expired state token"})
			return
		}
		code := c.Query("code")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		token, err := g.o.Exchange(ctx, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		client := g.o.Client(ctx, token)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to get user info"})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "failed to get user info"})
			return
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "failed to get user info"})
			return
		}
		var user domain.GoogleUser
		if err := json.Unmarshal(data, &user); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !user.VerifiedEmail {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Email provided is unverified"})
			return
		}
		u, err := store.GetUserByProvider(c.Request.Context(), cfg, provider, user.ID)
		if err != nil {
			if err == pgx.ErrNoRows {
				u, err = store.CreateOrLinkOauth(c.Request.Context(), cfg, user.ID, user.Email, user.Name, provider)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
				return
			}
		}
		sessionIdCookie, err := utils.PerformLoginActivity(c.Request.Context(), store, cfg, u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong, couldn't sign you in"})
			return
		}
		c.SetCookieData(sessionIdCookie)
		c.Redirect(http.StatusFound, cfg.CLIENT_DASHBOARD)
	}
}

type GithubOauth struct {
	o *oauth2.Config
}

func InitGithubOauth(cfg *configs.EnvData) *GithubOauth {
	c := oauth2.Config{
		ClientID:     cfg.GITHUB_CLIENT_ID,
		ClientSecret: cfg.GITHUB_CLIENT_SECRET,
		RedirectURL:  cfg.GITHUB_REDIRECT_URL,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
	return &GithubOauth{
		o: &c,
	}
}
func (g *GithubOauth) HandleRedirectToGithub(cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		// logger := utils.GetLogger(c)
		state := utils.GenUUID()
		exp := time.Now().Add(15 * time.Minute)
		cookie := utils.SetAndGetCookieDetails("state", state, cfg.PRODUCTION, exp)
		c.SetCookieData(cookie)
		c.Redirect(http.StatusFound, g.o.AuthCodeURL(state))
	}
}
