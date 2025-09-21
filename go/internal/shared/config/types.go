package config

type HttpConfig struct {
	AppEnv       string `env_default:"production"`
	Host         string `env_default:"localhost"`
	Secret       string `env_default:"super_secret"`
	Port         string `env_default:"8080"`
	DB           DatabaseConfig
	LogLevel     string `env_default:"info"`
	WallpaperDir string `env_default:"/wallpapers"`
	OIDC         *OIDCConfig
}

type GrpcConfig struct {
	AppEnv       string `env_default:"development"`
	Host         string `env_default:"localhost"`
	Port         string `env_default:"9090"`
	DB           DatabaseConfig
	LogLevel     string `env_default:"info"`
	WallpaperDir string `env_default:"/wallpapers"`
	DJHost       string `env_default:"localhost"`
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
	Server            string
	LogLevel          string `env_default:"info"`
	DryRun            bool   `env_default:"true"`
	GreeterConfigPath string `env_default:"/etc/lightdm/lightdm-qt5-greeter.conf"`
}

type DatabaseConfig struct {
	Host     string `env_default:"localhost"`
	Port     string `env_default:"5432"`
	User     string `env_default:"postgres"`
	Password string `env_default:"postgres"`
	Name     string `env_default:"fogistration"`
	SSLMode  string `env_default:"disable"`
}

type OIDCConfig struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}
