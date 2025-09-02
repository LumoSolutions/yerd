package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/lumosolutions/yerd/server/internal/config"
	"github.com/lumosolutions/yerd/server/internal/core/http"
	"github.com/lumosolutions/yerd/server/internal/core/services"
)

func main() {
	if err := services.Boot(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := config.WriteConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "YERD",
		AppName:       "YERD Service API v2.0.0",
	})

	app.Get("/api/php/list", http.PhpGetList)
	app.Post("/api/php/install", http.PhpInstall)

	listenOn := fmt.Sprintf("127.0.0.1:%d", cfg.Yerd.Port)
	fmt.Printf("Listening on %s\n", listenOn)

	if err = app.Listen(listenOn); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
