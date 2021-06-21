package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors" // for TraceBack

	. "github.com/rehearsal-open/rehearsal/util"
)

//
type MissionConfig struct {
	currentDir string
	configFile string
}

func BuildMissionConfig() (*MissionConfig, error) {

	numArgs := len(os.Args)

	// find config file
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	configFile, err := filepath.Abs(os.Args[1])
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !FileExist(configFile) {
		return nil, errors.New("Cannot find config file: " + configFile)
	}

	missionConf := MissionConfig{
		currentDir: currentDir,
		configFile: configFile,
	}

	// read config file
	_, result, err := MarkUpFileLoad(configFile)
	if err != nil {
		return nil, errors.WithMessage(err, "Cannot read config file: "+configFile)
	}

	err = missionConf.parseConfigMap(result)
	if err != nil {
		return nil, errors.WithMessage(err, "Cannot understand config file: "+configFile)
	}

	for iArgs := 2; iArgs < numArgs; iArgs++ {
		switch os.Args[iArgs] {
		// Todo: Overwrite Settings with command arguments
		}
	}

	return &missionConf, nil
}

func (conf *MissionConfig) parseConfigMap(mapping map[string]interface{}) error {

	// todo: parsing Config
	return nil
}
