package pod

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ishanmadhav/supernetes/api"
)

// Name will be different for every pod. Use some unique name-gen library
// Selector will be the one for which we create replicas
// ID is the container ID inside of docker
type Pod struct {
	Name      string `json:"name"`
	Selector  string `json:"selector"`
	Container *types.ContainerJSON
	ID        string
	Spec      api.PodSpec
}

type PodRecipe struct {
	Name     string `json:"name"`
	Selector string `json:"selector"`
	Image    string `json:"image"`
	Spec     api.PodSpec
	Action   string `json:"action"`
	Port     string `json:"port"`
}

func (p *Pod) IsPodRunning(cli *client.Client) (bool, error) {
	ctx := context.Background()
	containerJSON, err := cli.ContainerInspect(ctx, p.ID)
	if err != nil {
		return false, err
	}

	return containerJSON.State.Running, nil
}
