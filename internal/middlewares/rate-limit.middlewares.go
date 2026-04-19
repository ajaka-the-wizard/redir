package middlewares

import (
	"net/http"
	"strings"
	"time"

	rl "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

type Limiter struct{}

func (l *Limiter) getStore(limit uint, rate time.Duration) *rl.Store {
	store := rl.InMemoryStore(&rl.InMemoryOptions{
		Rate:  rate,
		Limit: limit,
	})
	return &store
}

func InitRateLimiter() *Limiter {
	return &Limiter{}
}

func (l *Limiter) GetLimiterForAuth(limit_per_rate uint) gin.HandlerFunc {
	store := l.getStore(limit_per_rate, time.Minute)
	mv := rl.RateLimiter(*store, &rl.Options{
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
		ErrorHandler: rlerror,
	})
	return mv
}

func (l *Limiter) GetLimiterForProductAndUser(limit_per_rate uint) gin.HandlerFunc {
	store := l.getStore(limit_per_rate, time.Minute)
	mv := rl.RateLimiter(*store, &rl.Options{
		KeyFunc: func(c *gin.Context) string {
			sessionId, _ := c.Cookie("sessionId")
			return sessionId
		},
		ErrorHandler: rlerror,
	})
	return mv
}

func (l *Limiter) GetLimiterForClient(limit_per_rate uint) gin.HandlerFunc {
	store := l.getStore(limit_per_rate, time.Minute)
	mv := rl.RateLimiter(*store, &rl.Options{
		KeyFunc: func(c *gin.Context) string {
			api_key := c.GetHeader("Authorization")
			api_key = strings.TrimPrefix(api_key, "Bearer ")
			return api_key
		},
		ErrorHandler: rlerror,
	})
	return mv
}

var RL = InitRateLimiter()

func rlerror(c *gin.Context, i rl.Info) {
	retryafter := int(time.Until(i.ResetTime).Seconds())
	if retryafter < 0 {
		retryafter = 0
	}
	c.JSON(http.StatusTooManyRequests, gin.H{"success": false, "message": "Too many attempts, please try again later", "retryafter": retryafter})
	c.Abort()
}
