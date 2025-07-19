package database

import (
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
)

func GetUrl(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)
}
