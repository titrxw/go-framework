package config

type Config struct {
	App         App                 `mapstructure:"app" json:"app" yaml:"app"`
	DatabaseMap map[string]Database `mapstructure:"database" json:"database" yaml:"database"`
	RedisMap    map[string]Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Session     Session             `mapstructure:"session" json:"session" yaml:"session"`
	Cookie      Cookie              `mapstructure:"cookie" json:"cookie" yaml:"cookie"`
	LogMap      map[string]Log      `mapstructure:"log" json:"log" yaml:"log"`
}
