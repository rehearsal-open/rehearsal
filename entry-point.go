package main

import (
	"log"

	v0 "github.com/rehearsal-open/rehearsal/engine/v0"
	"github.com/rehearsal-open/rehearsal/load"
)

func main() {
	engine := v0.RehearsalEngine{}
	if entity, err := load.Load(); err != nil {
		log.Fatal(err)
	} else if err := engine.AssignConfig(entity); err != nil {
		log.Fatal(err)
	}
	engine.Finalize()
}
