package supercontroller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishanmadhav/supernetes/supercache"
)

func (s *SuperController) DeploymentRoutes() {
	s.app.Post("/update/deployments", s.UpdateDeployments)
}

func (s *SuperController) UpdateDeployments(c *fiber.Ctx) error {
	var deployments supercache.Deployments
	err := c.BodyParser(&deployments)
	if err != nil {
		return c.JSON(err)
	}
	s.state.Deployments = deployments.Deployments
	return c.JSON(deployments)
}
