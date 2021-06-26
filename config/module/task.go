package module

type TaskConfig struct {
	Name  string   `yaml:"name" json:"name"`
	Path  string   `yaml:"execPath" json:"execPath"`
	Args  []string `yaml:"args,omitempty" json:"args,omitempty"`
	Color string   `yaml:"color,omitempty" json:"color,omitempty"`
}
