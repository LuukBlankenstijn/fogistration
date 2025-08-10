package config

type HttpConfig struct {
	AppEnv   string `env_default:"development"`
	Host     string `env_default:"localhost"`
	Port     string `env_default:"8080"`
	DB       DatabaseConfig
	LogLevel string `env_default:"info"`
}

type GrpcConfig struct {
	AppEnv   string `env_default:"development"`
	Host     string `env_default:"localhost"`
	Port     string `env_default:"9090"`
	DB       DatabaseConfig
	LogLevel string `env_default:"info"`
}

type DomJudgeConfig struct {
	AppEnv       string `env_default:"development"`
	DJHost       string `env_default:"localhost"`
	DB           DatabaseConfig
	LogLevel     string `env_default:"info"`
	Username     string
	Password     string
	SyncInterval string `env_default:"15m"`
}

type ClientConfig struct {
	Server   string
	LogLevel string `env_default:"info"`
	DryRun   bool   `env_default:"true"`
}

type DatabaseConfig struct {
	Host     string `env_default:"localhost"`
	Port     string `env_default:"5432"`
	User     string `env_default:"postgres"`
	Password string `env_default:"postgres"`
	Name     string `env_default:"fogistration"`
	SSLMode  string `env_default:"disable"`
}
