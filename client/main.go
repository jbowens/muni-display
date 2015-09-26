package main

import "github.com/octavore/naga/service"

func main() {
	s := service.New(&Module{})
	if err := s.RunCommand("start"); err != nil {
		panic(err)
	}
}
