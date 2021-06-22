package main

import (
	"fmt"
	"log"

	"github.com/rehearsal-open/rehearsal/config/mission"
)

func main() {
	conf, err := mission.BuildMissionConfig()
	if conf == nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println(conf)
	}

}
