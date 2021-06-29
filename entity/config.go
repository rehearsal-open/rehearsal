package entity

type Config struct {
	Command  string       `yaml:"-"`
	YamlPath string       `yaml:"-"`
	Dir      string       `yaml:"directory"`
	SyncMs   int          `yaml:"syncms"`
	TaskConf []TaskConfig `yaml:"tasks"`
}

type TaskConfig struct {
	Name     string   `yaml:"name"`
	Type     string   `yaml:"type"`
	ExecPath string   `yaml:"execPath"`
	Args     []string `yaml:"args"`
	ShowOut  bool     `yaml:"show"`
}
