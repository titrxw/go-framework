package config

type App struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	Env  string `mapstructure:"env" json:"env" yaml:"env"`
	Lang string `mapstructure:"lang" json:"lang" yaml:"lang"`
}
