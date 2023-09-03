package cmdapi

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/ishanmadhav/supernetes/api"
	"github.com/ishanmadhav/supernetes/internals/constants"
)

func CreateDeploymentAPI(deployment api.Deployment) error {
	var resp api.Deployment
	err := requests.
		URL("http://localhost" + constants.SUPERAPI_PORT + "/deployment").
		BodyJSON(&deployment).
		ToJSON(&resp).
		Fetch(context.Background())

	if err != nil {
		return err
	}

	fmt.Print(resp)
	return nil
}
