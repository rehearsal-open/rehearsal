package main

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
	// "github.com/rehearsal-open/rehearsal/runner"
)

func LoadRunner(path string) (runner.Runner, error) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var data map[string]interface{}

	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err)
	}

	// Todo: Convert from map to run
}
