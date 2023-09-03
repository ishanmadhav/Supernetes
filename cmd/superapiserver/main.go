package main

import "github.com/ishanmadhav/supernetes/superapiserver"

func main() {
	s, err := superapiserver.NewSuperAPIServer()
	if err != nil {
		panic(err)
	}
	s.Run()
}
