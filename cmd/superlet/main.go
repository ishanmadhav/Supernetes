package main

import (
	"log"

	"github.com/ishanmadhav/supernetes/superlet"
)

func main() {
	s, err := superlet.NewSuperlet()
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
