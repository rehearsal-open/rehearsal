package entity

import (
	cli "github.com/rehearsal-open/rehearsal/entity/cli-color"
)

type Color cli.CliColor

// Global configuration
type Conf struct {
	Dir    string              `yaml:"directory" json:"directory"`
	SyncMs int                 `yaml:"syncms" json:"syncms"`
	Tasks  map[string]TaskConf `yaml:"tasks" json:"tasks"`
}

// Each task configuration
type TaskConf struct {
	Kind     string   `yaml:"kind" json:"kind"`
	Name     string   `yaml:"-" json:"-"`
	Path     string   `yaml:"execPath" json:"execPath"`
	Args     []string `yaml:"args,omitempty" json:"args,omitempty"`
	ColorStr string   `yaml:"color,omitempty" json:"color,omitempty"`
	SendTo   []string `yaml:"sendTo,omitempty" json:"sendTo,omitEmpty"`
	SyncMs   int      `yaml:"syncms" json:"syncms"`
	Color    Color    `yaml:"-" json:"-"`
}
