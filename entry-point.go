package main

import (
	"fmt"
	"log"

	v0 "github.com/rehearsal-open/rehearsal/engine/v0"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/load"
)

func main() {
	engine := v0.RehearsalEngine{}
	if en, err := load.Load(); err != nil {
		log.Fatal(err)
	} else {
		switch en.Command {
		case "run":
			if err := engine.AssignConfig(en); err != nil {
				log.Fatal(err)
			} else if err := engine.Run(); err != nil {
				log.Fatal(err)
			}
		case "about":
			fmt.Println(entity.AboutInfomation)
		}
	}
	engine.Finalize()
}
