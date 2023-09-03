package main

import "github.com/ishanmadhav/supernetes/supercontroller"

func main() {
	s, err := supercontroller.NewSuperController()
	if err != nil {
		panic(err)
	}

	err = s.Start()
	if err != nil {
		panic(err)
	}
}
