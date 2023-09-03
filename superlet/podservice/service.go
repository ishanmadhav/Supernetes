package podservice

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/ishanmadhav/supernetes/superlet/pod"
)

type PodService struct {
	client   *client.Client
	pods     []*pod.Pod
	pipeline chan pod.PodRecipe
	mutex    *sync.Mutex
}

func NewPodService() (*PodService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &PodService{client: cli, pipeline: make(chan pod.PodRecipe), pods: make([]*pod.Pod, 0), mutex: &sync.Mutex{}}, nil
}

// Create Pod from PodRecipe
// This function is called from the PodSink() function of the pipeline
func (p *PodService) CreatePod(spec pod.PodRecipe) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	fmt.Println("Created pod was hit")
	ctx := context.Background()
	out, err := p.client.ImagePull(ctx, spec.Image, types.ImagePullOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(spec.Port + "/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: spec.Port,
				},
			},
		},
	}

	resp, err := p.client.ContainerCreate(ctx, &container.Config{
		Image: spec.Image,
		ExposedPorts: nat.PortSet{
			nat.Port(spec.Port + "/tcp"): struct{}{},
		},
	}, hostConfig, &network.NetworkingConfig{}, nil, spec.Name)
	if err != nil {
		return err
	}

	fmt.Print(resp)
	// err = p.client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	// if err != nil {
	// 	fmt.Println("Couldn't remove container")
	// 	return err
	// }

	if err := p.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	containerJSON, err := p.GetContainer(resp.ID)
	if err != nil {
		fmt.Println("Couldn't get container which suggests it was created")
		return err
	}
	createdPod := &pod.Pod{
		Name:      spec.Name,
		Selector:  spec.Selector,
		Spec:      spec.Spec,
		Container: containerJSON,
		ID:        resp.ID,
	}

	p.pods = append(p.pods, createdPod)

	return nil
}

// Get Container By ID
func (p *PodService) GetContainer(id string) (*types.ContainerJSON, error) {
	ctx := context.Background()
	container, err := p.client.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}
	return &container, nil
}

// Gets all the docker containers, running or not	running
func (p *PodService) GetAllContainers() ([]types.Container, error) {
	ctx := context.Background()
	containers, err := p.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// Get All Pods in service
func (p *PodService) GetAllPods() []*pod.Pod {
	return p.pods
}

// Get Pods by Selector
func (p *PodService) GetPodsBySelector(selector string) []*pod.Pod {
	pods := make([]*pod.Pod, 0)
	for _, pod := range p.pods {
		if pod.Selector == selector {
			pods = append(pods, pod)
		}
	}
	return pods
}

// Get Pod By ID
func (p *PodService) GetPodById(id string) *pod.Pod {
	for _, pod := range p.pods {
		if pod.ID == id {
			return pod
		}
	}
	return nil
}

// Gets the first pod in the pod list with given selector
func (p *PodService) GetFirstPodWithSelector(selector string) *pod.Pod {
	for _, pod := range p.pods {
		if pod.Selector == selector {
			return pod
		}
	}
	return nil
}

func (p *PodService) GetPodIndex(pod *pod.Pod) int {
	for i, p := range p.pods {
		if p.ID == pod.ID {
			return i
		}
	}
	return -1
}

// Get Pod by Name
func (p *PodService) GetPodByName(name string) *pod.Pod {
	for _, pod := range p.pods {
		if pod.Name == name {
			return pod
		}
	}
	return nil
}

// Delete Pod by ID
// NOTE: Removal from pod list yet to be imolement
func (p *PodService) DeletePodById(id string) error {

	ctx := context.Background()
	err := p.client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Delete All Pods
func (p *PodService) DeleteAllPods() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	ctx := context.Background()
	for _, pod := range p.pods {
		fmt.Println(pod.ID)
		err := p.StopPodByID(pod.ID)
		if err != nil {
			fmt.Print(err)
			return err
		}
		err = p.client.ContainerRemove(ctx, pod.ID, types.ContainerRemoveOptions{})
		if err != nil {
			fmt.Print(err)
			return err
		}
	}
	return nil
}

// Delete specified number of pods by selector
func (p *PodService) DeletePodsBySelector(selector string, count uint) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	pods := p.GetPodsBySelector(selector)
	if len(pods) < int(count) {
		return fmt.Errorf("Not enough pods to delete")
	}

	for i := 0; i < int(count); i++ {
		err := p.StopPodByID(pods[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PodService) DeletePodBySelector(spec pod.PodRecipe) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	fmt.Println("Delete pod by selector was hit")
	pod := p.GetFirstPodWithSelector(spec.Selector)
	if pod == nil {
		return fmt.Errorf("No pod with selector %s found", spec.Selector)
	}
	err := p.StopPodByID(pod.ID)
	if err != nil {
		return err
	}
	err = p.RemovePodByID(pod.ID)
	if err != nil {
		return err
	}
	podIndex := p.GetPodIndex(pod)
	if podIndex == -1 {
		return fmt.Errorf("Pod not found in pod list")
	}
	p.pods = removePod(p.pods, podIndex)
	return nil
}

func removePod(slice []*pod.Pod, index int) []*pod.Pod {
	return append(slice[:index], slice[index+1:]...)
}

// func removeElement(slice []int, index int) []int {
// 	return append(slice[:index], slice[index+1:]...)
// }

// Run Pod by Selector
func (p *PodService) RunPodBySelector(selector string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	pods := p.GetPodsBySelector(selector)
	for _, pod := range pods {
		err := p.client.ContainerStart(context.Background(), pod.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// Run All Pods
func (p *PodService) RunAllPods() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, pod := range p.pods {
		err := p.client.ContainerStart(context.Background(), pod.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop Pod By Selector Name
func (p *PodService) StopPodBySelector(selector string) error {
	pods := p.GetPodsBySelector(selector)
	for _, pod := range pods {
		err := p.client.ContainerStop(context.Background(), pod.ID, container.StopOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop Pod by ID
func (p *PodService) StopPodByID(id string) error {
	err := p.client.ContainerStop(context.Background(), id, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Remove Pod by ID
func (p *PodService) RemovePodByID(id string) error {
	err := p.client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

// List logs of pod by ID
func (p *PodService) ListPodLogsByID(id string) (io.ReadCloser, error) {
	ctx := context.Background()
	return p.client.ContainerLogs(ctx, id, types.ContainerLogsOptions{ShowStdout: true})
}
