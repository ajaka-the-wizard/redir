package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
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
	return hex.EncodeToString(id[:])
}

func GenUUID() uuid.UUID {
	return uuid.New()
}

func SetAndGetCookieDetails(v string, s bool) *http.Cookie {
	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    v,
		HttpOnly: true,
		Secure:   s,
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
	sd := hex.EncodeToString(d)
	return sd, nil
}
func VerifyMultipStepHash(k string, h string) error {
	hash := HashPrivateKey(k)
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(hash))
	if err != nil {
		return err
	}
	return nil
}
