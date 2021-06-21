package main

import (
	"fmt"
	"log"

	"github.com/rehearsal-open/rehearsal/config"
)

func main() {
	conf, err := config.BuildMissionConfig()
	if conf == nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println(conf)
	}

}
