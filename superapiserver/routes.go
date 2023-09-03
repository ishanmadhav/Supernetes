package superapiserver

func (s *SuperAPIServer) DeploymentRoutes() {

	s.app.Get("/deployment/:name", s.GetDeploymentByName)
	s.app.Get("/pods/:selector", s.GetPodsBySelector)
	s.app.Get("/pods", s.GetAllPods)
	s.app.Post("/deployment", s.CreateDeployment)
	s.app.Delete("/pods", s.DeleteAllPods)

	// s.app.Put("/deployment/:name", s.EditDeploymentByName)
	// s.app.Delete("/deployment/:name", s.DeleteDeploymentByName)
}
