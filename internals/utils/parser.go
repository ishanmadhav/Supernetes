package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ishanmadhav/supernetes/api"
)

func ParseDeploymentFileJSON(fileName string) (api.Deployment, error) {
	jsonData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return api.Deployment{}, err
	}

	var deployment api.Deployment
	err = json.Unmarshal(jsonData, &deployment)
	if err != nil {
		return api.Deployment{}, err
	}

	return deployment, nil
}
