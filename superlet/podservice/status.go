package podservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/ishanmadhav/supernetes/internals/constants"
	"github.com/ishanmadhav/supernetes/supercache"
)

func (p *PodService) GetPodsStatusBySelector(selector string) uint {
	pods := p.GetPodsBySelector(selector)
	runningPodCount := uint(0)
	for _, pod := range pods {
		isRunning, err := pod.IsPodRunning(p.client)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Is Running")
		fmt.Println(isRunning)
		if !isRunning {
			p.DeletePodById(pod.ID)
		} else {
			runningPodCount++
		}
	}
	return runningPodCount
}

// Checks the status of all the pods
// When some pod status changes from running to stop
// we ping the deployment controller about the same
// and then we remove that pod from our list
// the deployment controller will then send a new request to create a new pod
// NOTE: This function might cause lots of issues without a mutex
// find out why its crashing on addition of mutex and fix this on priority basis
func (p *PodService) CheckAllPodsStatus() {
	for {
		if len(p.pipeline) > 0 {
			time.Sleep(3 * time.Second)
			continue
		}
		p.mutex.Lock()
		deployments, err := fetchDeployments()
		if err != nil {
			fmt.Println(err)
		}

		for ind, deployment := range deployments.Deployments {
			selector := deployment.Selector
			cnt := p.GetPodsStatusBySelector(selector)
			deployments.Deployments[ind].Replicas = cnt
			fmt.Println(deployments.Deployments[ind])
		}
		fmt.Print(deployments)
		err = updateDeployments(deployments)
		if err != nil {
			fmt.Println(err)
		}
		p.mutex.Unlock()
		time.Sleep(10 * time.Second)

	}
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
	var resp supercache.Deployments
	err := requests.
		URL("http://localhost" + constants.SUPERCONTROLLER_PORT + "/update/deployments").
		BodyJSON(&deployments).
		ToJSON(&resp).
		Fetch(context.Background())
	if err != nil {
		return err
	}
	return nil
}
