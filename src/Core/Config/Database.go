package config

type Database struct {
	Driver   string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DbName   string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	Charset  string `mapstructure:"charset" json:"charset" yaml:"charset"`
	Prefix   string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	PoolSize int    `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
}
