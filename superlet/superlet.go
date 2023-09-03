package superlet

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ishanmadhav/supernetes/internals/constants"
	"github.com/ishanmadhav/supernetes/superlet/podservice"
)

type Superlet struct {
	service *podservice.PodService
	app     *fiber.App
}

type SuperletInterface struct {
}

func NewSuperlet() (*Superlet, error) {
	service, err := podservice.NewPodService()
	if err != nil {
		log.Fatal(err)
	}
	return &Superlet{service: service, app: fiber.New()}, nil
}

func (s *Superlet) Run() error {
	go s.service.PodSink()
	go s.service.CheckAllPodsStatus()
	s.app.Get("/docker/containers", s.GetAllDockerContainers)
	s.app.Post("/pod", s.CreatePod)
	s.app.Get("/pods", s.GetPods)
	s.app.Get("/pods/:id", s.GetPodById)
	s.app.Get("/pods/selector/:selector", s.GetPodsBySelector)
	s.app.Delete("/pod/:id", s.DeletePodById)
	s.app.Delete("/pods", s.DeleteAllPods)
	s.app.Delete("/pod/select/:selector", s.DeletePodsBySelector)
	s.app.Get("/run/pod/:selector", s.RunPodBySelector)
	s.app.Get("/run/pods", s.RunAllPods)

	err := s.app.Listen(constants.SUPERLET_PORT)
	if err != nil {
		return err
	}
	return nil
}
