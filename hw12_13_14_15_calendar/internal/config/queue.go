package config

type QueueConf struct {
	Type         string
	Dsn          string
	Exchange     string
	ExchangeType string
	Queue        string
}
