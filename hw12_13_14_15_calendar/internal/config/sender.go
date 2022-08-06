package config

type SenderConfig struct {
	Logger  LoggerConf
	Storage StorageConf
	Queue   QueueConf
}

func NewSenderConfig() SenderConfig {
	return SenderConfig{}
}
