package config

type Log struct {
	Path    string `mapstructure:"path" json:"path" yaml:"path"`
	Level   string `mapstructure:"level" json:"level" yaml:"level"`
	MaxDays int    `mapstructure:"max_days" json:"max_days" yaml:"max_days"`
	MaxSize int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
}
