package api

type PodSpec struct {
	Selector string `json:"selector"`
	Image    string `json:"image"`
	Replicas uint   `json:"replicas"`
	Port     string `json:"port"`
	Action   string `json:"action"`
}
