package module

type GlobalConfig struct {
	Dir   string       `yaml:"directory" json:"directory"`
	Tasks []TaskConfig `yaml:"tasks" json:"tasks"`
}
