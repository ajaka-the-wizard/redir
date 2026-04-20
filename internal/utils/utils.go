package utils

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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
	val, ok := c.Get("user")
	if !ok {
		return nil, false
	}
	user, ok := val.(*domain.
		LightUser)
	return user, ok
}

func GetMedia(c *gin.Context) (*models.Media, bool) {
	val, ok := c.Get("media")
	if !ok {
		return nil, false
	}
	media, ok := val.(*models.Media)
	return media, ok
}

func GetId(c *gin.Context) (int, bool) {
	val, ok := c.Get("id")
	if !ok {
		return 0, false
	}
	id, ok := val.(int)
	return id, ok
}

func GetProduct(c *gin.Context) (*models.Product, bool) {
	if p, ok := c.Get("product"); ok {
		product, ok := p.(*models.Product)
		return product, ok
	}
	return nil, false
}

func genTwoUUIDsSeparatedBySomething(sep string) string {
	id1 := GenCleanedUpUUid()
	id2 := GenCleanedUpUUid()
	return id1 + sep + id2
}

func GeneratePrivateKey() string {
	prefix := "rp_live_"
	uuids := genTwoUUIDsSeparatedBySomething("")
	key := prefix + uuids
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
	exp := lastAccessedTime.Add(24 * time.Hour)
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
		return time.Time{}, err
	}
	return time.Unix(unixTime, 0), nil
}

func StringifyTime(t time.Time) string {
	str := strconv.FormatInt(t.Unix(), 10)
	return str
}

func GenerateKeyForUpload(cfg *configs.EnvData, productId int) (string, string) {
	return genInnerKey(cfg, productId), genPublicKey(cfg, productId)
}

func genInnerKey(cfg *configs.EnvData, productId int) string {
	id := GenUUID()
	return fmt.Sprintf("%s/%d/%s", cfg.BUCKET_ROOT, productId, id)
}

func genPublicKey(cfg *configs.EnvData, productId int) string {
	uuids := genTwoUUIDsSeparatedBySomething("-")
	return fmt.Sprintf("%s/%d/%s", cfg.DATA_GET_PATH, productId, uuids)
}

func ValidateAndReturnUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func ValidateAssetId(s string) (int, bool) {
	res := strings.Split(s, "/")
	if len(res) < 2 {
		return 0, false
	}
	productId, err := strconv.Atoi(res[0])
	if err != nil {
		return 0, false
	}
	for u := range strings.SplitSeq(res[1], "-") {
		err = uuid.Validate(u)
		if err != nil {
			return 0, false
		}
	}
	return productId, true
}
