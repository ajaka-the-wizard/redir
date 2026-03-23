package handlers

import (
	"fmt"
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUser(pool *pgxpool.Pool, cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		var response domain.
			GetMeResponse
		response.Success = false
		response.Message = "Something went wrong"
		u, ok := utils.GetUser(c)
		if !ok {
			response.Message = "Couldnt identify user"
			c.JSON(http.StatusBadRequest, &response)
			return
		}
		user, err := repository.GetUserById(pool, cfg, u.Id)
		if err != nil {
			if err == pgx.ErrNoRows {
				response.Message = fmt.Sprintf("Couldnt find user with id of %v", u.Id)
				c.JSON(http.StatusNotFound, &response)
				return
			}
			c.JSON(http.StatusInternalServerError, &response)
			return
		}
		response.Success = true
		response.Message = "User retrieved successfully"
		response.User = *user
		c.JSON(http.StatusOK, &response)
	}
}
