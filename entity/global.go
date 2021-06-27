package entity

import (
	color "github.com/rehearsal-open/rehearsal/entity/cli-color"
)

// Global configuration
type Conf struct {
	ConfigDir  string               `yaml:"-" json:"-"`
	Dir        string               `yaml:"directory" json:"directory"`
	SyncMs     int                  `yaml:"syncms" json:"syncms"`
	Tasks      []TaskConf           `yaml:"tasks" json:"tasks"`
	MaxNameLen int                  `yaml:"-" json:"-"`
	TasksMap   map[string]*TaskConf `yaml:"-" json:"-"`
}

// Each task configuration
type TaskConf struct {
	Name     string         `yaml:"name" json:"name"`
	Type     string         `yaml:"type" json:"type"`
	Path     string         `yaml:"execPath" json:"execPath"`
	Args     []string       `yaml:"args,omitempty" json:"args,omitempty"`
	ColorStr string         `yaml:"color,omitempty" json:"color,omitempty"`
	SendTo   []string       `yaml:"sendTo,omitempty" json:"sendTo,omitEmpty"`
	SyncMs   int            `yaml:"syncms" json:"syncms"`
	Color    color.CliColor `yaml:"-" json:"-"`
}

const (
	CLITask string = "CLI"
)
