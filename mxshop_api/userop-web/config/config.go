package config

type GoodsSrvConfig struct {
	Name string `mapstructure:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type ConsulConfig struct {
	Host string `mapstructuer:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name          string         `mapstructure:"name"`
	Address       string         `mapstructure:"address"`
	Port          int            `structure:"port"`
	Tags          []string       `mapstructure:"tags"`
	GoodsInfo     GoodsSrvConfig `mapstructure:"goods_info"`
	JWTInfo       JWTConfig      `mapstructure:"jwt"`
	ConsulInfo    ConsulConfig   `mapstructure:"consul"`
	UserOpSrvInfo GoodsSrvConfig `mapstructure:"userop_info"`
}

type NacosConfig struct {
	Host        string `mapstructuer:"host"`
	Port        uint64 `mapstructure:"port"`
	NamespaceId string `mapstructure:"namespace_id"`
	DataId      string `mapstructure:"data_id"`
	Group       string `mapstructure:"group"`
}
