package config

type Config struct {
	Logger     LoggerConf
	Storage    StorageConf
	HTTPServer HTTPServerConf
	GRPCServer GRPCServerConf
}

func NewConfig() Config {
	return Config{}
}
