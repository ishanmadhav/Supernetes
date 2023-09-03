package podservice

import (
	"fmt"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/ishanmadhav/supernetes/api"
	"github.com/ishanmadhav/supernetes/superlet/pod"
)

func waitForTime() {
	time.Sleep(10 * time.Second)
	fmt.Println("Waited for 10 seconds")
}

// This will act like a sink/consumer for our pipeline of pod creation
// Instead of a queue, we use a channel that acts like a queue
func (p *PodService) PodSink() {
	for {
		select {
		case newPod := <-p.pipeline:
			fmt.Println(newPod)
			if newPod.Action == "CREATE" {
				p.CreatePod(newPod)
			} else if newPod.Action == "DELETE" {
				fmt.Println("We reached the pod sink for deletion")
				p.DeletePodBySelector(newPod)
			} else if newPod.Action == "UPDATE" {
				//p.UpdatePod(newPod)
			} else if newPod.Action == "RESTART" {
				//p.RestartPod(newPod)
			} else if newPod.Action == "STOP" {
				//p.StopPod(newPod)
			} else if newPod.Action == "RUN" {
				//p.RunPod(newPod)
			} else if newPod.Action == "RUNALL" {
				//p.RunAllPods(newPod)
			} else if newPod.Action == "DELETEALL" {
				//p.DeleteAllPods(newPod)
			} else if newPod.Action == "DELETEBYSELECTOR" {
				//p.DeletePodsBySelector(newPod)
			} else if newPod.Action == "RUNBYSELECTOR" {
				//p.RunPodBySelector(newPod)
			} else {
				fmt.Println("Invalid action")
			}

		}
	}
}

// Acts as source for pod
func (p *PodService) PodSource(spec api.PodSpec) error {
	fmt.Println("PodSource was hit")
	fmt.Print(spec)
	preparedPodList := PodMaker(spec)
	fmt.Print(preparedPodList)
	for _, pod := range preparedPodList {
		p.AddPod(pod)
	}
	return nil
}

//

func PodMaker(spec api.PodSpec) []pod.PodRecipe {

	podList := make([]pod.PodRecipe, 0)

	for i := 0; i < int(spec.Replicas); i++ {
		podName := namesgenerator.GetRandomName(0)
		tempPod := pod.PodRecipe{
			Name:     podName,
			Selector: spec.Selector,
			Image:    spec.Image,
			Spec:     spec,
			Action:   spec.Action,
			Port:     spec.Port,
		}
		podList = append(podList, tempPod)
	}
	return podList
}

// This will add one pod to the pipeline
func (p *PodService) AddPod(newPod pod.PodRecipe) {
	//fmt.Println(newPod)
	p.pipeline <- newPod
}
