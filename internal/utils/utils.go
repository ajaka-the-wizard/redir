package utils

import (
	"context"
	"crypto/sha256"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/domain"
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

func SetAndGetCookieDetails(n string, v string, s bool) *http.Cookie {
	cookie := http.Cookie{
		Name:     n,
		Value:    v,
		HttpOnly: true,
		Secure:   s,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
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
	if l, e := c.Get("logger"); e {
		return l.(*slog.Logger)
	}
	return slog.Default()
}
