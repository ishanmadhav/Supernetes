package api

type Deployment struct {
	Name     string `json:"name"`
	Selector string `json:"selector"`
	Replicas uint   `json:"replicas"`
	Image    string `json:"image"`
	Port     string `json:"port"`
}
