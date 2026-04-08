package utils

import (
	"context"
	"crypto/sha256"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateContextWithStatedTime(seconds string) (context.Context, context.CancelFunc) {
	IntSecond, err := strconv.Atoi(seconds)
	if err != nil {
		IntSecond = 2
	}
	return context.WithTimeout(context.Background(), (time.Duration(IntSecond))*time.Second)
}
func GenCleanedUpUUid() string {
	id := GenUUID()
	return strings.ReplaceAll(id, "-", "")
}

func GenUUID() string {
	return uuid.New().String()
}

func SetAndGetCookieDetails(n string, v string, s bool, exp time.Time) *http.Cookie {
	cookie := http.Cookie{
		Name:     n,
		Value:    v,
		HttpOnly: true,
		Secure:   s,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  exp,
	}
	return &cookie
}

func GetUser(c *gin.Context) (*domain.
	LightUser, bool) {
	val, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	user, ok := val.(*domain.
		LightUser)
	return user, ok
}

func GetId(c *gin.Context) (int, bool) {
	val, exists := c.Get("id")
	if !exists {
		return 0, false
	}
	id, ok := val.(int)
	return id, ok
}

func GeneratePrivateKey() string {
	prefix := "rp_live_"
	id1 := GenCleanedUpUUid()
	id2 := GenCleanedUpUUid()
	key := prefix + id1 + id2
	return key
}

func HashPrivateKey(k string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(k))
	hash := hasher.Sum(nil)
	return hash
}

func PerformMultiStepHash(k string) (string, error) {
	hash := HashPrivateKey(k)
	d, err := bcrypt.GenerateFromPassword(hash, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	sd := string(d)
	return sd, nil
}
func VerifyMultipStepHash(k string, h string) error {
	hash := HashPrivateKey(k)
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(hash))
}

func GetLogger(c *gin.Context) *slog.Logger {
	if l, ok := c.Get("logger"); ok {
		if logger, ok := l.(*slog.Logger); ok && logger != nil {
			return logger
		}
	}
	return slog.Default()
}

func PerformLoginActivity(mmap *memory.AuthMemoryMap, cfg *configs.EnvData, user *domain.LightUser) (*http.Cookie, *http.Cookie) {
	id := GenCleanedUpUUid()
	lastAccessedTime := mmap.SetUserOnline(id, user)
	exp := time.Now().Add(24 * time.Hour)
	sessionIdCookie := SetAndGetCookieDetails("sessionId", id, cfg.PRODUCTION, exp)
	lastAccessTimeCookie := SetAndGetCookieDetails("lastUpdateTime", StringifyTime(lastAccessedTime), cfg.PRODUCTION, exp)
	return sessionIdCookie, lastAccessTimeCookie
}

func GetTimeFromCookie(c *gin.Context) (time.Time, error) {
	cookie, err := c.Cookie("lastUpdateTime")
	if err != nil {
		return time.Time{}, err
	}
	unixTime, err := strconv.ParseInt(cookie, 10, 64)
	if err != nil {
		return time.Time{}, nil
	}
	return time.Unix(unixTime, 0), nil
}

func StringifyTime(t time.Time) string {
	str := strconv.FormatInt(t.Unix(), 10)
	return str
}
