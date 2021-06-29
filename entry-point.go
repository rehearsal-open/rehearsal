package main

import (
	"log"

	"github.com/rehearsal-open/rehearsal/load"
)

func main() {
	if entity, err := load.Load(); err != nil {
		log.Fatal(err)
	}
}
