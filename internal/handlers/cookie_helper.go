package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func setCookieFromHttpCookie(c *gin.Context, cookie *http.Cookie) {
	maxAge := int(time.Until(cookie.Expires).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	c.SetCookie(cookie.Name, cookie.Value, maxAge, cookie.Path, "", cookie.Secure, cookie.HttpOnly)
}
