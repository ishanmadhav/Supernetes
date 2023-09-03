package superlet

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ishanmadhav/supernetes/api"
)

func (s *Superlet) GetAllDockerContainers(c *fiber.Ctx) error {
	containers, err := s.service.GetAllContainers()
	if err != nil {
		return c.JSON(err)
	}
	fmt.Print(err)
	return c.JSON(containers)

}

func (s *Superlet) CreatePod(c *fiber.Ctx) error {
	// someSpec := api.PodSpec{
	// 	Selector: "demoContainer",
	// 	Image:    "bfirsh/reticulate-splines",
	// 	Replicas: 3,
	// }
	var newSpec api.PodSpec
	err := c.BodyParser(&newSpec)
	if err != nil {
		return c.JSON(err)
	}
	err = s.service.PodSource(newSpec)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(newSpec)
}

func (s *Superlet) GetPods(c *fiber.Ctx) error {
	pods := s.service.GetAllPods()
	return c.JSON(pods)
}

func (s *Superlet) GetPodById(c *fiber.Ctx) error {
	return c.JSON(s.service.GetPodById(c.Params("id")))
}

func (s *Superlet) GetPodsBySelector(c *fiber.Ctx) error {
	return c.JSON(s.service.GetPodsBySelector(c.Params("selector")))
}

func (s *Superlet) DeletePodById(c *fiber.Ctx) error {
	err := s.service.DeletePodById(c.Params("id"))
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON("Pod Deleted")
}

func (s *Superlet) DeleteAllPods(c *fiber.Ctx) error {
	err := s.service.DeleteAllPods()
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON("All Pods Deleted")
}

func (s *Superlet) DeletePodsBySelector(c *fiber.Ctx) error {
	fmt.Println("Delete Pods route was hit")
	var spec api.PodSpec
	err := c.BodyParser(&spec)
	if err != nil {
		return c.JSON(err)
	}

	err = s.service.PodSource(spec)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(spec)
}

func (s *Superlet) RunPodBySelector(c *fiber.Ctx) error {
	err := s.service.RunPodBySelector(c.Params("selector"))
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON("Pods Running")
}

func (s *Superlet) RunAllPods(c *fiber.Ctx) error {
	err := s.service.RunAllPods()
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON("All Pods Running")
}
