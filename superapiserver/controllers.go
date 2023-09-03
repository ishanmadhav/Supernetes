package superapiserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/ishanmadhav/supernetes/api"
	"github.com/ishanmadhav/supernetes/internals/constants"
	"github.com/ishanmadhav/supernetes/supercache"
)

//Deployment Controlers

func (s *SuperAPIServer) GetDeploymentByName(c *fiber.Ctx) error {
	return c.JSON("Deployment")
}

func (s *SuperAPIServer) GetPodsBySelector(c *fiber.Ctx) error {
	return c.JSON("Pods")
}

func (s *SuperAPIServer) GetAllPods(c *fiber.Ctx) error {
	return c.JSON("All Pods")
}

func (s *SuperAPIServer) DeleteAllPods(c *fiber.Ctx) error {
	return c.JSON("Delete all pods")
}

// Create a deployement in the cache
// Supercontroller will monitor the deployments key in the cache
// using some sort of loop
// the alternative to this would have been the use of a message queue, but we would have to transmit the
// crashing and creation of new pods as well
func (s *SuperAPIServer) CreateDeployment(c *fiber.Ctx) error {
	var deployment api.Deployment
	err := c.BodyParser(&deployment)
	if err != nil {
		return c.JSON(err)
	}
	err = s.ProcessDeployment(deployment)
	return c.JSON(deployment)
}

func (s *SuperAPIServer) ProcessDeployment(deployment api.Deployment) error {

	deployments, err := fetchDeployments()
	if err != nil {
		return err
	}
	fmt.Print(deployments)

	//check if deployment exists, if it does exist, update it
	if deploymentExists(deployment, deployments) {
		fmt.Print("Deployment exists")
		deployments = updateDeployment(deployment, deployments)
	} else {
		//if it doesn't exist, create it
		fmt.Print("Deployment does not exist")
		deployments = createDeployment(deployment, deployments)
	}

	err = updateDeployments(deployments)
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

func updateDeployments(deployments supercache.Deployments) error {
	var reqBody supercache.Body
	var resp supercache.Response
	reqBody.Key = "deployments"
	str, err := json.Marshal(deployments)
	if err != nil {
		return err
	}
	reqBody.Value = string(str)
	err = requests.
		URL("http://localhost" + constants.SUPERCACHE_PORT + "/set").
		BodyJSON(&reqBody).
		ToJSON(&resp).
		Fetch(context.Background())
	return nil
}

func deploymentExists(deployment api.Deployment, deployments supercache.Deployments) bool {
	for _, d := range deployments.Deployments {
		if d.Name == deployment.Name {
			return true
		}
	}
	return false
}

func updateDeployment(deployment api.Deployment, deployments supercache.Deployments) supercache.Deployments {
	for i, d := range deployments.Deployments {
		if d.Name == deployment.Name {
			deployments.Deployments[i] = deployment
		}
	}
	return deployments
}

func createDeployment(deployment api.Deployment, deployments supercache.Deployments) supercache.Deployments {
	deployments.Deployments = append(deployments.Deployments, deployment)
	return deployments
}
