package main

import (
	"log"

	engine "github.com/rehearsal-open/rehearsal/engine/v0"

	"github.com/rehearsal-open/rehearsal/parser"
)

func main() {
	rehearsalEngine := engine.RehearsalEngine{}

	if config, err := parser.Load(); err != nil {
		log.Fatal(err.Error())
	} else if err := rehearsalEngine.AssignConfig(config); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println(rehearsalEngine.Conf)
		// engine.Execute()
	}
}
