package appanalytics

import (
	"github.com/gofiber/fiber/v3"
	transporttcp "github.com/tudorhulban/analytics77/infra/transport-tcp"
)

type App struct {
	transportHTTP *fiber.App
	transportTCP  *transporttcp.TransportTCP
}

func NewAppAnalytics() (*App, error) {
	return &App{
			transportHTTP: fiber.New(
				fiber.Config{
					BodyLimit: 1 * 1024 * 1024, // in mb
				},
			),
		},
		nil
}
