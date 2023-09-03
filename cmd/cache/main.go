package main

import "github.com/ishanmadhav/supernetes/supercache"

func main() {
	s := supercache.NewSuperCacheServer()
	s.Start()
}
