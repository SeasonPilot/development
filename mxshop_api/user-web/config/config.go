package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type AliSmsConfig struct {
	ApiKey    string `mapstructure:"key"`
	ApiSecret string `mapstructure:"secret"`
}

type RedisConfig struct {
	Host   string `mapstructuer:"host"`
	Port   int    `mapstructure:"port"`
	Expire int    `mapstruture:"expire"`
}

type ServerConfig struct {
	Name       string        `mapstructure:"name"`
	Port       int           `structure:"port"`
	UserInfo   UserSrvConfig `mapstructure:"user-srv"`
	JWTInfo    JWTConfig     `mapstructure:"jwt"`
	AliSmsInfo AliSmsConfig  `mapstructure:"sms"`
	RedisInfo  RedisConfig   `mapstructure:"redis"`
}
