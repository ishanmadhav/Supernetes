package superapiserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishanmadhav/supernetes/internals/constants"
)

type SuperAPIServer struct {
	app *fiber.App
}

func NewSuperAPIServer() (*SuperAPIServer, error) {
	return &SuperAPIServer{app: fiber.New()}, nil
}

func (s *SuperAPIServer) Run() error {
	s.DeploymentRoutes()
	err := s.app.Listen(constants.SUPERAPI_PORT)
	if err != nil {
		return err
	}
	return nil
}
