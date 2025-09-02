package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lumosolutions/yerd/server/internal/config"
	"github.com/lumosolutions/yerd/server/internal/core/services"
)

type PhpListResponse struct {
	Error    bool              `json:"error"`
	Message  string            `json:"message"`
	Versions []*config.PhpInfo `json:"versions"`
	LogFile  string            `json:"log_file"`
}

func PhpGetList(c *fiber.Ctx) error {
	logger, err := services.NewLogger("php-list")
	if err != nil {
		return fmt.Errorf("unable to create logger, %v", err)
	}

	logger.Info("list", "Getting Versions")

	versions, err := ListInstalledPhpVersions()

	if err != nil {
		return c.JSON(PhpListResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	logger.Close()

	return c.JSON(PhpListResponse{
		Error:    false,
		Versions: versions,
		LogFile:  logger.GetLogFile(),
	})
}

func ListInstalledPhpVersions() ([]*config.PhpInfo, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load configuration, %v", err)
	}

	if cfg.Php == nil {
		return nil, nil
	}

	phpVersions := make([]*config.PhpInfo, 0, len(cfg.Php))
	for _, phpInfo := range cfg.Php {
		info := phpInfo
		phpVersions = append(phpVersions, info)
	}

	return phpVersions, nil
}
