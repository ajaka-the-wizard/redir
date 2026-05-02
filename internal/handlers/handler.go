package handlers

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/models"
)

func HydrateMedias(cfg *configs.EnvData, m []models.Media) {
	for i := range m {
		m[i].PublicKey = cfg.DATA_GET_PATH + m[i].PublicKey
	}
}
