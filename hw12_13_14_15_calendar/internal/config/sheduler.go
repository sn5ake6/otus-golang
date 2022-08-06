package config

type SchedulerConfig struct {
	Logger    LoggerConf
	Storage   StorageConf
	Queue     QueueConf
	Intervals IntervalsConf
}

type IntervalsConf struct {
	Notify string
	Delete string
}

func NewSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{}
}
