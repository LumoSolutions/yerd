package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lumosolutions/yerd/server/internal/core/http/output"
	"github.com/lumosolutions/yerd/server/internal/core/manager"
	"github.com/lumosolutions/yerd/server/internal/core/services"
)

type InstallRequest struct {
	Version string `json:"version"`
}

func PhpInstall(c *fiber.Ctx) error {
	var req InstallRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to parse request: " + err.Error()})
	}
	
	version := req.Version

	return output.UseStream(c, func(writer *output.StreamWriter) error {
		if version == "" {
			writer.WriteError("Version parameter is required")
			return nil
		}

		logger, err := services.NewLogger("php-install")
		if err != nil {
			writer.WriteError("Failed to create logger: %s", err.Error())
			return nil
		}

		i, err := manager.NewPhpManager(version, writer, logger)
		if err != nil {
			writer.WriteError("Failed to PHP Manager: %s", err.Error())
			return nil
		}

		err = i.Install()
		if err != nil {
			writer.WriteError("Failed to install PHP: %s", err.Error())
			return nil
		}

		return nil
	})
}
