package api

type Job struct {
	Name   string `json:"name"`
	IsCron bool   `json:"is_cron"`
}
