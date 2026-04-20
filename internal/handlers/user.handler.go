package handlers

import (
	"fmt"
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUser(pool *pgxpool.Pool, cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := utils.GetUser(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Couldn't identify user"})
			return
		}
		user, err := repository.GetUserById(c.Request.Context(), pool, cfg, u.Id)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": fmt.Sprintf("Couldn't find user with id of %v", u.Id)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "User retrieved successfully", "user": *user})
	}
}
