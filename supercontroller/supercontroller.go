package supercontroller

import (
	"context"
	"encoding/json"

	"github.com/carlmjohnson/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/ishanmadhav/supernetes/api"
	"github.com/ishanmadhav/supernetes/internals/constants"
	"github.com/ishanmadhav/supernetes/supercache"
)

type NodeState struct {
	Deployments []api.Deployment `json:"deployments"`
	Services    []api.Service    `json:"services"`
	Job         []api.Job        `json:"jobs"`
}

type SuperController struct {
	app   *fiber.App
	state NodeState
}

func NewSuperController() (*SuperController, error) {
	app := fiber.New()
	return &SuperController{app: app}, nil
}

func (s *SuperController) Start() error {
	go s.Watch()
	s.DeploymentRoutes()
	err := s.app.Listen(constants.SUPERCONTROLLER_PORT)
	if err != nil {
		return err
	}
	return nil
}

func fetchDeployments() (supercache.Deployments, error) {
	ctx := context.Background()
	var resp supercache.Response
	err := requests.
		URL("http://localhost" + constants.SUPERCACHE_PORT + "/get").
		Method("GET").
		BodyJSON(&supercache.Body{Key: "deployments"}).
		ToJSON(&resp).
		Fetch(ctx)

	if err != nil {
		return supercache.Deployments{}, err
	}

	var deployments supercache.Deployments
	err = json.Unmarshal([]byte(resp.Message), &deployments)

	if err != nil {
		return supercache.Deployments{}, err
	}
	return deployments, nil
}
