package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ajaka-the-wizard/redir/internal/cache"
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/database"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/routes"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type authResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func setupAuthRouter(t *testing.T) (*gin.Engine, *cache.Sredis, *pgxpool.Pool, *configs.EnvData) {
	t.Helper()
	err := godotenv.Load()
	if err != nil {
		panic("Env file not found")
	}

	cfg := &configs.EnvData{
		DATABASE_URL:   os.Getenv("DATABASE_URL"),
		REDIS_ADDR:     os.Getenv("REDIS_ADDR"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
		PRODUCTION:     false,
	}

	if cfg.DATABASE_URL == "" {
		t.Fatal("DATABASE_URL environment variable must be set for integration tests")
	}
	if cfg.REDIS_ADDR == "" {
		t.Fatal("REDIS_ADDR environment variable must be set for integration tests")
	}

	logger := slog.Default()
	ctx := context.Background()

	rdb := cache.InitializeRedis(ctx, cfg, logger)
	t.Cleanup(func() {
		rdb.Clean()
	})

	pool := database.ConnectDB(ctx, logger, cfg.DATABASE_URL)
	t.Cleanup(func() {
		pool.Close()
	})

	repo := repository.InitializeRepository(pool)
	store := store.InitializeStore(rdb, repo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.GenAndAttachRequestIdMiddleware())
	router.Use(middlewares.AttachLoggerToContext())
	router.Use(middlewares.PerformBasicRequestCycleCalculations())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	routes.AuthRoutes(v1, cfg, store)

	return router, rdb, pool, cfg
}

func TestRegisterRouteIntegration(t *testing.T) {
	router, _, _, _ := setupAuthRouter(t)

	email := fmt.Sprintf("test-register-%s@example.com", uuid.NewString())
	body, err := json.Marshal(map[string]string{
		"full_name": "Integration Test User",
		"email":     email,
		"password":  "password123",
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var resp authResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("expected success=true, got false, body=%s", rec.Body.String())
	}
	if resp.Message == "" {
		t.Fatal("expected non-empty success message")
	}
}
