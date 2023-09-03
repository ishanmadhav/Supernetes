package supercache

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/ishanmadhav/supernetes/api"
	"github.com/ishanmadhav/supernetes/internals/constants"
)

type SuperCacheServer struct {
	app *fiber.App
}

type Body struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BooleanResponse struct {
	Exists bool `json:"exists"`
}

type Response struct {
	Message string `json:"message"`
}

func NewSuperCacheServer() *SuperCacheServer {
	return &SuperCacheServer{app: fiber.New()}
}

func (s *SuperCacheServer) Start() {
	cache, err := NewSuperCache()
	if err != nil {
		panic(err)
	}
	err = Hydrate(cache)
	if err != nil {
		panic(err)
	}

	s.app.Get("/get/", func(c *fiber.Ctx) error {
		var body Body
		err := c.BodyParser(&body)
		if err != nil {
			return c.JSON(err)
		}
		str, err := cache.Get(body.Key)
		if err != nil {
			return c.JSON(err)
		}
		msg := Response{Message: str}
		return c.JSON(msg)
	})

	s.app.Post("/set", func(c *fiber.Ctx) error {
		var body Body
		err := c.BodyParser(&body)
		if err != nil {
			return c.JSON(err)
		}
		err = cache.Set(body.Key, []byte(body.Value))
		if err != nil {
			return c.JSON(err)
		}
		msg := Response{Message: "Successfully set key"}
		return c.JSON(msg)
	})

	s.app.Delete("/delete", func(c *fiber.Ctx) error {
		var body Body
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}
		exists := cache.Exists(body.Key)
		if !exists {
			msg := Response{Message: "Key does not exist"}
			return c.JSON(msg)
		}

		err = cache.Delete(body.Key)
		if err != nil {
			return c.JSON(err)
		}
		msg := Response{Message: "Successfully deleted key"}
		return c.JSON(msg)
	})

	s.app.Get("/exists", func(c *fiber.Ctx) error {
		var body Body
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}
		exists := cache.Exists(body.Key)
		return c.JSON(BooleanResponse{Exists: exists})
	})

	s.app.Post("/reset", func(c *fiber.Ctx) error {
		err := cache.Reset()
		if err != nil {
			return c.JSON(err)
		}
		msg := Response{Message: "Successfully reset cache"}
		return c.JSON(msg)
	})

	s.app.Get("/len", func(c *fiber.Ctx) error {
		return c.JSON(cache.Len())
	})

	s.app.Listen(constants.SUPERCACHE_PORT)
}

type Deployments struct {
	Deployments []api.Deployment `json:"deployments"`
}

type Services struct {
	Services []api.Service `json:"services"`
}

func Hydrate(c *SuperCache) error {
	deployments := Deployments{Deployments: []api.Deployment{}}
	str, err := json.Marshal(deployments)
	if err != nil {
		return err
	}
	err = c.Set(constants.DEPLOYMENTS_KEY, str)
	if err != nil {
		return err
	}

	services := Services{Services: []api.Service{}}
	str, err = json.Marshal(services)
	if err != nil {
		return err
	}
	err = c.Set(constants.SERVICES_KEY, str)
	if err != nil {
		return err
	}

	return nil
}
