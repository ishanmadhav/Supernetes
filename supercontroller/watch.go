package supercontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/ishanmadhav/supernetes/api"
	"github.com/ishanmadhav/supernetes/internals/constants"
)

func (s *SuperController) Watch() {
	go s.WatchDeployments()
	go s.WatchServices()
}

// Whenever we register a difference between desired state and actual state
// we need to update the pods in the superlet for the node
// now this update can be one of the following:
// 1. Create a new pod
// 2. Delete an existing pod
// 3. Update an existing pod
// 4. Increase the number of replicas of the pod
// 5. Stop a running pod
// 6. Restart a pod
// 7. Update the image of a pod
// We can create something like a job for things like restarting, deletion, image updation, etc.
// Could use a message queue instead of the cache for that.
// Or for consistency use the cache and
func (s *SuperController) WatchDeployments() {
	for {
		deployments, err := fetchDeployments()
		fmt.Println(deployments)
		if err != nil {
			fmt.Print(err)
		}

		for _, deployment := range deployments.Deployments {
			isDeploymentPresent, ind := deploymentExists(deployment, s.state.Deployments)
			if isDeploymentPresent {
				fmt.Println("Deployment is present")
				if !areDeploymentsSame(deployment, s.state.Deployments[ind]) {
					//update the deployment
					prev := s.state.Deployments[ind]
					s.state.Deployments[ind] = deployment

					if prev.Image != deployment.Image || prev.Port != deployment.Port || prev.Selector != deployment.Selector {
						//Create a job and add to queue
						//This queue could be a channel or rabbitmq queue
						//Instead of adding job to a queue, we could also just create a different kind of pod spec
						//Since, we already have something akin to a queue running on the superlet, we can technically just use to it dish out actions
						//Spec will have Action="UPDATE"
					} else if prev.Replicas != deployment.Replicas && prev.Image == deployment.Image && prev.Port == deployment.Port && prev.Selector == deployment.Selector {
						if prev.Replicas < deployment.Replicas {
							//Create more pods for this selector
							diff := deployment.Replicas - prev.Replicas
							spec := api.PodSpec{
								Selector: deployment.Selector,
								Image:    deployment.Image,
								Replicas: diff,
								Port:     deployment.Port,
								Action:   "CREATE",
							}
							var resp api.PodSpec
							//Send a HTTP request for pod creation
							err := requests.
								URL("http://localhost" + constants.SUPERLET_PORT + "/pod").
								Method("POST").
								BodyJSON(&spec).
								ToJSON(&resp).
								Fetch(context.Background())
							if err != nil {
								fmt.Println(err)
							}
						} else {
							//Delete some of the pods for this selector
							diff := prev.Replicas - deployment.Replicas
							spec := api.PodSpec{
								Selector: deployment.Selector,
								Image:    deployment.Image,
								Replicas: diff,
								Port:     deployment.Port,
								Action:   "DELETE",
							}
							var resp api.PodSpec
							//Send a HTTP request for pod deletion
							err := requests.
								URL("http://localhost" + constants.SUPERLET_PORT + "/pod/select/" + spec.Selector).
								Method("DELETE").
								BodyJSON(&spec).
								ToJSON(&resp).
								Fetch(context.Background())
							if err != nil {
								fmt.Println(err)
							}
						}

					} else {
						fmt.Println("Deployment are same")
					}
				}
			} else {
				//create the deployment
				fmt.Print("Creating deployment")
				s.state.Deployments = append(s.state.Deployments, deployment)
				spec := api.PodSpec{
					Selector: deployment.Selector,
					Image:    deployment.Image,
					Replicas: deployment.Replicas,
					Port:     deployment.Port,
					Action:   "CREATE",
				}
				var resp api.PodSpec
				//Send a HTTP request for pod creation
				err := requests.
					URL("http://localhost" + constants.SUPERLET_PORT + "/pod").
					Method("POST").
					BodyJSON(&spec).
					ToJSON(&resp).
					Fetch(context.Background())
				if err != nil {
					fmt.Println(err)
				}
			}
			//check if the deployment has changed
			//if it has changed, update the deployment
			//else continue
		}
		time.Sleep(10 * time.Second)
	}

}

// Adds the required number of pods
func (s *SuperController) AddPods(deployment api.Deployment, count uint) error {
	return nil
}

// Removes the given number of pods
func (s *SuperController) RemovePods(deployment api.Deployment, count uint) error {
	return nil
}

func (s *SuperController) WatchServices() {
	for {
		fmt.Println("Services")
		time.Sleep(10 * time.Second)
	}
}

func deploymentExists(deployment api.Deployment, deployments []api.Deployment) (bool, int) {
	for i, d := range deployments {
		if d.Name == deployment.Name {
			return true, i
		}
	}
	return false, -1
}

func areDeploymentsSame(deploymentA api.Deployment, deploymentB api.Deployment) bool {
	return deploymentA.Name == deploymentB.Name &&
		deploymentA.Image == deploymentB.Image &&
		deploymentA.Replicas == deploymentB.Replicas
}
